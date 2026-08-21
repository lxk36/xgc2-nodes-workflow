package signalwait

import (
	"testing"
	"time"

	"github.com/lxk36/xgc2-orchestration-core/kernel/canonicaljson"
	protocol "github.com/lxk36/xgc2-orchestration-core/kernel/node"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
)

func TestSignalWaitRequiresExactRecordedEventResolution(t *testing.T) {
	executor := New()
	registry := protocol.NewRegistry()
	if err := registry.Register(executor); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Seal(); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 13, 11, 0, 0, 0, time.UTC)
	input := map[string]any{"signalId": "stop-experiment"}
	inputDigest, _ := canonicaljson.DigestValue(input)
	request := contracts.NodeInvocationRequest{
		SchemaVersion: protocol.InvocationSchemaVersion,
		InvocationID:  "signal-invocation", RunID: "signal-run", NodeID: "hold",
		TypeRef: executor.Descriptor().TypeRef, DescriptorDigest: executor.Descriptor().DescriptorDigest,
		AttemptID: "signal-attempt", AttemptOrdinal: 1, Input: input, InputDigest: inputDigest,
		RequestedAt: now, Deadline: now.Add(time.Hour),
	}
	waiting, err := registry.Execute(t.Context(), request)
	if err != nil || waiting.Status != contracts.NodeResultWaiting || waiting.Wait == nil ||
		waiting.Wait.Kind != contracts.NodeWaitEvent || waiting.Wait.SubjectRef != "stop-experiment" {
		t.Fatalf("waiting result = %+v err=%v", waiting, err)
	}
	payload := map[string]any{"signalId": "stop-experiment"}
	payloadDigest, _ := canonicaljson.DigestValue(payload)
	resumed, err := registry.Resume(t.Context(), contracts.NodeResumeRequest{
		InvocationID: request.InvocationID, RunID: request.RunID, NodeID: request.NodeID,
		TypeRef: request.TypeRef, DescriptorDigest: request.DescriptorDigest,
		AttemptID: request.AttemptID, AttemptOrdinal: request.AttemptOrdinal,
		Input: input, InputDigest: inputDigest, Wait: *waiting.Wait,
		Resolution: contracts.NodeWaitResolution{
			Kind: waiting.Wait.Kind, SubjectRef: waiting.Wait.SubjectRef,
			ConditionDigest: waiting.Wait.ConditionDigest, Status: contracts.NodeWaitResolvedSucceeded,
			Payload: payload, PayloadDigest: payloadDigest, ObservedAt: now.Add(time.Minute),
		},
		RequestedAt: now.Add(time.Minute),
	})
	if err != nil || resumed.Status != contracts.NodeResultSucceeded || resumed.Output["signalId"] != "stop-experiment" {
		t.Fatalf("resumed result = %+v err=%v", resumed, err)
	}
}
