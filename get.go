package jsonvalue

// Len returns length of an object or array type JSON value.
//
// Len 返回当前对象类型或数组类型的 JSON 的成员长度。如果不是这两种类型，那么会返回 0。
func (v *V) Len() int {
	if v.impl == nil {
		return 0
	}
	return v.impl.Len()
}

// Get returns JSON value in specified position. Param formats are like At().
//
// Get 返回按照参数指定的位置的 JSON 成员值。参数格式与 At() 函数相同
func (v *V) Get(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	if v.impl == nil {
		return &V{}, ErrValueUninitialized
	}
	return v.impl.get(false, firstParam, otherParams...)
}

func (v *V) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	if v.impl == nil {
		return &V{}, ErrValueUninitialized
	}
	return v.impl.get(caseless, firstParam, otherParams...)
}

// MustGet is same as Get(), but does not return error. If error occurs, a JSON value with
// NotExist type will be returned.
//
// MustGet 与 Get() 函数相同，不过不返回错误。如果发生错误了，那么会返回一个 ValueType() 返回值为 NotExist
// 的 JSON 值对象。
func (v *V) MustGet(firstParam interface{}, otherParams ...interface{}) *V {
	res, _ := v.Get(firstParam, otherParams...)
	return res
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
	if ret.ValueType() != String {
		return []byte{}, ErrTypeNotMatch
	}
	b, err := b64.DecodeString(ret.String())
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
	if ret.ValueType() != String {
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
	i, err := v.getInt64(caseless, firstParam, otherParams...)
	return int(i), err
}

// GetUint is equalivent to v, err := Get(...); v.Uint(). If error occurs, returns 0.
//
// GetUint 等效于 v, err := Get(...); v.Uint()。如果发生错误，则返回 0。
func (v *V) GetUint(firstParam interface{}, otherParams ...interface{}) (uint, error) {
	return v.getUint(false, firstParam, otherParams...)
}

func (v *V) getUint(caseless bool, firstParam interface{}, otherParams ...interface{}) (uint, error) {
	u, err := v.getUint64(caseless, firstParam, otherParams...)
	return uint(u), err
}

// GetInt64 is equalivent to v, err := Get(...); v.Int64(). If error occurs, returns 0.
//
// GetInt64 等效于 v, err := Get(...); v.Int64()。如果发生错误，则返回 0。
func (v *V) GetInt64(firstParam interface{}, otherParams ...interface{}) (int64, error) {
	return v.getInt64(false, firstParam, otherParams...)
}

func (v *V) getInt64(caseless bool, firstParam interface{}, otherParams ...interface{}) (int64, error) {
	if v.impl == nil {
		return 0, ErrValueUninitialized
	}

	v, err := v.impl.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	return v.impl.Int64()
}

// GetUint64 is equalivent to v, err := Get(...); v.Unt64(). If error occurs, returns 0.
//
// GetUint64 等效于 v, err := Get(...); v.Unt64()。如果发生错误，则返回 0。
func (v *V) GetUint64(firstParam interface{}, otherParams ...interface{}) (uint64, error) {
	return v.getUint64(false, firstParam, otherParams...)
}

func (v *V) getUint64(caseless bool, firstParam interface{}, otherParams ...interface{}) (uint64, error) {
	if v.impl == nil {
		return 0, ErrValueUninitialized
	}

	v, err := v.impl.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	return v.impl.Uint64()
}

// GetInt32 is equalivent to v, err := Get(...); v.Int32(). If error occurs, returns 0.
//
// GetInt32 等效于 v, err := Get(...); v.Int32()。如果发生错误，则返回 0。
func (v *V) GetInt32(firstParam interface{}, otherParams ...interface{}) (int32, error) {
	return v.getInt32(false, firstParam, otherParams...)
}

func (v *V) getInt32(caseless bool, firstParam interface{}, otherParams ...interface{}) (int32, error) {
	i, err := v.getInt64(caseless, firstParam, otherParams...)
	return int32(i), err
}

// GetUint32 is equalivent to v, err := Get(...); v.Uint32(). If error occurs, returns 0.
//
// GetUint32 等效于 v, err := Get(...); v.Uint32()。如果发生错误，则返回 0。
func (v *V) GetUint32(firstParam interface{}, otherParams ...interface{}) (uint32, error) {
	return v.getUint32(false, firstParam, otherParams...)
}

func (v *V) getUint32(caseless bool, firstParam interface{}, otherParams ...interface{}) (uint32, error) {
	u, err := v.getUint64(caseless, firstParam, otherParams...)
	return uint32(u), err
}

// GetFloat64 is equalivent to v, err := Get(...); v.Float64(). If error occurs, returns 0.0.
//
// GetFloat64 等效于 v, err := Get(...); v.Float64()。如果发生错误，则返回 0.0。
func (v *V) GetFloat64(firstParam interface{}, otherParams ...interface{}) (float64, error) {
	return v.getFloat64(false, firstParam, otherParams...)
}

func (v *V) getFloat64(caseless bool, firstParam interface{}, otherParams ...interface{}) (float64, error) {
	if v.impl == nil {
		return 0, ErrValueUninitialized
	}

	v, err := v.impl.get(caseless, firstParam, otherParams...)
	if err != nil {
		return 0, err
	}
	return v.impl.Float64()
}

// GetFloat32 is equalivent to v, err := Get(...); v.Float32(). If error occurs, returns 0.0.
//
// GetFloat32 等效于 v, err := Get(...); v.Float32()。如果发生错误，则返回 0.0。
func (v *V) GetFloat32(firstParam interface{}, otherParams ...interface{}) (float32, error) {
	return v.getFloat32(false, firstParam, otherParams...)
}

func (v *V) getFloat32(caseless bool, firstParam interface{}, otherParams ...interface{}) (float32, error) {
	f, err := v.getFloat64(caseless, firstParam, otherParams...)
	return float32(f), err
}

// GetBool is equalivent to v, err := Get(...); v.Bool(). If error occurs, returns false.
//
// GetBool 等效于 v, err := Get(...); v.Bool()。如果发生错误，则返回 false。
func (v *V) GetBool(firstParam interface{}, otherParams ...interface{}) (bool, error) {
	return v.getBool(false, firstParam, otherParams...)
}

func (v *V) getBool(caseless bool, firstParam interface{}, otherParams ...interface{}) (bool, error) {
	if v.impl == nil {
		return false, ErrValueUninitialized
	}

	v, err := v.impl.get(caseless, firstParam, otherParams...)
	if err != nil {
		return false, err
	}
	return v.impl.Bool()
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
	if ret.ValueType() != Null {
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
	if ret.ValueType() != Object {
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
	if ret.ValueType() != Array {
		return &V{}, ErrTypeNotMatch
	}
	return ret, nil
}
