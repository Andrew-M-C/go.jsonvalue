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
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/buger/jsonparser"
)

var (
	b64 = base64.StdEncoding
)

// test:
// go test -v -failfast -cover -coverprofile ./cover.out && go tool cover -html=./cover.out -o ./cover.html && open ./cover.html

// V is the main type of jsonvalue, representing a JSON value.
//
// V 是 jsonvalue 的主类型，表示一个 JSON 值。
type V struct {
	valueType  jsonparser.ValueType
	valueBytes []byte
	parsed     bool

	num       num
	valueStr  string
	valueBool bool
	children  children
}

type num struct {
	negative bool
	floated  bool
	i64      int64
	u64      uint64
	f64      float64
}

type children struct {
	object map[string]*V
	array  *list.List

	// As official json package supports caseless key accessing, I decide to di it as well
	lowerCaseKeys map[string]map[string]struct{}
}

func new(t jsonparser.ValueType) *V {
	v := V{}
	v.valueType = t
	return &v
}

func newObject() *V {
	v := new(jsonparser.Object)
	v.children = children{
		object:        make(map[string]*V),
		lowerCaseKeys: make(map[string]map[string]struct{}),
	}
	return v
}

func newArray() *V {
	v := new(jsonparser.Array)
	v.children = children{
		array: list.New(),
	}
	return v
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

// UnmarshalString is equavilent to Unmarshal(unsafeBtoS(b)), but much more efficient.
//
// UnmarshalString 等效于 Unmarshal(unsafeBtoS(b))，但效率更高。
func UnmarshalString(s string) (*V, error) {
	// reference: https://stackoverflow.com/questions/41591097/slice-bounds-out-of-range-when-using-unsafe-pointer
	// sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// bh := reflect.SliceHeader{
	// 	Data: sh.Data,
	// 	Len:  sh.Len,
	// 	Cap:  sh.Len,
	// }
	// b := *(*[]byte)(unsafe.Pointer(&bh))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sh.Cap = sh.Len
	b := *(*[]byte)(unsafe.Pointer(sh))

	trueB := make([]byte, len(b))
	copy(trueB, b)
	return Unmarshal(trueB)
}

// TODO:
func unmarshalWithIter(it *iter, offset, end int) (*V, error) {
	offset, reachEnd := it.skipBlanks(offset, end)
	if reachEnd {
		return nil, fmt.Errorf("%w, cannot find any symbol characters found", ErrRawBytesUnrecignized)
	}

	chr := it.b[offset]
	switch chr {
	case '{':
		// TODO: object
	case '[':
		// TODO: array
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		// TODO: number
	case '"':
		sectLenWithoutQuote, _, err := it.parseStrFromBytesForwardWithQuote(offset)
		if err != nil {
			return nil, err
		}
		return NewString(unsafeBtoS(it.b[offset+1 : offset+1+sectLenWithoutQuote])), nil

	case 't':
		_, err := it.parseTrue(offset)
		if err != nil {
			return nil, err
		}
		return newTrue(), nil

	case 'f':
		_, err := it.parseFalse(offset)
		if err != nil {
			return nil, err
		}
		return newFalse(), nil

	case 'n':
		// TODO: null
	default:
		return nil, fmt.Errorf("%w, invalid character \\u%04X at Position %d", ErrRawBytesUnrecignized, chr, offset)
	}
	return nil, ErrRawBytesUnrecignized // TODO: 后续去掉
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
			return newFromObject(b[i:])
		case '[':
			return newFromArray(b[i:])
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			ret, err = newFromNumber(b[i:])
			if err != nil {
				return
			}
			err = ret.parseNumber()
			if err != nil {
				return nil, err
			}
			return ret, nil

		case '"':
			ret = new(jsonparser.String)
			ret.valueStr, ret.valueBytes, err = parseString(b[i:])
			if err != nil {
				return nil, err
			}
			ret.parsed = true
			return ret, nil

		case 't':
			return newFromTrue(b[i:])
		case 'f':
			return newFromFalse(b[i:])
		case 'n':
			return newFromNull(b[i:])
		default:
			return nil, ErrRawBytesUnrecignized
		}
	}

	return nil, ErrRawBytesUnrecignized
}

