package jsonvalue

// Caseless is returned by Caseless(). operations of Caseless type are same as (*V).Get(), but are via caseless key.
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

// Caseless mark current value to be caseless mode
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
