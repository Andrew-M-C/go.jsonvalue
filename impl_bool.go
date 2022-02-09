package jsonvalue

import (
	"bytes"
	"fmt"
)

var (
	trueBytes  = []byte{'t', 'r', 'u', 'e'}
	falseBytes = []byte{'f', 'a', 'l', 's', 'e'}
)

type boolValue bool

// NewBool returns an initialied boolean jsonvalue object
//
// NewBool 用给定的 bool 返回一个初始化好的布尔类型的 jsonvalue 值
func NewBool(b bool) *V {
	return &V{
		impl: boolValue(b),
	}
}

// ======== deleter interface ========

func (v boolValue) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	return ErrNotFound
}

// ======== typper interface ========

func (v boolValue) ValueType() ValueType {
	return Boolean
}

// ======== getter interface ========

func (v boolValue) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return &V{}, ErrNotFound
}

// ======== setter interface ========

func (v boolValue) setAt(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Set()", v.ValueType())
}

// ======== iterater interface ========

func (v boolValue) RangeObjects(callback func(k string, v *V) bool) {
	// do nothing
}

func (v boolValue) RangeArray(callback func(i int, v *V) bool) {
	// do nothing
}

func (v boolValue) ForRangeObj() map[string]*V {
	return map[string]*V{}
}

func (v boolValue) ForRangeArr() []*V {
	return nil
}

func (v boolValue) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter)
	close(ch)
	return ch
}

func (v boolValue) IterArray() <-chan *ArrayIter {
	ch := make(chan *ArrayIter)
	close(ch)
	return ch
}

// ======== marshaler interface ========

func (v boolValue) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	if bool(v) {
		buf.Write(trueBytes)
	} else {
		buf.Write(falseBytes)
	}
	return nil
}

// ======== valuer interface ========

func (v boolValue) Bool() (bool, error) {
	return bool(v), nil
}

func (v boolValue) Int64() (int64, error) {
	return 0, ErrTypeNotMatch
}

func (v boolValue) Uint64() (uint64, error) {
	return 0, ErrTypeNotMatch
}

func (v boolValue) Float64() (float64, error) {
	return 0, ErrTypeNotMatch
}

func (v boolValue) String() string {
	if bool(v) {
		return "true"
	}
	return "false"
}

func (v boolValue) Len() int {
	return 0
}

// ======== inserter interface ========

func (v boolValue) insertBefore(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

func (v boolValue) insertAfter(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

// ======= appender interface ========

func (v boolValue) appendInTheBeginning(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}

func (v boolValue) appendInTheEnd(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}
