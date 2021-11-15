package jsonvalue

import (
	"fmt"
	"strings"
)

// Len returns length of an object or array type JSON value.
//
// Len 返回当前对象类型或数组类型的 JSON 的成员长度。如果不是这两种类型，那么会返回 0。
func (v *V) Len() int {
	switch v.valueType {
	case Array:
		return len(v.children.array)
	case Object:
		return len(v.children.object)
	default:
		return 0
	}
}

// Get returns JSON value in specified position. Param formats are like At().
//
// Get 返回按照参数指定的位置的 JSON 成员值。参数格式与 At() 函数相同
func (v *V) Get(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return v.get(false, firstParam, otherParams...)
}

// MustGet is same as Get(), but does not return error. If error occurs, a JSON value with
// NotExist type will be returned.
//
// MustGet 与 Get() 函数相同，不过不返回错误。如果发生错误了，那么会返回一个 ValueType() 返回值为 NotExist
// 的 JSON 值对象。
func (v *V) MustGet(firstParam interface{}, otherParams ...interface{}) *V {
	res, _ := v.get(false, firstParam, otherParams...)
	return res
}

func (v *V) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	child, err := v.getInCurrValue(caseless, firstParam)
	if err != nil {
		return &V{}, err
	}

	if len(otherParams) == 0 {
		return child, nil
	}
	return child.get(caseless, otherParams[0], otherParams[1:]...)
}

func (v *V) initCaselessStorage() {
	if v.children.lowerCaseKeys != nil {
		return
	}
	v.children.lowerCaseKeys = make(map[string]map[string]struct{}, len(v.children.object))
	for k := range v.children.object {
		v.addCaselessKey(k)
	}
}

func (v *V) getFromObjectChildren(caseless bool, key string) (child *V, exist bool) {
	child, exist = v.children.object[key]
	if exist {
		return child, true
	}

	if !caseless {
		return &V{}, false
	}

	v.initCaselessStorage()

	lowerCaseKey := strings.ToLower(key)
	keys, exist := v.children.lowerCaseKeys[lowerCaseKey]
	if !exist {
		return &V{}, false
	}

	for actualKey := range keys {
		child, exist = v.children.object[actualKey]
		if exist {
			return child, true
		}
	}

	return &V{}, false
}

func (v *V) getInCurrValue(caseless bool, param interface{}) (*V, error) {
	if v.valueType == Array {
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

	} else if v.valueType == Object {
		// string expected
		key, err := intfToString(param)
		if err != nil {
			return &V{}, err
		}
		child, exist := v.getFromObjectChildren(caseless, key)
		if !exist {
			return &V{}, ErrNotFound
		}
		return child, nil

	} else {
		return &V{}, fmt.Errorf("%v type does not supports Get()", v.valueType)
	}
}

// GetBytes is similar with v, err := Get(...); v.Bytes(). But if error occurs or Base64 decode error, returns error.
//
// GetBytes 类似于 v, err := Get(...); v.Bytes()，但如果查询中发生错误，或者 base64 解码错误，则返回错误。
func (v *V) GetBytes(firstParam interface{}, otherParams ...interface{}) ([]byte, error) {
	return v.getBytes(false, firstParam, otherParams...)
}

func (v *V) getBytes(caseless bool, firstParam interface{}, otherParams ...interface{}) ([]byte, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return []byte{}, err
	}
	if ret.valueType != String {
		return []byte{}, ErrTypeNotMatch
	}
	b, err := b64.DecodeString(ret.valueStr)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// GetString is equalivent to v, err := Get(...); v.String(). If error occurs, returns "".
//
// GetString 等效于 v, err := Get(...); v.String()。如果发生错误，则返回 ""。
func (v *V) GetString(firstParam interface{}, otherParams ...interface{}) (string, error) {
	return v.getString(false, firstParam, otherParams...)
}

func (v *V) getString(caseless bool, firstParam interface{}, otherParams ...interface{}) (string, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return "", err
	}
	if ret.valueType != String {
		return "", ErrTypeNotMatch
	}
	return ret.String(), nil
}

// GetInt is equalivent to v, err := Get(...); v.Int(). If error occurs, returns 0.
//
// GetInt 等效于 v, err := Get(...); v.Int()。如果发生错误，则返回 0。
func (v *V) GetInt(firstParam interface{}, otherParams ...interface{}) (int, error) {
	return v.getInt(false, firstParam, otherParams...)
}

func (v *V) getInt(caseless bool, firstParam interface{}, otherParams ...interface{}) (int, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Int(), err
}

// GetUint is equalivent to v, err := Get(...); v.Uint(). If error occurs, returns 0.
//
// GetUint 等效于 v, err := Get(...); v.Uint()。如果发生错误，则返回 0。
func (v *V) GetUint(firstParam interface{}, otherParams ...interface{}) (uint, error) {
	return v.getUint(false, firstParam, otherParams...)
}

func (v *V) getUint(caseless bool, firstParam interface{}, otherParams ...interface{}) (uint, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Uint(), err
}

