// Package actioncall publishes the product-neutral child Action control node.
// The public orchestration controller recognizes its schema mode and owns the
// complete durable child Run lifecycle; Execute must therefore never be used.
package actioncall

import (
	"context"
	"errors"

	protocol "github.com/lxk36/xgc2-orchestration-core/kernel/node"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
)

const packageDigest = "sha256:37b6a36922a60d990d85ee3659013d0a3331b66e6ac7b8fdfe94af4d26056584"

type Executor struct{ descriptor contracts.NodeDescriptor }

func New() *Executor {
	empty := contracts.Schema{Type: contracts.TypeObject, Properties: map[string]contracts.Schema{}}
	descriptor := contracts.NodeDescriptor{
		SchemaVersion:  protocol.DescriptorSchemaVersion,
		TypeRef:        "xgc.workflow.action-call/v1",
		DisplayName:    "Durable child Action",
		PackageRef:     "xgc2-nodes-workflow",
		PackageDigest:  packageDigest,
		InputSchema:    empty,
		OutputSchema:   empty,
		Mode:           contracts.NodeWaiting,
		Determinism:    contracts.NodeRecorded,
		SchemaMode:     contracts.NodeSchemaCallAction,
		MaxInputBytes:  65536,
		MaxOutputBytes: 1048576,
	}
	descriptor.DescriptorDigest, _ = protocol.DescriptorDigest(descriptor)
	return &Executor{descriptor: descriptor}
}

func (executor *Executor) Descriptor() contracts.NodeDescriptor { return executor.descriptor }

func (*Executor) Execute(context.Context, contracts.NodeInvocationRequest) (contracts.NodeResult, error) {
	return contracts.NodeResult{}, errors.New("child Action control nodes are executed only by the orchestration kernel")
}
