package jsonvalue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Opt is the option of jsonvalue in marshaling.
//
// Opt 表示序列化当前 jsonvalue 类型时的参数
type Opt struct {
	// OmitNull tells how to handle null json value. The default value is false.
	// If OmitNull is true, null value will be omitted when marshaling.
	//
	// OmitNull 表示是否忽略 JSON 中的 null 类型值。默认为 false.
	OmitNull bool

	// MarshalLessFunc is used to handle sequences of marshaling. Since object is
	// implemented by hash map, the sequence of keys is unexpectable. For situations
	// those need settled JSON key-value sequence, please use MarshalLessFunc.
	//
	// Note: Elements in an array value would NOT trigger this function as they are
	// already sorted.
	//
	// We provides a example DefaultStringSequence. It is quite useful when calculating
	// idempotence of a JSON text, as key-value sequences should be fixed.
	//
	// MarshalLessFunc 用于处理序列化 JSON 对象类型时，键值对的顺序。由于 object 类型是采用 go 原生的 map 类型，采用哈希算法实现，
	// 因此其键值对的顺序是不可控的。而为了提高效率，jsonvalue 的内部实现中并不会刻意保存键值对的顺序。如果有必要在序列化时固定键值对顺序的话，
	// 可以使用这个函数。
	//
	// 注意：array 类型中键值对的顺序不受这个函数的影响
	//
	// 此外，我们提供了一个例子: DefaultStringSequence。当需要计算 JSON 文本的幂等值时，
	// 由于需要不变的键值对顺序，因此这个函数是非常有用的。
	MarshalLessFunc MarshalLessFunc

	// MarshalKeySequence is used to handle sequance of marshaling. This is much simpler
	// than MarshalLessFunc, just pass a string slice identifying key sequence. For keys
	// those are not in this slice, they would be appended in the end according to result
	// of Go string comparing. Therefore this parameter is useful for ensure idempotence.
	//
	// MarshalKeySequence 也用于处理序列化时的键值对顺序。与 MarshalLessFunc 不同，这个只需要用字符串切片的形式指定键的顺序即可，
	// 实现上更为简易和直观。对于那些不在指定切片中的键，那么将会统一放在结尾，并且按照 go 字符串对比的结果排序。也可以保证幂等。
	MarshalKeySequence []string
	keySequence        map[string]int // generated from MarshalKeySequence

	// FloatNaNHandleType tells what to deal with float NaN.
	//
	// FloatNaNHandleType 表示当处理 float 的时候，如果遇到了 NaN 的话，要如何处理。
	FloatNaNHandleType FloatNaNHandleType
	// FloatNaNToString works with FloatNaNHandleType = FloatNaNConvertToString. It tells what string to replace
	// to with NaN. If not specified, NaN will be set as string "NaN".
	//
	// FloatNaNToString 搭配 FloatNaNHandleType = FloatNaNConvertToString 使用，表示将 NaN 映射为哪个字符串。
	// 这个值如果不指定，则默认会被设置为字符串 "NaN"
	FloatNaNToString string
	// FloatNaNToFloat works with FloatNaNHandleType = FloatNaNConvertToFloat. It tells what float number will
	// be mapped to as for NaN. NaN, +Inf or -Inf are not allowed for this option.
	//
	// FloatNaNToFloat 搭配 FloatNaNHandleType = FloatNaNConvertToFloat 使用，表示将 NaN 映射为哪个 float64 值。
	// 不允许指定为 NaN, +Inf 或 -Inf。如果不指定，则映射为 0
	FloatNaNToFloat float64

	// FloatInfHandleType tells what to deal with float +Inf and -Inf.
	//
	// FloatInfHandleType 表示当处理 float 的时候，如果遇到了 +Inf 和 -Inf 的话，要如何处理。
	FloatInfHandleType FloatInfHandleType
	// FloatInfPositiveToString works with FloatInfHandleType = FloatInfConvertToFloat. It tells what float number will
	// be mapped to as for +Inf. If not specified, +Inf will be set as string "+Inf"
	//
	// FloatInfPositiveToString 搭配 FloatInfHandleType = FloatInfConvertToFloat 使用，表示将 NaN 映射为哪个字符串。
	// 这个值如果不指定，则默认会被设置为字符串 "+Inf"
	FloatInfPositiveToString string
	// FloatInfNegativeToString works with FloatInfHandleType = FloatInfConvertToFloat. It tells what float number will
	// be mapped to as for -Inf. If not specified, -Inf will be set as string "-" + strings.TrimLeft(FloatInfPositiveToString, "+").
	//
	// FloatInfNegativeToString 搭配 FloatInfHandleType = FloatInfConvertToFloat 使用，表示将 NaN 映射为哪个字符串。
	// 这个值如果不指定，则默认会被设置为字符串 "-" + strings.TrimLeft(FloatInfPositiveToString, "+")。
	FloatInfNegativeToString string
	// FloatInfToFloat works with FloatInfHandleType = FloatInfConvertToFloat. It tells what float numbers will be
	// mapped to as for +Inf. And -Inf will be specified as the negative value of this option.
	// +Inf or -Inf are not allowed for this option.
	//
	// FloatInfToFloat 搭配 FloatInfHandleType = FloatInfConvertToFloat 使用，表示将 +Inf 映射为哪个 float64 值。而 -Inf
	// 则会被映射为这个值的负数。
	// 不允许指定为 NaN, +Inf 或 -Inf。如果不指定，则映射为 0
	FloatInfToFloat float64
}