// GetInt64 is equalivent to v, err := Get(...); v.Int64(). If error occurs, returns 0.
//
// GetInt64 等效于 v, err := Get(...); v.Int64()。如果发生错误，则返回 0。
func (v *V) GetInt64(firstParam interface{}, otherParams ...interface{}) (int64, error) {
	return v.getInt64(false, firstParam, otherParams...)
}

func (v *V) getInt64(caseless bool, firstParam interface{}, otherParams ...interface{}) (int64, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Int64(), err
}

// GetUint64 is equalivent to v, err := Get(...); v.Unt64(). If error occurs, returns 0.
//
// GetUint64 等效于 v, err := Get(...); v.Unt64()。如果发生错误，则返回 0。
func (v *V) GetUint64(firstParam interface{}, otherParams ...interface{}) (uint64, error) {
	return v.getUint64(false, firstParam, otherParams...)
}

func (v *V) getUint64(caseless bool, firstParam interface{}, otherParams ...interface{}) (uint64, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Uint64(), err
}

// GetInt32 is equalivent to v, err := Get(...); v.Int32(). If error occurs, returns 0.
//
// GetInt32 等效于 v, err := Get(...); v.Int32()。如果发生错误，则返回 0。
func (v *V) GetInt32(firstParam interface{}, otherParams ...interface{}) (int32, error) {
	return v.getInt32(false, firstParam, otherParams...)
}

func (v *V) getInt32(caseless bool, firstParam interface{}, otherParams ...interface{}) (int32, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Int32(), err
}

// GetUint32 is equalivent to v, err := Get(...); v.Uint32(). If error occurs, returns 0.
//
// GetUint32 等效于 v, err := Get(...); v.Uint32()。如果发生错误，则返回 0。
func (v *V) GetUint32(firstParam interface{}, otherParams ...interface{}) (uint32, error) {
	return v.getUint32(false, firstParam, otherParams...)
}

func (v *V) getUint32(caseless bool, firstParam interface{}, otherParams ...interface{}) (uint32, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Uint32(), err
}

// GetFloat64 is equalivent to v, err := Get(...); v.Float64(). If error occurs, returns 0.0.
//
// GetFloat64 等效于 v, err := Get(...); v.Float64()。如果发生错误，则返回 0.0。
func (v *V) GetFloat64(firstParam interface{}, otherParams ...interface{}) (float64, error) {
	return v.getFloat64(false, firstParam, otherParams...)
}

func (v *V) getFloat64(caseless bool, firstParam interface{}, otherParams ...interface{}) (float64, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Float64(), err
}

// GetFloat32 is equalivent to v, err := Get(...); v.Float32(). If error occurs, returns 0.0.
//
// GetFloat32 等效于 v, err := Get(...); v.Float32()。如果发生错误，则返回 0.0。
func (v *V) GetFloat32(firstParam interface{}, otherParams ...interface{}) (float32, error) {
	return v.getFloat32(false, firstParam, otherParams...)
}

func (v *V) getFloat32(caseless bool, firstParam interface{}, otherParams ...interface{}) (float32, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Float32(), err
}

// GetBool is equalivent to v, err := Get(...); v.Bool(). If error occurs, returns false.
//
// GetBool 等效于 v, err := Get(...); v.Bool()。如果发生错误，则返回 false。
func (v *V) GetBool(firstParam interface{}, otherParams ...interface{}) (bool, error) {
	return v.getBool(false, firstParam, otherParams...)
}

func (v *V) getBool(caseless bool, firstParam interface{}, otherParams ...interface{}) (bool, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return false, err
	}
	ret, err = getBoolAndErrorFromValue(ret)
	return ret.Bool(), err
}

// GetNull is equalivent to v, err := Get(...); raise err if error occurs or v.IsNull() == false.
//
// GetNull 等效于 v, err := Get(...);，如果发生错误或者 v.IsNull() == false 则返回错误。
func (v *V) GetNull(firstParam interface{}, otherParams ...interface{}) error {
	return v.getNull(false, firstParam, otherParams...)
}

func (v *V) getNull(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return err
	}
	if ret.valueType != Null {
		return ErrTypeNotMatch
	}
	return nil
}

// GetObject is equalivent to v, err := Get(...); raise err if error occurs or v.IsObject() == false.
//
// GetObject 等效于 v, err := Get(...);，如果发生错误或者 v.IsObject() == false 则返回错误。
func (v *V) GetObject(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return v.getObject(false, firstParam, otherParams...)
}

func (v *V) getObject(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return &V{}, err
	}
	if ret.valueType != Object {
		return &V{}, ErrTypeNotMatch
	}
	return ret, nil
}

// GetArray is equalivent to v, err := Get(...); raise err if or v.IsArray() == false.
//
// GetArray 等效于 v, err := Get(...);，如果发生错误或者 v.IsArray() == false 则返回错误。
func (v *V) GetArray(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return v.getArray(false, firstParam, otherParams...)
}

func (v *V) getArray(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return &V{}, err
	}
	if ret.valueType != Array {
		return &V{}, ErrTypeNotMatch
	}
	return ret, nil
}
