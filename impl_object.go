package jsonvalue

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type objectValue struct {
	children      map[string]*V
	lowerCaseKeys map[string]map[string]struct{}
}

// NewObject returns an object-typed jsonvalue object. If keyValues is specified, it will also create some key-values in
// the object. Now we supports basic types only. Such as int/uint, int/int8/int16/int32/int64,
// uint/uint8/uint16/uint32/uint64 series, string, bool, nil.
//
// NewObject 返回一个初始化好的 object 类型的 jsonvalue 值。可以使用可选的 map[string]interface{} 类型参数初始化该 object 的下一级键值对，
// 不过目前只支持基础类型，也就是: int/uint, int/int8/int16/int32/int64, uint/uint8/uint16/uint32/uint64, string, bool, nil。
func NewObject(keyValues ...M) *V {
	impl := &objectValue{}

	if len(keyValues) == 0 {
		impl.children = map[string]*V{}
	} else {
		kv := keyValues[0]
		impl.children = make(map[string]*V, len(kv))
		if kv != nil {
			impl.parseNewObjectKV(kv)
		}
	}

	return &V{
		impl: impl,
	}
}

func newObject() (*V, *objectValue) {
	impl := &objectValue{
		children: make(map[string]*V, 128),
	}

	return &V{impl: impl}, impl
}

// M is the alias of map[string]interface{}
type M map[string]interface{}

func (v *objectValue) parseNewObjectKV(kv M) {
	for k, val := range kv {
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Invalid:
			v.children[k] = NewNull()
		case reflect.Bool:
			v.children[k] = NewBool(rv.Bool())
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			v.children[k] = NewInt64(rv.Int())
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			v.children[k] = NewUint64(rv.Uint())
		case reflect.Float32, reflect.Float64:
			v.children[k] = NewFloat64(rv.Float())
		case reflect.String:
			v.children[k] = NewString(rv.String())
		// case reflect.Map:
		// 	if rv.Type().Key().Kind() == reflect.String && rv.Type().Elem().Kind() == reflect.Interface {
		// 		if m, ok := rv.Interface().(M); ok {
		// 			sub := NewObject(m)
		// 			if sub != nil {
		// 				v.Set(sub).At(k)
		// 			}
		// 		}
		// 	}
		default:
			// continue
		}
	}
}

// ======== deleter interface ========

