package jsonvalue

import (
	"strconv"
	"testing"
)

func TestRange(t *testing.T) {
	var checkCount int
	var topic string
	var err error
	var raw string
	var v *V
	var rangeCount int
	checkErr := func() {
		checkCount++
		if err != nil {
			t.Errorf("%02d - %s - error occurred: %v", checkCount, topic, err)
			return
		}
	}

	topic = "invalid array range"
	v = NewString("")
	v.RangeArray(func(i int, c *V) bool {
		return true
	}) // just do not panic

	topic = "nil array callback"
	v = NewArray()
	v.AppendNull().InTheEnd()
	v.RangeArray(nil) // just do not panic

	for iter := range v.IterArray() {
		_ = iter.V.String() // just do not panic
	}

	topic = "unmarshal array"
	raw = `[1,2,3,4,5,6,7,8,9,10,"11","12",13.0]`
	v, err = UnmarshalString(raw)
	checkErr()

	topic = "complete array range"
	rangeCount = 0
	v.RangeArray(func(i int, c *V) bool {
		rangeCount++
		element := 0
		if c.IsNumber() {
			element = c.Int()
		} else if c.IsString() {
			s := c.String()
			element, err = strconv.Atoi(s)
			checkErr()
		} else {
			t.Errorf("invalid jsonvalue type")
			return false
		}

		if element != i+1 {
			t.Errorf("unexpected element %d, index %d", element, i)
			return false
		}

		return true
	})

	if rangeCount != 13 {
		t.Errorf("unexpected rangeCount: %d", rangeCount)
		return
	}

	topic = "broken array range"
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

	// object
	topic = "invalid object range"
	v = NewString("")
	v.RangeObjects(func(k string, c *V) bool {
		return true
	}) // just do not panic

	topic = "nil object callback"
	v = NewObject()
	v.SetString("world").At("hello")
	v.RangeObjects(nil) // just do not panic

	for iter := range v.IterObjects() {
		_ = iter.V.String()
		// just do not panic
	}

	topic = "unmarshal object"
	raw = `{"number":12345,"bool":true,"string":"hello, world","null":null}`
	v, err = UnmarshalString(raw)
	checkErr()

	topic = "complete object range"
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

	topic = "broken object range"
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
}
