package actioncall

import (
	"testing"

	protocol "github.com/lxk36/xgc2-orchestration-core/kernel/node"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
)

func TestDescriptorIsASealedKernelOwnedChildActionControl(t *testing.T) {
	executor := New()
	descriptor := executor.Descriptor()
	if descriptor.TypeRef != "xgc.workflow.action-call/v1" ||
		descriptor.Mode != contracts.NodeWaiting ||
		descriptor.Determinism != contracts.NodeRecorded ||
		descriptor.SchemaMode != contracts.NodeSchemaCallAction ||
		len(descriptor.AllowedEffectKinds) != 0 || descriptor.CompensationTypeRef != "" {
		t.Fatalf("descriptor = %+v", descriptor)
	}
	registry := protocol.NewRegistry()
	if err := registry.Register(executor); err != nil {
		t.Fatal(err)
	}
	if digest, err := registry.Seal(); err != nil || !contracts.ValidDigest(digest) {
		t.Fatalf("catalog digest=%q err=%v", digest, err)
	}
}
