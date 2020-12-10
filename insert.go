package jsonvalue

import (
	"github.com/buger/jsonparser"
)

// Insert type is for After() and Before() function. Please refer for realated function.
//
// Should be generated ONLY BY V.Insert function!
type Insert struct {
	v *V
	c *V // child
}

// Insert starts setting a child JSON value
func (v *V) Insert(child *V) *Insert {
	if nil == child {
		child = NewNull()
	}
	return &Insert{
		v: v,
		c: child,
	}
}

// InsertString is equivalent to Insert(jsonvalue.NewString(s))
func (v *V) InsertString(s string) *Insert {
	return v.Insert(NewString(s))
}

// InsertBool is equivalent to Insert(jsonvalue.NewBool(b))
func (v *V) InsertBool(b bool) *Insert {
	return v.Insert(NewBool(b))
}

// InsertInt is equivalent to Insert(jsonvalue.NewInt(b))
func (v *V) InsertInt(i int) *Insert {
	return v.Insert(NewInt(i))
}

// InsertInt64 is equivalent to Insert(jsonvalue.NewInt64(b))
func (v *V) InsertInt64(i int64) *Insert {
	return v.Insert(NewInt64(i))
}

// InsertInt32 is equivalent to Insert(jsonvalue.NewInt32(b))
func (v *V) InsertInt32(i int32) *Insert {
	return v.Insert(NewInt32(i))
}

// InsertUint is equivalent to Insert(jsonvalue.NewUint(b))
func (v *V) InsertUint(u uint) *Insert {
	return v.Insert(NewUint(u))
}

// InsertUint64 is equivalent to Insert(jsonvalue.NewUint64(b))
func (v *V) InsertUint64(u uint64) *Insert {
	return v.Insert(NewUint64(u))
}

// InsertUint32 is equivalent to Insert(jsonvalue.NewUint32(b))
func (v *V) InsertUint32(u uint32) *Insert {
	return v.Insert(NewUint32(u))
}

// InsertFloat64 is equivalent to Insert(jsonvalue.NewFloat64(b))
func (v *V) InsertFloat64(f float64, prec int) *Insert {
	return v.Insert(NewFloat64(f, prec))
}

// InsertFloat32 is equivalent to Insert(jsonvalue.NewFloat32(b))
func (v *V) InsertFloat32(f float32, prec int) *Insert {
	return v.Insert(NewFloat32(f, prec))
}

// InsertNull is equivalent to Insert(jsonvalue.NewNull())
func (v *V) InsertNull() *Insert {
	return v.Insert(NewNull())
}

// InsertObject is equivalent to Insert(jsonvalue.NewObject())
func (v *V) InsertObject() *Insert {
	return v.Insert(NewObject())
}

// InsertArray is equivalent to Insert(jsonvalue.NewArray())
func (v *V) InsertArray() *Insert {
	return v.Insert(NewArray())
}

// Before completes the following operation of Insert().
//
// The last parameter identifies the postion where a new JSON is inserted after, it should ba an interger, no matter signed or unsigned.
// If the position is zero or positive interger, it tells the index of an array. If the position is negative, it tells the backward index of an array.
//
// For example, 0 represents the first, and -2 represents the second last.
func (ins *Insert) Before(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	v := ins.v
	c := ins.c
	if v.valueType == jsonparser.NotExist {
		return nil, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(otherParams)
	if 0 == paramCount {
		if v.valueType != jsonparser.Array {
			return nil, ErrNotArrayValue
		}

		pos, err := intfToInt(firstParam)
		if err != nil {
			return nil, err
		}

		e := v.elementAtIndex(pos)
		if nil == e {
			return nil, ErrOutOfRange
		}
		v.children.array.InsertBefore(c, e)
		return c, nil
	}

	// this is not the last iterarion
	child, err := v.GetArray(firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return nil, err
	}

	childIns := Insert{
		v: child,
		c: c,
	}
	return childIns.Before(otherParams[paramCount-1])
}

// After completes the following operation of Insert().
//
// The last parameter identifies the postion where a new JSON is inserted after, it should ba an interger, no matter signed or unsigned.
// If the position is zero or positive interger, it tells the index of an array. If the position is negative, it tells the backward index of an array.
//
// For example, 0 represents the first, and -2 represents the second last.
func (ins *Insert) After(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	v := ins.v
	c := ins.c
	if nil == v || v.valueType == jsonparser.NotExist {
		return nil, ErrValueUninitialized
	}
	if nil == c || c.valueType == jsonparser.NotExist {
		return nil, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(otherParams)
	if 0 == paramCount {
		if v.valueType != jsonparser.Array {
			return nil, ErrNotArrayValue
		}

		pos, err := intfToInt(firstParam)
		if err != nil {
			return nil, err
		}

		e := v.elementAtIndex(pos)
		if nil == e {
			return nil, ErrOutOfRange
		}
		v.children.array.InsertAfter(c, e)
		return c, nil
	}

	// this is not the last iterarion
	child, err := v.GetArray(firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return nil, err
	}

	childIns := Insert{
		v: child,
		c: c,
	}
	return childIns.After(otherParams[paramCount-1])
}
