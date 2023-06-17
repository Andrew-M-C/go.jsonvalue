package jsonvalue

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/Andrew-M-C/go.jsonvalue/internal/buffer"
	"github.com/Andrew-M-C/go.jsonvalue/internal/unsafe"
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
		return nil, ErrValueUninitialized
	}

	buf := buffer.NewBuffer()
	opt := combineOptions(opts)

	err = v.marshalToBuffer(nil, buf, opt)
	if err != nil {
		return nil, err
	}

	b = buf.Bytes()
	return b, nil
}

// MarshalString is same with Marshal, but returns string. It is much more efficient than string(b).
//
// MarshalString 与 Marshal 相同, 不同的是返回 string 类型。它比 string(b) 操作更高效。
func (v *V) MarshalString(opts ...Option) (s string, err error) {
	b, err := v.Marshal(opts...)
	if err != nil {
		return "", err
	}
	return unsafe.BtoS(b), nil
}

func (v *V) marshalToBuffer(parentInfo *ParentInfo, buf buffer.Buffer, opt *Opt) (err error) {
	switch v.valueType {
	default:
		// do nothing
	case String:
		v.marshalString(buf, opt)
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

func (v *V) marshalString(buf buffer.Buffer, opt *Opt) {
	_ = buf.WriteByte('"')
	escapeStringToBuff(v.valueStr, buf, opt)
	_ = buf.WriteByte('"')
}

func (v *V) marshalBoolean(buf buffer.Buffer) {
	_, _ = buf.WriteString(formatBool(v.valueBool))
}

func (v *V) marshalNumber(buf buffer.Buffer, opt *Opt) error {
	if b := v.srcByte; len(b) > 0 {
		_, _ = buf.Write(b)
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

func marshalNaN(buf buffer.Buffer, opt *Opt) error {
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
		_, _ = buf.Write(b)

	case FloatNaNNull:
		_, _ = buf.WriteString("null")

	case FloatNaNConvertToString:
		if s := opt.FloatNaNToString; s == "" {
			_, _ = buf.WriteString(`"NaN"`)
		} else {
			_ = buf.WriteByte('"')
			escapeStringToBuff(s, buf, opt)
			_ = buf.WriteByte('"')
		}
	}

	return nil
}

func marshalInfP(buf buffer.Buffer, opt *Opt) error {
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
		_, _ = buf.Write(b)

	case FloatInfNull:
		_, _ = buf.WriteString("null")

	case FloatInfConvertToString:
		if s := opt.FloatInfPositiveToString; s == "" {
			_, _ = buf.WriteString(`"+Inf"`)
		} else {
			_ = buf.WriteByte('"')
			escapeStringToBuff(s, buf, opt)
			_ = buf.WriteByte('"')
		}
	}

	return nil
}

func marshalInfN(buf buffer.Buffer, opt *Opt) error {
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
		_, _ = buf.Write(b)

	case FloatInfNull:
		_, _ = buf.WriteString("null")

	case FloatInfConvertToString:
		_ = buf.WriteByte('"')
		if s := opt.FloatInfNegativeToString; s != "" {
			escapeStringToBuff(s, buf, opt)
		} else if opt.FloatInfPositiveToString != "" {
			s = "-" + strings.TrimLeft(opt.FloatInfPositiveToString, "+")
			escapeStringToBuff(s, buf, opt)
		} else {
			_, _ = buf.WriteString(`-Inf`)
		}
		_ = buf.WriteByte('"')
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

func (v *V) marshalNull(buf buffer.Buffer) {
	_, _ = buf.WriteString("null")
}

func (v *V) marshalObject(parentInfo *ParentInfo, buf buffer.Buffer, opt *Opt) {
	if len(v.children.object) == 0 {
		_, _ = buf.WriteString("{}")
		return
	}

	opt.indent.cnt++
	_ = buf.WriteByte('{')

	if opt.MarshalLessFunc != nil {
		sov := v.newSortObjectV(parentInfo, opt)
		sov.marshalObjectWithLessFunc(buf, opt)
	} else if len(opt.MarshalKeySequence) > 0 {
		sssv := v.newSortStringSliceV(opt)
		sssv.marshalObjectWithStringSlice(buf, opt)
	} else if opt.marshalBySetSequence {
		sssv := v.newSortStringSliceVBySetSeq(opt)
		sssv.marshalObjectWithStringSlice(buf, opt)
	} else {
		firstWritten := false
		for k, child := range v.children.object {
			firstWritten = writeObjectChildren(nil, buf, !firstWritten, k, child.v, opt)
		}
	}

	opt.indent.cnt--
	if opt.indent.enabled {
		_ = buf.WriteByte('\n')
		writeIndent(buf, opt)
	}
	_ = buf.WriteByte('}')
}

func writeObjectChildren(
	parentInfo *ParentInfo, buf buffer.Buffer, isFirstOne bool, key string, child *V, opt *Opt,
) (written bool) {
	if child.IsNull() && opt.OmitNull {
		return false
	}
	if !isFirstOne {
		_ = buf.WriteByte(',')
	}

	if opt.indent.enabled {
		_ = buf.WriteByte('\n')
		writeIndent(buf, opt)
	}

	_ = buf.WriteByte('"')
	escapeStringToBuff(key, buf, opt)

	if opt.indent.enabled {
		_, _ = buf.WriteString("\": ")
	} else {
		_, _ = buf.WriteString("\":")
	}

	_ = child.marshalToBuffer(parentInfo, buf, opt)
	return true
}

func writeIndent(buf buffer.Buffer, opt *Opt) {
	_, _ = buf.WriteString(opt.indent.prefix)
	for i := 0; i < opt.indent.cnt; i++ {
		_, _ = buf.WriteString(opt.indent.indent)
	}
}

func (v *V) marshalArray(parentInfo *ParentInfo, buf buffer.Buffer, opt *Opt) {
	if len(v.children.arr) == 0 {
		_, _ = buf.WriteString("[]")
		return
	}

	opt.indent.cnt++
	_ = buf.WriteByte('[')

	v.RangeArray(func(i int, child *V) bool {
		if i > 0 {
			_ = buf.WriteByte(',')
		}
		if opt.indent.enabled {
			_ = buf.WriteByte('\n')
			writeIndent(buf, opt)
		}
		if opt.MarshalLessFunc == nil {
			_ = child.marshalToBuffer(nil, buf, opt)
		} else {
			_ = child.marshalToBuffer(v.newParentInfo(parentInfo, intKey(i)), buf, opt)
		}
		return true
	})

	opt.indent.cnt--
	if opt.indent.enabled {
		_ = buf.WriteByte('\n')
		writeIndent(buf, opt)
	}
	_ = buf.WriteByte(']')
}
