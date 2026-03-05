package jsonvalue

import "testing"

func testInternal(t *testing.T) {
	cv("test predict numbers", func() { testInternalPredict(t) })
	cv("test anyToInt function", func() { testAnyToInt(t) })
	cv("test anyToString function", func() { testAnyToString(t) })
}

func testInternalPredict(t *testing.T) {
	const maxUint32 = 0xFFFFFFFF
	v := internal.predict.calcStorage
	p := internal.predict.bytesPerValue

	defer func(v uint64, p uint64) {
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
}

func testAnyToInt(*testing.T) {
	// Test normal cases (should work)
	result, err := anyToInt(int(42))
	so(err, isNil)
	so(result, eq, 42)

	result, err = anyToInt(int8(5))
	so(err, isNil)
	so(result, eq, 5)

	result, err = anyToInt(uint16(100))
	so(err, isNil)
	so(result, eq, 100)

	// Test nil parameter (should error)
	_, err = anyToInt(nil)
	so(err, isErr)

	// Test non-number types (should error)
	_, err = anyToInt("not a number")
	so(err, isErr)

	_, err = anyToInt([]int{1, 2, 3})
	so(err, isErr)

	_, err = anyToInt(map[string]int{"a": 1})
	so(err, isErr)

	_, err = anyToInt(struct{ x int }{x: 1})
	so(err, isErr)

	// Test PathItem with valid Idx (>= 0)
	result, err = anyToInt(PathItem{Idx: 0, Key: ""})
	so(err, isNil)
	so(result, eq, 0)

	result, err = anyToInt(PathItem{Idx: 5, Key: "ignored"})
	so(err, isNil)
	so(result, eq, 5)

	// Test PathItem with invalid Idx (< 0, represents a key-only item)
	_, err = anyToInt(PathItem{Idx: -1, Key: "someKey"})
	so(err, isErr)
}

func testAnyToString(*testing.T) {
	// Test normal string
	result, err := anyToString("hello")
	so(err, isNil)
	so(result, eq, "hello")

	// Test nil parameter (should error)
	_, err = anyToString(nil)
	so(err, isErr)

	// Test non-string types (should error)
	_, err = anyToString(42)
	so(err, isErr)

	_, err = anyToString(3.14)
	so(err, isErr)

	_, err = anyToString([]string{"a", "b"})
	so(err, isErr)

	// Test PathItem with valid Key
	result, err = anyToString(PathItem{Key: "myKey", Idx: -1})
	so(err, isNil)
	so(result, eq, "myKey")

	result, err = anyToString(PathItem{Key: "anotherKey", Idx: 0})
	so(err, isNil)
	so(result, eq, "anotherKey")

	// Test PathItem with empty Key (represents an index-only item)
	_, err = anyToString(PathItem{Key: "", Idx: 3})
	so(err, isErr)
}
