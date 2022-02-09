package jsonvalue

import (
	"bytes"
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
	if v.impl == nil {
		return []byte{}, ErrValueUninitialized
	}

	buf := bytes.Buffer{}
	opt := combineOptions(opts)
	err = v.impl.marshalToBuffer(v, nil, &buf, opt)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// MarshalString is same with Marshal, but returns string. It is much more efficient than string(b).
//
// MarshalString 与 Marshal 相同, 不同的是返回 string 类型。它比 string(b) 操作更高效。
func (v *V) MarshalString(opts ...Option) (s string, err error) {
	if v.impl == nil {
		return "", ErrValueUninitialized
	}

	buf := bytes.Buffer{}
	opt := combineOptions(opts)
	err = v.impl.marshalToBuffer(v, nil, &buf, opt)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
