// Package signalwait publishes a product-neutral durable signal boundary.
// Execute records an event wait; only the orchestration controller's exact
// external-wait ingress can supply the matching observation and invoke Resume.
package signalwait

import (
	"context"
	"errors"

	"github.com/lxk36/xgc2-orchestration-core/kernel/canonicaljson"
	protocol "github.com/lxk36/xgc2-orchestration-core/kernel/node"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
)

const packageDigest = "sha256:d66a909c6570eaf553950e819d991e5306c872de294b86ca7a18b214f58a2ec0"

type Executor struct{ descriptor contracts.NodeDescriptor }

func New() *Executor {
	stringSchema := contracts.Schema{Type: contracts.TypeString}
	schema := contracts.Schema{
		Type:       contracts.TypeObject,
		Properties: map[string]contracts.Schema{"signalId": stringSchema},
		Required:   []string{"signalId"},
	}
	descriptor := contracts.NodeDescriptor{
		SchemaVersion: protocol.DescriptorSchemaVersion,
		TypeRef:       "xgc.workflow.signal-wait/v1", DisplayName: "Durable workflow signal wait",
		PackageRef: "xgc2-nodes-workflow", PackageDigest: packageDigest,
		InputSchema: schema, OutputSchema: schema,
		Mode: contracts.NodeWaiting, Determinism: contracts.NodeRecorded,
		MaxInputBytes: 65536, MaxOutputBytes: 65536,
	}
	descriptor.DescriptorDigest, _ = protocol.DescriptorDigest(descriptor)
	return &Executor{descriptor: descriptor}
}

func (executor *Executor) Descriptor() contracts.NodeDescriptor { return executor.descriptor }

func (executor *Executor) Execute(_ context.Context, request contracts.NodeInvocationRequest) (contracts.NodeResult, error) {
	signalID, _ := request.Input["signalId"].(string)
	if !contracts.ValidIdentifier(signalID) {
		return contracts.NodeResult{}, errors.New("signal wait requires a portable signal identity")
	}
	conditionDigest, err := canonicaljson.DigestValue(map[string]any{"signalId": signalID})
	if err != nil {
		return contracts.NodeResult{}, err
	}
	evidenceDigest, err := canonicaljson.DigestValue(map[string]any{
		"kind": contracts.NodeWaitEvent, "signalId": signalID, "conditionDigest": conditionDigest,
	})
	if err != nil {
		return contracts.NodeResult{}, err
	}
	return contracts.NodeResult{
		Status: contracts.NodeResultWaiting,
		Wait: &contracts.NodeWait{
			Kind: contracts.NodeWaitEvent, SubjectRef: signalID, ConditionDigest: conditionDigest,
		},
		EvidenceDigest: evidenceDigest,
	}, nil
}

func (executor *Executor) Resume(_ context.Context, request contracts.NodeResumeRequest) (contracts.NodeResult, error) {
	signalID, _ := request.Input["signalId"].(string)
	resolvedSignal, _ := request.Resolution.Payload["signalId"].(string)
	if request.Resolution.Kind != contracts.NodeWaitEvent ||
		request.Resolution.Status != contracts.NodeWaitResolvedSucceeded ||
		request.Resolution.SubjectRef != signalID || resolvedSignal != signalID {
		return contracts.NodeResult{}, errors.New("signal observation does not match the exact durable wait")
	}
	output := map[string]any{"signalId": signalID}
	outputDigest, err := canonicaljson.DigestValue(output)
	if err != nil {
		return contracts.NodeResult{}, err
	}
	evidenceDigest, err := canonicaljson.DigestValue(map[string]any{
		"outputDigest": outputDigest, "resolutionDigest": request.Resolution.PayloadDigest,
	})
	if err != nil {
		return contracts.NodeResult{}, err
	}
	return contracts.NodeResult{
		Status: contracts.NodeResultSucceeded, Output: output,
		OutputDigest: outputDigest, EvidenceDigest: evidenceDigest,
	}, nil
}
