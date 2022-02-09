package jsonvalue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type float64Value struct {
	f float64
	s string
	b []byte
}

func newFloat64ByRaw(f float64, b []byte) *V {
	res := &float64Value{
		f: f,
		b: b,
	}
	return &V{
		impl: res,
	}
}

func newFloat64f(f float64, format byte, prec, bitsize int) *V {
	res := &float64Value{
		f: f,
	}
	if isValidFloat(f) {
		res.s = strconv.FormatFloat(f, format, prec, bitsize)
	}
	return &V{
		impl: res,
	}
}

// NewFloat64 returns an initialied num jsonvalue value by float64 type. The format and precision control is the same
// with encoding/json: https://github.com/golang/go/blob/master/src/encoding/json/encode.go#L575
//
// NewFloat64 根据指定的 flout64 类型返回一个初始化好的数字类型的 jsonvalue 值。数字转出来的字符串格式参照 encoding/json 中的逻辑。
func NewFloat64(f float64) *V {
	abs := math.Abs(f)
	format := byte('f')
	if abs < 1e-6 || abs >= 1e21 {
		format = byte('e')
	}

	return newFloat64f(f, format, -1, 64)
}

// NewFloat64f returns an initialied num jsonvalue value by float64 type. The format and prec parameter are used in
// strconv.FormatFloat(). Only 'f', 'E', 'e', 'G', 'g' formats are supported, other formats will be mapped to 'g'.
//
// NewFloat64f 根据指定的 float64 类型返回一个初始化好的数字类型的 jsonvalue 值。其中参数 format 和 prec 分别用于
// strconv.FormatFloat() 函数. 只有 'f'、'E'、'e'、'G'、'g' 格式是支持的，其他配置均统一映射为 'g'。
func NewFloat64f(f float64, format byte, prec int) *V {
	if err := validateFloat64Format(format); err != nil {
		format = 'g'
	}
	return newFloat64f(f, format, prec, 64)
}

// NewFloat32 returns an initialied num jsonvalue value by float32 type. The format and precision control is the same
// with encoding/json: https://github.com/golang/go/blob/master/src/encoding/json/encode.go#L575
//
// NewFloat32 根据指定的 float32 类型返回一个初始化好的数字类型的 jsonvalue 值。数字转出来的字符串格式参照 encoding/json 中的逻辑。
func NewFloat32(f float32) *V {
	f64 := float64(f)
	abs := math.Abs(f64)
	format := byte('f')
	if abs < 1e-6 || abs >= 1e21 {
		format = byte('e')
	}

	return newFloat64f(f64, format, -1, 32)
}

// NewFloat32f returns an initialied num jsonvalue value by float64 type. The format and prec parameter are used in
// strconv.FormatFloat(). Only 'f', 'E', 'e', 'G', 'g' formats are supported, other formats will be mapped to 'g'.
//
// NewFloat32f 根据指定的 float64 类型返回一个初始化好的数字类型的 jsonvalue 值。其中参数 format 和 prec 分别用于
// strconv.FormatFloat() 函数. 只有 'f'、'E'、'e'、'G'、'g' 格式是支持的，其他配置均统一映射为 'g'。
func NewFloat32f(f float32, format byte, prec int) *V {
	if err := validateFloat64Format(format); err != nil {
		format = 'g'
	}
	return newFloat64f(float64(f), format, prec, 64)
}

func isValidFloat(f float64) bool {
	if math.IsNaN(f) {
		return false
	}
	if math.IsInf(f, 0) {
		return false
	}
	return true
}

func validateFloat64Format(f byte) error {
	switch f {
	case 'f', 'E', 'e', 'G', 'g':
		return nil
	default:
		return fmt.Errorf("unsupported float value in option: %c", f)
	}
}

// ======== deleter interface ========

func (v *float64Value) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	return ErrNotFound
}

// ======== typper interface ========

func (v *float64Value) ValueType() ValueType {
	return Number
}

// ======== getter interface ========

func (v *float64Value) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	return &V{}, ErrNotFound
}

// ======== setter interface ========

func (v *float64Value) setAt(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Set()", v.ValueType())
}

// ======== iterater interface ========

func (v *float64Value) RangeObjects(callback func(k string, v *V) bool) {
	// do nothing
}

func (v *float64Value) RangeArray(callback func(i int, v *V) bool) {
	// do nothing
}

func (v *float64Value) ForRangeObj() map[string]*V {
	return map[string]*V{}
}

func (v *float64Value) ForRangeArr() []*V {
	return nil
}

func (v *float64Value) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter)
	close(ch)
	return ch
}

func (v *float64Value) IterArray() <-chan *ArrayIter {
	ch := make(chan *ArrayIter)
	close(ch)
	return ch
}

//  ======== marshaler interface ========

