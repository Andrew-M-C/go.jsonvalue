package jsonvalue

import (
	"errors"
	"fmt"
	"strings"
)

// ================ INSERT ================

// MARK: INSERT

// Inserter type is for After() and Before() methods.
//
// # Should be generated ONLY BY V.Insert() !
//
// Insert 类型适用于 After() 和 Before() 方法。请注意：该类型仅应由 V.Insert 函数生成！
type Inserter interface {
	// After completes the following operation of Insert(). It inserts value AFTER
	//  specified position.
	//
	// The last parameter identifies the position where a new JSON is inserted after,
	//  it should ba an integer, no matter signed or unsigned. If the position is
	// zero or positive integer, it tells the index of an array. If the position
	// is negative, it tells the backward index of an array.
	//
	// For example, 0 represents the first, and -2 represents the second last.
	//
	// After 结束并完成 Insert() 函数的后续插入操作，表示插入到指定位置的前面。
	//
	// 在 Before 函数的最后一个参数指定了被插入的 JSON 数组的位置，这个参数应当是一个整型
	// （有无符号类型均可）。
	// 如果这个值等于0或者正整数，那么它指定的是在 JSON 数组中的位置（从0开始）。如果这个值是负数,
	// 那么它指定的是 JSON 数组中从最后一个位置开始算起的位置。
	//
	// 举例说明：0 表示第一个位置，而 -2 表示倒数第二个位置。
	After(firstParam any, otherParams ...any) (*V, error)

	// Before completes the following operation of Insert(). It inserts value BEFORE
	// specified position.
	//
	// The last parameter identifies the position where a new JSON is inserted after,
	// it should ba an integer, no matter signed or unsigned.
	// If the position is zero or positive integer, it tells the index of an array.
	// If the position is negative, it tells the backward index of an array.
	//
	// For example, 0 represents the first, and -2 represents the second last.
	//
	// Before 结束并完成 Insert() 函数的后续插入操作，表示插入到指定位置的后面。
	//
	// 在 Before 函数的最后一个参数指定了被插入的 JSON 数组的位置，这个参数应当是一个整型
	//（有无符号类型均可）。
	// 如果这个值等于0或者正整数，那么它指定的是在 JSON 数组中的位置（从0开始）。如果这个值是负数,
	// 那么它指定的是 JSON 数组中从最后一个位置开始算起的位置。
	//
	// 举例说明：0 表示第一个位置，而 -2 表示倒数第二个位置。
	Before(firstParam any, otherParams ...any) (*V, error)
}

type insert struct {
	v *V
	c *V // child

	err error
}

// Insert starts inserting a child JSON value
//
// Insert 开启一个 JSON 数组成员的插入操作.
func (v *V) Insert(child any) Inserter {
	var ch *V
	var err error

	if child == nil {
		ch = NewNull()
	} else if childV, ok := child.(*V); ok {
		ch = childV
	} else {
		ch, err = Import(child)
	}

	return &insert{
		v:   v,
		c:   ch,
		err: err,
	}
}

// InsertString is equivalent to Insert(jsonvalue.NewString(s))
//
// InsertString 等效于 Insert(jsonvalue.NewString(s))
func (v *V) InsertString(s string) Inserter {
	return v.Insert(NewString(s))
}

// InsertBool is equivalent to Insert(jsonvalue.NewBool(b))
//
// InsertBool 等效于 Insert(jsonvalue.NewBool(b))
func (v *V) InsertBool(b bool) Inserter {
	return v.Insert(NewBool(b))
}

// InsertInt is equivalent to Insert(jsonvalue.NewInt(b))
//
// InsertInt 等效于 Insert(jsonvalue.NewInt(b))
func (v *V) InsertInt(i int) Inserter {
	return v.Insert(NewInt(i))
}

// InsertInt64 is equivalent to Insert(jsonvalue.NewInt64(b))
//
// InsertInt64 等效于 Insert(jsonvalue.NewInt64(b))
func (v *V) InsertInt64(i int64) Inserter {
	return v.Insert(NewInt64(i))
}

// InsertInt32 is equivalent to Insert(jsonvalue.NewInt32(b))
//
// InsertInt32 等效于 Insert(jsonvalue.NewInt32(b))
func (v *V) InsertInt32(i int32) Inserter {
	return v.Insert(NewInt32(i))
}

