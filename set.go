package jsonvalue

import (
	"container/list"
	"fmt"
	"reflect"

	"github.com/buger/jsonparser"
)

type set struct {
	v *V
	c *V // child
}

// Set starts setting a child JSON value
func (v *V) Set(child *V) *set {
	if nil == child {
		child = NewNull()
	}
	return &set{
		v: v,
		c: child,
	}
}

// At
func (s *set) At(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	v := s.v
	c := s.c
	if v.valueType == jsonparser.NotExist {
		return nil, ErrValueUninitialized
	}

	// this is the last iteration
	if 0 == len(otherParams) {
		switch v.valueType {
		default:
			return nil, fmt.Errorf("%v type does not supports Set()", v.valueType)

		case jsonparser.Object:
			var k string
			k, err := intfToString(firstParam)
			if err != nil {
				return nil, err
			}
			v.objectChildren[k] = c
			return c, nil

		case jsonparser.Array:
			pos, err := intfToInt(firstParam)
			if err != nil {
				return nil, err
			}
			err = v.insertAtIndex(c, pos)
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	}

	// this is not the last iterarion
	if v.valueType == jsonparser.Object {
		k, err := intfToString(firstParam)
		if err != nil {
			return nil, err
		}
		child, exist := v.objectChildren[k]
		if false == exist {
			if _, err := intfToString(otherParams[0]); err == nil {
				child = NewObject()
			} else if _, err := intfToInt(otherParams[0]); err == nil {
				child = NewArray()
			} else {
				return nil, fmt.Errorf("unexpected type %v for Set()", reflect.TypeOf(otherParams[0]))
			}
		}
		next := set{
			v: child,
			c: c,
		}
		_, err = next.At(otherParams[0], otherParams[1:]...)
		if err != nil {
			return nil, err
		}
		if false == exist {
			v.objectChildren[k] = child
		}
		return c, nil
	}

	// array type
	if v.valueType == jsonparser.Array {
		pos, err := intfToInt(firstParam)
		if err != nil {
			return nil, err
		}
		child, _ := v.childAtIndex(pos)
		isNewChild := false
		if nil == child {
			isNewChild = true
			if _, err := intfToString(otherParams[0]); err == nil {
				child = NewObject()
			} else if _, err := intfToInt(otherParams[0]); err == nil {
				child = NewArray()
			} else {
				return nil, fmt.Errorf("unexpected type %v for Set()", reflect.TypeOf(otherParams[0]))
			}
		}
		next := set{
			v: child,
			c: c,
		}
		_, err = next.At(otherParams[0], otherParams[1:]...)
		if err != nil {
			return nil, err
		}
		// OK to add this object
		if isNewChild {
			v.arrayChildren.PushBack(child)
		}
		return c, nil
	}

	// illegal type
	return nil, fmt.Errorf("%v type does not supports Set()", v.valueType)
}

func (v *V) childAtIndex(pos int) (*V, bool) { // if nil returned, means that just push
	if pos < 0 {
		return nil, false
	}
	if 0 == v.arrayChildren.Len() {
		return nil, false
	}

	// find element at pos
	var e *list.Element
	i := 0
	for e = v.arrayChildren.Front(); e != nil && i < pos; e = e.Next() {
		i++
	}

	if nil == e {
		return nil, true
	}
	return e.Value.(*V), false
}

func (v *V) insertAtIndex(child *V, pos int) (err error) {
	if pos < 0 {
		// append in the end
		v.arrayChildren.PushBack(child)
		return
	}

	if 0 == v.arrayChildren.Len() {
		v.arrayChildren.PushBack(child)
		return
	}

	// find element at pos
	var e *list.Element
	i := 0
	for e = v.arrayChildren.Front(); e != nil && i < pos; e = e.Next() {
		i++
	}

	// exceeds array length
	if nil == e {
		v.arrayChildren.PushBack(child)
		return
	}

	// find position correct
	v.arrayChildren.InsertAfter(child, e)
	return
}

// ==== short cuts ====
