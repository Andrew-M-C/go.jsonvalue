package jsonvalue

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// Len returns length of an object or array type JSON value
func (v *V) Len() int {
	switch v.valueType {
	case jsonparser.Array:
		return v.arrayChildren.Len()
	case jsonparser.Object:
		return len(v.objectChildren)
	default:
		return 0
	}
}

// Get returns JSON value in specified position. Param formats are like At()
func (v *V) Get(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	child, err := v.getInCurrValue(firstParam)
	if err != nil {
		return nil, err
	}

	if 0 == len(otherParams) {
		return child, nil
	}
	return child.Get(otherParams[0], otherParams[1:]...)
}

func (v *V) getFromObjectChildren(key string) (child *V, exist bool) {
	child, exist = v.objectChildren[key]
	if exist {
		return child, true
	}

	lowerCaseKey := strings.ToLower(key)
	keys, exist := v.lowerCaseKeys[lowerCaseKey]
	if !exist {
		return nil, false
	}

	for actualKey := range keys {
		child, exist = v.objectChildren[actualKey]
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

// GetString is equalivent to v, err := Get(...); v.String()
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

// GetInt is equalivent to v, err := Get(...); v.Int()
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

// GetUint is equalivent to v, err := Get(...); v.Int()
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

// GetInt64 is equalivent to v, err := Get(...); v.Int()
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

// GetUint64 is equalivent to v, err := Get(...); v.Int()
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

// GetInt32 is equalivent to v, err := Get(...); v.Int()
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

// GetUint32 is equalivent to v, err := Get(...); v.Int()
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

// GetFloat64 is equalivent to v, err := Get(...); v.Int()
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

// GetFloat32 is equalivent to v, err := Get(...); v.Int()
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

// GetBool is equalivent to v, err := Get(...); v.Bool()
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

// GetNull is equalivent to v, err := Get(...); raise err if v.IsNull() == false
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

// GetObject is equalivent to v, err := Get(...); raise err if v.IsObject() == false
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

// GetArray is equalivent to v, err := Get(...); raise err if v.IsArray() == false
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
