package jsonvalue

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

// ---------------- array sorting ----------------

// ArrayLessFunc is used in SortArray(), identifying which member is ahead.
//
// ArrayLessFunc 用于 SortArray() 函数中，指定两个成员谁在前面。
type ArrayLessFunc func(v1, v2 *V) bool

// SortArray is used to re-arrange sequence of the array. Invokers should pass less function for sorting.
// Nothing would happens either lessFunc is nil or v is not an array.
//
// SortArray 用于对 array 类型的 JSON 的子成员进行重新排序。基本逻辑与 sort.Sort 函数相同。当 lessFunc 为 nil，或者当前 JSON 不是一个
// array 类型时，什么变化都不会发生。
func (v *V) SortArray(lessFunc ArrayLessFunc) {
	if v.impl == nil {
		return
	}
	if nil == lessFunc {
		return
	}

	if !v.IsArray() {
		return
	}

	arr := (v.impl).(*arrayValue)
	sort.Slice(arr.children, func(i, j int) bool {
		vI := arr.children[i]
		vJ := arr.children[j]
		return lessFunc(vI, vJ)
	})
}

// ---------------- marshal sorting ----------------

// Key is the element of KeyPath
//
// Key 是 KeyPath 类型的成员
type Key struct {
	s string
	i int
}

func intKey(i int) Key {
	return Key{i: i}
}

func stringKey(s string) Key {
	return Key{s: s}
}

// String returns string value of a key
//
// String 返回当前键值对的键的描述
func (k *Key) String() string {
	if k.s != "" {
		return k.s
	}
	return strconv.Itoa(k.i)
}

// IsString tells if current key is a string, which indicates a child of an object.
//
// IsString 判断当前的键是不是一个 string 类型，如果是的话，那么它是一个 object JSON 的子成员。
func (k *Key) IsString() bool {
	return k.s != ""
}

// Int returns int value of a key.
//
// Int 返回当前键值对的 int 值。
func (k *Key) Int() int {
	if k.s == "" {
		return k.i
	}
	return 0
}

// IsInt tells if current key is a integer, which indicates a child of an array
//
// IsInt 判断当前的键是不是一个整型类型，如果是的话，那么它是一个 array JSON 的子成员。
func (k *Key) IsInt() bool {
	return k.s == ""
}

// KeyPath identifies a full path of keys of object in jsonvalue.
//
// KeyPath 表示一个对象在指定 jsonvalue 中的完整的键路径。
type KeyPath []*Key

// String returns last element of key path.
//
// String 返回 KeyPath 的最后一个成员的描述。
func (p KeyPath) String() (s string) {
	buff := bytes.Buffer{}
	buff.WriteRune('[')

	defer func() {
		buff.WriteRune(']')
		s = buff.String()
	}()

	for i, k := range p {
		if i > 0 {
			buff.WriteRune(' ')
		}
		if k.IsInt() {
			s := strconv.Itoa(k.Int())
			buff.WriteString(s)
		} else {
			buff.WriteRune('"')
			escapeStringToBuff(k.String(), &buff, &Opt{})
			buff.WriteRune('"')
		}
	}

	return
}

// ParentInfo show informations of parent of a JSON value.
//
// ParentInfo 表示一个 JSON 值的父节点信息。
type ParentInfo struct {
	Parent  *V
	KeyPath KeyPath
}

func (v *V) newParentInfo(nilableParentInfo *ParentInfo, key Key) *ParentInfo {
	if nil == nilableParentInfo {
		return &ParentInfo{
			Parent:  v,
			KeyPath: KeyPath{&key},
		}
	}

	return &ParentInfo{
		Parent:  v,
		KeyPath: append(nilableParentInfo.KeyPath, &key),
	}
}

// MarshalLessFunc is used in marshaling, for sorting marshaled data.
//
// MarshalLessFunc 用于序列化，指定 object 类型的 JSON 的键值对顺序。
type MarshalLessFunc func(nilableParent *ParentInfo, key1, key2 string, v1, v2 *V) bool

// DefaultStringSequence simply use strings.Compare() to define the sequence
// of various key-value pairs of an object value.
// This function is used in Opt.MarshalLessFunc.
//
// DefaultStringSequence 使用 strings.Compare() 函数来判断键值对的顺序。用于 Opt.MarshalLessFunc。
func DefaultStringSequence(parent *ParentInfo, key1, key2 string, v1, v2 *V) bool {
	return strings.Compare(key1, key2) <= 0
}
