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
	return get(v, false, firstParam, otherParams...)
}

// MustGet is same as Get(), but does not return error. If error occurs, a JSON value with
// NotExist type will be returned.
//
// MustGet 与 Get() 函数相同，不过不返回错误。如果发生错误了，那么会返回一个 ValueType() 返回值为 NotExist
// 的 JSON 值对象。
func (v *V) MustGet(firstParam any, otherParams ...any) *V {
	res, _ := get(v, false, firstParam, otherParams...)
	return res
}

func get(v *V, caseless bool, firstParam any, otherParams ...any) (*V, error) {
	if ok, p1, p2 := isSliceAndExtractDividedParams(firstParam); ok {
		if len(otherParams) > 0 {
			return &V{}, ErrMultipleParamNotSupportedWithIfSliceOrArrayGiven
		}
		return get(v, caseless, p1, p2...)
	}
	child, err := getInCurrentValue(v, caseless, firstParam)
	if err != nil {
		return &V{}, err
	}

	if len(otherParams) == 0 {
		return child, nil
	}
	return get(child, caseless, otherParams[0], otherParams[1:]...)
}

func initCaselessStorage(v *V) {
	if v.children.lowerCaseKeys != nil {
		return
	}
	v.children.lowerCaseKeys = make(map[string]map[string]struct{}, len(v.children.object))
	for k := range v.children.object {
		addCaselessKey(v, k)
	}
}

