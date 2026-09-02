package nansen

import "testing"

func TestPointerHelpers(t *testing.T) {
	if got := *StringPtr("value"); got != "value" {
		t.Errorf("StringPtr() = %q", got)
	}
	if got := *IntPtr(42); got != 42 {
		t.Errorf("IntPtr() = %d", got)
	}
	if got := *BoolPtr(true); !got {
		t.Error("BoolPtr() = false")
	}
	if got := *Float64Ptr(1.5); got != 1.5 {
		t.Errorf("Float64Ptr() = %v", got)
	}
}
