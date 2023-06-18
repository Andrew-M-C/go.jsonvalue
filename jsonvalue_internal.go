package jsonvalue

import (
	"encoding"
	"encoding/base64"
	"encoding/json"
	"reflect"
)

var internal = struct {
	b64 *base64.Encoding

	defaultMarshalOption *Opt

	predict struct {
		bytesPerValue int
		calcStorage   uint64 // upper 32 bits - size; lower 32 bits - value count
	}

	types struct {
		JSONMarshaler reflect.Type
		TextMarshaler reflect.Type
	}
}{}

func init() {
	internal.b64 = base64.StdEncoding
	internal.defaultMarshalOption = emptyOptions()
	internalAddPredictSizePerValue(16, 1)

	internal.types.JSONMarshaler = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	internal.types.TextMarshaler = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
}

func internalLoadPredictSizePerValue() int {
	if n := internal.predict.bytesPerValue; n > 0 {
		return int(n)
	}

	v := internal.predict.calcStorage
	total := v >> 32
	num := v & 0xFFFFFFFF
	return int(total / num)
}

func internalAddPredictSizePerValue(total, num int) {
	v := internal.predict.calcStorage
	preTotal := v >> 32
	preNum := v & 0xFFFFFFFF

	nextTotal := uint64(total) + preTotal
	nextNum := uint64(num) + preNum

	if nextTotal < 0x7FFFFFFF {
		v := (nextTotal << 32) + nextNum
		internal.predict.calcStorage = v
		return
	}

	per := nextTotal / nextNum
	internal.predict.bytesPerValue = int(per)
	internal.predict.calcStorage = (per << 32) + 1
}

type pool interface {
	get() *V
}

type globalPool struct{}

func (globalPool) get() *V {
	return &V{}
}

type poolImpl struct {
	pool []V

	count  int
	actual int // actual counted values

	rawSize int
}

func newPool(rawSize int) *poolImpl {
	per := internalLoadPredictSizePerValue()
	cnt := rawSize / per

	p := &poolImpl{
		pool:    make([]V, cnt),
		count:   cnt,
		actual:  0,
		rawSize: rawSize,
	}

	return p
}

func (p *poolImpl) get() *V {
	if p.actual < p.count {
		v := &p.pool[p.actual]
		p.actual++
		return v
	}

	p.actual++
	return globalPool{}.get()
}

func (p *poolImpl) release() {
	internalAddPredictSizePerValue(p.rawSize, p.actual)
}