var dot = []byte{'.'}

func (v *V) parseNumber() (err error) {
	b := v.valueBytes

	// if v.num == nil {
	// 	v.num = &num{}
	// }

	if bytes.Contains(b, dot) {
		v.num.floated = true
		v.num.f64, err = parseFloat(b)
		if err != nil {
			return
		}

		v.parsed = true
		v.num.negative = (v.num.f64 < 0)
		v.num.i64 = int64(v.num.f64)
		v.num.u64 = uint64(v.num.f64)

	} else if b[0] == '-' {
		v.num.negative = true
		v.num.i64, err = parseInt(b)
		if err != nil {
			return
		}

		v.parsed = true
		v.num.u64 = uint64(v.num.i64)
		v.num.f64 = float64(v.num.i64)

	} else {
		v.num.negative = false
		v.num.u64, err = parseUint(b)
		if err != nil {
			return
		}

		v.parsed = true
		v.num.i64 = int64(v.num.u64)
		v.num.f64 = float64(v.num.u64)
	}

	return nil
}

// ==== simple object parsing ====
func newFromNumber(b []byte) (ret *V, err error) {
	v := new(jsonparser.Number)
	v.valueBytes = b
	return v, nil
}

// func newFromString(b []byte) (ret *V, err error) {
// 	v := new()
// 	v.valueType = jsonparser.String
// 	v.valueBytes = b
// 	return v, nil
// }

var (
	trueBytes  = []byte{'t', 'r', 'u', 'e'}
	falseBytes = []byte{'f', 'a', 'l', 's', 'e'}
)

func newTrue() *V {
	v := new(jsonparser.Boolean)
	v.parsed = true
	v.valueBytes = trueBytes
	v.valueBool = true
	return v
}

func newFalse() *V {
	v := new(jsonparser.Boolean)
	v.parsed = true
	v.valueBytes = falseBytes
	v.valueBool = false
	return v
}

func newFromTrue(b []byte) (ret *V, err error) {
	if len(b) != 4 || unsafeBtoS(b) != "true" {
		return nil, ErrNotValidBoolValue
	}
	v := new(jsonparser.Boolean)
	v.parsed = true
	v.valueBytes = []byte{'t', 'r', 'u', 'e'}
	v.valueBool = true
	return v, nil
}

func newFromFalse(b []byte) (ret *V, err error) {
	if len(b) != 5 || unsafeBtoS(b) != "false" {
		return nil, ErrNotValidBoolValue
	}
	v := new(jsonparser.Boolean)
	v.parsed = true
	v.valueBytes = []byte{'f', 'a', 'l', 's', 'e'}
	v.valueBool = false
	return v, nil
}

func newFromBool(b []byte) (ret *V, err error) {
	v := new(jsonparser.Boolean)

	switch unsafeBtoS(b) {
	case "true":
		v.parsed = true
		v.valueBytes = []byte{'t', 'r', 'u', 'e'}
		v.valueBool = true
	case "false":
		v.parsed = true
		v.valueBytes = []byte{'f', 'a', 'l', 's', 'e'}
		v.valueBool = false
	default:
		return nil, ErrNotValidBoolValue
	}

	return v, nil
}

func newFromNull(b []byte) (ret *V, err error) {
	if len(b) != 4 || unsafeBtoS(b) != "null" {
		return nil, ErrNotValidBoolValue
	}
	v := new(jsonparser.Null)
	v.parsed = true
	return v, nil
}

