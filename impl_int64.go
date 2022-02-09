package jsonvalue

import (
	"bytes"
	"fmt"
	"strconv"
)

type int64Value int64

// NewInt64 returns an initialied num jsonvalue object by int64 type
//
// NewInt64 用给定的 int64 返回一个初始化好的数字类型的 jsonvalue 值
func NewInt64(i int64) *V {
	if i > 0 {
		return NewUint64(uint64(i))
	}
	return &V{
		impl: int64Value(i),
	}
}

// NewInt32 returns an initialied num jsonvalue object by int32 type
//
// NewInt32 用给定的 int32 返回一个初始化好的数字类型的 jsonvalue 值
func NewInt32(i int32) *V {
	return NewInt64(int64(i))
}

// NewInt returns an initialied num jsonvalue object by int type
//
// NewInt 用给定的 int 返回一个初始化好的数字类型的 jsonvalue 值
func NewInt(i int) *V {
	return NewInt64(int64(i))
}

// ======== deleter interface ========

func (v int64Value) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	return ErrNotFound
}

// ======== typper interface ========

func (v int64Value) ValueType() ValueType {
	return Number
}

// ======== getter interface ========

func (v int64Value) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return &V{}, ErrNotFound
}

// ======== setter interface ========

func (v int64Value) setAt(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Set()", v.ValueType())
}

// ======== iterater interface ========

func (v int64Value) RangeObjects(callback func(k string, v *V) bool) {
	// do nothing
}

func (v int64Value) RangeArray(callback func(i int, v *V) bool) {
	// do nothing
}

func (v int64Value) ForRangeObj() map[string]*V {
	return map[string]*V{}
}

func (v int64Value) ForRangeArr() []*V {
	return nil
}

func (v int64Value) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter)
	close(ch)
	return ch
}

func (v int64Value) IterArray() <-chan *ArrayIter {
	ch := make(chan *ArrayIter)
	close(ch)
	return ch
}

// ======== marshaler interface ========

func (v int64Value) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	buf.WriteString(strconv.FormatInt(int64(v), 10))
	return nil
}

// ======== valuer interface ========

func (v int64Value) Bool() (bool, error) {
	return v != 0, ErrTypeNotMatch
}

func (v int64Value) Int64() (int64, error) {
	return int64(v), nil
}

func (v int64Value) Uint64() (uint64, error) {
	return uint64(v), nil
}

func (v int64Value) Float64() (float64, error) {
	return float64(v), nil
}

func (v int64Value) String() string {
	return strconv.FormatUint(uint64(v), 10)
}

func (v int64Value) Len() int {
	return 0
}

// ======== numberAsserter interface ========

func (v int64Value) IsFloat() bool {
	return false
}

func (v int64Value) IsInteger() bool {
	return true
}

func (v int64Value) IsNegative() bool {
	return true
}

func (v int64Value) IsPositive() bool {
	return false
}

func (v int64Value) GreaterThanInt64Max() bool {
	return false
}

// ======== inserter interface ========

func (v int64Value) insertBefore(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

func (v int64Value) insertAfter(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

// ======= appender interface ========

func (v int64Value) appendInTheBeginning(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}

func (v int64Value) appendInTheEnd(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}
