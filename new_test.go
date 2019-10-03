package jsonvalue

import "testing"

func TestNewString(t *testing.T) {
	s := "你好，世界"
	v := NewString(s)
	if s != v.String() {
		t.Errorf("test NewString Failed")
	}
}

func TestNewBool(t *testing.T) {
	v := NewBool(true)
	if v.Bool() != true {
		t.Errorf("test NewBool - true failed")
		return
	}

	v = NewBool(false)
	if v.Bool() != false {
		t.Errorf("test NewBool - false failed")
		return
	}
}

func TestNewNull(t *testing.T) {
	v := NewNull()
	if false == v.IsNull() {
		t.Errorf("test NewNull failed")
		return
	}
}

func TestNewInteger(t *testing.T) {
	i := int64(-1234567)

	v := NewInt(int(i))
	if v.Int() != int(i) {
		t.Errorf("Test NewInt failed")
		return
	}

	v = NewUint(uint(i))
	if v.Uint() != uint(i) {
		t.Errorf("Test NewUint failed")
		return
	}

	v = NewInt32(int32(i))
	if v.Int32() != int32(i) {
		t.Errorf("Test NewInt32 failed")
		return
	}

	v = NewUint32(uint32(i))
	if v.Uint32() != uint32(i) {
		t.Errorf("Test NewUint32 failed")
		return
	}

	v = NewInt64(int64(i))
	if v.Int64() != int64(i) {
		t.Errorf("Test NewInt64 failed")
		return
	}

	v = NewUint64(uint64(i))
	if v.Uint64() != uint64(i) {
		t.Errorf("Test NewUint64 failed")
		return
	}
}

func TestNewFloat(t *testing.T) {
	s := "3.1415926535"
	f := 3.1415926535

	v := NewFloat64(f, -1)
	if v.String() != s {
		t.Errorf("Test NewFloat64 failed")
		return
	}

	v = NewFloat64(f, 2)
	if v.String() != "3.14" {
		t.Errorf("Test NewFloat64 failed")
		return
	}

	s = "3.1415927"
	v = NewFloat32(float32(f), -1)
	if v.String() != s {
		t.Errorf("Test NewFloat32 failed: %s", v.String())
		return
	}

	v = NewFloat32(float32(f), 5)
	if v.String() != "3.14159" {
		t.Errorf("Test NewFloat32 failed")
		return
	}
}

func TestEmptyObject(t *testing.T) {
	v := NewObject()
	b, _ := v.Marshal()
	if string(b) != "{}" {
		t.Errorf("Test NewObject failed")
		return
	}
}

func TestEmptyArray(t *testing.T) {
	v := NewArray()
	b, _ := v.Marshal()
	if string(b) != "[]" {
		t.Errorf("TestNewObject failed")
		return
	}
}
