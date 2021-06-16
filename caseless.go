package jsonvalue

// Caseless is returned by Caseless(). operations of Caseless type are same as (*V).Get(), but are via caseless key.
//
// Caseless 类型通过 Caseless() 函数返回。通过 Caseless 接口操作的所有操作均与 (*v).Get() 相同，但是对 key 进行读取的时候，
// 不区分大小写。
type Caseless interface {
	Get(firstParam interface{}, otherParams ...interface{}) (*V, error)
	GetBytes(firstParam interface{}, otherParams ...interface{}) ([]byte, error)
	GetString(firstParam interface{}, otherParams ...interface{}) (string, error)
	GetInt(firstParam interface{}, otherParams ...interface{}) (int, error)
	GetUint(firstParam interface{}, otherParams ...interface{}) (uint, error)
	GetInt64(firstParam interface{}, otherParams ...interface{}) (int64, error)
	GetUint64(firstParam interface{}, otherParams ...interface{}) (uint64, error)
	GetInt32(firstParam interface{}, otherParams ...interface{}) (int32, error)
	GetUint32(firstParam interface{}, otherParams ...interface{}) (uint32, error)
	GetFloat64(firstParam interface{}, otherParams ...interface{}) (float64, error)
	GetFloat32(firstParam interface{}, otherParams ...interface{}) (float32, error)
	GetBool(firstParam interface{}, otherParams ...interface{}) (bool, error)
	GetNull(firstParam interface{}, otherParams ...interface{}) error
	GetObject(firstParam interface{}, otherParams ...interface{}) (*V, error)
	GetArray(firstParam interface{}, otherParams ...interface{}) (*V, error)

	Delete(firstParam interface{}, otherParams ...interface{}) error
}

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

func (g *caselessOper) Get(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return g.v.get(true, firstParam, otherParams...)
}

func (g *caselessOper) GetBytes(firstParam interface{}, otherParams ...interface{}) ([]byte, error) {
	return g.v.getBytes(true, firstParam, otherParams...)
}

func (g *caselessOper) GetString(firstParam interface{}, otherParams ...interface{}) (string, error) {
	return g.v.getString(true, firstParam, otherParams...)
}

func (g *caselessOper) GetInt(firstParam interface{}, otherParams ...interface{}) (int, error) {
	return g.v.getInt(true, firstParam, otherParams...)
}

func (g *caselessOper) GetUint(firstParam interface{}, otherParams ...interface{}) (uint, error) {
	return g.v.getUint(true, firstParam, otherParams...)
}

func (g *caselessOper) GetInt64(firstParam interface{}, otherParams ...interface{}) (int64, error) {
	return g.v.getInt64(true, firstParam, otherParams...)
}

func (g *caselessOper) GetUint64(firstParam interface{}, otherParams ...interface{}) (uint64, error) {
	return g.v.getUint64(true, firstParam, otherParams...)
}

func (g *caselessOper) GetInt32(firstParam interface{}, otherParams ...interface{}) (int32, error) {
	return g.v.getInt32(true, firstParam, otherParams...)
}

func (g *caselessOper) GetUint32(firstParam interface{}, otherParams ...interface{}) (uint32, error) {
	return g.v.getUint32(true, firstParam, otherParams...)
}

func (g *caselessOper) GetFloat64(firstParam interface{}, otherParams ...interface{}) (float64, error) {
	return g.v.getFloat64(true, firstParam, otherParams...)
}

func (g *caselessOper) GetFloat32(firstParam interface{}, otherParams ...interface{}) (float32, error) {
	return g.v.getFloat32(true, firstParam, otherParams...)
}

func (g *caselessOper) GetBool(firstParam interface{}, otherParams ...interface{}) (bool, error) {
	return g.v.getBool(true, firstParam, otherParams...)
}

func (g *caselessOper) GetNull(firstParam interface{}, otherParams ...interface{}) error {
	return g.v.getNull(true, firstParam, otherParams...)
}

func (g *caselessOper) GetObject(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return g.v.getObject(true, firstParam, otherParams...)
}

func (g *caselessOper) GetArray(firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return g.v.getArray(true, firstParam, otherParams...)
}

func (g *caselessOper) Delete(firstParam interface{}, otherParams ...interface{}) error {
	return g.v.delete(true, firstParam, otherParams...)
}