// InsertUint is equivalent to Insert(jsonvalue.NewUint(b))
//
// InsertUint 等效于 Insert(jsonvalue.NewUint(b))
func (v *V) InsertUint(u uint) Inserter {
	return v.Insert(NewUint(u))
}

// InsertUint64 is equivalent to Insert(jsonvalue.NewUint64(b))
//
// InsertUint64 等效于 Insert(jsonvalue.NewUint64(b))
func (v *V) InsertUint64(u uint64) Inserter {
	return v.Insert(NewUint64(u))
}

// InsertUint32 is equivalent to Insert(jsonvalue.NewUint32(b))
//
// InsertUint32 等效于 Insert(jsonvalue.NewUint32(b))
func (v *V) InsertUint32(u uint32) Inserter {
	return v.Insert(NewUint32(u))
}

// InsertFloat64 is equivalent to Insert(jsonvalue.NewFloat64(b))
//
// InsertFloat64 等效于 Insert(jsonvalue.NewFloat64(b))
func (v *V) InsertFloat64(f float64) Inserter {
	return v.Insert(NewFloat64(f))
}

// InsertFloat32 is equivalent to Insert(jsonvalue.NewFloat32(b))
//
// InsertFloat32 等效于 Insert(jsonvalue.NewFloat32(b))
func (v *V) InsertFloat32(f float32) Inserter {
	return v.Insert(NewFloat32(f))
}

// InsertNull is equivalent to Insert(jsonvalue.NewNull())
//
// InsertNull 等效于 Insert(jsonvalue.NewNull())
func (v *V) InsertNull() Inserter {
	return v.Insert(NewNull())
}

// InsertObject is equivalent to Insert(jsonvalue.NewObject())
//
// InsertObject 等效于 Insert(jsonvalue.NewObject())
func (v *V) InsertObject() Inserter {
	return v.Insert(NewObject())
}

// InsertArray is equivalent to Insert(jsonvalue.NewArray())
//
// InsertArray 等效于 Insert(jsonvalue.NewArray())
func (v *V) InsertArray() Inserter {
	return v.Insert(NewArray())
}

func (ins *insert) Before(firstParam any, otherParams ...any) (*V, error) {
	if ins.err != nil {
		return &V{}, ins.err
	}
	if ok, p1, p2 := isSliceAndExtractDividedParams(firstParam); ok {
		if len(otherParams) > 0 {
			return &V{}, ErrMultipleParamNotSupportedWithIfSliceOrArrayGiven
		}
		return ins.Before(p1, p2...)
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

		pos, err := anyToInt(firstParam)
		if err != nil {
			return &V{}, err
		}

		pos = posAtIndexForInsertBefore(v, pos)
		if pos < 0 {
			return &V{}, ErrOutOfRange
		}
		insertToArr(v, pos, c)
		return c, nil
	}

	// this is not the last iteration
	child, err := v.GetArray(firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return &V{}, err
	}

	childIns := &insert{
		v: child,
		c: c,
	}
	return childIns.Before(otherParams[paramCount-1])
}

func (ins *insert) After(firstParam any, otherParams ...any) (*V, error) {
	if ins.err != nil {
		return &V{}, ins.err
	}
	if ok, p1, p2 := isSliceAndExtractDividedParams(firstParam); ok {
		if len(otherParams) > 0 {
			return &V{}, ErrMultipleParamNotSupportedWithIfSliceOrArrayGiven
		}
		return ins.After(p1, p2...)
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

		pos, err := anyToInt(firstParam)
		if err != nil {
			return &V{}, err
		}

		pos, appendToEnd := posAtIndexForInsertAfter(v, pos)
		if pos < 0 {
			return &V{}, ErrOutOfRange
		}
		if appendToEnd {
			appendToArr(v, c)
		} else {
			insertToArr(v, pos, c)
		}
		return c, nil
	}

	// this is not the last iteration
	child, err := v.GetArray(firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return &V{}, err
	}

	childIns := &insert{
		v: child,
		c: c,
	}
	return childIns.After(otherParams[paramCount-1])
}

func insertToArr(v *V, pos int, child *V) {
	v.children.arr = append(v.children.arr, nil)
	copy(v.children.arr[pos+1:], v.children.arr[pos:])
	v.children.arr[pos] = child
}

// ================ APPEND ================

// MARK: APPEND

