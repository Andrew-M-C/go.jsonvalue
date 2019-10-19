package jsonvalue

import (
	"bytes"

	"github.com/buger/jsonparser"
)

// Opt is the option of jsonvalue.
type Opt struct {
	// OmitNull tells how to handle null json value. The default value is false.
	// If OmitNull is true, null value will be omitted whan marshaling.
	OmitNull bool
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
		v.marshalToBuffer(&buf, &defaultOption)
	} else {
		v.marshalToBuffer(&buf, &opt[0])
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
		v.marshalToBuffer(&buf, &defaultOption)
	} else {
		v.marshalToBuffer(&buf, &opt[0])
	}

	return buf.String(), err
}

func (v *V) marshalToBuffer(buf *bytes.Buffer, opt *Opt) {
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
		v.marshalObject(buf, opt)
	case jsonparser.Array:
		v.marshalArray(buf, opt)
	}
	return
}

func (v *V) marshalString(buf *bytes.Buffer) {
	buf.WriteRune('"')
	escapeStringToBuff(v.stringValue, buf)
	buf.WriteRune('"')
	return
}

func (v *V) marshalBoolean(buf *bytes.Buffer) {
	buf.WriteString(formatBool(v.boolValue))
	return
}

func (v *V) marshalNumber(buf *bytes.Buffer) {
	buf.Write(v.rawNumBytes)
	return
}

func (v *V) marshalNull(buf *bytes.Buffer) {
	buf.WriteString("null")
}

func (v *V) marshalObject(buf *bytes.Buffer, opt *Opt) {
	buf.WriteRune('{')
	defer buf.WriteRune('}')

	i := 0

	for k, child := range v.objectChildren {
		if child.IsNull() && opt.OmitNull {
			continue
		}
		if i > 0 {
			buf.WriteRune(',')
		}

		buf.WriteRune('"')
		escapeStringToBuff(k, buf)
		buf.WriteString("\":")

		child.marshalToBuffer(buf, opt)
		i++
	}

	return
}

func (v *V) marshalArray(buf *bytes.Buffer, opt *Opt) {
	buf.WriteRune('[')
	defer buf.WriteRune(']')

	v.RangeArray(func(i int, child *V) bool {
		if i > 0 {
			buf.WriteRune(',')
		}
		child.marshalToBuffer(buf, opt)
		return true
	})

	return
}
