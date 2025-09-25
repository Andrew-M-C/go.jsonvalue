package jsonvalue

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/Andrew-M-C/go.jsonvalue/internal/buffer"
	"github.com/Andrew-M-C/go.jsonvalue/internal/unsafe"
)

// MustMarshal is the same as Marshal. If error occurs, an empty byte slice will be returned.
//
// MustMarshal 与 Marshal 相同，但是当错误发生时，什么都不做，直接返回空数据
func (v *V) MustMarshal(opts ...Option) []byte {
	ret, err := v.Marshal(opts...)
	if err != nil {
		return []byte{}
	}
	return ret
}

// MustMarshalString is the same as MarshalString, If error occurs, an empty byte slice will be returned.
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
	if v == nil || v.valueType == NotExist {
		return nil, ErrValueUninitialized
	}

	buf := buffer.NewBuffer()
	opt := combineOptions(opts)

	err = marshalToBuffer(v, nil, buf, opt)
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

// MarshalWrite is the same as Marshal, but writes to a io.Writer.
//
// MarshalWrite 与 Marshal 相同，但是将序列化后的 JSON 写入到 io.Writer 中。
func (v *V) MarshalWrite(w io.Writer, opts ...Option) error {
	if v == nil || v.valueType == NotExist {
		return ErrValueUninitialized
	}
	if w == nil {
		return ErrNilParameter
	}
	opt := combineOptions(opts)
	return marshalToBuffer(v, nil, w, opt)
}

func marshalToBuffer(v *V, parentInfo *ParentInfo, buf io.Writer, opt *Opt) (err error) {
	switch v.valueType {
	default:
		// do nothing
	case String:
		marshalString(v, buf, opt)
	case Boolean:
		marshalBoolean(v, buf)
	case Number:
		err = marshalNumber(v, buf, opt)
	case Null:
		marshalNull(buf)
	case Object:
		marshalObject(v, parentInfo, buf, opt)
	case Array:
		marshalArray(v, parentInfo, buf, opt)
	}
	return err
}

func marshalString(v *V, buf io.Writer, opt *Opt) {
	_, _ = buf.Write([]byte{'"'})
	escapeStringToBuff(v.valueStr, buf, opt)
	_, _ = buf.Write([]byte{'"'})
}

func marshalBoolean(v *V, buf io.Writer) {
	if v.valueBool {
		_, _ = buf.Write([]byte{'t', 'r', 'u', 'e'})
	} else {
		_, _ = buf.Write([]byte{'f', 'a', 'l', 's', 'e'})
	}
}

func marshalNumber(v *V, buf io.Writer, opt *Opt) error {
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

func marshalNaN(buf io.Writer, opt *Opt) error {
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
		_, _ = buf.Write([]byte{'n', 'u', 'l', 'l'})

	case FloatNaNConvertToString:
		if s := opt.FloatNaNToString; s == "" {
			_, _ = buf.Write([]byte{'"', 'N', 'a', 'N', '"'})
		} else {
			_, _ = buf.Write([]byte{'"'})
			escapeStringToBuff(s, buf, opt)
			_, _ = buf.Write([]byte{'"'})
		}
	}

	return nil
}

func marshalInfP(buf io.Writer, opt *Opt) error {
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
		_, _ = buf.Write([]byte{'n', 'u', 'l', 'l'})

	case FloatInfConvertToString:
		if s := opt.FloatInfPositiveToString; s == "" {
			_, _ = buf.Write([]byte{'"', '+', 'I', 'n', 'f', '"'})
		} else {
			_, _ = buf.Write([]byte{'"'})
			escapeStringToBuff(s, buf, opt)
			_, _ = buf.Write([]byte{'"'})
		}
	}

	return nil
}

