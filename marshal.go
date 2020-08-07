package jsonvalue

import (
	"bytes"

	"github.com/buger/jsonparser"
)

// Opt is the option of jsonvalue.
type Opt struct {
	// OmitNull tells how to handle null json value. The default value is false.
	// If OmitNull is true, null value will be omitted when marshaling.
	OmitNull bool

	// MarshalLessFunc is used to handle sequences of marshaling. Since object is
	// implemented by hash map, the sequence of keys is unexpectable. For situations
	// those need settled JSON key-value sequence, please use MarshalLessFunc.
	//
	// Note: Elements in an array value would NOT trigger this function as they are
	// already sorted.
	//
	// We provides a example DefaultStringSequence.
	MarshalLessFunc MarshalLessFunc

	// MarshalKeySequence is used to handle sequance of marshaling. This is much simpler
	// than MarshalLessFunc, just pass a string slice identifying key sequence. For keys
	// those are not in this slice, they would be appended in the end according to result
	// of Go string comparing.
	MarshalKeySequence []string
	keySequence        map[string]int // generated from MarshalKeySequence
}

var defaultOption = Opt{
	OmitNull: false,
}

// MustMarshal is the same as Marshal, but panics if error pccurred
func (v *V) MustMarshal(opt ...Opt) []byte {
	ret, err := v.Marshal(opt...)
	if err != nil {
		panic(err)
	}
	return ret
}

// MustMarshalString is the same as MarshalString, but panics if error pccurred
func (v *V) MustMarshalString(opt ...Opt) string {
	ret, err := v.MarshalString(opt...)
	if err != nil {
		panic(err)
	}
	return ret
}

// Marshal returns marshaled bytes
func (v *V) Marshal(opt ...Opt) (b []byte, err error) {
	if jsonparser.NotExist == v.valueType {
		return nil, ErrValueUninitialized
	}

	buf := bytes.Buffer{}

	if 0 == len(opt) {
		v.marshalToBuffer(nil, &buf, &defaultOption)
	} else {
		v.marshalToBuffer(nil, &buf, &opt[0])
	}

	return buf.Bytes(), err
}

// MarshalString is same with Marshal, but returns string
func (v *V) MarshalString(opt ...Opt) (s string, err error) {
	if jsonparser.NotExist == v.valueType {
		return "", ErrValueUninitialized
	}

	buf := bytes.Buffer{}

	if 0 == len(opt) {
		v.marshalToBuffer(nil, &buf, &defaultOption)
	} else {
		v.marshalToBuffer(nil, &buf, &opt[0])
	}

	return buf.String(), err
}

func (v *V) marshalToBuffer(parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) {
	switch v.valueType {
	default:
		// do nothing
	case jsonparser.String:
		v.marshalString(buf)
	case jsonparser.Boolean:
		v.marshalBoolean(buf)
	case jsonparser.Number:
		v.marshalNumber(buf)
	case jsonparser.Null:
		v.marshalNull(buf)
	case jsonparser.Object:
		v.marshalObject(parentInfo, buf, opt)
	case jsonparser.Array:
		v.marshalArray(parentInfo, buf, opt)
	}
	return
}

func (v *V) marshalString(buf *bytes.Buffer) {
	if v.valueBytes != nil {
		buf.Write(v.valueBytes)
	} else {
		buf.WriteByte('"')
		escapeStringToBuff(v.stringValue, buf)
		buf.WriteByte('"')
	}
	return
}

func (v *V) marshalBoolean(buf *bytes.Buffer) {
	buf.WriteString(formatBool(v.boolValue))
	return
}

func (v *V) marshalNumber(buf *bytes.Buffer) {
	buf.Write(v.valueBytes)
	return
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

	for k, child := range v.objectChildren {
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

	return
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

	return
}