type FloatNaNHandleType uint8

const (
	// FloatNaNTreatAsError indicates that error will be returned when a float number is NaN when marshaling.
	//
	// FloatNaNTreatAsError 表示当 marshal 遇到 NaN 时，返回错误。这是默认选项。
	FloatNaNTreatAsError FloatNaNHandleType = 0
	// FloatNaNConvertToFloat indicates that NaN will be replaced as another float number when marshaling. This option
	// works with option FloatNaNToFloat.
	//
	// FloatNaNConvertToFloat 表示当 marshal 遇到 NaN 时，将值置为另一个数。搭配 FloatNaNToFloat 选项使用。
	FloatNaNConvertToFloat FloatNaNHandleType = 1
	// FloatNaNNull indicates that NaN key-value pair will be set as null when marshaling.
	//
	// FloatNaNNull 表示当 marshal 遇到 NaN 时，则将值设置为 null
	FloatNaNNull FloatNaNHandleType = 2
	// FloatNaNConvertToString indicates that NaN will be replaced as a string when marshaling. This option
	// works with option FloatNaNToString.
	//
	// FloatNaNConvertToString 表示当 marshal 遇到 NaN 时，将值设置为一个字符串。搭配 FloatNaNToString 选项使用。
	FloatNaNConvertToString FloatNaNHandleType = 3
)

type FloatInfHandleType uint8

const (
	// FloatInfTreatAsError indicates that error will be returned when a float number is Inf or -Inf when marshaling.
	//
	// FloatInfTreatAsError 表示当 marshal 遇到 Inf 或 -Inf 时，返回错误。这是默认选项。
	FloatInfTreatAsError FloatInfHandleType = 0
	// FloatInfConvertToFloat indicates that Inf and -Inf will be replaced as another float number when marshaling.
	// This option works with option FloatInfToFloat.
	//
	// FloatInfConvertToFloat 表示当 marshal 遇到 Inf 或 -Inf 时，将值置为另一个数。搭配 FloatInfToFloat 选项使用。
	FloatInfConvertToFloat FloatInfHandleType = 1
	// FloatInfNull indicates that Inf or -Inf key-value pair will be set as null when marshaling.
	//
	// FloatInfNull 表示当 marshal 遇到 Inf 和 -Inf 时，则将值设置为 null
	FloatInfNull FloatInfHandleType = 2
	// FloatInfConvertToString indicates that Inf anf -Inf will be replaced as a string when marshaling. This option
	// works with option FloatInfPositiveToString and FloatInfNegativeToString.
	//
	// FloatInfConvertToString 表示当 marshal 遇到 Inf 和 -Inf 时，将值设置为一个字符串。搭配 FloatInfPositiveToString
	// FloatInfNegativeToString 选项使用。
	FloatInfConvertToString FloatInfHandleType = 3
)

var defaultOption = Opt{
	OmitNull: false,
}

// MustMarshal is the same as Marshal. If error pccurred, an empty byte slice will be returned.
//
// MustMarshal 与 Marshal 相同，但是当错误发生时，什么都不做，直接返回空数据
func (v *V) MustMarshal(opt ...Opt) []byte {
	ret, err := v.Marshal(opt...)
	if err != nil {
		return []byte{}
	}
	return ret
}

