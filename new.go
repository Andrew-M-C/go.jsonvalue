package jsonvalue

import (
	"encoding/json"
	"strconv"

	"github.com/buger/jsonparser"
)

// NewString returns an initialied string jsonvalue object
//
// NewString 用给定的 string 返回一个初始化好的字符串类型的 jsonvalue 值
func NewString(s string) *V {
	v := new(jsonparser.String)
	v.valueStr = s
	v.parsed = true
	return v
}

// NewInt64 returns an initialied num jsonvalue object by int64 type
//
// NewInt64 用给定的 int64 返回一个初始化好的数字类型的 jsonvalue 值
func NewInt64(i int64) *V {
	v := new(jsonparser.Number)
	// v.num = &num{}
	v.num.floated = false
	v.num.negative = i < 0
	v.num.f64 = float64(i)
	v.num.i64 = i
	v.num.u64 = uint64(i)
	s := strconv.FormatInt(v.num.i64, 10)
	v.srcByte = []byte(s)
	v.srcOffset, v.srcEnd = 0, len(s)
	v.parsed = true
	return v
}

// NewUint64 returns an initialied num jsonvalue object by uint64 type
//
// NewUint64 用给定的 uint64 返回一个初始化好的数字类型的 jsonvalue 值
func NewUint64(u uint64) *V {
	v := new(jsonparser.Number)
	// v.num = &num{}
	v.num.floated = false
	v.num.negative = false
	v.num.f64 = float64(u)
	v.num.i64 = int64(u)
	v.num.u64 = u
	s := strconv.FormatUint(v.num.u64, 10)
	v.srcByte = []byte(s)
	v.srcOffset, v.srcEnd = 0, len(s)
	v.parsed = true
	return v
}

// NewInt returns an initialied num jsonvalue object by int type
//
// NewInt 用给定的 int 返回一个初始化好的数字类型的 jsonvalue 值
func NewInt(i int) *V {
	return NewInt64(int64(i))
}

// NewUint returns an initialied num jsonvalue object by uint type
//
// NewUint 用给定的 uint 返回一个初始化好的数字类型的 jsonvalue 值
func NewUint(u uint) *V {
	return NewUint64(uint64(u))
}

// NewInt32 returns an initialied num jsonvalue object by int32 type
//
// NewInt32 用给定的 int32 返回一个初始化好的数字类型的 jsonvalue 值
func NewInt32(i int32) *V {
	return NewInt64(int64(i))
}

// NewUint32 returns an initialied num jsonvalue object by uint32 type
//
// NewUint32 用给定的 uint32 返回一个初始化好的数字类型的 jsonvalue 值
func NewUint32(u uint32) *V {
	return NewUint64(uint64(u))
}

// NewBool returns an initialied boolean jsonvalue object
//
// NewBool 用给定的 bool 返回一个初始化好的布尔类型的 jsonvalue 值
func NewBool(b bool) *V {
	v := new(jsonparser.Boolean)
	v.valueBool = b
	v.parsed = true
	return v
}

// NewNull returns an initialied null jsonvalue object
//
// NewNull 返回一个初始化好的 null 类型的 jsonvalue 值
func NewNull() *V {
	v := new(jsonparser.Null)
	v.parsed = true
	return v
}

// NewObject returns an object-typed jsonvalue object. If keyValues is specified, it will also create some key-values in
// the object. Now we supports basic types only. Such as int/uint, int/int8/int16/int32/int64,
// uint/uint8/uint16/uint32/uint64 series, string, bool, nil.
//
// NewObject 返回一个初始化好的 object 类型的 jsonvalue 值。可以使用可选的 map[string]interface{} 类型参数初始化该 object 的下一级键值对，
// 不过目前只支持基础类型，也就是: int/uint, int/int8/int16/int32/int64, uint/uint8/uint16/uint32/uint64, string, bool, nil。
func NewObject(keyValues ...map[string]interface{}) *V {
	v := newObject()
	v.parsed = true

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
		switch val := val.(type) {
		case nil:
			v.SetNull().At(k)
		case string:
			v.SetString(val).At(k)
		case bool:
			v.SetBool(val).At(k)
		case int:
			v.SetInt(val).At(k)
		case uint:
			v.SetUint(val).At(k)
		case int8:
			v.SetInt32(int32(val)).At(k)
		case uint8:
			v.SetUint32(uint32(val)).At(k)
		case int16:
			v.SetInt32(int32(val)).At(k)
		case uint16:
			v.SetUint32(uint32(val)).At(k)
		case int32:
			v.SetInt32(val).At(k)
		case uint32:
			v.SetUint32(val).At(k)
		case int64:
			v.SetInt64(val).At(k)
		case uint64:
			v.SetUint64(val).At(k)
		case float32:
			v.SetFloat32(val, -1).At(k)
		case float64:
			v.SetFloat64(val, -1).At(k)
		default:
			// continue
		}
	}
}

// NewArray returns an emty array-typed jsonvalue object
//
// NewArray 返回一个初始化好的 array 类型的 jsonvalue 值。
func NewArray() *V {
	v := newArray()
	v.parsed = true
	return v
}

// NewFloat64 returns an initialied num jsonvalue object by float64 type. The precision prec controls the number of
// digits. Use -1 in prec for automatically precision.
//
// NewFloat64 根据指定的 flout64 类型返回一个初始化好的数字类型的 jsonvalue 值。
// 参数 precision prec 指定需要编码的小数点后的位数。使用 -1 则交给编译器自行判断。
func NewFloat64(f float64, prec int) *V {
	v := new(jsonparser.Number)
	// v.num = &num{}
	v.num.negative = f < 0
	v.num.f64 = f
	v.num.i64 = int64(f)
	v.num.u64 = uint64(f)
	if prec >= 0 {
		s := strconv.FormatFloat(f, 'f', prec, 64)
		v.srcByte = []byte(s)
		v.srcOffset, v.srcEnd = 0, len(s)
	} else {
		// s := fmt.Sprintf("%f", f)
		// b := []byte(s)
		b, _ := json.Marshal(&f)
		v.srcByte = b
		v.srcOffset, v.srcEnd = 0, len(b)
	}
	v.parsed = true
	return v
}

// NewFloat32 returns an initialied num jsonvalue object by float32 type. The precision prec controls the number of
// digits. Use -1 in prec for automatically precision.
//
// NewFloat32 根据指定的 float32 类型返回一个初始化好的数字类型的 jsonvalue 值。
// 参数 precision prec 指定需要编码的小数点后的位数。使用 -1 则交给编译器自行判断。
func NewFloat32(f float32, prec int) *V {
	v := new(jsonparser.Number)
	// v.num = &num{}
	v.num.negative = f < 0
	v.num.f64 = float64(f)
	v.num.i64 = int64(f)
	v.num.u64 = uint64(f)
	if prec >= 0 {
		s := strconv.FormatFloat(v.num.f64, 'f', prec, 64)
		v.srcByte = []byte(s)
		v.srcOffset, v.srcEnd = 0, len(s)
	} else {
		// s := fmt.Sprintf("%f", f)
		// b := []byte(s)
		b, _ := json.Marshal(&f)
		v.srcByte = b
		v.srcOffset, v.srcEnd = 0, len(b)
	}
	v.parsed = true
	return v
}
