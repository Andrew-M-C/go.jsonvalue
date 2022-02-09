package jsonvalue

import (
	"bytes"
	"fmt"
)

type nullValue int

// NewNull returns an initialied null jsonvalue object
//
// NewNull 返回一个初始化好的 null 类型的 jsonvalue 值
func NewNull() *V {
	return &V{
		impl: nullValue(0),
	}
}

// ======== deleter interface ========

func (v nullValue) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	return ErrNotFound
}

// ======== typper interface ========

func (v nullValue) ValueType() ValueType {
	return Null
}

// ======== getter interface ========

func (v nullValue) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return &V{}, ErrNotFound
}

// ======== setter interface ========

func (v nullValue) setAt(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Set()", v.ValueType())
}

// ======== iterater interface ========

func (v nullValue) RangeObjects(callback func(k string, v *V) bool) {
	// do nothing
}

func (v nullValue) RangeArray(callback func(i int, v *V) bool) {
	// do nothing
}

func (v nullValue) ForRangeObj() map[string]*V {
	return map[string]*V{}
}

func (v nullValue) ForRangeArr() []*V {
	return nil
}

func (v nullValue) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter)
	close(ch)
	return ch
}

func (v nullValue) IterArray() <-chan *ArrayIter {
	ch := make(chan *ArrayIter)
	close(ch)
	return ch
}

//  ======== marshaler interface ========

func (v nullValue) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	buf.WriteString("null")
	return nil
}

// ======== valuer interface ========

func (v nullValue) Bool() (bool, error) {
	return false, ErrTypeNotMatch
}

func (v nullValue) Int64() (int64, error) {
	return 0, ErrTypeNotMatch
}

func (v nullValue) Uint64() (uint64, error) {
	return 0, ErrTypeNotMatch
}

func (v nullValue) Float64() (float64, error) {
	return 0, ErrTypeNotMatch
}

func (v nullValue) String() string {
	return "null"
}

func (v nullValue) Len() int {
	return 0
}

// ======== inserter interface ========

func (v nullValue) insertBefore(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

func (v nullValue) insertAfter(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

// ======= appender interface ========

func (v nullValue) appendInTheBeginning(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}

func (v nullValue) appendInTheEnd(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}
