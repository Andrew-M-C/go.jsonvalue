package jsonvalue

import (
	"strconv"
	"testing"
)

func testIteration(t *testing.T) {
	cv("Range Array", func() { testRangeArray(t) })
	cv("Range Object", func() { testRangeObject(t) })
	cv("Range Object by seq", func() { testRangeObjectsBySetSequence(t) })
}

func testRangeArray(t *testing.T) {
	cv("invalid array range", func() {
		v := NewString("")
		v.RangeArray(func(_ int, _ *V) bool {
			t.Errorf("should NOT iter here!!!")
			return true
		}) // just do not panic

		for range MustUnmarshalString("invalid").IterArray() {
			t.Errorf("should NOT iter here!!!")
		}

		for range MustUnmarshalString("invalid").ForRangeArr() {
			t.Errorf("should NOT iter here!!!")
		}

		invalidV, _ := MustUnmarshalString("invalid").Get("another invalid", 1, 2, 3, "opps")
		for range invalidV.IterArray() {
			t.Errorf("should NOT iter here!!!")
		}

		for range invalidV.ForRangeArr() {
			t.Errorf("should NOT iter here!!!")
		}
	})

	cv("nil array callback", func() {
		v := NewArray()
		v.MustAppendNull().InTheEnd()
		v.RangeArray(nil) // just do not panic

		for iter := range v.IterArray() {
			_ = iter.V.String() // just do not panic
		}

		for _, v := range v.ForRangeArr() {
			_ = v.String() // just do not panic
		}
	})

	cv("array range", func() {
		raw := `[1,2,3,4,5,6,7,8,9,10,"11","12",13.0]`
		v, err := UnmarshalString(raw)
		so(err, isNil)

		rangeCount := 0
		v.RangeArray(func(i int, c *V) bool {
			rangeCount++
			element := 0
			if c.IsNumber() {
				element = c.Int()
			} else if c.IsString() {
				s := c.String()
				element, err = strconv.Atoi(s)
				so(err, isNil)
			} else {
				t.Errorf("invalid jsonvalue type")
				return false
			}

			so(element, eq, i+1)

			if element != i+1 {
				t.Errorf("unexpected element %d, index %d", element, i)
				return false
			}

			return true
		})

		so(rangeCount, eq, 13)

		// broken array range
		rangeCount = 0
		brokenCount := 4
		v.RangeArray(func(i int, _ *V) bool {
			if i < brokenCount {
				rangeCount++
				return true
			}
			return false
		})
		if rangeCount != brokenCount {
			t.Errorf("expected rangeCount %d but %d got", brokenCount, rangeCount)
			return
		}
	})
}

func testRangeObject(t *testing.T) {
	cv("invalid object range", func() {
		v := NewString("")
		v.RangeObjects(func(_ string, _ *V) bool {
			t.Errorf("should NOT iter here!!!")
			return true
		}) // just do not panic

		for range MustUnmarshalString("invalid").IterObjects() {
			t.Errorf("should NOT iter here!!!")
		}

		for range MustUnmarshalString("invalid").ForRangeObj() {
			t.Errorf("should NOT iter here!!!")
		}

		invalidV, _ := MustUnmarshalString("invalid").Get("another invalid", 1, 2, 3, "opps")
		for range invalidV.IterObjects() {
			t.Errorf("should NOT iter here!!!")
		}

		for range invalidV.ForRangeObj() {
			t.Errorf("should NOT iter here!!!")
		}
	})

	cv("nil object callback", func() {
		v := NewObject()
		v.MustSetString("world").At("hello")
		v.RangeObjects(nil) // just do not panic

		iterCount := 0
		for iter := range v.IterObjects() {
			iterCount++
			_ = iter.V.String() // just do not panic
		}
		so(iterCount, eq, 1)

		iterCount = 0
		for _, v := range v.ForRangeObj() {
			iterCount++
			_ = v.String() // just do not panic
		}
		so(iterCount, eq, 1)
	})

	cv("unmarshal object and whole object range", func() {
		raw := `{"number":12345,"bool":true,"string":"hello, world","null":null}`
		v, err := UnmarshalString(raw)
		so(err, isNil)

		checkKeys := map[string]bool{
			"number": true,
			"bool":   true,
			"string": true,
			"null":   true,
		}

		expectedIterCount := len(checkKeys)
		iterCount := 0
		v.RangeObjects(func(k string, _ *V) bool {
			iterCount++
			delete(checkKeys, k)
			return true
		})
		so(iterCount, eq, expectedIterCount)
		so(len(checkKeys), isZero)

		// broken object range
		caughtValue := ""
		caughtKey := ""
		v.RangeObjects(func(k string, v *V) bool {
			caughtKey = k
			caughtValue = v.String()
			return k != "string"
		})

		if caughtKey == "string" && caughtValue == "hello, world" {
			// OK
		} else {
			t.Errorf("unexpected K-V: %s - %s", caughtKey, caughtValue)
			return
		}
	})
}

func testRangeObjectsBySetSequence(t *testing.T) {
	cv("invalid object range", func() {
		v := NewString("")
		v.RangeObjectsBySetSequence(func(_ string, _ *V) bool {
			t.Errorf("should NOT iter here!!!")
			return true
		}) // just do not panic
	})

	cv("nil callback", func() {
		v := NewObject()
		v.RangeObjectsBySetSequence(nil)
		// just do not panic
		so(true, isTrue)
	})

	cv("unmarshal object and whole object range", func() {
		v := NewObject()
		const size = 50
		const iterate = 100
		for i := 0; i < size; i++ {
			v.MustSet(i).At(strconv.FormatInt(int64(i), 10))
		}
		so(v.Len(), eq, size)

		for i := 0; i < iterate; i++ {
			lastNum := int64(-1)
			v.RangeObjectsBySetSequence(func(key string, v *V) bool {
				k, err := strconv.ParseInt(key, 10, 64)
				so(err, isNil)
				so(k, eq, v.Int64())

				so(k, eq, lastNum+1)

				lastNum = k

				// last one? exist
				return k+2 < size
			})
			so(lastNum, ne, -1)
			so(lastNum, eq, size-2)
		}
	})
}
