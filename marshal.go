package jsonvalue

import (
	"bytes"
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
}

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
		v.marshalToBuffer(nil, &buf, &defaultOption)
	} else {
		v.marshalToBuffer(nil, &buf, &opt[0])
	}

	return buf.Bytes(), err
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
	case String:
		v.marshalString(buf)
	case Boolean:
		v.marshalBoolean(buf)
	case Number:
		v.marshalNumber(buf)
	case Null:
		v.marshalNull(buf)
	case Object:
		v.marshalObject(parentInfo, buf, opt)
	case Array:
		v.marshalArray(parentInfo, buf, opt)
	}
}

func (v *V) marshalString(buf *bytes.Buffer) {
	buf.WriteByte('"')
	escapeStringToBuff(v.valueStr, buf)
	buf.WriteByte('"')
}

func (v *V) marshalBoolean(buf *bytes.Buffer) {
	buf.WriteString(formatBool(v.valueBool))
}

func (v *V) marshalNumber(buf *bytes.Buffer) {
	buf.Write(v.valueBytes())
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
