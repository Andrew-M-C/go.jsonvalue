package jsonvalue

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// Len returns length of an object or array type JSON value.
//
// Len 返回当前对象类型或数组类型的 JSON 的成员长度。如果不是这两种类型，那么会返回 0。
func (v *V) Len() int {
	switch v.valueType {
	case jsonparser.Array:
		return v.children.array.Len()
	case jsonparser.Object:
		return len(v.children.object)
	default:
		return 0
	}
}

// Get returns JSON value in specified position. Param formats are like At().
//
// Get 返回按照参数指定的位置的 JSON 成员值。参数格式与 At() 函数相同
func (v *V) Get(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	child, err := v.getInCurrValue(firstParam)
	if err != nil {
		return nil, err
	}

	if len(otherParams) == 0 {
		return child, nil
	}
	return child.Get(otherParams[0], otherParams[1:]...)
}

func (v *V) getFromObjectChildren(key string) (child *V, exist bool) {
	child, exist = v.children.object[key]
	if exist {
		return child, true
	}

	lowerCaseKey := strings.ToLower(key)
	keys, exist := v.children.lowerCaseKeys[lowerCaseKey]
	if !exist {
		return nil, false
	}

	for actualKey := range keys {
		child, exist = v.children.object[actualKey]
		if exist {
			return child, true
		}
	}

	return nil, false
}

func (v *V) getInCurrValue(param interface{}) (*V, error) {
	if v.valueType == jsonparser.Array {
		// integer expected
		pos, err := intfToInt(param)
		if err != nil {
			return nil, err
		}
		child, _ := v.childAtIndex(pos)
		if nil == child {
			return nil, ErrOutOfRange
		}
		return child, nil

	} else if v.valueType == jsonparser.Object {
		// string expected
		key, err := intfToString(param)
		if err != nil {
			return nil, err
		}
		child, exist := v.getFromObjectChildren(key)
		if !exist {
			return nil, ErrNotFound
		}
		return child, nil

	} else {
		return nil, fmt.Errorf("%v type does not supports Get()", v.valueType)
	}
}

// GetString is equalivent to v, err := Get(...); v.String(). If error occurs, returns "".
//
// GetString 等效于 v, err := Get(...); v.String()。如果发生错误，则返回 ""。
func (v *V) GetString(firstParam interface{}, otherParams ...interface{}) (string, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return "", err
	}
	if ret.valueType != jsonparser.String {
		return "", ErrTypeNotMatch
	}
	return ret.String(), nil
}

// GetInt is equalivent to v, err := Get(...); v.Int(). If error occurs, returns 0.
//
// GetInt 等效于 v, err := Get(...); v.Int()。如果发生错误，则返回 0。
func (v *V) GetInt(firstParam interface{}, otherParams ...interface{}) (int, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Int(), nil
}

// GetUint is equalivent to v, err := Get(...); v.Uint(). If error occurs, returns 0.
//
// GetUint 等效于 v, err := Get(...); v.Uint()。如果发生错误，则返回 0。
func (v *V) GetUint(firstParam interface{}, otherParams ...interface{}) (uint, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Uint(), nil
}

// GetInt64 is equalivent to v, err := Get(...); v.Int64(). If error occurs, returns 0.
//
// GetInt64 等效于 v, err := Get(...); v.Int64()。如果发生错误，则返回 0。
func (v *V) GetInt64(firstParam interface{}, otherParams ...interface{}) (int64, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Int64(), nil
}

// GetUint64 is equalivent to v, err := Get(...); v.Unt64(). If error occurs, returns 0.
//
// GetUint64 等效于 v, err := Get(...); v.Unt64()。如果发生错误，则返回 0。
func (v *V) GetUint64(firstParam interface{}, otherParams ...interface{}) (uint64, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Uint64(), nil
}

// GetInt32 is equalivent to v, err := Get(...); v.Int32(). If error occurs, returns 0.
//
// GetInt32 等效于 v, err := Get(...); v.Int32()。如果发生错误，则返回 0。
func (v *V) GetInt32(firstParam interface{}, otherParams ...interface{}) (int32, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Int32(), nil
}

// GetUint32 is equalivent to v, err := Get(...); v.Uint32(). If error occurs, returns 0.
//
// GetUint32 等效于 v, err := Get(...); v.Uint32()。如果发生错误，则返回 0。
func (v *V) GetUint32(firstParam interface{}, otherParams ...interface{}) (uint32, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Uint32(), nil
}

// GetFloat64 is equalivent to v, err := Get(...); v.Float64(). If error occurs, returns 0.0.
//
// GetFloat64 等效于 v, err := Get(...); v.Float64()。如果发生错误，则返回 0.0。
func (v *V) GetFloat64(firstParam interface{}, otherParams ...interface{}) (float64, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Float64(), nil
}

// GetFloat32 is equalivent to v, err := Get(...); v.Float32(). If error occurs, returns 0.0.
//
// GetFloat32 等效于 v, err := Get(...); v.Float32()。如果发生错误，则返回 0.0。
func (v *V) GetFloat32(firstParam interface{}, otherParams ...interface{}) (float32, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	if ret.valueType != jsonparser.Number {
		return 0, ErrTypeNotMatch
	}
	return ret.Float32(), nil
}

// GetBool is equalivent to v, err := Get(...); v.Bool(). If error occurs, returns false.
//
// GetBool 等效于 v, err := Get(...); v.Bool()。如果发生错误，则返回 false。
func (v *V) GetBool(firstParam interface{}, otherParams ...interface{}) (bool, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return false, err
	}
	if ret.valueType != jsonparser.Boolean {
		return false, ErrTypeNotMatch
	}
	return ret.Bool(), nil
}

// GetNull is equalivent to v, err := Get(...); raise err if error occurs or v.IsNull() == false.
//
// GetNull 等效于 v, err := Get(...);，如果发生错误或者 v.IsNull() == false 则返回错误。
func (v *V) GetNull(firstParam interface{}, otherParams ...interface{}) error {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return err
	}
	if ret.valueType != jsonparser.Null {
		return ErrTypeNotMatch
	}
	return nil
}

// GetObject is equalivent to v, err := Get(...); raise err if error occurs or v.IsObject() == false.
//
// GetObject 等效于 v, err := Get(...);，如果发生错误或者 v.IsObject() == false 则返回错误。
func (v *V) GetObject(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return nil, err
	}
	if ret.valueType != jsonparser.Object {
		return nil, ErrTypeNotMatch
	}
	return ret, nil
}

// GetArray is equalivent to v, err := Get(...); raise err if or v.IsArray() == false.
//
// GetArray 等效于 v, err := Get(...);，如果发生错误或者 v.IsArray() == false 则返回错误。
func (v *V) GetArray(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	ret, err := v.Get(firstParam, otherParams...)
	if err != nil {
		return nil, err
	}
	if ret.valueType != jsonparser.Array {
		return nil, ErrTypeNotMatch
	}
	return ret, nil
}
