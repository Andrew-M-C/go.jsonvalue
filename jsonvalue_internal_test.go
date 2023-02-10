package jsonvalue

import "testing"

func testInternal(t *testing.T) {
	cv("test predict numbers", func() { testInternalPredict(t) })
}

func testInternalPredict(t *testing.T) {
	const maxUint32 = 0xFFFFFFFF
	v := internal.predict.calcStorage
	p := internal.predict.bytesPerValue

	defer func(v uint64, p int) {
		// cancel mocking
		internal.predict.calcStorage = v
		internal.predict.bytesPerValue = p
	}(v, p)

	size := internalLoadPredictSizePerValue()

	t.Logf("total: %d, count %d, calculated per size: %d", v>>32, v&maxUint32, size)
	so(size, eq, (v>>32)/(v&maxUint32))

	var mockTotal uint64 = maxUint32
	var mockCount uint64 = maxUint32 / 16
	internal.predict.calcStorage = (mockTotal << 32) + mockCount
	internal.predict.bytesPerValue = 0

	size = internalLoadPredictSizePerValue()
	so(size, eq, 16)

	// try unmarshal a new jsonvalue, which will make total overflowing uint32
	const raw = `1234567890123456`
	so(len(raw), eq, 16)
	jv, err := UnmarshalString(raw)
	so(err, isNil)
	so(jv.String(), eq, raw)

	size = internalLoadPredictSizePerValue()
	so(size, eq, 16)

	v = internal.predict.calcStorage
	p = internal.predict.bytesPerValue
	t.Logf("total: %d, count %d, calculated per size: %d", v>>32, v&maxUint32, size)
	so(size, eq, (v>>32)/(v&maxUint32))
	so(size, eq, p)

	// TODO:
}
