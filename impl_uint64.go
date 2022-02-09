package jsonvalue

import (
	"bytes"
	"fmt"
	"strconv"
)

type uint64Value uint64

// NewUint64 returns an initialied num jsonvalue object by uint64 type
//
// NewUint64 用给定的 uint64 返回一个初始化好的数字类型的 jsonvalue 值
func NewUint64(u uint64) *V {
	return &V{
		impl: uint64Value(u),
	}
}

// NewUint returns an initialied num jsonvalue object by uint type
//
// NewUint 用给定的 uint 返回一个初始化好的数字类型的 jsonvalue 值
func NewUint(u uint) *V {
	return NewUint64(uint64(u))
}

// NewUint32 returns an initialied num jsonvalue object by uint32 type
//
// NewUint32 用给定的 uint32 返回一个初始化好的数字类型的 jsonvalue 值
func NewUint32(u uint32) *V {
	return NewUint64(uint64(u))
}

// ======== deleter interface ========

func (v uint64Value) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	return ErrNotFound
}

// ======== typper interface ========

func (v uint64Value) ValueType() ValueType {
	return Number
}

// ======== getter interface ========

func (v uint64Value) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return &V{}, ErrNotFound
}

// ======== setter interface ========

func (v uint64Value) setAt(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Set()", v.ValueType())
}

// ======== iterater interface ========

func (v uint64Value) RangeObjects(callback func(k string, v *V) bool) {
	// do nothing
}

func (v uint64Value) RangeArray(callback func(i int, v *V) bool) {
	// do nothing
}

func (v uint64Value) ForRangeObj() map[string]*V {
	return map[string]*V{}
}

func (v uint64Value) ForRangeArr() []*V {
	return nil
}

func (v uint64Value) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter)
	close(ch)
	return ch
}

func (v uint64Value) IterArray() <-chan *ArrayIter {
	ch := make(chan *ArrayIter)
	close(ch)
	return ch
}

// ======== marshaler interface ========

func (v uint64Value) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	buf.WriteString(strconv.FormatUint(uint64(v), 10))
	return nil
}

// ======== valuer interface ========

func (v uint64Value) Bool() (bool, error) {
	return v != 0, ErrTypeNotMatch
}

func (v uint64Value) Int64() (int64, error) {
	return int64(v), nil
}

func (v uint64Value) Uint64() (uint64, error) {
	return uint64(v), nil
}

func (v uint64Value) Float64() (float64, error) {
	return float64(v), nil
}

func (v uint64Value) String() string {
	return strconv.FormatUint(uint64(v), 10)
}

func (v uint64Value) Len() int {
	return 0
}

// ======== numberAsserter interface ========

func (v uint64Value) IsFloat() bool {
	return false
}

func (v uint64Value) IsInteger() bool {
	return true
}

func (v uint64Value) IsNegative() bool {
	return false
}

func (v uint64Value) IsPositive() bool {
	return true
}

func (v uint64Value) GreaterThanInt64Max() bool {
	return v > 0x7FFFFFFFFFFFFFFF
}

// ======== inserter interface ========

func (v uint64Value) insertBefore(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

func (v uint64Value) insertAfter(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

// ======= appender interface ========

func (v uint64Value) appendInTheBeginning(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}

func (v uint64Value) appendInTheEnd(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}
