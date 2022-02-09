package jsonvalue

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type stringValue string

// NewString returns an initialied string jsonvalue object
//
// NewString 用给定的 string 返回一个初始化好的字符串类型的 jsonvalue 值
func NewString(s string) *V {
	return &V{
		impl: stringValue(s),
	}
}

// NewBytes returns an initialized string with Base64 string by given bytes
//
// NewBytes 用给定的字节串，返回一个初始化好的字符串类型的 jsonvalue 值，内容是字节串 Base64 之后的字符串。
func NewBytes(b []byte) *V {
	s := base64.StdEncoding.EncodeToString(b)
	return NewString(s)
}

// ======== deleter interface ========

func (v stringValue) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	return ErrNotFound
}

// ======== typper interface ========

func (v stringValue) ValueType() ValueType {
	return String
}

// ======== getter interface ========

func (v stringValue) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return &V{}, ErrNotFound
}

// ======== setter interface ========

func (v stringValue) setAt(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Set()", v.ValueType())
}

// ======== iterater interface ========

func (v stringValue) RangeObjects(callback func(k string, v *V) bool) {
	// do nothing
}

func (v stringValue) RangeArray(callback func(i int, v *V) bool) {
	// do nothing
}

func (v stringValue) ForRangeObj() map[string]*V {
	return map[string]*V{}
}

func (v stringValue) ForRangeArr() []*V {
	return nil
}

func (v stringValue) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter)
	close(ch)
	return ch
}

func (v stringValue) IterArray() <-chan *ArrayIter {
	ch := make(chan *ArrayIter)
	close(ch)
	return ch
}

// ======== marshaler interface ========

func (v stringValue) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	buf.WriteByte('"')
	escapeStringToBuff(string(v), buf, opt)
	buf.WriteByte('"')
	return nil
}

// ======== valuer interface ========

func (v stringValue) Bool() (bool, error) {
	return v == "true", ErrTypeNotMatch
}

func (v stringValue) Int64() (int64, error) {
	s := strings.TrimSpace(string(v))
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, ErrTypeNotMatch
	}
	f, err := v.Float64()
	return int64(f), err
}

func (v stringValue) Uint64() (uint64, error) {
	s := strings.TrimSpace(string(v))
	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return u, ErrTypeNotMatch
	}
	f, err := v.Float64()
	return uint64(f), err
}

func (v stringValue) Float64() (float64, error) {
	s := strings.TrimSpace(string(v))
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f, ErrTypeNotMatch
	}
	return 0, fmt.Errorf("%w: %v", ErrParseNumberFromString, err)
}

func (v stringValue) String() string {
	return string(v)
}

func (v stringValue) Len() int {
	return 0
}

// ======== inserter interface ========

func (v stringValue) insertBefore(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

func (v stringValue) insertAfter(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

// ======= appender interface ========

func (v stringValue) appendInTheBeginning(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}

func (v stringValue) appendInTheEnd(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}
