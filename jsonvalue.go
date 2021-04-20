// Package jsonvalue is for JSON parsing and setting. It is used in situations those Go structures cannot achieve, or "map[string]interface{}" could not do properbally.
//
// As a quick start:
// 	v := jsonvalue.NewObject()
// 	v.SetString("Hello, JSON").At("someObject", "someObject", "someObject", "message")  // automatically create sub objects
// 	fmt.Println(v.MustMarshalString())                                                  // marshal to string type. Use MustMarshal if you want []byte instead.
// 	// Output:
// 	// {"someObject":{"someObject":{"someObject":{"message":"Hello, JSON!"}}}
//
// If you want to parse raw JSON data, use Unmarshal()
// 	raw := []byte(`{"message":"hello, world"}`)
// 	v, err := jsonvalue.Unmarshal(raw)
// 	s, _ := v.GetString("message")
// 	fmt.Println(s)
// 	// Output:
// 	// hello, world
//
// jsonvalue 包用于 JSON 的解析（反序列化）和编码（序列化）。通常情况下我们用 struct 来处理结构化的 JSON，但是有时候使用 struct 不方便或者是功能不足的时候，
// go 一般而言使用的是 "map[string]interface{}"，但是后者也有很多不方便的地方。本包即是用于替代这些不方便的情况的。
//
// 快速上手：
// 	v := jsonvalue.NewObject()
// 	v.SetString("Hello, JSON").At("someObject", "someObject", "someObject", "message")  // 自动创建子成员
// 	fmt.Println(v.MustMarshalString())                                                  // 序列化为 string 类型，如果你要 []byte 类型，则使用 MustMarshal 函数。
// 	// 输出:
// 	// {"someObject":{"someObject":{"someObject":{"message":"Hello, JSON!"}}}
//
// 如果要反序列化原始的 JSON 文本，则使用 Unmarshal():
// 	raw := []byte(`{"message":"hello, world"}`)
// 	v, err := jsonvalue.Unmarshal(raw)
// 	s, _ := v.GetString("message")
// 	fmt.Println(s)
// 	// 输出:
// 	// hello, world
package jsonvalue

