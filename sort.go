package jsonvalue

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

// ---------------- array sorting ----------------

// ArrayLessFunc is used in SortArray()
type ArrayLessFunc func(v1, v2 *V) bool

// SortArray is used to re-arrange sequence of the array. Invokers should pass less function for sorting.
// Nothing would happens either lessFunc is nil or v is not an array.
func (v *V) SortArray(lessFunc ArrayLessFunc) {
	if nil == lessFunc {
		return
	}
	if false == v.IsArray() {
		return
	}

	sav := newSortV(v, lessFunc)
	sav.Sort()
	return
}

type sortArrayV struct {
	v        *V
	lessFunc ArrayLessFunc
	children []*V
}

func newSortV(v *V, lessFunc ArrayLessFunc) *sortArrayV {
	sav := sortArrayV{
		v:        v,
		lessFunc: lessFunc,
		children: make([]*V, 0, v.Len()),
	}

	for e := v.arrayChildren.Front(); e != nil; e = e.Next() {
		sav.children = append(sav.children, e.Value.(*V))
	}
	return &sav
}

func (v *sortArrayV) Sort() {
	// sort
	sort.Sort(v)

	// re-arrange children
	v.v.arrayChildren.Init()
	for _, child := range v.children {
		v.v.arrayChildren.PushBack(child)
	}
	return
}

func (v *sortArrayV) Len() int {
	return len(v.children)
}

func (v *sortArrayV) Less(i, j int) bool {
	v1 := v.children[i]
	v2 := v.children[j]
	return v.lessFunc(v1, v2)
}

func (v *sortArrayV) Swap(i, j int) {
	v.children[i], v.children[j] = v.children[j], v.children[i]
}

// ---------------- marshal sorting ----------------

// Key is the element of KeyPath
type Key struct {
	v interface{}
}

func intKey(i int) Key {
	return Key{v: i}
}

func stringKey(s string) Key {
	return Key{v: s}
}

// String returns string value of a key
func (k *Key) String() string {
	if s, ok := k.v.(string); ok {
		return s
	}
	if i, ok := k.v.(int); ok {
		return strconv.Itoa(i)
	}
	return ""
}

// IsString tells if current key is a string, which indicates a child of an object
func (k *Key) IsString() bool {
	_, ok := k.v.(string)
	return ok
}

// Int returns int value of a key
func (k *Key) Int() int {
	if i, ok := k.v.(int); ok {
		return i
	}
	return 0
}

// IsInt tells if current key is a integer, which indicates a child of an array
func (k *Key) IsInt() bool {
	_, ok := k.v.(int)
	return ok
}

// KeyPath identifies a full path of keys of object in jsonvalue
type KeyPath []*Key

// String returns last element of key path
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
			escapeStringToBuff(k.String(), &buff)
			buff.WriteRune('"')
		}
	}

	return
}

// ParentInfo show informations of parent of a json value
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

// MarshalLessFunc is used in marshaling, for sorting marshaled data
type MarshalLessFunc func(nilableParent *ParentInfo, key1, key2 string, v1, v2 *V) bool

// DefaultStringSequence simply use strings.Compare() to define the sequence
// of various key-value pairs of an object value.
// This function is used in Opt.MarshalLessFunc.
func DefaultStringSequence(parent *ParentInfo, key1, key2 string, v1, v2 *V) bool {
	return strings.Compare(key1, key2) <= 0
}

func (sov *sortObjectV) marshalObjectWithLessFunc(buf *bytes.Buffer, opt *Opt) {
	buf.WriteRune('{')
	defer buf.WriteRune('}')

	// sort
	sort.Sort(sov)

	// marshal
	marshaledCount := 0
	for i, key := range sov.keys {
		child := sov.values[i]
		if child.IsNull() && opt.OmitNull {
			continue
		}
		if marshaledCount > 0 {
			buf.WriteRune(',')
		}

		buf.WriteRune('"')
		escapeStringToBuff(key, buf)
		buf.WriteString("\":")

		child.marshalToBuffer(child.newParentInfo(sov.parentInfo, stringKey(key)), buf, opt)
		marshaledCount++
	}

	return
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
		keys:       make([]string, 0, len(v.objectChildren)),
		values:     make([]*V, 0, len(v.objectChildren)),
	}
	for k, child := range v.objectChildren {
		sov.keys = append(sov.keys, k)
		sov.values = append(sov.values, child)
	}

	return &sov
}

// marshalObjectWithStringSlice use a slice to determine sequence of object
func (sssv *sortStringSliceV) marshalObjectWithStringSlice(buf *bytes.Buffer, opt *Opt) {
	buf.WriteRune('{')
	defer buf.WriteRune('}')

	// sort
	sort.Sort(sssv)

	// marshal
	marshaledCount := 0
	for i, key := range sssv.keys {
		child := sssv.values[i]
		if child.IsNull() && opt.OmitNull {
			continue
		}
		if marshaledCount > 0 {
			buf.WriteRune(',')
		}

		buf.WriteRune('"')
		escapeStringToBuff(key, buf)
		buf.WriteString("\":")

		child.marshalToBuffer(nil, buf, opt)
		marshaledCount++
	}

	return
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
	for k, child := range v.objectChildren {
		sssv.keys = append(sssv.keys, k)
		sssv.values = append(sssv.values, child)
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