package jsonvalue

import (
	"bytes"
	"fmt"
	"strings"
)

// ================ GET ================

// Len returns length of an object or array type JSON value.
//
// Len 返回当前对象类型或数组类型的 JSON 的成员长度。如果不是这两种类型，那么会返回 0。
func (v *V) Len() int {
	switch v.valueType {
	case Array:
		return len(v.children.arr)
	case Object:
		return len(v.children.object)
	default:
		return 0
	}
}

// Get returns JSON value in specified position. Param formats are like At().
//
// Get 返回按照参数指定的位置的 JSON 成员值。参数格式与 At() 函数相同
func (v *V) Get(firstParam any, otherParams ...any) (*V, error) {
	return v.get(false, firstParam, otherParams...)
}

// MustGet is same as Get(), but does not return error. If error occurs, a JSON value with
// NotExist type will be returned.
//
// MustGet 与 Get() 函数相同，不过不返回错误。如果发生错误了，那么会返回一个 ValueType() 返回值为 NotExist
// 的 JSON 值对象。
func (v *V) MustGet(firstParam any, otherParams ...any) *V {
	res, _ := v.get(false, firstParam, otherParams...)
	return res
}

func (v *V) get(caseless bool, firstParam any, otherParams ...any) (*V, error) {
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
	childProperty, exist := v.children.object[key]
	if exist {
		return childProperty.v, true
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
		childProperty, exist = v.children.object[actualKey]
		if exist {
			return childProperty.v, true
		}
	}

	return &V{}, false
}