import (
	"bytes"
	"container/list"
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// V is the main type of jsonvalue, representing a JSON value.
//
// V 是 jsonvalue 的主类型，表示一个 JSON 值。
type V struct {
	valueType  jsoniter.ValueType
	valueBytes []byte

	status struct {
		parsed   bool
		negative bool
		floated  bool
	}

	value struct {
		str     string
		i64     int64
		u64     uint64
		boolean bool
		f64     float64
	}

	children struct {
		object map[string]*V
		array  *list.List

		// As official json package supports caseless key accessing, I decide to do it as well
		lowerCaseKeys map[string]map[string]struct{}
	}
}

func new() *V {
	v := V{}
	v.valueType = jsoniter.InvalidValue
	return &v
}

func newObject() *V {
	v := V{}
	v.valueType = jsoniter.ObjectValue
	v.children.object = make(map[string]*V)
	v.children.lowerCaseKeys = make(map[string]map[string]struct{})
	return &v
}

func newArray() *V {
	v := V{}
	v.valueType = jsoniter.ArrayValue
	v.children.array = list.New()
	return &v
}

func (v *V) addCaselessKey(k string) {
	lowerK := strings.ToLower(k)
	keys, exist := v.children.lowerCaseKeys[lowerK]
	if !exist {
		keys = make(map[string]struct{})
		v.children.lowerCaseKeys[lowerK] = keys
	}
	keys[k] = struct{}{}
}

func (v *V) delCaselessKey(k string) {
	lowerK := strings.ToLower(k)
	keys, exist := v.children.lowerCaseKeys[lowerK]
	if !exist {
		return
	}

	delete(keys, k)

	if len(keys) == 0 {
		delete(v.children.lowerCaseKeys, lowerK)
	}
}

// UnmarshalString is equavilent to Unmarshal(string(b)), but much more efficient.
//
// UnmarshalString 等效于 Unmarshal(string(b))，但效率更高。
func UnmarshalString(s string) (*V, error) {
	// reference: https://stackoverflow.com/questions/41591097/slice-bounds-out-of-range-when-using-unsafe-pointer
	// sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// bh := reflect.SliceHeader{
	// 	Data: sh.Data,
	// 	Len:  sh.Len,
	// 	Cap:  sh.Len,
	// }
	// b := *(*[]byte)(unsafe.Pointer(&bh))
	return Unmarshal(unsafeStoB(s))
}

// Unmarshal parse raw bytes(encoded in UTF-8 or pure AscII) and returns a *V instance.
//
// Unmarshal 解析原始的字节类型数据（以 UTF-8 或纯 AscII 编码），并返回一个 *V 对象。
func Unmarshal(b []byte) (ret *V, err error) {
	if nil == b || len(b) == 0 {
		return nil, ErrNilParameter
	}

	for i, c := range b {
		switch c {
		case ' ', '\r', '\n', '\t', '\b':
			// continue
		case '{':
			// object start
			return newFromObject(jsoniter.Get(b[i:]))
		case '[':
			return newFromArray(jsoniter.Get(b[i:]))
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			ret, err = newFromNumberBytes(b[i:])
			if err != nil {
				return
			}
			err = ret.parseNumber()
			if err != nil {
				return nil, err
			}
			return ret, nil

		case '"':
			ret = new()
			ret.valueType = jsoniter.StringValue
			ret.value.str, _, err = parseString(b[i:])
			if err != nil {
				return nil, err
			}
			ret.status.parsed = true
			return ret, nil

		case 't':
			return newFromTrue(b[i:])
		case 'f':
			return newFromFalse(b[i:])
		case 'n':
			if bytes.HasPrefix(b[i:], []byte("null")) {
				return newFromNull()
			}
			return nil, ErrNotValidNulllValue
		default:
			return nil, ErrRawBytesUnrecignized
		}
	}

	return nil, ErrRawBytesUnrecignized
}

func (v *V) parseNumber() (err error) {
	b := v.valueBytes

	if bytes.Contains(b, []byte(".")) {
		v.status.floated = true
		v.value.f64, err = parseFloat(b)
		if err != nil {
			return
		}

		v.status.parsed = true
		v.status.negative = (v.value.f64 < 0)
		v.value.i64 = int64(v.value.f64)
		v.value.u64 = uint64(v.value.f64)

	} else if b[0] == '-' {
		v.status.negative = true
		v.value.i64, err = parseInt(b)
		if err != nil {
			return
		}

		v.status.parsed = true
		v.value.u64 = uint64(v.value.i64)
		v.value.f64 = float64(v.value.i64)

	} else {
		v.status.negative = false
		v.value.u64, err = parseUint(b)
		if err != nil {
			return
		}

		v.status.parsed = true
		v.value.i64 = int64(v.value.u64)
		v.value.f64 = float64(v.value.u64)
	}

	return nil
}

// ==== simple object parsing ====
func newFromNumber(any jsoniter.Any) (ret *V, err error) {
	v := new()
	v.valueType = jsoniter.NumberValue
	v.valueBytes = unsafeStoB(any.ToString())
	return v, nil
}

func newFromNumberBytes(b []byte) (ret *V, err error) {
	v := new()
	v.valueType = jsoniter.NumberValue
	v.valueBytes = b
	return v, nil
}

// func newFromString(b []byte) (ret *V, err error) {
// 	v := new()
// 	v.valueType = jsoniter.StringValue
// 	v.valueBytes = b
// 	return v, nil
// }

func newFromTrue(b []byte) (ret *V, err error) {
	if len(b) != 4 || string(b) != "true" {
		return nil, ErrNotValidBoolValue
	}
	v := new()
	v.status.parsed = true
	v.valueType = jsoniter.BoolValue
	v.valueBytes = []byte{'t', 'r', 'u', 'e'}
	v.value.boolean = true
	return v, nil
}

func newFromFalse(b []byte) (ret *V, err error) {
	if len(b) != 5 || string(b) != "false" {
		return nil, ErrNotValidBoolValue
	}
	v := new()
	v.status.parsed = true
	v.valueType = jsoniter.BoolValue
	v.valueBytes = []byte{'f', 'a', 'l', 's', 'e'}
	v.value.boolean = false
	return v, nil
}

func newFromBool(b []byte) (ret *V, err error) {
	v := new()
	v.valueType = jsoniter.BoolValue

	switch string(b) {
	case "true":
		v.status.parsed = true
		v.valueBytes = []byte{'t', 'r', 'u', 'e'}
		v.value.boolean = true
	case "false":
		v.status.parsed = true
		v.valueBytes = []byte{'f', 'a', 'l', 's', 'e'}
		v.value.boolean = false
	default:
		return nil, ErrNotValidBoolValue
	}

	return v, nil
}

func newFromNull() (ret *V, err error) {
	v := new()
	v.status.parsed = true
	v.valueType = jsoniter.NilValue
	return v, nil
}

// ====
func newFromArray(any jsoniter.Any) (ret *V, err error) {
	o := newArray()
	le := any.Size()

	for i := 0; i < le; i++ {
		ch := any.Get(i)
		var child *V
		switch ch.ValueType() {
		default:
			err = fmt.Errorf("unmarshal error: %w", ch.LastError())
			// err = fmt.Errorf("invalid value type: %v", ch.ValueType())
			return nil, err
		case jsoniter.ObjectValue:
			child, err = newFromObject(ch)
		case jsoniter.ArrayValue:
			child, err = newFromArray(ch)
		case jsoniter.NumberValue:
			child, err = newFromNumber(ch)
		case jsoniter.BoolValue:
			child, err = newFromBool(unsafeStoB(ch.ToString()))
		case jsoniter.NilValue:
			child, err = newFromNull()
		case jsoniter.StringValue:
			child = new()
			child.status.parsed = true
			child.valueType = jsoniter.StringValue
			child.value.str = ch.ToString()
		}

		if err != nil {
			return nil, err
		}
		o.children.array.PushBack(child)
	}

	// done
	return o, nil
}

// ==== object parsing ====
func newFromObject(any jsoniter.Any) (ret *V, err error) {
	o := newObject()
	keys := any.Keys()

	for _, k := range keys {
		ch := any.Get(k)
		var child *V
		switch ch.ValueType() {
		default:
			err = fmt.Errorf("unmarshal error: %w", ch.LastError())
			// err = fmt.Errorf("invalid value type: %v", ch.ValueType())
			return nil, err
		case jsoniter.ObjectValue:
			child, err = newFromObject(ch)
		case jsoniter.ArrayValue:
			child, err = newFromArray(ch)
		case jsoniter.NumberValue:
			child, err = newFromNumber(ch)
		case jsoniter.BoolValue:
			child, err = newFromBool(unsafeStoB(ch.ToString()))
		case jsoniter.NilValue:
			child, err = newFromNull()
		case jsoniter.StringValue:
			child = new()
			child.status.parsed = true
			child.valueType = jsoniter.StringValue
			child.value.str = ch.ToString()
		}

		if err != nil {
			return nil, err
		}
		o.setToObjectChildren(k, child)
	}

	// done
	return o, nil
}

// ==== type access ====

// IsObject tells whether value is an object
//
// IsObject 判断当前值是不是一个对象类型
func (v *V) IsObject() bool {
	return v.valueType == jsoniter.ObjectValue
}

// IsArray tells whether value is an array
//
// IsArray 判断当前值是不是一个数组类型
func (v *V) IsArray() bool {
	return v.valueType == jsoniter.ArrayValue
}

// IsString tells whether value is a string
//
// IsString 判断当前值是不是一个字符串类型
func (v *V) IsString() bool {
	return v.valueType == jsoniter.StringValue
}

// IsNumber tells whether value is a number
//
// IsNumber 判断当前值是不是一个数字类型
func (v *V) IsNumber() bool {
	return v.valueType == jsoniter.NumberValue
}

// IsFloat tells whether value is a float point number. If there is no decimal point in original text, it returns false
// while IsNumber returns true.
//
// IsFloat 判断当前值是不是一个浮点数类型。如果给定的数不包含小数点，那么即便是数字类型，该函数也会返回 false.
func (v *V) IsFloat() bool {
	if v.valueType != jsoniter.NumberValue {
		return false
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return v.status.floated
}

// IsInteger tells whether value is a fix point interger
//
// IsNumber 判断当前值是不是一个定点数整型
func (v *V) IsInteger() bool {
	if v.valueType != jsoniter.NumberValue {
		return false
	}
	if !v.status.parsed {
		err := v.parseNumber()
		if err != nil {
			return false
		}
	}
	return !(v.status.floated)
}

// IsNegative tells whether value is a negative number
//
// IsNegative 判断当前值是不是一个负数
func (v *V) IsNegative() bool {
	if v.valueType != jsoniter.NumberValue {
		return false
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return v.status.negative
}

// IsPositive tells whether value is a positive number
//
// IsPositive 判断当前值是不是一个正数
func (v *V) IsPositive() bool {
	if v.valueType != jsoniter.NumberValue {
		return false
	}
	if !v.status.parsed {
		err := v.parseNumber()
		if err != nil {
			return false
		}
	}
	return !(v.status.negative)
}

// GreaterThanInt64Max return true when ALL conditions below are met:
// 	1. It is a number value.
// 	2. It is a positive interger.
// 	3. Its value is greater than 0x7fffffffffffffff.
//
// GreaterThanInt64Max 判断当前值是否超出 int64 可表示的范围。当以下条件均成立时，返回 true，否则返回 false：
// 	1. 是一个数字类型值.
// 	2. 是一个正整型数字.
// 	3. 该正整数的值大于 0x7fffffffffffffff.
func (v *V) GreaterThanInt64Max() bool {
	if v.valueType != jsoniter.NumberValue {
		return false
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	if v.status.negative {
		return false
	}
	return v.value.u64 > 0x7fffffffffffffff
}

// IsBoolean tells whether value is a boolean
//
// IsBoolean 判断当前值是不是一个布尔类型
func (v *V) IsBoolean() bool {
	return v.valueType == jsoniter.BoolValue
}

// IsNull tells whether value is a null
//
// IsBoolean 判断当前值是不是一个空类型
func (v *V) IsNull() bool {
	return v.valueType == jsoniter.NilValue
}

// ==== value access ====

func getNumberFromNotNumberValue(v *V) *V {
	if !v.IsString() {
		return NewInt(0)
	}
	ret, _ := newFromNumberBytes([]byte(v.value.str))
	err := ret.parseNumber()
	if err != nil {
		return NewInt64(0)
	}
	return ret
}

// Bool returns represented bool value. If value is not boolean, returns false.
//
// Bool 返回布尔类型值。如果当前值不是布尔类型，则返回 false。
func (v *V) Bool() bool {
	return v.value.boolean
}

// Int returns represented int value. If value is not a number, returns zero.
//
// Int 返回 int 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int() int {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Int()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return int(v.value.i64)
}

// Uint returns represented uint value. If value is not a number, returns zero.
//
// Uint 返回 uint 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint() uint {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Uint()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return uint(v.value.u64)
}

// Int64 returns represented int64 value. If value is not a number, returns zero.
//
// Int64 返回 int64 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int64() int64 {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Int64()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return int64(v.value.i64)
}

// Uint64 returns represented uint64 value. If value is not a number, returns zero.
//
// Uint64 返回 uint64 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint64() uint64 {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Uint64()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return uint64(v.value.u64)
}