func marshalInfN(buf io.Writer, opt *Opt) error {
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
		_, _ = buf.Write([]byte{'n', 'u', 'l', 'l'})

	case FloatInfConvertToString:
		_, _ = buf.Write([]byte{'"'})
		if s := opt.FloatInfNegativeToString; s != "" {
			escapeStringToBuff(s, buf, opt)
		} else if opt.FloatInfPositiveToString != "" {
			s = "-" + strings.TrimLeft(opt.FloatInfPositiveToString, "+")
			escapeStringToBuff(s, buf, opt)
		} else {
			_, _ = buf.Write([]byte{'-', 'I', 'n', 'f'})
		}
		_, _ = buf.Write([]byte{'"'})
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

func marshalNull(buf io.Writer) {
	_, _ = buf.Write([]byte{'n', 'u', 'l', 'l'})
}

func marshalObject(v *V, parentInfo *ParentInfo, buf io.Writer, opt *Opt) {
	if len(v.children.object) == 0 {
		_, _ = buf.Write([]byte{'{', '}'})
		return
	}

	opt.indent.cnt++
	_, _ = buf.Write([]byte{'{'})

	if opt.MarshalLessFunc != nil {
		sov := newSortObjectV(v, parentInfo, opt)
		sov.marshalObjectWithLessFunc(buf, opt)

	} else if len(opt.MarshalKeySequence) > 0 {
		sssv := newSortStringSliceV(v, opt)
		sssv.marshalObjectWithStringSlice(buf, opt)

	} else if opt.marshalBySetSequence {
		sssv := newSortStringSliceVBySetSeq(v)
		sssv.marshalObjectWithStringSlice(buf, opt)

	} else {
		writeObjectKVInRandomizedSequence(v, buf, opt)
	}

	opt.indent.cnt--
	if opt.indent.enabled {
		_, _ = buf.Write([]byte{'\n'})
		writeIndent(buf, opt)
	}
	_, _ = buf.Write([]byte{'}'})
}

func writeObjectKVInRandomizedSequence(v *V, buf io.Writer, opt *Opt) {
	firstWritten := false
	for k, child := range v.children.object {
		firstWritten = writeObjectChildren(nil, buf, !firstWritten, k, child.v, opt)
	}
}

func writeObjectChildren(
	parentInfo *ParentInfo, buf io.Writer, isFirstOne bool, key string, child *V, opt *Opt,
) (written bool) {
	if child.IsNull() && opt.OmitNull {
		return false
	}
	if !isFirstOne {
		_, _ = buf.Write([]byte{','})
	}

	if opt.indent.enabled {
		_, _ = buf.Write([]byte{'\n'})
		writeIndent(buf, opt)
	}

	_, _ = buf.Write([]byte{'"'})
	escapeStringToBuff(key, buf, opt)

	if opt.indent.enabled {
		_, _ = buf.Write([]byte{'"', ':', ' '})
	} else {
		_, _ = buf.Write([]byte{'"', ':'})
	}

	_ = marshalToBuffer(child, parentInfo, buf, opt)
	return true
}

func writeIndent(buf io.Writer, opt *Opt) {
	_, _ = buf.Write(unsafe.StoB(opt.indent.prefix))
	for i := 0; i < opt.indent.cnt; i++ {
		_, _ = buf.Write(unsafe.StoB(opt.indent.indent))
	}
}

func marshalArray(v *V, parentInfo *ParentInfo, buf io.Writer, opt *Opt) {
	if len(v.children.arr) == 0 {
		_, _ = buf.Write([]byte{'[', ']'})
		return
	}

	opt.indent.cnt++
	_, _ = buf.Write([]byte{'['})

	v.RangeArray(func(i int, child *V) bool {
		if i > 0 {
			_, _ = buf.Write([]byte{','})
		}
		if opt.indent.enabled {
			_, _ = buf.Write([]byte{'\n'})
			writeIndent(buf, opt)
		}
		if opt.MarshalLessFunc == nil {
			_ = marshalToBuffer(child, nil, buf, opt)
		} else {
			_ = marshalToBuffer(child, newParentInfo(v, parentInfo, intKey(i)), buf, opt)
		}
		return true
	})

	opt.indent.cnt--
	if opt.indent.enabled {
		_, _ = buf.Write([]byte{'\n'})
		writeIndent(buf, opt)
	}
	_, _ = buf.Write([]byte{']'})
}