// ====
func newFromArray(b []byte) (ret *V, err error) {
	o := newArray()

	jsonparser.ArrayEach(b, func(v []byte, t jsonparser.ValueType, _ int, _ error) {
		if err != nil {
			return
		}

		var child *V

		switch t {
		default:
			err = fmt.Errorf("invalid value type: %v", t)
		case jsonparser.Object:
			child, err = newFromObject(v)
		case jsonparser.Array:
			child, err = newFromArray(v)
		case jsonparser.Number:
			child, err = newFromNumber(v)
		case jsonparser.Boolean:
			child, err = newFromBool(v)
		case jsonparser.Null:
			child, err = newFromNull(v)
		case jsonparser.String:
			s, err := parseStringNoQuote(v)
			if err != nil {
				return
			}
			child = new(jsonparser.String)
			child.parsed = true
			child.valueStr = s
		}

		if err != nil {
			return
		}
		o.children.array.PushBack(child)
	})

	// done
	if err != nil {
		return
	}
	return o, nil
}

// ==== object parsing ====
func newFromObject(b []byte) (ret *V, err error) {
	o := newObject()

	err = jsonparser.ObjectEach(b, func(k, v []byte, t jsonparser.ValueType, _ int) error {
		// key
		var child *V
		key, err := parseStringNoQuote(k)
		if err != nil {
			return err
		}

		switch t {
		default:
			return fmt.Errorf("invalid value type: %v", t)
		case jsonparser.Object:
			child, err = newFromObject(v)
		case jsonparser.Array:
			child, err = newFromArray(v)
		case jsonparser.Number:
			child, err = newFromNumber(v)
		case jsonparser.Boolean:
			child, err = newFromBool(v)
		case jsonparser.Null:
			child, err = newFromNull(v)
		case jsonparser.String:
			s, err := parseStringNoQuote(v)
			if err != nil {
				return err
			}
			child = new(jsonparser.String)
			child.parsed = true
			child.valueStr = s
		}

		if err != nil {
			return err
		}
		o.setToObjectChildren(key, child)
		return nil
	})

	// done
	if err != nil {
		return
	}
	return o, nil
}

// ==== type access ====

// IsObject tells whether value is an object
//
// IsObject 判断当前值是不是一个对象类型
func (v *V) IsObject() bool {
	return v.valueType == jsonparser.Object
}

// IsArray tells whether value is an array
//
// IsArray 判断当前值是不是一个数组类型
func (v *V) IsArray() bool {
	return v.valueType == jsonparser.Array
}

// IsString tells whether value is a string
//
// IsString 判断当前值是不是一个字符串类型
func (v *V) IsString() bool {
	return v.valueType == jsonparser.String
}

// IsNumber tells whether value is a number
//
// IsNumber 判断当前值是不是一个数字类型
func (v *V) IsNumber() bool {
	return v.valueType == jsonparser.Number
}

// IsFloat tells whether value is a float point number. If there is no decimal point in original text, it returns false
// while IsNumber returns true.
//
// IsFloat 判断当前值是不是一个浮点数类型。如果给定的数不包含小数点，那么即便是数字类型，该函数也会返回 false.
func (v *V) IsFloat() bool {
	if v.valueType != jsonparser.Number {
		return false
	}
	if !v.parsed {
		v.parseNumber()
	}
	return v.num.floated
}

// IsInteger tells whether value is a fix point interger
//
// IsNumber 判断当前值是不是一个定点数整型
func (v *V) IsInteger() bool {
	if v.valueType != jsonparser.Number {
		return false
	}
	if !v.parsed {
		err := v.parseNumber()
		if err != nil {
			return false
		}
	}
	return !(v.num.floated)
}

// IsNegative tells whether value is a negative number
//
// IsNegative 判断当前值是不是一个负数
func (v *V) IsNegative() bool {
	if v.valueType != jsonparser.Number {
		return false
	}
	if !v.parsed {
		v.parseNumber()
	}
	return v.num.negative
}