func getFromObjectChildren(v *V, caseless bool, key string) (child *V, exist bool) {
	childProperty, exist := v.children.object[key]
	if exist {
		return childProperty.v, true
	}

	if !caseless {
		return &V{}, false
	}

	initCaselessStorage(v)

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

func getInCurrentValue(v *V, caseless bool, param any) (*V, error) {
	switch v.valueType {
	case Array:
		return getInCurrentArray(v, param)
	case Object:
		return getInCurrentObject(v, caseless, param)
	default:
		return &V{}, fmt.Errorf("%v type does not supports Get()", v.valueType)
	}
}

func getInCurrentArray(v *V, param any) (*V, error) {
	// integer expected
	pos, err := anyToInt(param)
	if err != nil {
		return &V{}, err
	}
	child, ok := childAtIndex(v, pos)
	if !ok {
		return &V{}, ErrOutOfRange
	}
	return child, nil
}

func getInCurrentObject(v *V, caseless bool, param any) (*V, error) {
	// string expected
	key, err := anyToString(param)
	if err != nil {
		return &V{}, err
	}
	child, exist := getFromObjectChildren(v, caseless, key)
	if !exist {
		return &V{}, ErrNotFound
	}
	return child, nil
}

// GetBytes is similar with v, err := Get(...); v.Bytes(). But if error occurs or Base64 decode error, returns error.
//
// GetBytes 类似于 v, err := Get(...); v.Bytes()，但如果查询中发生错误，或者 base64 解码错误，则返回错误。
func (v *V) GetBytes(firstParam any, otherParams ...any) ([]byte, error) {
	return getBytes(v, false, firstParam, otherParams...)
}

func getBytes(v *V, caseless bool, firstParam any, otherParams ...any) ([]byte, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
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

// GetString is equivalent to v, err := Get(...); v.String(). If error occurs, returns "".
//
// GetString 等效于 v, err := Get(...); v.String()。如果发生错误，则返回 ""。
func (v *V) GetString(firstParam any, otherParams ...any) (string, error) {
	return getString(v, false, firstParam, otherParams...)
}

func getString(v *V, caseless bool, firstParam any, otherParams ...any) (string, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return "", err
	}
	if ret.valueType != String {
		return "", ErrTypeNotMatch
	}
	return ret.String(), nil
}

// GetInt is equivalent to v, err := Get(...); v.Int(). If error occurs, returns 0.
//
// GetInt 等效于 v, err := Get(...); v.Int()。如果发生错误，则返回 0。
func (v *V) GetInt(firstParam any, otherParams ...any) (int, error) {
	return getInt(v, false, firstParam, otherParams...)
}

func getInt(v *V, caseless bool, firstParam any, otherParams ...any) (int, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Int(), err
}

// GetUint is equivalent to v, err := Get(...); v.Uint(). If error occurs, returns 0.
//
// GetUint 等效于 v, err := Get(...); v.Uint()。如果发生错误，则返回 0。
func (v *V) GetUint(firstParam any, otherParams ...any) (uint, error) {
	return getUint(v, false, firstParam, otherParams...)
}

func getUint(v *V, caseless bool, firstParam any, otherParams ...any) (uint, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Uint(), err
}

// GetInt64 is equivalent to v, err := Get(...); v.Int64(). If error occurs, returns 0.
//
// GetInt64 等效于 v, err := Get(...); v.Int64()。如果发生错误，则返回 0。
func (v *V) GetInt64(firstParam any, otherParams ...any) (int64, error) {
	return getInt64(v, false, firstParam, otherParams...)
}

func getInt64(v *V, caseless bool, firstParam any, otherParams ...any) (int64, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Int64(), err
}

// GetUint64 is equivalent to v, err := Get(...); v.Unt64(). If error occurs, returns 0.
//
// GetUint64 等效于 v, err := Get(...); v.Unt64()。如果发生错误，则返回 0。
func (v *V) GetUint64(firstParam any, otherParams ...any) (uint64, error) {
	return getUint64(v, false, firstParam, otherParams...)
}

func getUint64(v *V, caseless bool, firstParam any, otherParams ...any) (uint64, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Uint64(), err
}

// GetInt32 is equivalent to v, err := Get(...); v.Int32(). If error occurs, returns 0.
//
// GetInt32 等效于 v, err := Get(...); v.Int32()。如果发生错误，则返回 0。
func (v *V) GetInt32(firstParam any, otherParams ...any) (int32, error) {
	return getInt32(v, false, firstParam, otherParams...)
}

func getInt32(v *V, caseless bool, firstParam any, otherParams ...any) (int32, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Int32(), err
}

// GetUint32 is equivalent to v, err := Get(...); v.Uint32(). If error occurs, returns 0.
//
// GetUint32 等效于 v, err := Get(...); v.Uint32()。如果发生错误，则返回 0。
func (v *V) GetUint32(firstParam any, otherParams ...any) (uint32, error) {
	return getUint32(v, false, firstParam, otherParams...)
}

func getUint32(v *V, caseless bool, firstParam any, otherParams ...any) (uint32, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Uint32(), err
}

// GetFloat64 is equivalent to v, err := Get(...); v.Float64(). If error occurs, returns 0.0.
//
// GetFloat64 等效于 v, err := Get(...); v.Float64()。如果发生错误，则返回 0.0。
func (v *V) GetFloat64(firstParam any, otherParams ...any) (float64, error) {
	return getFloat64(v, false, firstParam, otherParams...)
}

func getFloat64(v *V, caseless bool, firstParam any, otherParams ...any) (float64, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Float64(), err
}

// GetFloat32 is equivalent to v, err := Get(...); v.Float32(). If error occurs, returns 0.0.
//
// GetFloat32 等效于 v, err := Get(...); v.Float32()。如果发生错误，则返回 0.0。
func (v *V) GetFloat32(firstParam any, otherParams ...any) (float32, error) {
	return getFloat32(v, false, firstParam, otherParams...)
}

func getFloat32(v *V, caseless bool, firstParam any, otherParams ...any) (float32, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	ret, err = getNumberAndErrorFromValue(ret)
	return ret.Float32(), err
}

// GetBool is equivalent to v, err := Get(...); v.Bool(). If error occurs, returns false.
//
// GetBool 等效于 v, err := Get(...); v.Bool()。如果发生错误，则返回 false。
func (v *V) GetBool(firstParam any, otherParams ...any) (bool, error) {
	return getBool(v, false, firstParam, otherParams...)
}

func getBool(v *V, caseless bool, firstParam any, otherParams ...any) (bool, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return false, err
	}
	ret, err = getBoolAndErrorFromValue(ret)
	return ret.Bool(), err
}

// GetNull is equivalent to v, err := Get(...); raise err if error occurs or v.IsNull() == false.
//
// GetNull 等效于 v, err := Get(...);，如果发生错误或者 v.IsNull() == false 则返回错误。
func (v *V) GetNull(firstParam any, otherParams ...any) error {
	return getNull(v, false, firstParam, otherParams...)
}

