package jsonvalue

import (
	"testing"
)

func TestSortArray(t *testing.T) {
	arr := NewArray()

	arr.AppendInt(0).InTheEnd()
	arr.AppendInt(1).InTheEnd()
	arr.AppendInt(2).InTheEnd()
	arr.AppendInt(3).InTheEnd()
	arr.AppendInt(4).InTheEnd()
	arr.AppendInt(5).InTheEnd()
	arr.AppendInt(6).InTheEnd()
	arr.AppendInt(7).InTheEnd()
	arr.AppendInt(8).InTheEnd()
	arr.AppendInt(9).InTheEnd()

	t.Logf("pre-sorted: '%s'", arr.MustMarshalString())

	lessFunc := func(v1, v2 *V) bool {
		return v1.Int() > v2.Int()
	}
	arr.SortArray(lessFunc)

	res := arr.MustMarshalString()
	t.Logf("sorted res: '%s'", res)

	if res != `[9,8,7,6,5,4,3,2,1,0]` {
		t.Errorf("array sort failed")
		return
	}

	return
}

func TestSortArrayError(t *testing.T) {
	// simple test, should not panic
	v := NewInt(1)
	v.SortArray(func(v1, v2 *V) bool { return false })

	v = NewArray()
	v.SortArray(nil)

	return
}
