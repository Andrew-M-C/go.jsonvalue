package jsonvalue

import (
	"encoding/base64"
)

// ======== internal global variable ========

var internal = g{}

type g struct {
	b64 *base64.Encoding

	defaultMarshalOption *Opt

	predict struct {
		bytesPerValue int
		calcStorage   uint64 // upper 32 bits - size; lower 32 bits - value count
	}
}

func init() {
	internal.b64 = base64.StdEncoding
	internal.defaultMarshalOption = emptyOptions()
	internal.predict.calcStorage = (10 << 32) + 1
}

func (internal *g) loadPredictSizePerValue() int {
	if n := internal.predict.bytesPerValue; n > 0 {
		return int(n)
	}

	v := internal.predict.calcStorage
	total := int(v >> 32)
	num := int(v & 0xFFFFFFFF)
	return total / num
}

func (internal *g) addPredictSizePerValue(total, num int) {
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

// ======== value pool ========

type valuePool interface {
	get() *V
}

type globalValuePool struct{}

func (globalValuePool) get() *V {
	return &V{}
}

type poolImpl struct {
	pool []V

	count  int
	actual int // actual counted values

	rawSize int
}

func newPool(rawSize int) *poolImpl {
	per := internal.loadPredictSizePerValue()
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
	return globalValuePool{}.get()
}

func (p *poolImpl) release() {
	internal.addPredictSizePerValue(p.rawSize, p.actual)
}

// ======== bytes buffer ========
