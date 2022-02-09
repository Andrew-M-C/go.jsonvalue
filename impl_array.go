package jsonvalue

import (
	"bytes"
	"fmt"
	"reflect"
)

type arrayValue struct {
	children []*V
}

// NewArray returns an emty array-typed jsonvalue object
//
// NewArray 返回一个初始化好的 array 类型的 jsonvalue 值。
func NewArray() *V {
	return &V{
		impl: &arrayValue{},
	}
}

func newArray() (*V, *arrayValue) {
	impl := &arrayValue{}
	v := &V{
		impl: impl,
	}
	return v, impl
}

// ======== deleter interface ========

func (v *arrayValue) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	paramCount := len(otherParams)
	if paramCount == 0 {
		return v.deleteInCurrValue(caseless, firstParam)
	}

	child, err := v.get(caseless, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	// if child == nil {
	// 	return ErrNotFound
	// }

	return child.impl.delete(caseless, otherParams[paramCount-1])
}

func (v *arrayValue) deleteInCurrValue(caseless bool, param interface{}) error {
	// interger expected
	pos, err := intfToInt(param)
	if err != nil {
		return err
	}

	pos = v.posAtIndexForRead(pos)
	if pos < 0 {
		return ErrOutOfRange
	}
	v.deleteInArr(pos)
	return nil
}

func (v *arrayValue) deleteInArr(pos int) {
	le := len(v.children)
	v.children[pos] = nil
	copy(v.children[pos:], v.children[pos+1:])
	v.children = v.children[:le-1]
}

// ======== typper interface ========

func (v *arrayValue) ValueType() ValueType {
	return Array
}

// ======== getter interface ========

func (v *arrayValue) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	child, err := v.getInCurrValue(firstParam)
	if err != nil {
		return &V{}, err
	}

	if len(otherParams) == 0 {
		return child, nil
	}
	return child.impl.get(caseless, otherParams[0], otherParams[1:]...)
}

func (v *arrayValue) getInCurrValue(param interface{}) (*V, error) {
	// integer expected
	pos, err := intfToInt(param)
	if err != nil {
		return &V{}, err
	}
	child, ok := v.childAtIndex(pos)
	if !ok {
		return &V{}, ErrOutOfRange
	}
	return child, nil
}

func (v *arrayValue) childAtIndex(pos int) (*V, bool) { // if nil returned, means that just push
	pos = v.posAtIndexForRead(pos)
	if pos < 0 {
		return &V{}, false
	}
	return v.children[pos], true
}

func (v *arrayValue) posAtIndexForRead(pos int) int {
	le := len(v.children)
	if le == 0 {
		return -1
	}

	if pos < 0 {
		pos += le
		if pos < 0 {
			return -1
		}
		return pos
	}

	if pos >= le {
		return -1
	}

	return pos
}

// ======== setter interface ========

func (v *arrayValue) setAt(end *V, firstParam interface{}, otherParams ...interface{}) error {
	// this is the last iteration
	if len(otherParams) == 0 {
		pos, err := intfToInt(firstParam)
		if err != nil {
			return err
		}
		err = v.setAtIndex(end, pos)
		if err != nil {
			return err
		}
		return nil
	}

	// this is not the last param
	pos, err := intfToInt(firstParam)
	if err != nil {
		return err
	}

	isNewChild := false
	child, ok := v.childAtIndex(pos)
	if !ok {
		isNewChild = true
		if _, err := intfToString(otherParams[0]); err == nil {
			child = NewObject()
		} else if i, err := intfToInt(otherParams[0]); err == nil {
			if i != 0 {
				return ErrOutOfRange
			}
			child = NewArray()
		} else {
			return fmt.Errorf("unexpected type %v for Set()", reflect.TypeOf(otherParams[0]))
		}
	}

	// go deeper
	next := Set{
		v: child,
		c: end,
	}
	_, err = next.At(otherParams[0], otherParams[1:]...)
	if err != nil {
		return err
	}
	// OK to add this object
	if isNewChild {
		v.children = append(v.children, child)
	}
	return nil
}

func (v *arrayValue) setAtIndex(child *V, pos int) error {
	pos, appendToEnd := v.posAtIndexForSet(pos)
	if pos < 0 {
		return ErrOutOfRange
	}
	if appendToEnd {
		v.children = append(v.children, child)
	} else {
		v.children[pos] = child
	}
	return nil
}

func (v *arrayValue) posAtIndexForSet(pos int) (newPos int, appendToEnd bool) {
	if pos == len(v.children) {
		return pos, true
	}
	pos = v.posAtIndexForRead(pos)
	return pos, false
}

// ======== iterater interface ========

func (v *arrayValue) RangeObjects(callback func(k string, v *V) bool) {
	// do nothing
}

func (v *arrayValue) RangeArray(callback func(i int, v *V) bool) {
	if nil == callback {
		return
	}

	for i, child := range v.children {
		if ok := callback(i, child); !ok {
			break
		}
	}
}

func (v *arrayValue) ForRangeObj() map[string]*V {
	return map[string]*V{}
}

func (v *arrayValue) ForRangeArr() []*V {
	res := make([]*V, 0, len(v.children))
	return append(res, v.children...)
}

func (v *arrayValue) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter)
	close(ch)
	return ch
}

