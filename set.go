package jsonvalue

import (
	"fmt"
	"reflect"
)

// Set type is for At() only.
//
// Set 类型仅用于 At() 函数。
type Setter interface {
	// At completes the following operation of Set(). It defines position of value in Set() and return the new value set.
	//
	// The usage of At() is perhaps the most important. This function will recursively search for child value, and set the
	// new value specified by Set() or SetXxx() series functions. Please unfold and read the following examples, they are important.
	//
	// At 完成 Set() 函数的后续操作并设置相应的子成员。其参数指定了应该在哪个位置设置子成员，并且返回被设置的子成员对象。
	//
	// 该函数的用法恐怕是 jsonvalue 中最重要的内容了：该函数会按照给定的可变参数递归地一层一层查找 JSON 值的子成员，并且设置到指定的位置上。
	// 设置的逻辑说明起来比较抽象，请打开以下的例子以了解，这非常重要。
	At(firstParam interface{}, otherParams ...interface{}) (*V, error)
}

type setter struct {
	v *V
	c *V // child

	err error
}

// Set starts setting a child JSON value. Any legal JSON value typped parameter
// is accepted, such as string, int, float, bool, nil, *jsonvalue.V, or even
// a struct or map or slice.
//
// Please refer to examples of "func (set Setter) At(...)"
//
// https://godoc.org/github.com/Andrew-M-C/go.jsonvalue/#Set.At
//
// Set 开始设置一个 JSON 子成员。任何合法的 JSON 类型都可以作为参数, 比如 string, int,
// float, bool, nil, *jsonvalue.V 等类型, 甚至也支持结构体、map、切片、数组。
//
// 请参见 "func (set Setter) At(...)" 例子.
//
// https://godoc.org/github.com/Andrew-M-C/go.jsonvalue/#Set.At
func (v *V) Set(child any) Setter {
	var ch *V
	var err error

	if child == nil {
		ch = NewNull()
	} else if childV, ok := child.(*V); ok {
		ch = childV
	} else {
		ch, err = Import(child)
	}

	return &setter{
		v:   v,
		c:   ch,
		err: err,
	}
}

// SetString is equivalent to Set(jsonvalue.NewString(s))
//
// SetString 等效于 Set(jsonvalue.NewString(s))
func (v *V) SetString(s string) Setter {
	return v.Set(NewString(s))
}

// SetBytes is equivalent to Set(NewString(base64.StdEncoding.EncodeToString(b)))
//
// SetBytes 等效于 Set(NewString(base64.StdEncoding.EncodeToString(b)))
func (v *V) SetBytes(b []byte) Setter {
	s := internal.b64.EncodeToString(b)
	return v.SetString(s)
}

// SetBool is equivalent to Set(jsonvalue.NewBool(b))
//
// SetBool 等效于 Set(jsonvalue.NewBool(b))
func (v *V) SetBool(b bool) Setter {
	return v.Set(NewBool(b))
}

// SetInt is equivalent to Set(jsonvalue.NewInt(b))
//
// SetInt 等效于 Set(jsonvalue.NewInt(b))
func (v *V) SetInt(i int) Setter {
	return v.Set(NewInt(i))
}

// SetInt64 is equivalent to Set(jsonvalue.NewInt64(b))
//
// SetInt64 等效于 Set(jsonvalue.NewInt64(b))
func (v *V) SetInt64(i int64) Setter {
	return v.Set(NewInt64(i))
}

// SetInt32 is equivalent to Set(jsonvalue.NewInt32(b))
//
// SetInt32 等效于 Set(jsonvalue.NewInt32(b))
func (v *V) SetInt32(i int32) Setter {
	return v.Set(NewInt32(i))
}

// SetUint is equivalent to Set(jsonvalue.NewUint(b))
//
// SetUint 等效于 Set(jsonvalue.NewUint(b))
func (v *V) SetUint(u uint) Setter {
	return v.Set(NewUint(u))
}

// SetUint64 is equivalent to Set(jsonvalue.NewUint64(b))
//
// SetUint64 is equivalent to Set(jsonvalue.NewUint64(b))
func (v *V) SetUint64(u uint64) Setter {
	return v.Set(NewUint64(u))
}

// SetUint32 is equivalent to Set(jsonvalue.NewUint32(b))
//
// SetUint32 等效于 Set(jsonvalue.NewUint32(b))
func (v *V) SetUint32(u uint32) Setter {
	return v.Set(NewUint32(u))
}

// SetFloat64 is equivalent to Set(jsonvalue.NewFloat64(b))
//
// SetFloat64 等效于 Set(jsonvalue.NewFloat64(b))
func (v *V) SetFloat64(f float64) Setter {
	return v.Set(NewFloat64(f))
}

// SetFloat32 is equivalent to Set(jsonvalue.NewFloat32(b))
//
// SetFloat32 等效于 Set(jsonvalue.NewFloat32(b))
func (v *V) SetFloat32(f float32) Setter {
	return v.Set(NewFloat32(f))
}