// Appender type is for InTheEnd() or InTheBeginning() function.
//
// Appender 类型是用于 InTheEnd() 和 InTheBeginning() 函数的。使用者可以不用关注这个类型。
// 并且这个类型只应当由 V.Append() 产生。
type Appender interface {
	InTheBeginning(params ...any) (*V, error)
	InTheEnd(params ...any) (*V, error)
}

type appender struct {
	v *V
	c *V // child

	err error
}

// Append starts appending a child JSON value to a JSON array.
//
// Append 开始将一个 JSON 值添加到一个数组中。需结合 InTheEnd() 和 InTheBeginning() 函数使用。
func (v *V) Append(child any) Appender {
	var ch *V
	var err error

	if child == nil {
		ch = NewNull()
	} else if childV, ok := child.(*V); ok {
		ch = childV
	} else {
		ch, err = Import(child)
	}
	return &appender{
		v:   v,
		c:   ch,
		err: err,
	}
}

// AppendString is equivalent to Append(jsonvalue.NewString(s))
//
// AppendString 等价于 Append(jsonvalue.NewString(s))
func (v *V) AppendString(s string) Appender {
	return v.Append(NewString(s))
}

// AppendBytes is equivalent to Append(jsonvalue.NewBytes(b))
//
// AppendBytes 等价于 Append(jsonvalue.NewBytes(b))
func (v *V) AppendBytes(b []byte) Appender {
	return v.Append(NewBytes(b))
}

// AppendBool is equivalent to Append(jsonvalue.NewBool(b))
//
// AppendBool 等价于 Append(jsonvalue.NewBool(b))
func (v *V) AppendBool(b bool) Appender {
	return v.Append(NewBool(b))
}

// AppendInt is equivalent to Append(jsonvalue.NewInt(b))
//
// AppendInt 等价于 Append(jsonvalue.NewInt(b))
func (v *V) AppendInt(i int) Appender {
	return v.Append(NewInt(i))
}

// AppendInt64 is equivalent to Append(jsonvalue.NewInt64(b))
//
// AppendInt64 等价于 Append(jsonvalue.NewInt64(b))
func (v *V) AppendInt64(i int64) Appender {
	return v.Append(NewInt64(i))
}

// AppendInt32 is equivalent to Append(jsonvalue.NewInt32(b))
//
// AppendInt32 等价于 Append(jsonvalue.NewInt32(b))
func (v *V) AppendInt32(i int32) Appender {
	return v.Append(NewInt32(i))
}

// AppendUint is equivalent to Append(jsonvalue.NewUint(b))
//
// AppendUint 等价于 Append(jsonvalue.NewUint(b))
func (v *V) AppendUint(u uint) Appender {
	return v.Append(NewUint(u))
}

// AppendUint64 is equivalent to Append(jsonvalue.NewUint64(b))
//
// AppendUint64 等价于 Append(jsonvalue.NewUint64(b))
func (v *V) AppendUint64(u uint64) Appender {
	return v.Append(NewUint64(u))
}

// AppendUint32 is equivalent to Append(jsonvalue.NewUint32(b))
//
// AppendUint32 等价于 Append(jsonvalue.NewUint32(b))
func (v *V) AppendUint32(u uint32) Appender {
	return v.Append(NewUint32(u))
}

// AppendFloat64 is equivalent to Append(jsonvalue.NewFloat64(b))
//
// AppendUint32 等价于 Append(jsonvalue.NewUint32(b))
func (v *V) AppendFloat64(f float64) Appender {
	return v.Append(NewFloat64(f))
}

// AppendFloat32 is equivalent to Append(jsonvalue.NewFloat32(b))
//
// AppendFloat32 等价于 Append(jsonvalue.NewFloat32(b))
func (v *V) AppendFloat32(f float32) Appender {
	return v.Append(NewFloat32(f))
}

// AppendNull is equivalent to Append(jsonvalue.NewNull())
//
// AppendNull 等价于 Append(jsonvalue.NewNull())
func (v *V) AppendNull() Appender {
	return v.Append(NewNull())
}

// AppendObject is equivalent to Append(jsonvalue.NewObject())
//
// AppendObject 等价于 Append(jsonvalue.NewObject())
func (v *V) AppendObject() Appender {
	return v.Append(NewObject())
}

// AppendArray is equivalent to Append(jsonvalue.NewArray())
//
// AppendArray 等价于 Append(jsonvalue.NewArray())
func (v *V) AppendArray() Appender {
	return v.Append(NewArray())
}

