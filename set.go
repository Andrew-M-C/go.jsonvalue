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

// SetString is equivalent to Set(jsonvalue.NewString(s))
func (v *V) SetString(s string) *set {
	return v.Set(NewString(s))
}

// SetBool is equivalent to Set(jsonvalue.NewBool(b))
func (v *V) SetBool(b bool) *set {
	return v.Set(NewBool(b))
}

// SetInt is equivalent to Set(jsonvalue.NewInt(b))
func (v *V) SetInt(i int) *set {
	return v.Set(NewInt(i))
}

// SetInt64 is equivalent to Set(jsonvalue.NewInt64(b))
func (v *V) SetInt64(i int64) *set {
	return v.Set(NewInt64(i))
}

// SetInt32 is equivalent to Set(jsonvalue.NewInt32(b))
func (v *V) SetInt32(i int32) *set {
	return v.Set(NewInt32(i))
}

// SetUint is equivalent to Set(jsonvalue.NewUint(b))
func (v *V) SetUint(u uint) *set {
	return v.Set(NewUint(u))
}

// SetUint64 is equivalent to Set(jsonvalue.NewUint64(b))
func (v *V) SetUint64(u uint64) *set {
	return v.Set(NewUint64(u))
}

// SetUint32 is equivalent to Set(jsonvalue.NewUint32(b))
func (v *V) SetUint32(u uint32) *set {
	return v.Set(NewUint32(u))
}

// SetFloat64 is equivalent to Set(jsonvalue.NewFloat64(b))
func (v *V) SetFloat64(f float64, prec int) *set {
	return v.Set(NewFloat64(f, prec))
}

// SetFloat32 is equivalent to Set(jsonvalue.NewFloat32(b))
func (v *V) SetFloat32(f float32, prec int) *set {
	return v.Set(NewFloat32(f, prec))
}

// SetNull is equivalent to Set(jsonvalue.NewNull())
func (v *V) SetNull() *set {
	return v.Set(NewNull())
}

// SetObject is equivalent to Set(jsonvalue.NewObject())
func (v *V) SetObject() *set {
	return v.Set(NewObject())
}

// SetArray is equivalent to Set(jsonvalue.NewArray())
func (v *V) SetArray() *set {
	return v.Set(NewArray())
}

// At completes the following operation of Set(). It defines posttion of value in Set().
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
			err = v.setAtIndex(c, pos)
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

func (v *V) elementAtIndex(pos int) *list.Element {
	l := v.arrayChildren.Len()
	if 0 == l {
		return nil
	}
	if pos < 0 {
		pos = l + pos
		if pos < 0 {
			return nil
		}
	} else if pos >= l {
		return nil
	}

	// find element at pos
	var e *list.Element
	i := 0
	for e = v.arrayChildren.Front(); e != nil && i < pos; e = e.Next() {
		i++
	}
	return e
}

func (v *V) childAtIndex(pos int) (*V, bool) { // if nil returned, means that just push
	// find element at pos
	e := v.elementAtIndex(pos)
	if nil == e {
		return nil, true
	}
	return e.Value.(*V), false
}

func (v *V) setAtIndex(child *V, pos int) error {
	if 0 == v.arrayChildren.Len() {
		return ErrOutOfRange
	}
	if -1 == pos {
		pos = v.arrayChildren.Len() - 1
	}

	e := v.elementAtIndex(pos)
	if nil == e {
		return ErrOutOfRange
	}
	v.arrayChildren.InsertBefore(child, e)
	v.arrayChildren.Remove(e)
	return nil
}
