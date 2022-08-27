package jsonvalue

import (
	"errors"
	"fmt"
	"strings"
)

// ================ INSERT ================

// Insert type is for After() and Before() function. Please refer for realated function.
//
// Should be generated ONLY BY V.Insert function!
//
// Insert 类型适用于 After() 和 Before() 函数。请参见相关函数。请注意：该类型仅应由 V.Insert 函数生成！
type Insert struct {
	v *V
	c *V // child

	err error
}

// Insert starts inserting a child JSON value
//
// Insert 开启一个 JSON 数组成员的插入操作.
func (v *V) Insert(child any) *Insert {
	var ch *V
	var err error

	if child == nil {
		ch = NewNull()
	} else if childV, ok := child.(*V); ok {
		ch = childV
	} else {
		ch, err = Import(child)
	}

	return &Insert{
		v:   v,
		c:   ch,
		err: err,
	}
}

// InsertString is equivalent to Insert(jsonvalue.NewString(s))
//
// InsertString 等效于 Insert(jsonvalue.NewString(s))
func (v *V) InsertString(s string) *Insert {
	return v.Insert(NewString(s))
}

// InsertBool is equivalent to Insert(jsonvalue.NewBool(b))
//
// InsertBool 等效于 Insert(jsonvalue.NewBool(b))
func (v *V) InsertBool(b bool) *Insert {
	return v.Insert(NewBool(b))
}

// InsertInt is equivalent to Insert(jsonvalue.NewInt(b))
//
// InsertInt 等效于 Insert(jsonvalue.NewInt(b))
func (v *V) InsertInt(i int) *Insert {
	return v.Insert(NewInt(i))
}

// InsertInt64 is equivalent to Insert(jsonvalue.NewInt64(b))
//
// InsertInt64 等效于 Insert(jsonvalue.NewInt64(b))
func (v *V) InsertInt64(i int64) *Insert {
	return v.Insert(NewInt64(i))
}

// InsertInt32 is equivalent to Insert(jsonvalue.NewInt32(b))
//
// InsertInt32 等效于 Insert(jsonvalue.NewInt32(b))
func (v *V) InsertInt32(i int32) *Insert {
	return v.Insert(NewInt32(i))
}

// InsertUint is equivalent to Insert(jsonvalue.NewUint(b))
//
// InsertUint 等效于 Insert(jsonvalue.NewUint(b))
func (v *V) InsertUint(u uint) *Insert {
	return v.Insert(NewUint(u))
}

// InsertUint64 is equivalent to Insert(jsonvalue.NewUint64(b))
//
// InsertUint64 等效于 Insert(jsonvalue.NewUint64(b))
func (v *V) InsertUint64(u uint64) *Insert {
	return v.Insert(NewUint64(u))
}

// InsertUint32 is equivalent to Insert(jsonvalue.NewUint32(b))
//
// InsertUint32 等效于 Insert(jsonvalue.NewUint32(b))
func (v *V) InsertUint32(u uint32) *Insert {
	return v.Insert(NewUint32(u))
}

// InsertFloat64 is equivalent to Insert(jsonvalue.NewFloat64(b))
//
// InsertFloat64 等效于 Insert(jsonvalue.NewFloat64(b))
func (v *V) InsertFloat64(f float64) *Insert {
	return v.Insert(NewFloat64(f))
}

// InsertFloat32 is equivalent to Insert(jsonvalue.NewFloat32(b))
//
// InsertFloat32 等效于 Insert(jsonvalue.NewFloat32(b))
func (v *V) InsertFloat32(f float32) *Insert {
	return v.Insert(NewFloat32(f))
}

// InsertNull is equivalent to Insert(jsonvalue.NewNull())
//
// InsertNull 等效于 Insert(jsonvalue.NewNull())
func (v *V) InsertNull() *Insert {
	return v.Insert(NewNull())
}

// InsertObject is equivalent to Insert(jsonvalue.NewObject())
//
// InsertObject 等效于 Insert(jsonvalue.NewObject())
func (v *V) InsertObject() *Insert {
	return v.Insert(NewObject())
}