// InTheBeginning completes the following operation of Append().
//
// InTheBeginning 函数将 Append 函数指定的 JSON 值，添加到参数指定的数组的最前端
func (apd *appender) InTheBeginning(params ...any) (*V, error) {
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
		if v.Len() == 0 {
			appendToArr(v, c)
		} else {
			insertToArr(v, 0, c)
		}
		return c, nil
	}
	if ok, p := isSliceAndExtractJointParams(params[0]); ok {
		if len(params) > 1 {
			return &V{}, ErrMultipleParamNotSupportedWithIfSliceOrArrayGiven
		}
		return apd.InTheBeginning(p...)
	}

	// this is not the last iteration
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
		appendToArr(child, c)
	} else {
		insertToArr(child, 0, c)
	}
	return c, nil
}

// InTheEnd completes the following operation of Append().
//
// InTheEnd 函数将 Append 函数指定的 JSON 值，添加到参数指定的数组的最后面
func (apd *appender) InTheEnd(params ...any) (*V, error) {
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

		appendToArr(v, c)
		return c, nil
	}
	if ok, p := isSliceAndExtractJointParams(params[0]); ok {
		if len(params) > 1 {
			return &V{}, ErrMultipleParamNotSupportedWithIfSliceOrArrayGiven
		}
		return apd.InTheEnd(p...)
	}

	// this is not the last iteration
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

	appendToArr(child, c)
	return c, nil
}

// ================ DELETE ================

// MARK: DELETE

func delFromObjectChildren(v *V, caseless bool, key string) (exist bool) {
	_, exist = v.children.object[key]
	if exist {
		delete(v.children.object, key)
		delCaselessKey(v, key)
		return true
	}

	if !caseless {
		return false
	}

	initCaselessStorage(v)

	lowerKey := strings.ToLower(key)
	keys, exist := v.children.lowerCaseKeys[lowerKey]
	if !exist {
		return false
	}

	for actualKey := range keys {
		_, exist = v.children.object[actualKey]
		if exist {
			delete(v.children.object, actualKey)
			delCaselessKey(v, actualKey)
			return true
		}
	}

	return false
}

// Delete deletes specified JSON value. For example, parameters ("data", "list") identifies deleting value in data.list.
// While ("list", 1) means deleting the second element from the "list" array.
//
// Delete 从 JSON 中删除参数指定的对象。比如参数 ("data", "list") 表示删除 data.list 值；参数 ("list", 1) 则表示删除 list
// 数组的第2（从1算起）个值。
func (v *V) Delete(firstParam any, otherParams ...any) error {
	return v.delete(false, firstParam, otherParams...)
}

func (v *V) delete(caseless bool, firstParam any, otherParams ...any) error {
	if ok, p1, p2 := isSliceAndExtractDividedParams(firstParam); ok {
		if len(otherParams) > 0 {
			return ErrMultipleParamNotSupportedWithIfSliceOrArrayGiven
		}
		return v.delete(caseless, p1, p2...)
	}

	paramCount := len(otherParams)
	if paramCount == 0 {
		return deleteInCurrentValue(v, caseless, firstParam)
	}

	child, err := get(v, caseless, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	// if child == nil {
	// 	return ErrNotFound
	// }

	return child.delete(caseless, otherParams[paramCount-1])
}

func deleteInCurrentValue(v *V, caseless bool, param any) error {
	switch v.valueType {
	case Object:
		return deleteInCurrentObject(v, caseless, param)
	case Array:
		return deleteInCurrentArray(v, param)
	default:
		// else, this is an object value
		return fmt.Errorf("%v type does not supports Delete()", v.valueType)
	}
}

func deleteInCurrentObject(v *V, caseless bool, param any) error {
	// string expected
	key, err := anyToString(param)
	if err != nil {
		return err
	}
	if exist := delFromObjectChildren(v, caseless, key); !exist {
		return ErrNotFound
	}
	return nil
}

func deleteInCurrentArray(v *V, param any) error {
	// integer expected
	pos, err := anyToInt(param)
	if err != nil {
		return err
	}
	pos = posAtIndexForRead(v, pos)
	if pos < 0 {
		return ErrOutOfRange
	}
	deleteInArr(v, pos)
	return nil
}

func deleteInArr(v *V, pos int) {
	le := len(v.children.arr)
	v.children.arr[pos] = nil
	copy(v.children.arr[pos:], v.children.arr[pos+1:])
	v.children.arr = v.children.arr[:le-1]
}
