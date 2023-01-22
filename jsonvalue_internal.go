package jsonvalue

import (
	"encoding/base64"
	"sync"
)

var internal = struct {
	b64 *base64.Encoding

	defaultMarshalOption *Opt

	pool *sync.Pool

	predict struct {
		lenPerValue int
		calcStorage uint64
	}
}{}

func init() {
	internal.b64 = base64.StdEncoding
	internal.defaultMarshalOption = emptyOptions()

	internal.pool = &sync.Pool{
		New: func() any {
			return &V{}
		},
	}

	internal.predict.calcStorage = (10 << 32) + 1
}

func internalLoadPredictSizePerValue() int {
	if n := internal.predict.lenPerValue; n > 0 {
		return int(n)
	}

	v := internal.predict.calcStorage
	total := int(v >> 32)
	num := int(v & 0xFFFFFFFF)
	return total / num
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
	internal.predict.lenPerValue = int(per)
	internal.predict.calcStorage = (per << 32) + 1
}

type pool interface {
	get() *V
}

type globalPool struct{}

func (globalPool) get() *V {
	v, _ := internal.pool.Get().(*V)

	if v.valueType == Array && v.children.arr != nil {
		v.children.arr = v.children.arr[:0]

	} else if v.valueType == Object {
		for k := range v.children.object {
			delete(v.children.object, k)
		}
	}

	v.valueType = NotExist
	v.srcByte = nil
	v.valueStr = ""
	return v
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
	for i := p.actual; i < len(p.pool); i++ {
		v := &p.pool[i]
		internal.pool.Put(v)
	}
}