func (v *V) getInCurrValue(caseless bool, param any) (*V, error) {
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
func (v *V) GetBytes(firstParam any, otherParams ...any) ([]byte, error) {
	return v.getBytes(false, firstParam, otherParams...)
}

func (v *V) getBytes(caseless bool, firstParam any, otherParams ...any) ([]byte, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return []byte{}, err
	}
	if ret.valueType != String {
		return []byte{}, ErrTypeNotMatch
	}
	b, err := internal.b64.DecodeString(ret.valueStr)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// GetString is equalivent to v, err := Get(...); v.String(). If error occurs, returns "".
//
// GetString 等效于 v, err := Get(...); v.String()。如果发生错误，则返回 ""。
func (v *V) GetString(firstParam any, otherParams ...any) (string, error) {
	return v.getString(false, firstParam, otherParams...)
}

func (v *V) getString(caseless bool, firstParam any, otherParams ...any) (string, error) {
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
func (v *V) GetInt(firstParam any, otherParams ...any) (int, error) {
	return v.getInt(false, firstParam, otherParams...)
}

func (v *V) getInt(caseless bool, firstParam any, otherParams ...any) (int, error) {
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
func (v *V) GetUint(firstParam any, otherParams ...any) (uint, error) {
	return v.getUint(false, firstParam, otherParams...)
}

func (v *V) getUint(caseless bool, firstParam any, otherParams ...any) (uint, error) {
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
func (v *V) GetInt64(firstParam any, otherParams ...any) (int64, error) {
	return v.getInt64(false, firstParam, otherParams...)
}

func (v *V) getInt64(caseless bool, firstParam any, otherParams ...any) (int64, error) {
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
func (v *V) GetUint64(firstParam any, otherParams ...any) (uint64, error) {
	return v.getUint64(false, firstParam, otherParams...)
}

func (v *V) getUint64(caseless bool, firstParam any, otherParams ...any) (uint64, error) {
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
func (v *V) GetInt32(firstParam any, otherParams ...any) (int32, error) {
	return v.getInt32(false, firstParam, otherParams...)
}

func (v *V) getInt32(caseless bool, firstParam any, otherParams ...any) (int32, error) {
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
func (v *V) GetUint32(firstParam any, otherParams ...any) (uint32, error) {
	return v.getUint32(false, firstParam, otherParams...)
}

func (v *V) getUint32(caseless bool, firstParam any, otherParams ...any) (uint32, error) {
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
func (v *V) GetFloat64(firstParam any, otherParams ...any) (float64, error) {
	return v.getFloat64(false, firstParam, otherParams...)
}

func (v *V) getFloat64(caseless bool, firstParam any, otherParams ...any) (float64, error) {
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
func (v *V) GetFloat32(firstParam any, otherParams ...any) (float32, error) {
	return v.getFloat32(false, firstParam, otherParams...)
}

func (v *V) getFloat32(caseless bool, firstParam any, otherParams ...any) (float32, error) {
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
func (v *V) GetBool(firstParam any, otherParams ...any) (bool, error) {
	return v.getBool(false, firstParam, otherParams...)
}

func (v *V) getBool(caseless bool, firstParam any, otherParams ...any) (bool, error) {
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
func (v *V) GetNull(firstParam any, otherParams ...any) error {
	return v.getNull(false, firstParam, otherParams...)
}

func (v *V) getNull(caseless bool, firstParam any, otherParams ...any) error {
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
func (v *V) GetObject(firstParam any, otherParams ...any) (*V, error) {
	return v.getObject(false, firstParam, otherParams...)
}

func (v *V) getObject(caseless bool, firstParam any, otherParams ...any) (*V, error) {
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
func (v *V) GetArray(firstParam any, otherParams ...any) (*V, error) {
	return v.getArray(false, firstParam, otherParams...)
}

func (v *V) getArray(caseless bool, firstParam any, otherParams ...any) (*V, error) {
	ret, err := v.get(caseless, firstParam, otherParams...)
	if err != nil {
		return &V{}, err
	}
	if ret.valueType != Array {
		return &V{}, ErrTypeNotMatch
	}
	return ret, nil
}

// ================ CASELESS ================

// Caseless is returned by Caseless(). operations of Caseless type are same as (*V).Get(), but are via caseless key.
//
// Caseless 类型通过 Caseless() 函数返回。通过 Caseless 接口操作的所有操作均与 (*v).Get() 相同，但是对 key 进行读取的时候，
// 不区分大小写。
type Caseless interface {
	Get(firstParam any, otherParams ...any) (*V, error)
	MustGet(firstParam any, otherParams ...any) *V
	GetBytes(firstParam any, otherParams ...any) ([]byte, error)
	GetString(firstParam any, otherParams ...any) (string, error)
	GetInt(firstParam any, otherParams ...any) (int, error)
	GetUint(firstParam any, otherParams ...any) (uint, error)
	GetInt64(firstParam any, otherParams ...any) (int64, error)
	GetUint64(firstParam any, otherParams ...any) (uint64, error)
	GetInt32(firstParam any, otherParams ...any) (int32, error)
	GetUint32(firstParam any, otherParams ...any) (uint32, error)
	GetFloat64(firstParam any, otherParams ...any) (float64, error)
	GetFloat32(firstParam any, otherParams ...any) (float32, error)
	GetBool(firstParam any, otherParams ...any) (bool, error)
	GetNull(firstParam any, otherParams ...any) error
	GetObject(firstParam any, otherParams ...any) (*V, error)
	GetArray(firstParam any, otherParams ...any) (*V, error)

	Delete(firstParam any, otherParams ...any) error
	MustDelete(firstParam any, otherParams ...any)
}

var _ Caseless = (*V)(nil)

// Caseless returns Caseless interface to support caseless getting.
//
// IMPORTANT: This function is not gouroutine-safe. Write-mutex (instead of read-mutex) should be attached in cross-goroutine scenarios.
//
// Caseless 返回 Caseless 接口，从而实现不区分大小写的 Get 操作。
//
// 注意: 该函数不是协程安全的，如果在多协程场景下，调用该函数，需要加上写锁，而不能用读锁。
func (v *V) Caseless() Caseless {
	switch v.valueType {
	default:
		return v

	case Array, Object:
		return &caselessOper{
			v: v,
		}
	}
}

type caselessOper struct {
	v *V
}

func (g *caselessOper) Get(firstParam any, otherParams ...any) (*V, error) {
	return g.v.get(true, firstParam, otherParams...)
}

func (g *caselessOper) MustGet(firstParam any, otherParams ...any) *V {
	res, _ := g.v.get(true, firstParam, otherParams...)
	return res
}

func (g *caselessOper) GetBytes(firstParam any, otherParams ...any) ([]byte, error) {
	return g.v.getBytes(true, firstParam, otherParams...)
}

func (g *caselessOper) GetString(firstParam any, otherParams ...any) (string, error) {
	return g.v.getString(true, firstParam, otherParams...)
}

func (g *caselessOper) GetInt(firstParam any, otherParams ...any) (int, error) {
	return g.v.getInt(true, firstParam, otherParams...)
}

func (g *caselessOper) GetUint(firstParam any, otherParams ...any) (uint, error) {
	return g.v.getUint(true, firstParam, otherParams...)
}

func (g *caselessOper) GetInt64(firstParam any, otherParams ...any) (int64, error) {
	return g.v.getInt64(true, firstParam, otherParams...)
}

func (g *caselessOper) GetUint64(firstParam any, otherParams ...any) (uint64, error) {
	return g.v.getUint64(true, firstParam, otherParams...)
}

func (g *caselessOper) GetInt32(firstParam any, otherParams ...any) (int32, error) {
	return g.v.getInt32(true, firstParam, otherParams...)
}

func (g *caselessOper) GetUint32(firstParam any, otherParams ...any) (uint32, error) {
	return g.v.getUint32(true, firstParam, otherParams...)
}

func (g *caselessOper) GetFloat64(firstParam any, otherParams ...any) (float64, error) {
	return g.v.getFloat64(true, firstParam, otherParams...)
}

func (g *caselessOper) GetFloat32(firstParam any, otherParams ...any) (float32, error) {
	return g.v.getFloat32(true, firstParam, otherParams...)
}

func (g *caselessOper) GetBool(firstParam any, otherParams ...any) (bool, error) {
	return g.v.getBool(true, firstParam, otherParams...)
}

func (g *caselessOper) GetNull(firstParam any, otherParams ...any) error {
	return g.v.getNull(true, firstParam, otherParams...)
}

func (g *caselessOper) GetObject(firstParam any, otherParams ...any) (*V, error) {
	return g.v.getObject(true, firstParam, otherParams...)
}

func (g *caselessOper) GetArray(firstParam any, otherParams ...any) (*V, error) {
	return g.v.getArray(true, firstParam, otherParams...)
}

func (g *caselessOper) Delete(firstParam any, otherParams ...any) error {
	return g.v.delete(true, firstParam, otherParams...)
}

func (g *caselessOper) MustDelete(firstParam any, otherParams ...any) {
	_ = g.v.delete(true, firstParam, otherParams...)
}

// ==== internal value access functions ====

func getNumberFromNotNumberValue(v *V) *V {
	if !v.IsString() {
		return NewInt(0)
	}
	ret, _ := newFromNumber(globalPool{}, bytes.TrimSpace([]byte(v.valueStr)))
	err := ret.parseNumber(globalPool{})
	if err != nil {
		return NewInt64(0)
	}
	return ret
}

func getNumberAndErrorFromValue(v *V) (*V, error) {
	switch v.valueType {
	default:
		return NewInt(0), ErrTypeNotMatch

	case Number:
		return v, nil

	case String:
		ret, _ := newFromNumber(globalPool{}, bytes.TrimSpace([]byte(v.valueStr)))
		err := ret.parseNumber(globalPool{})
		if err != nil {
			return NewInt(0), fmt.Errorf("%w: %v", ErrParseNumberFromString, err)
		}
		return ret, ErrTypeNotMatch
	}
}

func getBoolAndErrorFromValue(v *V) (*V, error) {
	switch v.valueType {
	default:
		return NewBool(false), ErrTypeNotMatch

	case Number:
		return NewBool(v.Float64() != 0), ErrTypeNotMatch

	case String:
		if v.valueStr == "true" {
			return NewBool(true), ErrTypeNotMatch
		}
		return NewBool(false), ErrTypeNotMatch

	case Boolean:
		return v, nil
	}
}
