package jsonvalue

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIteration(t *testing.T) {
	test(t, "Range Array", testRangeArray)
	test(t, "Range Object", testRangeObject)
}

func testRangeArray(t *testing.T) {
	Convey("invalid array range", func() {
		v := NewString("")
		v.RangeArray(func(i int, c *V) bool {
			t.Errorf("should NOT iter here!!!")
			return true
		}) // just do not panic

		for range MustUnmarshalString("invalid").IterArray() {
			t.Errorf("should NOT iter here!!!")
		}

		invalidV, _ := MustUnmarshalString("invalid").Get("another invalid", 1, 2, 3, "opps")
		for range invalidV.IterArray() {
			t.Errorf("should NOT iter here!!!")
		}
	})

	Convey("nil array callback", func() {
		v := NewArray()
		v.AppendNull().InTheEnd()
		v.RangeArray(nil) // just do not panic

		for iter := range v.IterArray() {
			_ = iter.V.String() // just do not panic
		}
	})

	Convey("array range", func() {
		raw := `[1,2,3,4,5,6,7,8,9,10,"11","12",13.0]`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)

		rangeCount := 0
		v.RangeArray(func(i int, c *V) bool {
			rangeCount++
			element := 0
			if c.IsNumber() {
				element = c.Int()
			} else if c.IsString() {
				s := c.String()
				element, err = strconv.Atoi(s)
				So(err, ShouldBeNil)
			} else {
				t.Errorf("invalid jsonvalue type")
				return false
			}

			So(element, ShouldEqual, i+1)

			if element != i+1 {
				t.Errorf("unexpected element %d, index %d", element, i)
				return false
			}

			return true
		})

		So(rangeCount, ShouldEqual, 13)

		// broken array range
		rangeCount = 0
		brokenCount := 4
		v.RangeArray(func(i int, c *V) bool {
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
	Convey("invalid object range", func() {
		v := NewString("")
		v.RangeObjects(func(k string, c *V) bool {
			t.Errorf("should NOT iter here!!!")
			return true
		}) // just do not panic

		for range MustUnmarshalString("invalid").IterObjects() {
			t.Errorf("should NOT iter here!!!")
		}

		invalidV, _ := MustUnmarshalString("invalid").Get("another invalid", 1, 2, 3, "opps")
		for range invalidV.IterObjects() {
			t.Errorf("should NOT iter here!!!")
		}
	})

	Convey("nil object callback", func() {
		v := NewObject()
		v.SetString("world").At("hello")
		v.RangeObjects(nil) // just do not panic

		for iter := range v.IterObjects() {
			_ = iter.V.String() // just do not panic
		}
	})

	Convey("unmarshal object and whole object range", func() {
		raw := `{"number":12345,"bool":true,"string":"hello, world","null":null}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)

		checkKeys := map[string]bool{
			"number": true,
			"bool":   true,
			"string": true,
			"null":   true,
		}
		v.RangeObjects(func(k string, v *V) bool {
			delete(checkKeys, k)
			return true
		})
		if len(checkKeys) > 0 {
			t.Errorf("not all key checked, remains: %+v", checkKeys)
			return
		}

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
