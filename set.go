package jsonvalue

import (
	"container/list"
	"fmt"

	"github.com/buger/jsonparser"
)

// Set set recursively set child into a jsonvalue object
func (v *V) Set(firstParam interface{}, secondParam interface{}, otherParams ...interface{}) (value *V, err error) {
	if 0 == len(otherParams) {
		// This is the last value, now we should set
		switch v.valueType {
		default:
			err = fmt.Errorf("%v type does not supports Set()", v.valueType)
			return

		case jsonparser.Object:
			var k string
			var child *V
			k, err = intfToString(firstParam)
			if err != nil {
				return
			}
			child, err = intfToJsonvalue(secondParam)
			if err != nil {
				return
			}
			v.objectChildren[k] = child
			return child, nil

		case jsonparser.Array:
			var pos int
			var child *V
			pos, err = intfToInt(firstParam)
			if err != nil {
				return nil, err
			}
			child, err = intfToJsonvalue(secondParam)
			if err != nil {
				return
			}

			err = v.insertAtIndex(child, pos)
			if err != nil {
				return
			}
			return child, nil
		}
	}

	// Not last params? continue
	// object value
	if v.valueType == jsonparser.Object {
		// string key is expected
		var k string
		k, err = intfToString(firstParam)
		if err != nil {
			return nil, err
		}
		child, exist := v.objectChildren[k]
		if false == exist {
			// if next parameter is a string, create an object
			if _, strErr := intfToString(secondParam); strErr == nil {
				child = NewObject()
			} else if _, arrErr := intfToInt(secondParam); arrErr == nil {
				child = NewArray()
			} else {
				err = fmt.Errorf("unexpected type %v for Set()", v.valueType)
				return nil, err
			}
		}
		var grandChild *V
		grandChild, err = child.Set(secondParam, otherParams[0], otherParams[1:]...)
		if err != nil {
			return nil, err
		}
		// OK to add this object
		if false == exist {
			v.objectChildren[k] = child
		}
		return grandChild, nil
	}

	// array value
	if v.valueType == jsonparser.Array {
		// interger is expected
		var pos int
		pos, err = intfToInt(firstParam)
		if err != nil {
			return nil, err
		}
		child, _ := v.childAtIndex(pos)
		isNewChild := false
		if nil == child {
			isNewChild = true
			if _, strErr := intfToString(secondParam); strErr == nil {
				child = NewObject()
			} else if _, arrErr := intfToInt(secondParam); arrErr == nil {
				child = NewArray()
			} else {
				err = fmt.Errorf("unexpected type %v for Set()", v.valueType)
				return nil, err
			}
		}
		var grandChild *V
		grandChild, err = child.Set(secondParam, otherParams[0], otherParams[1:]...)
		if err != nil {
			return nil, err
		}
		// OK to add this object
		if isNewChild {
			v.arrayChildren.PushBack(child)
		}
		return grandChild, nil
	}

	return nil, fmt.Errorf("%s does not supports Set()", v.valueType)
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

// MustSet is same as Set(), but it panics if error occurred
func (v *V) MustSet(firstParam interface{}, secondParam interface{}, otherParams ...interface{}) *V {
	ret, err := v.Set(firstParam, secondParam, otherParams...)
	if err != nil {
		panic(err)
	}
	return ret
}

// SetString is short for Set(..., v). While the last parameter must be a string object
func (v *V) SetString(firstParam interface{}, secondParam interface{}, otherParams ...interface{}) (*V, error) {
	otherParamCount := len(otherParams)

	if 0 == otherParamCount {
		s, ok := secondParam.(string)
		if false == ok {
			return nil, fmt.Errorf("string value is not set in SetString()")
		}
		return v.Set(firstParam, NewString(s))
	}

	lastParam := otherParams[otherParamCount-1]
	s, ok := lastParam.(string)
	if false == ok {
		return nil, fmt.Errorf("string value is not set in SetString()")
	}
	otherParams[otherParamCount-1] = NewString(s)
	return v.Set(firstParam, secondParam, otherParams...)
}

// MustSetString is same as SetString(), but it panics if error occurred
func (v *V) MustSetString(firstParam interface{}, secondParam interface{}, otherParams ...interface{}) *V {
	ret, err := v.SetString(firstParam, secondParam, otherParams...)
	if err != nil {
		panic(err)
	}
	return ret
}

// func (v *V) SetInt(firstParam interface{}, secondParam interface{}, otherParams ...interface{}) (*V, error) {
// 	otherParamCount := len(otherParams)

// 	if 0 == otherParamCount {
// 		s, ok := secondParam.(string)
// 		if false == ok {
// 			return nil, fmt.Errorf("string value is not set in SetString()")
// 		}
// 		return v.Set(firstParam, NewString(s))
// 	}

// 	lastParam := otherParams[otherParamCount-1]
// 	s, ok := lastParam.(string)
// 	if false == ok {
// 		return nil, fmt.Errorf("string value is not set in SetString()")
// 	}
// 	otherParams[otherParamCount-1] = NewString(s)
// 	return v.Set(firstParam, secondParam, otherParams...)
// }