// MustMarshalString is the same as MarshalString, If error pccurred, an empty byte slice will be returned.
//
// MustMarshalString 与 MarshalString 相同，但是当错误发生时，什么都不做，直接返回空数据
func (v *V) MustMarshalString(opt ...Opt) string {
	ret, err := v.MarshalString(opt...)
	if err != nil {
		return ""
	}
	return ret
}

// Marshal returns marshaled bytes.
//
// Marshal 返回序列化后的 JSON 字节序列。
func (v *V) Marshal(opt ...Opt) (b []byte, err error) {
	if NotExist == v.valueType {
		return []byte{}, ErrValueUninitialized
	}

	buf := bytes.Buffer{}

	if len(opt) == 0 {
		err = v.marshalToBuffer(nil, &buf, &defaultOption)
	} else {
		err = v.marshalToBuffer(nil, &buf, &opt[0])
	}

	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// MarshalString is same with Marshal, but returns string. It is much more efficient than string(b).
//
// MarshalString 与 Marshal 相同, 不同的是返回 string 类型。它比 string(b) 操作更高效。
func (v *V) MarshalString(opt ...Opt) (s string, err error) {
	if NotExist == v.valueType {
		return "", ErrValueUninitialized
	}

	buf := bytes.Buffer{}

	if len(opt) == 0 {
		err = v.marshalToBuffer(nil, &buf, &defaultOption)
	} else {
		err = v.marshalToBuffer(nil, &buf, &opt[0])
	}

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (v *V) marshalToBuffer(parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	switch v.valueType {
	default:
		// do nothing
	case String:
		v.marshalString(buf)
	case Boolean:
		v.marshalBoolean(buf)
	case Number:
		err = v.marshalNumber(buf, opt)
	case Null:
		v.marshalNull(buf)
	case Object:
		v.marshalObject(parentInfo, buf, opt)
	case Array:
		v.marshalArray(parentInfo, buf, opt)
	}
	return err
}

func (v *V) marshalString(buf *bytes.Buffer) {
	buf.WriteByte('"')
	escapeStringToBuff(v.valueStr, buf)
	buf.WriteByte('"')
}

func (v *V) marshalBoolean(buf *bytes.Buffer) {
	buf.WriteString(formatBool(v.valueBool))
}

func (v *V) marshalNumber(buf *bytes.Buffer, opt *Opt) error {
	if b := v.valueBytes(); len(b) > 0 {
		buf.Write(v.valueBytes())
		return nil
	}
	// else, +Inf or -Inf or NaN
	if math.IsInf(v.num.f64, 1) { // +Inf
		return marshalInfP(buf, opt)
	}
	if math.IsInf(v.num.f64, -1) { // -Inf
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
			escapeStringToBuff(s, buf)
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
			escapeStringToBuff(s, buf)
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
			escapeStringToBuff(s, buf)
		} else if opt.FloatInfPositiveToString != "" {
			s = "-" + strings.TrimLeft(opt.FloatInfPositiveToString, "+")
			escapeStringToBuff(s, buf)
		} else {
			buf.WriteString(`-Inf`)
		}
		buf.WriteByte('"')
	}

	return nil
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

func (v *V) marshalNull(buf *bytes.Buffer) {
	buf.WriteString("null")
}

func (v *V) marshalObject(parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) {
	if opt.MarshalLessFunc != nil {
		sov := v.newSortObjectV(parentInfo, opt)
		sov.marshalObjectWithLessFunc(buf, opt)
		return
	}
	if len(opt.MarshalKeySequence) > 0 {
		sssv := v.newSortStringSliceV(opt)
		sssv.marshalObjectWithStringSlice(buf, opt)
		return
	}

	buf.WriteByte('{')
	defer buf.WriteByte('}')

	i := 0

	for k, child := range v.children.object {
		if child.IsNull() && opt.OmitNull {
			continue
		}
		if i > 0 {
			buf.WriteByte(',')
		}

		buf.WriteByte('"')
		escapeStringToBuff(k, buf)
		buf.WriteString("\":")

		child.marshalToBuffer(nil, buf, opt)
		i++
	}
}

func (v *V) marshalArray(parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) {
	buf.WriteByte('[')
	defer buf.WriteByte(']')

	v.RangeArray(func(i int, child *V) bool {
		if i > 0 {
			buf.WriteByte(',')
		}
		if opt.MarshalLessFunc == nil {
			child.marshalToBuffer(nil, buf, opt)
		} else {
			child.marshalToBuffer(v.newParentInfo(parentInfo, intKey(i)), buf, opt)
		}
		return true
	})
}
