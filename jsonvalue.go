// Package jsonvalue is for JSON parsing and setting. It is used in situations those
// Go structures cannot achieve, or "map[string]any" could not do properly.
//
// As a quick start:
//
//	v := jsonvalue.NewObject()
//	v.SetString("Hello, JSON").At("someObject", "someObject", "someObject", "message")  // automatically create sub objects
//	fmt.Println(v.MustMarshalString())                                                  // marshal to string type. Use MustMarshal if you want []byte instead.
//	// Output:
//	// {"someObject":{"someObject":{"someObject":{"message":"Hello, JSON!"}}}
//
// If you want to parse raw JSON data, use Unmarshal()
//
//	raw := []byte(`{"message":"hello, world"}`)
//	v, err := jsonvalue.Unmarshal(raw)
//	s, _ := v.GetString("message")
//	fmt.Println(s)
//	// Output:
//	// hello, world
//
// jsonvalue 包用于 JSON 的解析（反序列化）和编码（序列化）。通常情况下我们用 struct 来处理
// 结构化的 JSON，但是有时候使用 struct 不方便或者是功能不足的时候，go 一般而言使用的是
// "map[string]any"，但是后者也有很多不方便的地方。本包即是用于替代这些不方便的情况的。
//
// 快速上手：
//
//	v := jsonvalue.NewObject()
//	v.SetString("Hello, JSON").At("someObject", "someObject", "someObject", "message")  // 自动创建子成员
//	fmt.Println(v.MustMarshalString())                                                  // 序列化为 string 类型，如果你要 []byte 类型，则使用 MustMarshal 函数。
//	// 输出:
//	// {"someObject":{"someObject":{"someObject":{"message":"Hello, JSON!"}}}
//
// 如果要反序列化原始的 JSON 文本，则使用 Unmarshal():
//
//	raw := []byte(`{"message":"hello, world"}`)
//	v, err := jsonvalue.Unmarshal(raw)
//	s, _ := v.GetString("message")
//	fmt.Println(s)
//	// 输出:
//	// hello, world
package jsonvalue

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/Andrew-M-C/go.jsonvalue/internal/buffer"
	"github.com/Andrew-M-C/go.jsonvalue/internal/unsafe"
)

// ValueType identifying JSON value type
type ValueType int

const (
	// NotExist type tells that this JSON value is not exist or legal
	NotExist ValueType = iota
	// String JSON string type
	String
	// Number JSON number type
	Number
	// Object JSON object type
	Object
	// Array JSON array type
	Array
	// Boolean JSON boolean type
	Boolean
	// Null JSON null type
	Null
	// Unknown unknown JSON type
	Unknown
)

var typeStr = [Unknown + 1]string{
	"illegal",
	"string",
	"number",
	"object",
	"array",
	"boolean",
	"null",
	"unknown",
}

// String show the type name of JSON
func (t ValueType) String() string {
	if t > Unknown {
		t = NotExist
	} else if t < 0 {
		t = NotExist
	}
	return typeStr[int(t)]
}

// ValueType returns the type of this JSON value.
func (v *V) ValueType() ValueType {
	return v.valueType
}

// V is the main type of jsonvalue, representing a JSON value.
//
// V 是 jsonvalue 的主类型，表示一个 JSON 值。
type V struct {
	valueType ValueType

	srcByte []byte

	num       num
	valueStr  string
	valueBool bool
	children  children

	marshalToBuffer func(parentInfo *ParentInfo, buf buffer.Buffer, opt *Opt) error
}

type num struct {
	negative bool
	floated  bool
	i64      int64
	u64      uint64
	f64      float64
}

type childWithProperty struct {
	id uint32
	v  *V
}

type children struct {
	incrID uint32
	arr    []*V
	object map[string]childWithProperty

	// As official json package supports caseless key accessing, I decide to do it as well
	lowerCaseKeys map[string]map[string]struct{}
}