func (v *arrayValue) IterArray() <-chan *ArrayIter {
	c := make(chan *ArrayIter, len(v.children))

	go func() {
		for i, child := range v.children {
			c <- &ArrayIter{
				I: i,
				V: child,
			}
		}
		close(c)
	}()
	return c
}

// ======== marshaler interface ========

func (v *arrayValue) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	buf.WriteByte('[')

	for i, child := range v.children {
		if i > 0 {
			buf.WriteByte(',')
		}
		if opt.MarshalLessFunc == nil {
			child.impl.marshalToBuffer(child, nil, buf, opt)
		} else {
			child.impl.marshalToBuffer(child, curr.newParentInfo(parentInfo, intKey(i)), buf, opt)
		}
	}

	buf.WriteByte(']')
	return nil
}

// ======== valuer interface ========

func (v *arrayValue) Bool() (bool, error) {
	return false, ErrTypeNotMatch
}

func (v *arrayValue) Int64() (int64, error) {
	return 0, ErrTypeNotMatch
}

func (v *arrayValue) Uint64() (uint64, error) {
	return 0, ErrTypeNotMatch
}

func (v *arrayValue) Float64() (float64, error) {
	return 0, ErrTypeNotMatch
}

func (v *arrayValue) String() string {
	buff := bytes.Buffer{}
	buff.WriteByte('[')

	for i, child := range v.children {
		if i > 0 {
			buff.WriteByte(' ')
		}
		buff.WriteString(child.String())
	}

	buff.WriteByte(']')
	return buff.String()
}

func (v *arrayValue) Len() int {
	return len(v.children)
}

// ======== inserter interface ========

func (v *arrayValue) insertBefore(end *V, firstParam interface{}, otherParams ...interface{}) error {
	// this is the last iteration
	paramCount := len(otherParams)
	if paramCount == 0 {
		pos, err := intfToInt(firstParam)
		if err != nil {
			return err
		}

		pos = v.posAtIndexForInsertBefore(pos)
		if pos < 0 {
			return ErrOutOfRange
		}
		v.insertToArr(pos, end)
		return nil
	}

	// this is not the last iterarion
	child, err := v.get(false, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	if child.ValueType() != Array {
		return ErrTypeNotMatch
	}

	return child.impl.insertBefore(end, otherParams[paramCount-1])
}

func (v *arrayValue) insertAfter(end *V, firstParam interface{}, otherParams ...interface{}) error {
	// this is the last iteration
	paramCount := len(otherParams)
	if paramCount == 0 {
		pos, err := intfToInt(firstParam)
		if err != nil {
			return err
		}

		pos, appendToEnd := v.posAtIndexForInsertAfter(pos)
		if pos < 0 {
			return ErrOutOfRange
		}
		if appendToEnd {
			v.children = append(v.children, end)
		} else {
			v.insertToArr(pos, end)
		}
		return nil
	}

	// this is not the last iterarion
	child, err := v.get(false, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	if child.ValueType() != Array {
		return ErrTypeNotMatch
	}

	return child.impl.insertAfter(end, otherParams[paramCount-1])
}

func (v *arrayValue) insertToArr(pos int, child *V) {
	v.children = append(v.children, nil)
	copy(v.children[pos+1:], v.children[pos:])
	v.children[pos] = child
}

func (v *arrayValue) posAtIndexForInsertBefore(pos int) (newPos int) {
	le := len(v.children)
	if le == 0 {
		return -1
	}

	if pos == 0 {
		return 0
	}

	if pos < 0 {
		pos += le
		if pos < 0 {
			return -1
		}
		return pos
	}

	if pos >= le {
		return -1
	}

	return pos
}

func (v *arrayValue) posAtIndexForInsertAfter(pos int) (newPos int, appendToEnd bool) {
	le := len(v.children)
	if le == 0 {
		return -1, false
	}

	if pos == -1 {
		return le, true
	}

	if pos < 0 {
		pos += le
		if pos < 0 {
			return -1, false
		}
		return pos + 1, false
	}

	if pos >= le {
		return -1, false
	}

	return pos + 1, false
}

// ======= appender interface ========

func (v *arrayValue) appendInTheBeginning(end *V, params ...interface{}) error {
	// this is the last iteration
	var childV *arrayValue
	paramCount := len(params)
	if paramCount == 0 {
		childV = v
	} else {
		// this is not the last iterarion
		child, err := v.get(false, params[0], params[1:]...)
		if err != nil {
			return err
		}
		if child.ValueType() != Array {
			return ErrTypeNotMatch
		}
		childV = child.impl.(*arrayValue)
	}

	if len(childV.children) == 0 {
		childV.children = append(childV.children, end)
	} else {
		childV.insertToArr(0, end)
	}

	return nil
}

func (v *arrayValue) appendInTheEnd(end *V, params ...interface{}) error {
	// this is the last iteration
	var childV *arrayValue
	paramCount := len(params)
	if paramCount == 0 {
		childV = v
	} else {
		// this is not the last iterarion
		child, err := v.get(false, params[0], params[1:]...)
		if err != nil {
			return err
		}
		if child.ValueType() != Array {
			return ErrTypeNotMatch
		}
		childV = child.impl.(*arrayValue)
	}

	childV.children = append(childV.children, end)

	return nil
}
