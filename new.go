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
	v.stringValue = s
	v.parsed = true
	return v
}

// NewInt64 returns an initialied num jsonvalue object by int64 type
func NewInt64(i int64) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.floated = false
	v.negative = i < 0
	v.floatValue = float64(i)
	v.int64Value = i
	v.uint64Value = uint64(i)
	s := strconv.FormatInt(v.int64Value, 10)
	v.valueBytes = []byte(s)
	v.parsed = true
	return v
}

// NewUint64 returns an initialied num jsonvalue object by uint64 type
func NewUint64(u uint64) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.floated = false
	v.negative = false
	v.floatValue = float64(u)
	v.int64Value = int64(u)
	v.uint64Value = u
	s := strconv.FormatUint(v.uint64Value, 10)
	v.valueBytes = []byte(s)
	v.parsed = true
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
	v.boolValue = b
	v.parsed = true
	return v
}

// NewNull returns an initialied null jsonvalue object
func NewNull() *V {
	v := new()
	v.valueType = jsonparser.Null
	v.parsed = true
	return v
}

// NewObject returns an empty object-typed jsonvalue object
func NewObject() *V {
	v := newObject()
	v.parsed = true
	return v
}

// NewArray returns an emty array-typed jsonvalue object
func NewArray() *V {
	v := newArray()
	v.parsed = true
	return v
}

// NewFloat64 returns an initialied num jsonvalue object by float64 type. The precision prec controls the number of digits. Use -1 in prec for automatically precision.
func NewFloat64(f float64, prec int) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.negative = f < 0
	v.floatValue = f
	v.int64Value = int64(f)
	v.uint64Value = uint64(f)
	if prec >= 0 {
		s := strconv.FormatFloat(f, 'f', prec, 64)
		v.valueBytes = []byte(s)
	} else {
		b, _ := json.Marshal(&f)
		v.valueBytes = b
	}
	v.parsed = true
	return v
}

// NewFloat32 returns an initialied num jsonvalue object by float32 type. The precision prec controls the number of digits. Use -1 in prec for automatically precision.
func NewFloat32(f float32, prec int) *V {
	v := new()
	v.valueType = jsonparser.Number
	v.negative = f < 0
	v.floatValue = float64(f)
	v.int64Value = int64(f)
	v.uint64Value = uint64(f)
	if prec >= 0 {
		s := strconv.FormatFloat(v.floatValue, 'f', prec, 32)
		v.valueBytes = []byte(s)
	} else {
		b, _ := json.Marshal(&f)
		v.valueBytes = b
	}
	v.parsed = true
	return v
}