// InsertArray is equivalent to Insert(jsonvalue.NewArray())
//
// InsertArray 等效于 Insert(jsonvalue.NewArray())
func (v *V) InsertArray() *Insert {
	return v.Insert(NewArray())
}

// Before completes the following operation of Insert(). It inserts value BEFORE specified position.
//
// The last parameter identifies the postion where a new JSON is inserted after, it should ba an interger, no matter signed or unsigned.
// If the position is zero or positive interger, it tells the index of an array. If the position is negative, it tells the backward index of an array.
//
// For example, 0 represents the first, and -2 represents the second last.
//
// Before 结束并完成 Insert() 函数的后续插入操作，表示插入到指定位置的后面。
//
// 在 Before 函数的最后一个参数指定了被插入的 JSON 数组的位置，这个参数应当是一个整型（有无符号类型均可）。
// 如果这个值等于0或者正整数，那么它指定的是在 JSON 数组中的位置（从0开始）。如果这个值是负数，那么它指定的是 JSON 数组中从最后一个位置开始算起的位置。
//
// 举例说明：0 表示第一个位置，而 -2 表示倒数第二个位置。
func (ins *Insert) Before(firstParam any, otherParams ...any) (*V, error) {
	if ins.err != nil {
		return &V{}, ins.err
	}
	v := ins.v
	c := ins.c
	if v.valueType == NotExist {
		return &V{}, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(otherParams)
	if paramCount == 0 {
		if v.valueType != Array {
			return &V{}, ErrNotArrayValue
		}

		pos, err := intfToInt(firstParam)
		if err != nil {
			return &V{}, err
		}

		pos = v.posAtIndexForInsertBefore(pos)
		if pos < 0 {
			return &V{}, ErrOutOfRange
		}
		v.insertToArr(pos, c)
		return c, nil
	}

	// this is not the last iterarion
	child, err := v.GetArray(firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return &V{}, err
	}

	childIns := Insert{
		v: child,
		c: c,
	}
	return childIns.Before(otherParams[paramCount-1])
}

// After completes the following operation of Insert(). It inserts value AFTER specified position.
//
// The last parameter identifies the postion where a new JSON is inserted after, it should ba an interger, no matter signed or unsigned.
// If the position is zero or positive interger, it tells the index of an array. If the position is negative, it tells the backward index of an array.
//
// For example, 0 represents the first, and -2 represents the second last.
//
// After 结束并完成 Insert() 函数的后续插入操作，表示插入到指定位置的前面。
//
// 在 Before 函数的最后一个参数指定了被插入的 JSON 数组的位置，这个参数应当是一个整型（有无符号类型均可）。
// 如果这个值等于0或者正整数，那么它指定的是在 JSON 数组中的位置（从0开始）。如果这个值是负数，那么它指定的是 JSON 数组中从最后一个位置开始算起的位置。
//
// 举例说明：0 表示第一个位置，而 -2 表示倒数第二个位置。
func (ins *Insert) After(firstParam any, otherParams ...any) (*V, error) {
	if ins.err != nil {
		return &V{}, ins.err
	}
	v := ins.v
	c := ins.c
	if nil == v || v.valueType == NotExist {
		return &V{}, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(otherParams)
	if paramCount == 0 {
		if v.valueType != Array {
			return &V{}, ErrNotArrayValue
		}

		pos, err := intfToInt(firstParam)
		if err != nil {
			return &V{}, err
		}

		pos, appendToEnd := v.posAtIndexForInsertAfter(pos)
		if pos < 0 {
			return &V{}, ErrOutOfRange
		}
		if appendToEnd {
			v.appendToArr(c)
		} else {
			v.insertToArr(pos, c)
		}
		return c, nil
	}

	// this is not the last iterarion
	child, err := v.GetArray(firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return &V{}, err
	}

	childIns := Insert{
		v: child,
		c: c,
	}
	return childIns.After(otherParams[paramCount-1])
}

func (v *V) insertToArr(pos int, child *V) {
	v.children.arr = append(v.children.arr, nil)
	copy(v.children.arr[pos+1:], v.children.arr[pos:])
	v.children.arr[pos] = child
}

// ================ APPEND ================

// Append type is for InTheEnd() or InTheBeginning() function. Please refer to related functions.
//
// Should ONLY be generated by V.Append() function
//
// Append 类型是用于 InTheEnd() 和 InTheBeginning() 函数的。使用者可以不用关注这个类型。并且这个类型只应当由 V.Append() 产生。
type Append struct {
	v *V
	c *V // child

	err error
}

// Append starts appending a child JSON value to a JSON array.
//
// Append 开始将一个 JSON 值添加到一个数组中。需结合 InTheEnd() 和 InTheBeginning() 函数使用。
func (v *V) Append(child any) *Append {
	var ch *V
	var err error

	if child == nil {
		ch = NewNull()
	} else if childV, ok := child.(*V); ok {
		ch = childV
	} else {
		ch, err = Import(child)
	}
	return &Append{
		v:   v,
		c:   ch,
		err: err,
	}
}

// AppendString is equivalent to Append(jsonvalue.NewString(s))
//
// AppendString 等价于 Append(jsonvalue.NewString(s))
func (v *V) AppendString(s string) *Append {
	return v.Append(NewString(s))
}

// AppendBytes is equivalent to Append(jsonvalue.NewBytes(b))
//
// AppendBytes 等价于 Append(jsonvalue.NewBytes(b))
func (v *V) AppendBytes(b []byte) *Append {
	return v.Append(NewBytes(b))
}

// AppendBool is equivalent to Append(jsonvalue.NewBool(b))
//
// AppendBool 等价于 Append(jsonvalue.NewBool(b))
func (v *V) AppendBool(b bool) *Append {
	return v.Append(NewBool(b))
}

// AppendInt is equivalent to Append(jsonvalue.NewInt(b))
//
// AppendInt 等价于 Append(jsonvalue.NewInt(b))
func (v *V) AppendInt(i int) *Append {
	return v.Append(NewInt(i))
}

// AppendInt64 is equivalent to Append(jsonvalue.NewInt64(b))
//
// AppendInt64 等价于 Append(jsonvalue.NewInt64(b))
func (v *V) AppendInt64(i int64) *Append {
	return v.Append(NewInt64(i))
}

// AppendInt32 is equivalent to Append(jsonvalue.NewInt32(b))
//
// AppendInt32 等价于 Append(jsonvalue.NewInt32(b))
func (v *V) AppendInt32(i int32) *Append {
	return v.Append(NewInt32(i))
}

// AppendUint is equivalent to Append(jsonvalue.NewUint(b))
//
// AppendUint 等价于 Append(jsonvalue.NewUint(b))
func (v *V) AppendUint(u uint) *Append {
	return v.Append(NewUint(u))
}

// AppendUint64 is equivalent to Append(jsonvalue.NewUint64(b))
//
// AppendUint64 等价于 Append(jsonvalue.NewUint64(b))
func (v *V) AppendUint64(u uint64) *Append {
	return v.Append(NewUint64(u))
}

// AppendUint32 is equivalent to Append(jsonvalue.NewUint32(b))
//
// AppendUint32 等价于 Append(jsonvalue.NewUint32(b))
func (v *V) AppendUint32(u uint32) *Append {
	return v.Append(NewUint32(u))
}

// AppendFloat64 is equivalent to Append(jsonvalue.NewFloat64(b))
//
// AppendUint32 等价于 Append(jsonvalue.NewUint32(b))
func (v *V) AppendFloat64(f float64) *Append {
	return v.Append(NewFloat64(f))
}

// AppendFloat32 is equivalent to Append(jsonvalue.NewFloat32(b))
//
// AppendFloat32 等价于 Append(jsonvalue.NewFloat32(b))
func (v *V) AppendFloat32(f float32) *Append {
	return v.Append(NewFloat32(f))
}

// AppendNull is equivalent to Append(jsonvalue.NewNull())
//
// AppendNull 等价于 Append(jsonvalue.NewNull())
func (v *V) AppendNull() *Append {
	return v.Append(NewNull())
}

// AppendObject is equivalent to Append(jsonvalue.NewObject())
//
// AppendObject 等价于 Append(jsonvalue.NewObject())
func (v *V) AppendObject() *Append {
	return v.Append(NewObject())
}

// AppendArray is equivalent to Append(jsonvalue.NewArray())
//
// AppendArray 等价于 Append(jsonvalue.NewArray())
func (v *V) AppendArray() *Append {
	return v.Append(NewArray())
}

// InTheBeginning completes the following operation of Append().
//
// InTheBeginning 函数将 Append 函数指定的 JSON 值，添加到参数指定的数组的最前端
func (apd *Append) InTheBeginning(params ...any) (*V, error) {
	v := apd.v
	c := apd.c
	if nil == v || v.valueType == NotExist {
		return &V{}, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(params)
	if paramCount == 0 {
		if v.valueType != Array {
			return &V{}, ErrNotArrayValue
		}

		v.appendToArr(c)
		return c, nil
	}

	// this is not the last iterarion
	child, err := v.GetArray(params[0], params[1:paramCount]...)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return &V{}, err
		}
		child, err = v.SetArray().At(params[0], params[1:]...)
		if err != nil {
			return &V{}, err
		}
	}

	if child.Len() == 0 {
		child.appendToArr(c)
	} else {
		child.insertToArr(0, c)
	}
	return c, nil
}

// InTheEnd completes the following operation of Append().
//
// InTheEnd 函数将 Append 函数指定的 JSON 值，添加到参数指定的数组的最后面
func (apd *Append) InTheEnd(params ...any) (*V, error) {
	v := apd.v
	c := apd.c
	if v.valueType == NotExist {
		return &V{}, ErrValueUninitialized
	}

	// this is the last iteration
	paramCount := len(params)
	if paramCount == 0 {
		if v.valueType != Array {
			return &V{}, ErrNotArrayValue
		}

		v.appendToArr(c)
		return c, nil
	}

	// this is not the last iterarion
	child, err := v.GetArray(params[0], params[1:paramCount]...)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return &V{}, err
		}
		child, err = v.SetArray().At(params[0], params[1:]...)
		if err != nil {
			return &V{}, err
		}
	}

	child.appendToArr(c)
	return c, nil
}