func getNull(v *V, caseless bool, firstParam any, otherParams ...any) error {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return err
	}
	if ret.valueType != Null {
		return ErrTypeNotMatch
	}
	return nil
}

// GetObject is equivalent to v, err := Get(...); raise err if error occurs or v.IsObject() == false.
//
// GetObject 等效于 v, err := Get(...);，如果发生错误或者 v.IsObject() == false 则返回错误。
func (v *V) GetObject(firstParam any, otherParams ...any) (*V, error) {
	return getObject(v, false, firstParam, otherParams...)
}

func getObject(v *V, caseless bool, firstParam any, otherParams ...any) (*V, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
	if err != nil {
		return &V{}, err
	}
	if ret.valueType != Object {
		return &V{}, ErrTypeNotMatch
	}
	return ret, nil
}

// GetArray is equivalent to v, err := Get(...); raise err if or v.IsArray() == false.
//
// GetArray 等效于 v, err := Get(...);，如果发生错误或者 v.IsArray() == false 则返回错误。
func (v *V) GetArray(firstParam any, otherParams ...any) (*V, error) {
	return getArray(v, false, firstParam, otherParams...)
}

func getArray(v *V, caseless bool, firstParam any, otherParams ...any) (*V, error) {
	ret, err := get(v, caseless, firstParam, otherParams...)
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
// Caseless 类型通过 Caseless() 函数返回。通过 Caseless 接口操作的所有操作均与 (*v).Get()
// 相同，但是对 key 进行读取的时候，不区分大小写。
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
// IMPORTANT: This function is not gouroutine-safe. Write-mutex (instead of read-mutex)
// should be attached in cross-goroutine scenarios.
//
// Caseless 返回 Caseless 接口，从而实现不区分大小写的 Get 操作。
//
// 注意: 该函数不是协程安全的，如果在多协程场景下，调用该函数，需要加上写锁，而不能用读锁。
func (v *V) Caseless() Caseless {
	switch v.valueType {
	default:
		return v

	case Array, Object:
		return &caselessOp{
			v: v,
		}
	}
}

type caselessOp struct {
	v *V
}

func (g *caselessOp) Get(firstParam any, otherParams ...any) (*V, error) {
	return get(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) MustGet(firstParam any, otherParams ...any) *V {
	res, _ := get(g.v, true, firstParam, otherParams...)
	return res
}

func (g *caselessOp) GetBytes(firstParam any, otherParams ...any) ([]byte, error) {
	return getBytes(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetString(firstParam any, otherParams ...any) (string, error) {
	return getString(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetInt(firstParam any, otherParams ...any) (int, error) {
	return getInt(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetUint(firstParam any, otherParams ...any) (uint, error) {
	return getUint(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetInt64(firstParam any, otherParams ...any) (int64, error) {
	return getInt64(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetUint64(firstParam any, otherParams ...any) (uint64, error) {
	return getUint64(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetInt32(firstParam any, otherParams ...any) (int32, error) {
	return getInt32(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetUint32(firstParam any, otherParams ...any) (uint32, error) {
	return getUint32(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetFloat64(firstParam any, otherParams ...any) (float64, error) {
	return getFloat64(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetFloat32(firstParam any, otherParams ...any) (float32, error) {
	return getFloat32(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetBool(firstParam any, otherParams ...any) (bool, error) {
	return getBool(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetNull(firstParam any, otherParams ...any) error {
	return getNull(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetObject(firstParam any, otherParams ...any) (*V, error) {
	return getObject(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) GetArray(firstParam any, otherParams ...any) (*V, error) {
	return getArray(g.v, true, firstParam, otherParams...)
}

func (g *caselessOp) Delete(firstParam any, otherParams ...any) error {
	return g.v.delete(true, firstParam, otherParams...)
}

func (g *caselessOp) MustDelete(firstParam any, otherParams ...any) {
	_ = g.v.delete(true, firstParam, otherParams...)
}

// ==== internal value access functions ====

func getNumberFromNotNumberValue(v *V) *V {
	if !v.IsString() {
		return NewInt(0)
	}
	if v.valueStr == "" {
		return NewInt64(0)
	}
	ret, _ := newFromNumber(globalPool{}, bytes.TrimSpace([]byte(v.valueStr)))
	err := parseNumber(ret, globalPool{})
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
		err := parseNumber(ret, globalPool{})
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
