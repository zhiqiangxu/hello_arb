package tx

import (
	"sync/atomic"
	"testing"
)

func TestDispatcher(t *testing.T) {
	var pointer atomic.Pointer[chan struct{}]
	old := pointer.Swap(nil)
	if old != nil {
		t.Fatal("old != nil")
	}

	v := pointer.Load()
	if v != nil {
		t.Fatal("v != nil")
	}

	ch := make(chan struct{})
	pointer.Store(&ch)

	old = pointer.Swap(nil)
	if old == nil {
		t.Fatal("old == nil")
	}

	v = pointer.Load()
	if v != nil {
		t.Fatal("v != nil")
	}
}