// Int32 returns represented int32 value. If value is not a number, returns zero.
//
// Int32 返回 int32 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int32() int32 {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Int32()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return int32(v.value.i64)
}

// Uint32 returns represented uint32 value. If value is not a number, returns zero.
//
// Uint32 返回 uint32 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint32() uint32 {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Uint32()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return uint32(v.value.u64)
}

// Float64 returns represented float64 value. If value is not a number, returns zero.
//
// Float64 返回 float64 类型值。如果当前值不是数字类型，则返回 0.0。
func (v *V) Float64() float64 {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Float64()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return v.value.f64
}

// Float32 returns represented float32 value. If value is not a number, returns zero.
//
// Float32 返回 float32 类型值。如果当前值不是数字类型，则返回 0.0。
func (v *V) Float32() float32 {
	if v.valueType != jsoniter.NumberValue {
		return getNumberFromNotNumberValue(v).Float32()
	}
	if !v.status.parsed {
		v.parseNumber()
	}
	return float32(v.value.f64)
}

// String returns represented string value or the description for the jsonvalue.V instance if it is not a string.
//
// String 返回 string 类型值。如果当前值不是字符串类型，则返回当前 *V 类型的描述说明。
func (v *V) String() string {
	if v == nil {
		return ""
	}
	switch v.valueType {
	default:
		return ""
	case jsoniter.NilValue:
		return "null"
	case jsoniter.NumberValue:
		return string(v.valueBytes)
	case jsoniter.StringValue:
		if !v.status.parsed {
			var e error
			v.value.str, v.valueBytes, e = parseString(v.valueBytes)
			if nil == e {
				v.status.parsed = true
			}
		}
		return v.value.str
	case jsoniter.BoolValue:
		return formatBool(v.value.boolean)
	case jsoniter.ObjectValue:
		return v.packObjChildren()
	case jsoniter.ArrayValue:
		return v.packArrChildren()
	}
}

func (v *V) packObjChildren() string {
	buf := bytes.Buffer{}
	v.bufObjChildren(&buf)
	return buf.String()
}

func (v *V) bufObjChildren(buf *bytes.Buffer) {
	buf.WriteByte('{')
	i := 0
	for k, v := range v.children.object {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v.String())
		i++
	}
	buf.WriteByte('}')
}

func (v *V) packArrChildren() string {
	buf := bytes.Buffer{}
	v.bufArrChildren(&buf)
	return buf.String()
}

func (v *V) bufArrChildren(buf *bytes.Buffer) {
	buf.WriteByte('[')
	v.RangeArray(func(i int, v *V) bool {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
		i++
		return true
	})
	buf.WriteByte(']')
}