func (v *objectValue) delete(caseless bool, firstParam interface{}, otherParams ...interface{}) error {
	paramCount := len(otherParams)
	if paramCount == 0 {
		return v.deleteInCurrValue(caseless, firstParam)
	}

	child, err := v.get(caseless, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	// if child == nil {
	// 	return ErrNotFound
	// }

	return child.impl.delete(caseless, otherParams[paramCount-1])
}

func (v *objectValue) deleteInCurrValue(caseless bool, param interface{}) error {
	// string expected
	key, err := intfToString(param)
	if err != nil {
		return err
	}

	if exist := v.delFromObjectChildren(caseless, key); !exist {
		return ErrNotFound
	}
	return nil
}

func (v *objectValue) delFromObjectChildren(caseless bool, key string) (exist bool) {
	_, exist = v.children[key]
	if exist {
		delete(v.children, key)
		v.delCaselessKey(key)
		return true
	}

	if !caseless {
		return false
	}

	v.initCaselessStorage()

	lowerKey := strings.ToLower(key)
	keys, exist := v.lowerCaseKeys[lowerKey]
	if !exist {
		return false
	}

	for actualKey := range keys {
		_, exist = v.children[actualKey]
		if exist {
			delete(v.children, actualKey)
			v.delCaselessKey(actualKey)
			return true
		}
	}

	return false
}

func (v *objectValue) initCaselessStorage() {
	if v.lowerCaseKeys != nil {
		return
	}
	v.lowerCaseKeys = make(map[string]map[string]struct{}, len(v.children))
	for k := range v.children {
		v.addCaselessKey(k)
	}
}

func (v *objectValue) addCaselessKey(k string) {
	if v.lowerCaseKeys == nil {
		return
	}
	lowerK := strings.ToLower(k)
	keys, exist := v.lowerCaseKeys[lowerK]
	if !exist {
		keys = make(map[string]struct{})
		v.lowerCaseKeys[lowerK] = keys
	}
	keys[k] = struct{}{}
}

func (v *objectValue) delCaselessKey(k string) {
	if v.lowerCaseKeys == nil {
		return
	}
	lowerK := strings.ToLower(k)
	keys, exist := v.lowerCaseKeys[lowerK]
	if !exist {
		return
	}

	delete(keys, k)

	if len(keys) == 0 {
		delete(v.lowerCaseKeys, lowerK)
	}
}

// ======== typper interface ========

func (v *objectValue) ValueType() ValueType {
	return Object
}

// ======== getter interface ========

func (v *objectValue) get(caseless bool, firstParam interface{}, otherParams ...interface{}) (*V, error) {
	child, err := v.getInCurrValue(caseless, firstParam)
	if err != nil {
		return &V{}, err
	}

	if len(otherParams) == 0 {
		return child, nil
	}
	return child.impl.get(caseless, otherParams[0], otherParams[1:]...)
}

func (v *objectValue) getFromObjectChildren(caseless bool, key string) (child *V, exist bool) {
	child, exist = v.children[key]
	if exist {
		return child, true
	}

	if !caseless {
		return &V{}, false
	}

	v.initCaselessStorage()

	lowerCaseKey := strings.ToLower(key)
	keys, exist := v.lowerCaseKeys[lowerCaseKey]
	if !exist {
		return &V{}, false
	}

	for actualKey := range keys {
		child, exist = v.children[actualKey]
		if exist {
			return child, true
		}
	}

	return &V{}, false
}

func (v *objectValue) getInCurrValue(caseless bool, param interface{}) (*V, error) {
	// string expected
	key, err := intfToString(param)
	if err != nil {
		return &V{}, err
	}
	child, exist := v.getFromObjectChildren(caseless, key)
	if !exist {
		return &V{}, ErrNotFound
	}
	return child, nil
}

// ======== setter interface ========

func (v *objectValue) setAt(end *V, firstParam interface{}, otherParams ...interface{}) error {
	// this is the last iteration
	if len(otherParams) == 0 {
		var k string
		k, err := intfToString(firstParam)
		if err != nil {
			return err
		}
		v.setToObjectChildren(k, end)
		return nil
	}

	// this is not the last iterarion
	k, err := intfToString(firstParam)
	if err != nil {
		return err
	}
	child, exist := v.getFromObjectChildren(false, k)
	if !exist {
		if _, err := intfToString(otherParams[0]); err == nil {
			child = NewObject()
		} else if i, err := intfToInt(otherParams[0]); err == nil {
			if i != 0 {
				return ErrOutOfRange
			}
			child = NewArray()
		} else {
			return fmt.Errorf("unexpected type %v for Set()", reflect.TypeOf(otherParams[0]))
		}
	}
	next := Set{
		v: child,
		c: end,
	}
	_, err = next.At(otherParams[0], otherParams[1:]...)
	if err != nil {
		return err
	}
	if !exist {
		v.setToObjectChildren(k, child)
	}
	return nil
}

func (v *objectValue) setToObjectChildren(key string, child *V) {
	v.children[key] = child
	v.addCaselessKey(key)
}

// ======== iterater interface ========

func (v *objectValue) RangeObjects(callback func(k string, v *V) bool) {
	if nil == callback {
		return
	}

	for k, v := range v.children {
		ok := callback(k, v)
		if !ok {
			break
		}
	}
}

func (v *objectValue) RangeArray(callback func(i int, v *V) bool) {
	// do nothing
}

func (v *objectValue) ForRangeObj() map[string]*V {
	res := make(map[string]*V, len(v.children))
	for k, v := range v.children {
		res[k] = v
	}
	return res
}

func (v *objectValue) ForRangeArr() []*V {
	return nil
}

func (v *objectValue) IterObjects() <-chan *ObjectIter {
	c := make(chan *ObjectIter, len(v.children))

	go func() {
		for k, v := range v.children {
			c <- &ObjectIter{
				K: k,
				V: v,
			}
		}
		close(c)
	}()
	return c
}

func (v *objectValue) IterArray() <-chan *ArrayIter {
	ch := make(chan *ArrayIter)
	close(ch)
	return ch
}

//  ======== marshaler interface ========

func (v *objectValue) marshalToBuffer(curr *V, parentInfo *ParentInfo, buf *bytes.Buffer, opt *Opt) (err error) {
	if opt.MarshalLessFunc != nil {
		sov := v.newSortObjectV(parentInfo, opt)
		sov.marshalObjectWithLessFunc(buf, opt)
		return nil
	}
	if len(opt.MarshalKeySequence) > 0 {
		sssv := v.newSortStringSliceV(curr, opt)
		sssv.marshalObjectWithStringSlice(buf, opt)
		return nil
	}

	buf.WriteByte('{')

	i := 0
	for k, child := range v.children {
		if child.impl.ValueType() == Null && opt.OmitNull {
			continue
		}
		if i > 0 {
			buf.WriteByte(',')
		}

		buf.WriteByte('"')
		escapeStringToBuff(k, buf, opt)
		buf.WriteString("\":")

		child.impl.marshalToBuffer(child, nil, buf, opt)
		i++
	}

	buf.WriteByte('}')
	return nil
}

func (v *objectValue) newSortObjectV(parentInfo *ParentInfo, opt *Opt) *sortObjectV {
	sov := sortObjectV{
		parentInfo: parentInfo,
		lessFunc:   opt.MarshalLessFunc,
		keys:       make([]string, 0, len(v.children)),
		values:     make([]*V, 0, len(v.children)),
	}
	for k, child := range v.children {
		sov.keys = append(sov.keys, k)
		sov.values = append(sov.values, child)
	}

	return &sov
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
		escapeStringToBuff(key, buf, opt)
		buf.WriteString("\":")

		child.impl.marshalToBuffer(child, child.newParentInfo(sov.parentInfo, stringKey(key)), buf, opt)
		marshaledCount++
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

func (v *objectValue) newSortStringSliceV(curr *V, opt *Opt) *sortStringSliceV {
	if nil == opt.keySequence {
		opt.keySequence = make(map[string]int, len(opt.MarshalKeySequence))
		for i, str := range opt.MarshalKeySequence {
			opt.keySequence[str] = i
		}
	}

	sssv := sortStringSliceV{
		v:      curr,
		seq:    opt.keySequence,
		keys:   make([]string, 0, v.Len()),
		values: make([]*V, 0, v.Len()),
	}
	for k, child := range v.children {
		sssv.keys = append(sssv.keys, k)
		sssv.values = append(sssv.values, child)
	}

	return &sssv
}

type sortStringSliceV struct {
	v      *V
	seq    map[string]int
	keys   []string
	values []*V
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
		escapeStringToBuff(key, buf, opt)
		buf.WriteString("\":")

		child.impl.marshalToBuffer(sssv.v, nil, buf, opt)
		marshaledCount++
	}
}

// ======== valuer interface ========

func (v *objectValue) Bool() (bool, error) {
	return false, ErrTypeNotMatch
}

func (v *objectValue) Int64() (int64, error) {
	return 0, ErrTypeNotMatch
}

func (v *objectValue) Uint64() (uint64, error) {
	return 0, ErrTypeNotMatch
}

func (v *objectValue) Float64() (float64, error) {
	return 0, ErrTypeNotMatch
}

func (v *objectValue) String() string {
	buff := bytes.Buffer{}
	buff.WriteByte('{')

	firstWritten := false

	for k, child := range v.children {
		if !firstWritten {
			firstWritten = true
		} else {
			buff.WriteByte(' ')
		}
		buff.WriteString(k)
		buff.WriteByte(':')
		buff.WriteString(child.String())
	}

	buff.WriteByte('}')
	return buff.String()
}

func (v *objectValue) Len() int {
	return len(v.children)
}

// ======== inserter interface ========

func (v *objectValue) insertBefore(end *V, firstParam interface{}, otherParams ...interface{}) error {
	// this is the last iteration
	paramCount := len(otherParams)
	if paramCount == 0 {
		return ErrNotArrayValue
	}

	// this is not the last iterarion
	child, err := v.get(false, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	return child.impl.insertBefore(end, otherParams[paramCount-1])
}

func (v *objectValue) insertAfter(end *V, firstParam interface{}, otherParams ...interface{}) error {
	// this is the last iteration
	paramCount := len(otherParams)
	if paramCount == 0 {
		return ErrNotArrayValue
	}

	// this is not the last iterarion
	child, err := v.get(false, firstParam, otherParams[:paramCount-1]...)
	if err != nil {
		return err
	}
	return child.impl.insertAfter(end, otherParams[paramCount-1])
}

// ======= appender interface ========

func (v *objectValue) appendInTheBeginning(end *V, params ...interface{}) error {
	// this is the last iteration
	paramCount := len(params)
	if paramCount == 0 {
		return ErrNotArrayValue
	}

	// this is not the last iterarion
	child, err := v.get(false, params[0], params[1:]...)
	if err != nil {
		return err
	}
	return child.impl.appendInTheBeginning(end)
}

func (v *objectValue) appendInTheEnd(end *V, params ...interface{}) error {
	// this is the last iteration
	paramCount := len(params)
	if paramCount == 0 {
		return ErrNotArrayValue
	}

	// this is not the last iterarion
	child, err := v.get(false, params[0], params[1:]...)
	if err != nil {
		return err
	}
	return child.impl.appendInTheEnd(end)
}
