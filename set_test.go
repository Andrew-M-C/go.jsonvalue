package jsonvalue

import (
	"testing"
)

func TestSet(t *testing.T) {
	o := NewObject()
	child := NewString("Hello, world!")
	_, err := o.Set(child).At("data", "message", -1, "hello")
	if err != nil {
		t.Errorf("test Set failed: %v", err)
		return
	}

	b, _ := o.Marshal()
	t.Logf("after setting: %v", string(b))
	if string(b) != `{"data":{"message":[{"hello":"Hello, world!"}]}}` {
		t.Errorf("test Set() failed")
	}
	return
}