func (c *children) deepCopy() children {
	res := children{
		incrID: c.incrID,
	}

	// if length or arr > 0, this must be an array type
	if len(c.arr) > 0 {
		for _, v := range c.arr {
			res.arr = append(res.arr, v.deepCopy())
		}
		return res
	}

	// if this is an object?
	if c.object != nil {
		res.object = make(map[string]childWithProperty, len(c.object))
		for key, item := range c.object {
			res.object[key] = childWithProperty{
				id: item.id,
				v:  item.v.deepCopy(),
			}
		}
	}

	// no need to copy lowerCaseKeys because it could be rebuild after calling
	// Caseless() next time

	return res
}

func (v *V) addCaselessKey(k string) {
	if v.children.lowerCaseKeys == nil {
		return
	}
	lowerK := strings.ToLower(k)
	keys, exist := v.children.lowerCaseKeys[lowerK]
	if !exist {
		keys = make(map[string]struct{})
		v.children.lowerCaseKeys[lowerK] = keys
	}
	keys[k] = struct{}{}
}

func (v *V) delCaselessKey(k string) {
	if v.children.lowerCaseKeys == nil {
		return
	}
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

// MustUnmarshalString just like UnmarshalString(). If error occurs, a JSON value
// with "NotExist" type would be returned, which could do nothing and return nothing
// in later use. It is useful to shorten codes.
//
// MustUnmarshalString 的逻辑与 UnmarshalString() 相同，不过如果错误的话，会返回一个类型未
// "NotExist" 的 JSON 值，这个值在后续的操作中将无法返回有效的数据，或者是执行任何有效的操作。
// 但起码不会导致程序 panic，便于使用短代码实现一些默认逻辑。
func MustUnmarshalString(s string) *V {
	v, _ := UnmarshalString(s)
	return v
}

// UnmarshalString is equivalent to Unmarshal([]byte(b)), but more efficient.
//
// UnmarshalString 等效于 Unmarshal([]byte(b))，但效率更高。
func UnmarshalString(s string) (*V, error) {
	// reference: https://stackoverflow.com/questions/41591097/slice-bounds-out-of-range-when-using-unsafe-pointer
	// sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// bh := reflect.SliceHeader{
	// 	Data: sh.Data,
	// 	Len:  sh.Len,
	// 	Cap:  sh.Len,
	// }
	// b := *(*[]byte)(unsafe.Pointer(&bh))
	b := []byte(s)
	return UnmarshalNoCopy(b)
}

// MustUnmarshal just like Unmarshal(). If error occurres, a JSON value with "NotExist"
// type would be returned, which could do nothing and return nothing in later use. It
// is useful to shorten codes.
//
// MustUnmarshal 的逻辑与 Unmarshal() 相同，不过如果错误的话，会返回一个类型未 "NotExist" 的
// JSON 值，这个值在后续的操作中将无法返回有效的数据，或者是执行任何有效的操作。但起码不会导致程序
// panic，便于使用短代码实现一些默认逻辑。
func MustUnmarshal(b []byte) *V {
	v, _ := Unmarshal(b)
	return v
}

// Unmarshal parse raw bytes(encoded in UTF-8 or pure AscII) and returns a *V instance.
//
// Unmarshal 解析原始的字节类型数据（以 UTF-8 或纯 AscII 编码），并返回一个 *V 对象。
func Unmarshal(b []byte) (ret *V, err error) {
	le := len(b)
	if le == 0 {
		return nil, ErrNilParameter
	}

	trueB := make([]byte, len(b))
	copy(trueB, b)
	it := iter(trueB)
	p := newPool(len(b))
	ret, err = unmarshalWithIter(p, it, 0)
	p.release()
	return
}

// MustUnmarshalNoCopy just like UnmarshalNoCopy(). If error occurres, a JSON value
// with "NotExist" type would be returned, which could do nothing and return nothing
// in later use. It is useful to shorten codes.
//
// MustUnmarshalNoCopy 的逻辑与 UnmarshalNoCopy() 相同，不过如果错误的话，会返回一个类型未
// "NotExist" 的 JSON 值，这个值在后续的操作中将无法返回有效的数据，或者是执行任何有效的操作。
// 但起码不会导致程序 panic，便于使用短代码实现一些默认逻辑。
func MustUnmarshalNoCopy(b []byte) *V {
	v, _ := UnmarshalNoCopy(b)
	return v
}

// UnmarshalNoCopy is same as Unmarshal, but it does not copy another []byte instance
// for saving CPU time. But pay attention that the input []byte may be used as buffer
// by jsonvalue and mey be modified.
//
// UnmarshalNoCopy 与 Unmarshal 相同，但是这个函数在解析过程中不会重新复制一个 []byte，对于大
// json 的解析而言能够大大节省时间。但请注意传入的 []byte 变量可能会被 jsonvalue 用作缓冲区，并进行修改
func UnmarshalNoCopy(b []byte) (ret *V, err error) {
	le := len(b)
	if le == 0 {
		return &V{}, ErrNilParameter
	}
	p := newPool(len(b))
	ret, err = unmarshalWithIter(p, iter(b), 0)
	p.release()
	return
}

// ==== type access ====

// IsObject tells whether value is an object
//
// IsObject 判断当前值是不是一个对象类型
func (v *V) IsObject() bool {
	return v.valueType == Object
}

// IsArray tells whether value is an array
//
// IsArray 判断当前值是不是一个数组类型
func (v *V) IsArray() bool {
	return v.valueType == Array
}

// IsString tells whether value is a string
//
// IsString 判断当前值是不是一个字符串类型
func (v *V) IsString() bool {
	return v.valueType == String
}

// IsNumber tells whether value is a number
//
// IsNumber 判断当前值是不是一个数字类型
func (v *V) IsNumber() bool {
	return v.valueType == Number
}

// IsFloat tells whether value is a float point number. If there is no decimal point
// in original text, it returns false while IsNumber returns true.
//
// IsFloat 判断当前值是不是一个浮点数类型。如果给定的数不包含小数点，那么即便是数字类型，该函数也会返回 false.
func (v *V) IsFloat() bool {
	if v.valueType != Number {
		return false
	}
	return v.num.floated
}

// IsInteger tells whether value is a fix point interger
//
// IsNumber 判断当前值是不是一个定点数整型
func (v *V) IsInteger() bool {
	if v.valueType != Number {
		return false
	}
	return !(v.num.floated)
}

// IsNegative tells whether value is a negative number
//
// IsNegative 判断当前值是不是一个负数
func (v *V) IsNegative() bool {
	if v.valueType != Number {
		return false
	}
	return v.num.negative
}

// IsPositive tells whether value is a positive number
//
// IsPositive 判断当前值是不是一个正数
func (v *V) IsPositive() bool {
	if v.valueType != Number {
		return false
	}
	return !(v.num.negative)
}

// GreaterThanInt64Max return true when ALL conditions below are met:
//  1. It is a number value.
//  2. It is a positive integer.
//  3. Its value is greater than 0x7fffffffffffffff.
//
// GreaterThanInt64Max 判断当前值是否超出 int64 可表示的范围。当以下条件均成立时，返回 true,
// 否则返回 false：
//  1. 是一个数字类型值.
//  2. 是一个正整型数字.
//  3. 该正整数的值大于 0x7fffffffffffffff.
func (v *V) GreaterThanInt64Max() bool {
	if v.valueType != Number {
		return false
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
	return v.valueType == Boolean
}

// IsNull tells whether value is a null
//
// IsBoolean 判断当前值是不是一个空类型
func (v *V) IsNull() bool {
	return v.valueType == Null
}

// Bool returns represented bool value. If value is not boolean, returns false.
//
// Bool 返回布尔类型值。如果当前值不是布尔类型，则判断是否为 string，不是 string 返回 false;
// 是 string 的话则返回字面值是否等于 true
func (v *V) Bool() bool {
	if v.valueType == Boolean {
		return v.valueBool
	}
	b, _ := getBoolAndErrorFromValue(v)
	return b.valueBool
}

// Int returns represented int value. If value is not a number, returns zero.
//
// Int 返回 int 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int() int {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Int()
	}
	return int(v.num.i64)
}

