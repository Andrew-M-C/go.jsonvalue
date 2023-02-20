package jsonvalue

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/Andrew-M-C/go.jsonvalue/internal/buffer"
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
	if nil == lessFunc {
		return
	}
	if !v.IsArray() {
		return
	}

	sav := newSortV(v, lessFunc)
	sav.Sort()
}

type sortArrayV struct {
	v        *V
	lessFunc ArrayLessFunc
}

func newSortV(v *V, lessFunc ArrayLessFunc) *sortArrayV {
	sav := sortArrayV{
		v:        v,
		lessFunc: lessFunc,
	}
	return &sav
}

func (v *sortArrayV) Sort() {
	sort.Sort(v)
}

func (v *sortArrayV) Len() int {
	return len(v.v.children.arr)
}

func (v *sortArrayV) Less(i, j int) bool {
	v1 := v.v.children.arr[i]
	v2 := v.v.children.arr[j]
	return v.lessFunc(v1, v2)
}

func (v *sortArrayV) Swap(i, j int) {
	v.v.children.arr[i], v.v.children.arr[j] = v.v.children.arr[j], v.v.children.arr[i]
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
			escapeStringToBuff(k.String(), &buff, getDefaultOptions())
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

func (sov *sortObjectV) marshalObjectWithLessFunc(buf buffer.Buffer, opt *Opt) {
	// sort
	sort.Sort(sov)

	// marshal
	firstWritten := false
	for i, key := range sov.keys {
		child := sov.values[i]
		par := child.newParentInfo(sov.parentInfo, stringKey(key))
		firstWritten = writeObjectChildren(par, buf, !firstWritten, key, child, opt)
	}
}

type sortObjectV struct {
	parentInfo *ParentInfo
	lessFunc   MarshalLessFunc
	keys       []string
	values     []*V
}

func (sov *sortObjectV) Len() int {
	return len(sov.values)
}

func (sov *sortObjectV) Less(i, j int) bool {
	return sov.lessFunc(sov.parentInfo, sov.keys[i], sov.keys[j], sov.values[i], sov.values[j])
}

func (sov *sortObjectV) Swap(i, j int) {
	sov.keys[i], sov.keys[j] = sov.keys[j], sov.keys[i]
	sov.values[i], sov.values[j] = sov.values[j], sov.values[i]
}

func (v *V) newSortObjectV(parentInfo *ParentInfo, opt *Opt) *sortObjectV {
	sov := sortObjectV{
		parentInfo: parentInfo,
		lessFunc:   opt.MarshalLessFunc,
		keys:       make([]string, 0, len(v.children.object)),
		values:     make([]*V, 0, len(v.children.object)),
	}
	for k, child := range v.children.object {
		sov.keys = append(sov.keys, k)
		sov.values = append(sov.values, child.v)
	}

	return &sov
}

// marshalObjectWithStringSlice use a slice to determine sequence of object
func (sssv *sortStringSliceV) marshalObjectWithStringSlice(buf buffer.Buffer, opt *Opt) {
	// sort
	sort.Sort(sssv)

	// marshal
	firstWritten := false
	for i, key := range sssv.keys {
		child := sssv.values[i]
		firstWritten = writeObjectChildren(nil, buf, !firstWritten, key, child, opt)
	}
}

type sortStringSliceV struct {
	v      *V
	seq    map[string]int
	keys   []string
	values []*V
}

func (v *V) newSortStringSliceV(opt *Opt) *sortStringSliceV {
	if nil == opt.keySequence {
		opt.keySequence = make(map[string]int, len(opt.MarshalKeySequence))
		for i, str := range opt.MarshalKeySequence {
			opt.keySequence[str] = i
		}
	}

	sssv := sortStringSliceV{
		v:      v,
		seq:    opt.keySequence,
		keys:   make([]string, 0, v.Len()),
		values: make([]*V, 0, v.Len()),
	}
	for k, child := range v.children.object {
		sssv.keys = append(sssv.keys, k)
		sssv.values = append(sssv.values, child.v)
	}

	return &sssv
}

func (v *V) newSortStringSliceVBySetSeq(opt *Opt) *sortStringSliceV {
	keySequence := make(map[string]int, len(v.children.object))
	for k, child := range v.children.object {
		keySequence[k] = int(child.id)
	}

	sssv := sortStringSliceV{
		v:      v,
		seq:    keySequence,
		keys:   make([]string, 0, v.Len()),
		values: make([]*V, 0, v.Len()),
	}
	for k, child := range v.children.object {
		sssv.keys = append(sssv.keys, k)
		sssv.values = append(sssv.values, child.v)
	}

	return &sssv
}

func (sssv *sortStringSliceV) Len() int {
	return len(sssv.values)
}

func (sssv *sortStringSliceV) Less(i, j int) bool {
	k1 := sssv.keys[i]
	k2 := sssv.keys[j]

	seq1, exist1 := sssv.seq[k1]
	seq2, exist2 := sssv.seq[k2]

	if exist1 {
		if exist2 {
			return seq1 < seq2
		}
		return true
	}
	if exist2 {
		return false
	}

	return k1 <= k2
}

func (sssv *sortStringSliceV) Swap(i, j int) {
	sssv.keys[i], sssv.keys[j] = sssv.keys[j], sssv.keys[i]
	sssv.values[i], sssv.values[j] = sssv.values[j], sssv.values[i]
}