// SetNull is equivalent to Set(jsonvalue.NewNull())
//
// SetNull 等效于 Set(jsonvalue.NewNull())
func (v *V) SetNull() Setter {
	return v.Set(NewNull())
}

// SetObject is equivalent to Set(jsonvalue.NewObject())
//
// SetObject 等效于 Set(jsonvalue.NewObject())
func (v *V) SetObject() Setter {
	return v.Set(NewObject())
}

// SetArray is equivalent to Set(jsonvalue.NewArray())
//
// SetArray 等效于 Set(jsonvalue.NewArray())
func (v *V) SetArray() Setter {
	return v.Set(NewArray())
}

func (v *V) setToObjectChildren(key string, child *V) {
	v.children.incrID++
	v.children.object[key] = childWithProperty{
		id: v.children.incrID,
		v:  child,
	}
	v.addCaselessKey(key)
}

func (s *setter) At(firstParam any, otherParams ...any) (*V, error) {
	if s.err != nil {
		return &V{}, s.err
	}
	v := s.v
	c := s.c
	if nil == v || v.valueType == NotExist {
		return &V{}, ErrValueUninitialized
	}
	if nil == c || c.valueType == NotExist {
		return &V{}, ErrValueUninitialized
	}

	// this is the last iteration
	if len(otherParams) == 0 {
		switch v.valueType {
		default:
			return &V{}, fmt.Errorf("%v type does not supports Set()", v.valueType)

		case Object:
			var k string
			k, err := intfToString(firstParam)
			if err != nil {
				return &V{}, err
			}
			v.setToObjectChildren(k, c)
			return c, nil

		case Array:
			pos, err := intfToInt(firstParam)
			if err != nil {
				return &V{}, err
			}
			err = v.setAtIndex(c, pos)
			if err != nil {
				return &V{}, err
			}
			return c, nil
		}
	}

	// this is not the last iterarion
	if v.valueType == Object {
		k, err := intfToString(firstParam)
		if err != nil {
			return &V{}, err
		}
		child, exist := v.getFromObjectChildren(false, k)
		if !exist {
			if _, err := intfToString(otherParams[0]); err == nil {
				child = NewObject()
			} else if i, err := intfToInt(otherParams[0]); err == nil {
				if i != 0 {
					return &V{}, ErrOutOfRange
				}
				child = NewArray()
			} else {
				return &V{}, fmt.Errorf("unexpected type %v for Set()", reflect.TypeOf(otherParams[0]))
			}
		}
		next := &setter{
			v: child,
			c: c,
		}
		_, err = next.At(otherParams[0], otherParams[1:]...)
		if err != nil {
			return &V{}, err
		}
		if !exist {
			v.setToObjectChildren(k, child)
		}
		return c, nil
	}

	// array type
	if v.valueType == Array {
		pos, err := intfToInt(firstParam)
		if err != nil {
			return &V{}, err
		}
		child, ok := v.childAtIndex(pos)
		isNewChild := false
		if !ok {
			isNewChild = true
			if _, err := intfToString(otherParams[0]); err == nil {
				child = NewObject()
			} else if i, err := intfToInt(otherParams[0]); err == nil {
				if i != 0 {
					return &V{}, ErrOutOfRange
				}
				child = NewArray()
			} else {
				return &V{}, fmt.Errorf("unexpected type %v for Set()", reflect.TypeOf(otherParams[0]))
			}
		}
		next := &setter{
			v: child,
			c: c,
		}
		_, err = next.At(otherParams[0], otherParams[1:]...)
		if err != nil {
			return &V{}, err
		}
		// OK to add this object
		if isNewChild {
			v.appendToArr(child)
		}
		return c, nil
	}

	// illegal type
	return &V{}, fmt.Errorf("%v type does not supports Set()", v.valueType)
}

func (v *V) posAtIndexForSet(pos int) (newPos int, appendToEnd bool) {
	if pos == len(v.children.arr) {
		return pos, true
	}
	pos = v.posAtIndexForRead(pos)
	return pos, false
}

func (v *V) posAtIndexForInsertBefore(pos int) (newPos int) {
	le := len(v.children.arr)
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

func (v *V) posAtIndexForInsertAfter(pos int) (newPos int, appendToEnd bool) {
	le := len(v.children.arr)
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

func (v *V) posAtIndexForRead(pos int) int {
	le := len(v.children.arr)
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

func (v *V) childAtIndex(pos int) (*V, bool) { // if nil returned, means that just push
	pos = v.posAtIndexForRead(pos)
	if pos < 0 {
		return &V{}, false
	}
	return v.children.arr[pos], true
}

func (v *V) setAtIndex(child *V, pos int) error {
	pos, appendToEnd := v.posAtIndexForSet(pos)
	if pos < 0 {
		return ErrOutOfRange
	}
	if appendToEnd {
		v.children.arr = append(v.children.arr, child)
	} else {
		v.children.arr[pos] = child
	}
	return nil
}