// Uint returns represented uint value. If value is not a number, returns zero.
//
// Uint 返回 uint 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint() uint {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Uint()
	}
	return uint(v.num.u64)
}

// Int64 returns represented int64 value. If value is not a number, returns zero.
//
// Int64 返回 int64 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int64() int64 {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Int64()
	}
	return int64(v.num.i64)
}

// Uint64 returns represented uint64 value. If value is not a number, returns zero.
//
// Uint64 返回 uint64 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint64() uint64 {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Uint64()
	}
	return uint64(v.num.u64)
}

// Int32 returns represented int32 value. If value is not a number, returns zero.
//
// Int32 返回 int32 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Int32() int32 {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Int32()
	}
	return int32(v.num.i64)
}

// Uint32 returns represented uint32 value. If value is not a number, returns zero.
//
// Uint32 返回 uint32 类型值。如果当前值不是数字类型，则返回 0。
func (v *V) Uint32() uint32 {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Uint32()
	}
	return uint32(v.num.u64)
}

// Float64 returns represented float64 value. If value is not a number, returns zero.
//
// Float64 返回 float64 类型值。如果当前值不是数字类型，则返回 0.0。
func (v *V) Float64() float64 {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Float64()
	}
	return v.num.f64
}

// Float32 returns represented float32 value. If value is not a number, returns zero.
//
// Float32 返回 float32 类型值。如果当前值不是数字类型，则返回 0.0。
func (v *V) Float32() float32 {
	if v.valueType != Number {
		return getNumberFromNotNumberValue(v).Float32()
	}
	return float32(v.num.f64)
}

