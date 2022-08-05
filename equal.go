package jsonvalue

import (
	"bytes"

	"github.com/shopspring/decimal"
)

// Equal shows whether the content of two JSON values equal to each other.
//
// Equal 判断两个 JSON 的内容是否相等
func (v *V) Equal(another *V) bool {
	if v == nil || another == nil {
		return false
	}
	if v.valueType != another.valueType {
		return false
	}

	switch v.valueType {
	default: // including NotExist, Unknown
		return false
	case String:
		return v.valueStr == another.valueStr
	case Number:
		return numberEqual(v, another)
	case Object:
		return objectEqual(v, another)
	case Array:
		return arrayEqual(v, another)
	case Boolean:
		return v.valueBool == another.valueBool
	case Null:
		return true
	}
}

func numberEqual(left, right *V) bool {
	if bytes.Equal(left.srcByte, right.srcByte) {
		return true
	}

	d1, _ := decimal.NewFromString(string(left.srcByte))
	d2, _ := decimal.NewFromString(string(right.srcByte))
	return d1.Equal(d2)
}

func objectEqual(left, right *V) bool {
	if len(left.children.object) != len(right.children.object) {
		return false
	}

	for k, leftChild := range left.children.object {
		rightChild, exist := right.children.object[k]
		if !exist {
			return false
		}
		if !leftChild.v.Equal(rightChild.v) {
			return false
		}
	}
	return true
}

func arrayEqual(left, right *V) bool {
	if len(left.children.arr) != len(right.children.arr) {
		return false
	}

	for i, leftChild := range left.children.arr {
		rightChild := right.children.arr[i]
		if !leftChild.Equal(rightChild) {
			return false
		}
	}
	return true
}