// ================ DELETE ================

func (v *V) delFromObjectChildren(caseless bool, key string) (exist bool) {
	_, exist = v.children.object[key]
	if exist {
		delete(v.children.object, key)
		v.delCaselessKey(key)
		return true
	}

	if !caseless {
		return false
	}

	v.initCaselessStorage()

	lowerKey := strings.ToLower(key)
	keys, exist := v.children.lowerCaseKeys[lowerKey]
	if !exist {
		return false
	}

	for actualKey := range keys {
		_, exist = v.children.object[actualKey]
		if exist {
			delete(v.children.object, actualKey)
			v.delCaselessKey(actualKey)
			return true
		}
	}

	return false
}

// Delete deletes specified JSON value. Forexample, parameters ("data", "list") identifies deleting value in data.list.
// While ("list", 1) means deleting 2nd (count from one) element from the "list" array.
//
// Delete 从 JSON 中删除参数指定的对象。比如参数 ("data", "list") 表示删除 data.list 值；参数 ("list", 1) 则表示删除 list
// 数组的第2（从1算起）个值。
func (v *V) Delete(firstParam any, otherParams ...any) error {
	return v.delete(false, firstParam, otherParams...)
}

func (v *V) delete(caseless bool, firstParam any, otherParams ...any) error {
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

	return child.delete(caseless, otherParams[paramCount-1])
}

func (v *V) deleteInCurrValue(caseless bool, param any) error {
	if v.valueType == Object {
		// string expected
		key, err := intfToString(param)
		if err != nil {
			return err
		}

		if exist := v.delFromObjectChildren(caseless, key); !exist {
			return ErrNotFound
		}
		return nil
	}

	if v.valueType == Array {
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

	// else, this is an object value
	return fmt.Errorf("%v type does not supports Delete()", v.valueType)
}

func (v *V) deleteInArr(pos int) {
	le := len(v.children.arr)
	v.children.arr[pos] = nil
	copy(v.children.arr[pos:], v.children.arr[pos+1:])
	v.children.arr = v.children.arr[:le-1]
}