// Bytes returns represented binary data which is encoede as Base64 string. []byte{}
// would be returned if value is
// not a string type or base64 decode failed.
//
// Bytes 返回以 Base64 编码在 string 类型中的二进制数据。如果当前值不是字符串类型，或者是 base64
// 编码失败，则返回 []byte{}。
func (v *V) Bytes() []byte {
	if v.valueType != String {
		return []byte{}
	}
	b, err := internal.b64.DecodeString(v.valueStr)
	if err != nil {
		return []byte{}
	}
	return b
}

func (v *V) deepCopy() *V {
	if v == nil {
		return &V{}
	}

	switch v.valueType {
	default:
		// 	fallthrough
		// case NotExist, Unknown:
		return &V{}
	case String:
		return NewString(v.String())
	case Number:
		res := new(globalPool{}, Number)
		res.num = v.num
		res.srcByte = v.srcByte
		return res
	case Object:
		res := new(globalPool{}, Object)
		res.children = v.children.deepCopy()
		return res
	case Array:
		res := new(globalPool{}, Array)
		res.children = v.children.deepCopy()
		return res
	case Boolean:
		return NewBool(v.Bool())
	case Null:
		return NewNull()
	}
}

type deepCopier interface {
	deepCopy() *V
}

// String returns represented string value or the description for the jsonvalue.V
// instance if it is not a string.
//
// String 返回 string 类型值。如果当前值不是字符串类型，则返回当前 *V 类型的描述说明。
func (v *V) String() string {
	if v == nil {
		return "nil"
	}
	switch v.valueType {
	default:
		return ""
	case Null:
		return "null"
	case Number:
		if len(v.srcByte) > 0 {
			return unsafe.BtoS(v.srcByte)
		}
		return strconv.FormatFloat(v.num.f64, 'g', -1, 64)
	case String:
		return v.valueStr
	case Boolean:
		return formatBool(v.valueBool)
	case Object:
		return v.packObjChildren()
	case Array:
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
		buf.WriteString(v.v.String())
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
		return true
	})
	buf.WriteByte(']')
}