func (v *float64Value) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	if v.b != nil {
		buf.Write(v.b)
		return nil
	}
	if v.s != "" {
		buf.WriteString(v.s)
		return nil
	}

	// else, +Inf or -Inf or NaN
	if math.IsInf(v.f, 1) { // +Inf
		return marshalInfP(buf, opt)
	}
	if math.IsInf(v.f, -1) { // -Inf
		return marshalInfN(buf, opt)
	}

	return marshalNaN(buf, opt)
}

func marshalNaN(buf *bytes.Buffer, opt *Opt) error {
	switch opt.FloatNaNHandleType {
	default:
		fallthrough
	case FloatNaNTreatAsError:
		return fmt.Errorf("%w: %v", ErrUnsupportedFloat, math.NaN())

	case FloatNaNConvertToFloat:
		if !isValidFloat(opt.FloatNaNToFloat) {
			return fmt.Errorf("%w: %v", ErrUnsupportedFloatInOpt, opt.FloatNaNToFloat)
		}
		b, _ := json.Marshal(opt.FloatNaNToFloat)
		buf.Write(b)

	case FloatNaNNull:
		buf.WriteString("null")

	case FloatNaNConvertToString:
		if s := opt.FloatNaNToString; s == "" {
			buf.WriteString(`"NaN"`)
		} else {
			buf.WriteByte('"')
			escapeStringToBuff(s, buf, opt)
			buf.WriteByte('"')
		}
	}

	return nil
}

func marshalInfP(buf *bytes.Buffer, opt *Opt) error {
	switch opt.FloatInfHandleType {
	default:
		fallthrough
	case FloatInfTreatAsError:
		return fmt.Errorf("%w: %v", ErrUnsupportedFloat, math.Inf(1))

	case FloatInfConvertToFloat:
		if !isValidFloat(opt.FloatInfToFloat) {
			return fmt.Errorf("%w: %v", ErrUnsupportedFloatInOpt, opt.FloatInfToFloat)
		}
		b, _ := json.Marshal(opt.FloatInfToFloat)
		buf.Write(b)

	case FloatInfNull:
		buf.WriteString("null")

	case FloatInfConvertToString:
		if s := opt.FloatInfPositiveToString; s == "" {
			buf.WriteString(`"+Inf"`)
		} else {
			buf.WriteByte('"')
			escapeStringToBuff(s, buf, opt)
			buf.WriteByte('"')
		}
	}

	return nil
}

func marshalInfN(buf *bytes.Buffer, opt *Opt) error {
	switch opt.FloatInfHandleType {
	default:
		fallthrough
	case FloatInfTreatAsError:
		return fmt.Errorf("%w: %v", ErrUnsupportedFloat, math.Inf(-1))

	case FloatInfConvertToFloat:
		if !isValidFloat(opt.FloatInfToFloat) {
			return fmt.Errorf("%w: %v", ErrUnsupportedFloatInOpt, -opt.FloatInfToFloat)
		}
		b, _ := json.Marshal(-opt.FloatInfToFloat)
		buf.Write(b)

	case FloatInfNull:
		buf.WriteString("null")

	case FloatInfConvertToString:
		buf.WriteByte('"')
		if s := opt.FloatInfNegativeToString; s != "" {
			escapeStringToBuff(s, buf, opt)
		} else if opt.FloatInfPositiveToString != "" {
			s = "-" + strings.TrimLeft(opt.FloatInfPositiveToString, "+")
			escapeStringToBuff(s, buf, opt)
		} else {
			buf.WriteString(`-Inf`)
		}
		buf.WriteByte('"')
	}

	return nil
}

// ======== valuer interface ========

func (v *float64Value) Bool() (bool, error) {
	return v.f != 0, ErrTypeNotMatch
}

func (v *float64Value) Int64() (int64, error) {
	return int64(v.f), nil
}

func (v *float64Value) Uint64() (uint64, error) {
	return uint64(v.f), nil
}

func (v *float64Value) Float64() (float64, error) {
	return v.f, nil
}

func (v *float64Value) String() string {
	if v.s != "" {
		return v.s
	}
	if v.b != nil {
		return string(v.b)
	}
	return fmt.Sprint(v.f)
}

func (v *float64Value) Len() int {
	return 0
}

// ======== numberAsserter interface ========

func (v *float64Value) IsFloat() bool {
	return true
}

func (v *float64Value) IsInteger() bool {
	return false
}

func (v *float64Value) IsNegative() bool {
	return v.f < 0
}

func (v *float64Value) IsPositive() bool {
	return v.f >= 0
}

func (v *float64Value) GreaterThanInt64Max() bool {
	return false
}

// ======== inserter interface ========

func (v *float64Value) insertBefore(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

func (v *float64Value) insertAfter(child *V, firstParam interface{}, otherParams ...interface{}) error {
	return fmt.Errorf("%v type does not supports Insert()", v.ValueType())
}

// ======= appender interface ========

func (v *float64Value) appendInTheBeginning(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}

func (v *float64Value) appendInTheEnd(child *V, params ...interface{}) error {
	return fmt.Errorf("%v type does not supports Append()", v.ValueType())
}
