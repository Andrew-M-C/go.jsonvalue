package jsonvalue

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
