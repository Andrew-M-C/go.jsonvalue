package jsonvalue

import (
	"testing"
)

func TestSet(t *testing.T) {
	o := NewObject()
	child := NewString("Hello, world!")
	o.Set("data", "message", -1, "hello", child)

	b, _ := o.Marshal()
	t.Logf("after setting: %v", string(b))
	if string(b) != `{"data":{"message":[{"hello":"Hello, world!"}]}}` {
		t.Errorf("test Set() failed")
	}
	return
}
