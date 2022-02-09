package jsonvalue

import "bytes"

// deleter identifies a value supporting deletion operation.
//
// deleter 代表一个支持 delete 操作的值。
type deleter interface {
	delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error
}

// getter identifies a value supporting pure get operation only.
//
// getter 代表一个仅支持 get 操作的值。
type getter interface {
	get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error)
}

// setter identifies a value supporting set operation.
//
// setter 代表一个支持 set 操作的类型
type setter interface {
	setAt(child *V, firstParam interface{}, otherParams ...interface{}) error
}

// iterater identies iteration operation.
//
// iterater 代表迭代操作。
type iterater interface {
	RangeObjects(callback func(k string, v *V) bool)
	RangeArray(callback func(i int, v *V) bool)

	ForRangeObj() map[string]*V
	ForRangeArr() []*V

	IterObjects() <-chan *ObjectIter
	IterArray() <-chan *ArrayIter
}

// marshaler identifies marshaling operation.
//
// marshaler 代表序列化操作。
type marshaler interface {
	marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error)
}

// typper identifies ValueType operation.
type typper interface {
	ValueType() ValueType
}

// valuer identities a value getter
type valuer interface {
	Bool() (bool, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Float64() (float64, error)
	String() string

	Len() int
}

// inserter identifies a value inserter
type inserter interface {
	insertBefore(child *V, firstParam interface{}, otherParams ...interface{}) error
	insertAfter(child *V, firstParam interface{}, otherParams ...interface{}) error
}

// appender identifies a value appender
type appender interface {
	appendInTheBeginning(child *V, params ...interface{}) error
	appendInTheEnd(child *V, params ...interface{}) error
}

// numberAsserter identifies a number property asserter
type numberAsserter interface {
	IsFloat() bool
	IsInteger() bool
	IsNegative() bool
	IsPositive() bool
	GreaterThanInt64Max() bool
}

// value identifies a jsonvalue implememtation.
//
// value 代表一个完整的 json 值实现。
type value interface {
	typper
	getter
	setter
	deleter
	iterater
	marshaler
	valuer
	inserter
	appender

	// TODO: inserter and appender
	// TODO: 删除 numberAsserter
}