// IsPositive tells whether value is a positive number
//
// IsPositive 判断当前值是不是一个正数
func (v *V) IsPositive() bool {
	if v.valueType != jsonparser.Number {
		return false
	}
	if !v.parsed {
		err := v.parseNumber()
		if err != nil {
			return false
		}
	}
	return !(v.num.negative)
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
	if v.valueType != jsonparser.Number {
		return false
	}
	if !v.parsed {
		v.parseNumber()
	}
	if v.num.negative {
		return false
	}
	return v.num.u64 > 0x7fffffffffffffff
}

// IsBoolean tells whether value is a boolean
//
// IsBoolean 判断当前值是不是一个布尔类型
func (v *V) IsBoolean() bool {
	return v.valueType == jsonparser.Boolean
}

// IsNull tells whether value is a null
//
// IsBoolean 判断当前值是不是一个空类型
func (v *V) IsNull() bool {
	return v.valueType == jsonparser.Null
}

// ==== value access ====

func getNumberFromNotNumberValue(v *V) *V {
	if !v.IsString() {
		return NewInt(0)
	}
	ret, _ := newFromNumber([]byte(v.valueStr))
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
	return v.valueBool
}

// Int returns represented int value. If value is not a number, returns zero.
//
// Int 返回 int 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int() int {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Int()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return int(v.num.i64)
}

// Uint returns represented uint value. If value is not a number, returns zero.
//
// Uint 返回 uint 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint() uint {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Uint()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return uint(v.num.u64)
}

// Int64 returns represented int64 value. If value is not a number, returns zero.
//
// Int64 返回 int64 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int64() int64 {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Int64()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return int64(v.num.i64)
}

// Uint64 returns represented uint64 value. If value is not a number, returns zero.
//
// Uint64 返回 uint64 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint64() uint64 {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Uint64()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return uint64(v.num.u64)
}

// Int32 returns represented int32 value. If value is not a number, returns zero.
//
// Int32 返回 int32 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int32() int32 {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Int32()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return int32(v.num.i64)
}

// Uint32 returns represented uint32 value. If value is not a number, returns zero.
//
// Uint32 返回 uint32 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint32() uint32 {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Uint32()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return uint32(v.num.u64)
}

// Float64 returns represented float64 value. If value is not a number, returns zero.
//
// Float64 返回 float64 类型值。如果当前值不是数字类型，则返回 0.0。
func (v *V) Float64() float64 {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Float64()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return v.num.f64
}

// Float32 returns represented float32 value. If value is not a number, returns zero.
//
// Float32 返回 float32 类型值。如果当前值不是数字类型，则返回 0.0。
func (v *V) Float32() float32 {
	if v.valueType != jsonparser.Number {
		return getNumberFromNotNumberValue(v).Float32()
	}
	if !v.parsed {
		v.parseNumber()
	}
	return float32(v.num.f64)
}

// Bytes returns represented binary data which is encoede as Base64 string. []byte{} would be returned if value is
// not a string type or base64 decode failed.
//
// Bytes 返回以 Base64 编码在 string 类型中的二进制数据。如果当前值不是字符串类型，或者是 base64 编码失败，则返回 []byte{}。
func (v *V) Bytes() []byte {
	if v.valueType != jsonparser.String {
		return []byte{}
	}
	b, err := b64.DecodeString(v.valueStr)
	if err != nil {
		return []byte{}
	}
	return b
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
	case jsonparser.Null:
		return "null"
	case jsonparser.Number:
		return unsafeBtoS(v.valueBytes)
	case jsonparser.String:
		if !v.parsed {
			var e error
			v.valueStr, v.valueBytes, e = parseString(v.valueBytes)
			if nil == e {
				v.parsed = true
			}
		}
		return v.valueStr
	case jsonparser.Boolean:
		return formatBool(v.valueBool)
	case jsonparser.Object:
		return v.packObjChildren()
	case jsonparser.Array:
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
