package jsonvalue

import (
	"encoding/json"
	"strconv"

	"github.com/buger/jsonparser"
)

// NewString returns an initialied string jsonvalue object
func NewString(s string) *V {
	v := new()
	v.valueType = jsonparser.String
	v.value.str = s
	v.status.parsed = true
	return v
}

// NewInt64 returns an initialied num jsonvalue object by int64 type
func NewInt64(i int64) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.status.floated = false
	v.status.negative = i < 0
	v.value.f64 = float64(i)
	v.value.i64 = i
	v.value.u64 = uint64(i)
	s := strconv.FormatInt(v.value.i64, 10)
	v.valueBytes = []byte(s)
	v.status.parsed = true
	return v
}

// NewUint64 returns an initialied num jsonvalue object by uint64 type
func NewUint64(u uint64) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.status.floated = false
	v.status.negative = false
	v.value.f64 = float64(u)
	v.value.i64 = int64(u)
	v.value.u64 = u
	s := strconv.FormatUint(v.value.u64, 10)
	v.valueBytes = []byte(s)
	v.status.parsed = true
	return v
}

// NewInt returns an initialied num jsonvalue object by int type
func NewInt(i int) *V {
	return NewInt64(int64(i))
}

// NewUint returns an initialied num jsonvalue object by uint type
func NewUint(u uint) *V {
	return NewUint64(uint64(u))
}

// NewInt32 returns an initialied num jsonvalue object by int32 type
func NewInt32(i int32) *V {
	return NewInt64(int64(i))
}

// NewUint32 returns an initialied num jsonvalue object by uint32 type
func NewUint32(u uint32) *V {
	return NewUint64(uint64(u))
}

// NewBool returns an initialied boolean jsonvalue object
func NewBool(b bool) *V {
	v := new()
	v.valueType = jsonparser.Boolean
	v.value.boolean = b
	v.status.parsed = true
	return v
}

// NewNull returns an initialied null jsonvalue object
func NewNull() *V {
	v := new()
	v.valueType = jsonparser.Null
	v.status.parsed = true
	return v
}

// NewObject returns an object-typed jsonvalue object. If keyValues is specified, it will also create some key-values in
// the object. Now we supports basic types only. Such as int/uint series, string, bool, nil.
func NewObject(keyValues ...map[string]interface{}) *V {
	v := newObject()
	v.status.parsed = true

	if len(keyValues) > 0 {
		kv := keyValues[0]
		if kv != nil {
			v.parseNewObjectKV(kv)
		}
	}

	return v
}

func (v *V) parseNewObjectKV(kv map[string]interface{}) {
	for k, val := range kv {
		switch val.(type) {
		case nil:
			v.SetNull().At(k)
		case string:
			v.SetString(val.(string)).At(k)
		case bool:
			v.SetBool(val.(bool)).At(k)
		case int:
			v.SetInt(val.(int)).At(k)
		case uint:
			v.SetUint(val.(uint)).At(k)
		case int8:
			v.SetInt32(int32(val.(int8))).At(k)
		case uint8:
			v.SetUint32(uint32(val.(uint8))).At(k)
		case int16:
			v.SetInt32(int32(val.(int16))).At(k)
		case uint16:
			v.SetUint32(uint32(val.(uint16))).At(k)
		case int32:
			v.SetInt32(val.(int32)).At(k)
		case uint32:
			v.SetUint32(val.(uint32)).At(k)
		case int64:
			v.SetInt64(val.(int64)).At(k)
		case uint64:
			v.SetUint64(val.(uint64)).At(k)
		case float32:
			v.SetFloat32(val.(float32), -1).At(k)
		case float64:
			v.SetFloat64(val.(float64), -1).At(k)
		default:
			// continue
		}
	}
	return
}

// NewArray returns an emty array-typed jsonvalue object
func NewArray() *V {
	v := newArray()
	v.status.parsed = true
	return v
}

// NewFloat64 returns an initialied num jsonvalue object by float64 type. The precision prec controls the number of
// digits. Use -1 in prec for automatically precision.
func NewFloat64(f float64, prec int) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.status.negative = f < 0
	v.value.f64 = f
	v.value.i64 = int64(f)
	v.value.u64 = uint64(f)
	if prec >= 0 {
		s := strconv.FormatFloat(f, 'f', prec, 64)
		v.valueBytes = []byte(s)
	} else {
		b, _ := json.Marshal(&f)
		v.valueBytes = b
	}
	v.status.parsed = true
	return v
}

// NewFloat32 returns an initialied num jsonvalue object by float32 type. The precision prec controls the number of
// digits. Use -1 in prec for automatically precision.
func NewFloat32(f float32, prec int) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.status.negative = f < 0
	v.value.f64 = float64(f)
	v.value.i64 = int64(f)
	v.value.u64 = uint64(f)
	if prec >= 0 {
		s := strconv.FormatFloat(v.value.f64, 'f', prec, 32)
		v.valueBytes = []byte(s)
	} else {
		b, _ := json.Marshal(&f)
		v.valueBytes = b
	}
	v.status.parsed = true
	return v
}
