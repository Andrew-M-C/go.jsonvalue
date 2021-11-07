package jsonvalue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// MustMarshal is the same as Marshal. If error pccurred, an empty byte slice will be returned.
//
// MustMarshal 与 Marshal 相同，但是当错误发生时，什么都不做，直接返回空数据
func (v *V) MustMarshal(opts ...Option) []byte {
	ret, err := v.Marshal(opts...)
	if err != nil {
		return []byte{}
	}
	return ret
}

// MustMarshalString is the same as MarshalString, If error pccurred, an empty byte slice will be returned.
//
// MustMarshalString 与 MarshalString 相同，但是当错误发生时，什么都不做，直接返回空数据
func (v *V) MustMarshalString(opt ...Option) string {
	ret, err := v.MarshalString(opt...)
	if err != nil {
		return ""
	}
	return ret
}

// Marshal returns marshaled bytes.
//
// Marshal 返回序列化后的 JSON 字节序列。
func (v *V) Marshal(opts ...Option) (b []byte, err error) {
	if NotExist == v.valueType {
		return []byte{}, ErrValueUninitialized
	}

	buf := bytes.Buffer{}
	opt := combineOptions(opts)
	err = v.marshalToBuffer(nil, &buf, opt)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// MarshalString is same with Marshal, but returns string. It is much more efficient than string(b).
//
// MarshalString 与 Marshal 相同, 不同的是返回 string 类型。它比 string(b) 操作更高效。
func (v *V) MarshalString(opts ...Option) (s string, err error) {
	if NotExist == v.valueType {
		return "", ErrValueUninitialized
	}

	buf := bytes.Buffer{}
	opt := combineOptions(opts)
	err = v.marshalToBuffer(nil, &buf, opt)
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
