package jsonvalue

import (
	"fmt"
	"strconv"
	"testing"
)

func testIteration(t *testing.T) {
	cv("Range Array", func() { testRangeArray(t) })
	cv("Range Object", func() { testRangeObject(t) })
	cv("Range Object by seq", func() { testRangeObjectsBySetSequence(t) })
	cv("Walk", func() { testWalk(t) })
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

		invalidV, _ := MustUnmarshalString("invalid").Get("another invalid", 1, 2, 3, "oops")
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

		invalidV, _ := MustUnmarshalString("invalid").Get("another invalid", 1, 2, 3, "oops")
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

func testWalk(t *testing.T) {
	cv("nil callback function", func() { testWalkNilCallback(t) })
	cv("basic types", func() { testWalkBasicTypes(t) })
	cv("empty containers", func() { testWalkEmptyContainers(t) })
	cv("simple object", func() { testWalkSimpleObject(t) })
	cv("simple array", func() { testWalkSimpleArray(t) })
	cv("nested objects", func() { testWalkNestedObjects(t) })
	cv("nested arrays", func() { testWalkNestedArrays(t) })
	cv("mixed nesting", func() { testWalkMixedNesting(t) })
	cv("path validation", func() { testWalkPathValidation(t) })
	cv("complex nested structure", func() { testWalkComplexStructure(t) })
	cv("early termination", func() { testWalkEarlyTermination(t) })
}

func testWalkNilCallback(t *testing.T) {
	// Test nil callback - should not panic and do nothing
	v := NewString("test")
	v.Walk(nil) // Should not panic

	v2 := NewObject()
	v2.MustSetString("value").At("key")
	v2.Walk(nil) // Should not panic

	v3 := NewArray()
	v3.MustAppendString("item").InTheEnd()
	v3.Walk(nil) // Should not panic
}

func testWalkBasicTypes(t *testing.T) {
	// Test string type
	cv("string value", func() {
		v := NewString("hello")
		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			so(len(path), eq, 0) // Root level, no path
			so(val.ValueType(), eq, String)
			so(val.String(), eq, "hello")
			return true
		})
		so(walkCount, eq, 1)
	})

	// Test number type
	cv("number value", func() {
		v := NewInt64(42)
		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			lastItem := path.Last()
			so(len(path), eq, 0)
			so(lastItem.Idx, eq, -1)
			so(lastItem.Key, eq, "")
			so(val.ValueType(), eq, Number)
			so(val.Int64(), eq, 42)
			return true
		})
		so(walkCount, eq, 1)
	})

	// Test boolean type
	cv("boolean value", func() {
		v := NewBool(true)
		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			so(len(path), eq, 0)
			so(val.ValueType(), eq, Boolean)
			so(val.Bool(), eq, true)
			return true
		})
		so(walkCount, eq, 1)
	})

	// Test null type
	cv("null value", func() {
		v := NewNull()
		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			so(len(path), eq, 0)
			so(val.ValueType(), eq, Null)
			so(val.IsNull(), eq, true)
			return true
		})
		so(walkCount, eq, 1)
	})
}

func testWalkEmptyContainers(t *testing.T) {
	// Test empty object
	cv("empty object", func() {
		v := NewObject()
		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			t.Errorf("Should not walk through empty object")
			return true
		})
		so(walkCount, eq, 0)
	})

	// Test empty array
	cv("empty array", func() {
		v := NewArray()
		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			t.Errorf("Should not walk through empty array")
			return true
		})
		so(walkCount, eq, 0)
	})
}

func testWalkSimpleObject(t *testing.T) {
	v := NewObject()
	v.MustSetString("hello").At("greeting")
	v.MustSetInt64(123).At("number")
	v.MustSetBool(true).At("flag")
	v.MustSetNull().At("empty")

	expected := map[string]struct {
		valueType ValueType
		value     any
	}{
		"greeting": {String, "hello"},
		"number":   {Number, int64(123)},
		"flag":     {Boolean, true},
		"empty":    {Null, nil},
	}

	walkCount := 0
	v.Walk(func(path Path, val *V) bool {
		walkCount++
		so(len(path), eq, 1)

		key := path[0].Key
		so(path[0].Idx, eq, -1) // Not an array index
		so(key, ne, "")         // Should have a key

		expectedVal, exists := expected[key]
		so(exists, isTrue)
		so(val.ValueType(), eq, expectedVal.valueType)

		switch expectedVal.valueType {
		case String:
			so(val.String(), eq, expectedVal.value)
		case Number:
			so(val.Int64(), eq, expectedVal.value)
		case Boolean:
			so(val.Bool(), eq, expectedVal.value)
		case Null:
			so(val.IsNull(), isTrue)
		}

		delete(expected, key)
		return true
	})

	so(walkCount, eq, 4)
	so(len(expected), eq, 0) // All expected values should be found
}

func testWalkSimpleArray(t *testing.T) {
	v := NewArray()
	v.MustAppendString("first").InTheEnd()
	v.MustAppendInt64(456).InTheEnd()
	v.MustAppendBool(false).InTheEnd()
	v.MustAppendNull().InTheEnd()

	expectedValues := []struct {
		valueType ValueType
		value     any
	}{
		{String, "first"},
		{Number, int64(456)},
		{Boolean, false},
		{Null, nil},
	}

	walkCount := 0
	v.Walk(func(path Path, val *V) bool {
		so(len(path), eq, 1)

		idx := path[0].Idx
		so(path[0].Key, eq, "") // Not an object key
		so(idx, ne, -1)         // Should have an index
		so(idx >= 0, isTrue)
		so(idx < len(expectedValues), isTrue)

		expected := expectedValues[idx]
		so(val.ValueType(), eq, expected.valueType)

		switch expected.valueType {
		case String:
			so(val.String(), eq, expected.value)
		case Number:
			so(val.Int64(), eq, expected.value)
		case Boolean:
			so(val.Bool(), eq, expected.value)
		case Null:
			so(val.IsNull(), isTrue)
		}

		walkCount++
		return true
	})

	so(walkCount, eq, 4)
}

func testWalkNestedObjects(t *testing.T) {
	// Create nested object structure:
	// {
	//   "level1": {
	//     "level2": {
	//       "value": "deep"
	//     }
	//   }
	// }
	v := NewObject()
	level1 := NewObject()
	level2 := NewObject()
	level2.MustSetString("deep").At("value")
	level1.MustSet(level2).At("level2")
	v.MustSet(level1).At("level1")

	walkCount := 0
	v.Walk(func(path Path, val *V) bool {
		walkCount++

		if len(path) == 3 {
			// This should be the deepest value
			so(path[0].Key, eq, "level1")
			so(path[0].Idx, eq, -1)
			so(path[1].Key, eq, "level2")
			so(path[1].Idx, eq, -1)
			so(path[2].Key, eq, "value")
			so(path[2].Idx, eq, -1)
			so(val.ValueType(), eq, String)
			so(val.String(), eq, "deep")
		}

		return true
	})

	so(walkCount, eq, 1) // Only the leaf value should be visited
}

func testWalkNestedArrays(t *testing.T) {
	// Create nested array structure: [[[42]]]
	v := NewArray()
	level1 := NewArray()
	level2 := NewArray()
	level2.MustAppendInt64(42).InTheEnd()
	level1.MustAppend(level2).InTheEnd()
	v.MustAppend(level1).InTheEnd()

	walkCount := 0
	v.Walk(func(path Path, val *V) bool {
		walkCount++

		if len(path) == 3 {
			// This should be the deepest value
			so(path[0].Key, eq, "")
			so(path[0].Idx, eq, 0)
			so(path[1].Key, eq, "")
			so(path[1].Idx, eq, 0)
			so(path[2].Key, eq, "")
			so(path[2].Idx, eq, 0)
			so(val.ValueType(), eq, Number)
			so(val.Int64(), eq, 42)
		}

		return true
	})

	so(walkCount, eq, 1) // Only the leaf value should be visited
}

func testWalkMixedNesting(t *testing.T) {
	// Create mixed structure:
	// {
	//   "array": [
	//     "string_item",
	//     {
	//       "nested_key": 789
	//     }
	//   ],
	//   "simple": "value"
	// }
	v := NewObject()
	arr := NewArray()
	arr.MustAppendString("string_item").InTheEnd()
	nestedObj := NewObject()
	nestedObj.MustSetInt64(789).At("nested_key")
	arr.MustAppend(nestedObj).InTheEnd()
	v.MustSet(arr).At("array")
	v.MustSetString("value").At("simple")

	walkResults := make(map[string]any)
	walkCount := 0

	v.Walk(func(path Path, val *V) bool {
		walkCount++

		// Create a path string for identification using Path.String()
		pathStr := path.String()

		if val.ValueType() == String {
			walkResults[pathStr] = val.String()
		} else if val.ValueType() == Number {
			walkResults[pathStr] = val.Int64()
		}

		return true
	})

	so(walkCount, eq, 3) // "simple", "array[0]", "array[1].nested_key"
	so(walkResults["simple"], eq, "value")
	so(walkResults["array.[0]"], eq, "string_item")
	so(walkResults["array.[1].nested_key"], eq, int64(789))
}

func testWalkPathValidation(t *testing.T) {
	// Test that PathItem fields are correctly set
	v := NewObject()
	arr := NewArray()
	arr.MustAppendString("item0").InTheEnd()
	arr.MustAppendString("item1").InTheEnd()
	v.MustSet(arr).At("my_array")
	v.MustSetString("direct_value").At("my_string")

	pathValidations := make(map[string]bool)

	v.Walk(func(path Path, val *V) bool {
		pathStr := ""
		for _, p := range path {
			if p.Key != "" {
				// Object key
				so(p.Idx, eq, -1) // Idx should be -1 for object keys
				pathStr += "obj:" + p.Key + "/"
			} else {
				// Array index
				so(p.Idx, ne, -1)      // Idx should not be -1 for array indices
				so(p.Idx >= 0, isTrue) // Should be non-negative
				so(p.Key, eq, "")      // Key should be empty for array indices
				pathStr += fmt.Sprintf("arr:%d/", p.Idx)
			}
		}
		pathValidations[pathStr] = true
		return true
	})

	// Check that we visited all expected paths
	so(pathValidations["obj:my_array/arr:0/"], isTrue)
	so(pathValidations["obj:my_array/arr:1/"], isTrue)
	so(pathValidations["obj:my_string/"], isTrue)
	so(len(pathValidations), eq, 3)
}

func testWalkComplexStructure(t *testing.T) {
	// Create a complex nested structure to ensure comprehensive coverage
	raw := `{
		"users": [
			{
				"name": "John",
				"age": 30,
				"active": true
			},
			{
				"name": "Jane",
				"age": 25,
				"active": false,
				"metadata": null
			}
		],
		"config": {
			"debug": true,
			"settings": {
				"timeout": 5000,
				"retries": 3
			}
		},
		"empty_array": [],
		"empty_object": {}
	}`

	v, err := UnmarshalString(raw)
	so(err, isNil)

	leafValues := make(map[string]any)
	walkCount := 0

	v.Walk(func(path Path, val *V) bool {
		walkCount++

		// Store leaf values
		switch val.ValueType() {
		case String:
			leafValues[path.String()] = val.String()
		case Number:
			leafValues[path.String()] = val.Int64()
		case Boolean:
			leafValues[path.String()] = val.Bool()
		case Null:
			leafValues[path.String()] = nil
		}

		return true
	})

	// Verify we captured all leaf values correctly
	so(leafValues["users.[0].name"], eq, "John")
	so(leafValues["users.[0].age"], eq, int64(30))
	so(leafValues["users.[0].active"], eq, true)
	so(leafValues["users.[1].name"], eq, "Jane")
	so(leafValues["users.[1].age"], eq, int64(25))
	so(leafValues["users.[1].active"], eq, false)
	so(leafValues["users.[1].metadata"], isNil)
	so(leafValues["config.debug"], eq, true)
	so(leafValues["config.settings.timeout"], eq, int64(5000))
	so(leafValues["config.settings.retries"], eq, int64(3))

	// Should have walked through all leaf nodes (non-container values)
	so(walkCount, eq, 10) // All leaf values
}

func testWalkEarlyTermination(t *testing.T) {
	// Test early termination in simple array
	cv("early termination in array", func() {
		v := NewArray()
		v.MustAppendString("first").InTheEnd()
		v.MustAppendString("second").InTheEnd()
		v.MustAppendString("third").InTheEnd()
		v.MustAppendString("fourth").InTheEnd()

		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			// Stop at the second element
			return walkCount < 2
		})

		so(walkCount, eq, 2) // Should have stopped after 2 elements
	})

	// Test early termination in simple object
	cv("early termination in object", func() {
		v := NewObject()
		v.MustSetString("value1").At("key1")
		v.MustSetString("value2").At("key2")
		v.MustSetString("value3").At("key3")
		v.MustSetString("value4").At("key4")

		walkCount := 0
		var lastKey string
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			lastKey = path[0].Key
			// Stop after first element
			return false
		})

		so(walkCount, eq, 1) // Should have stopped after 1 element
		so(lastKey, ne, "")  // Should have captured a key
	})

	// Test early termination in nested structure
	cv("early termination in nested structure", func() {
		// Use array structure to ensure deterministic iteration order
		v := NewArray()
		v.MustAppendString("item1").InTheEnd()
		v.MustAppendString("item2").InTheEnd()
		v.MustAppendString("item3").InTheEnd()
		v.MustAppendString("item4").InTheEnd()
		v.MustAppendString("item5").InTheEnd()

		// This structure has 5 leaf nodes with deterministic iteration order
		totalWalkCount := 0
		lastValue := ""

		v.Walk(func(path Path, v *V) bool {
			totalWalkCount++
			// Build path string for debugging
			lastValue = v.String()

			// Stop after fourth element to allow some iteration
			return path.Last().Idx < 3
		})

		so(totalWalkCount, eq, 4)
		so(lastValue, eq, "item4")
	})

	// Test that returning true continues iteration normally
	cv("normal iteration with true return", func() {
		v := NewArray()
		v.MustAppendString("a").InTheEnd()
		v.MustAppendString("b").InTheEnd()
		v.MustAppendString("c").InTheEnd()

		walkCount := 0
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			return true // Always continue
		})

		so(walkCount, eq, 3) // Should visit all elements
	})

	// Test early termination with specific condition
	cv("conditional early termination", func() {
		v := NewArray()
		v.MustAppendString("apple").InTheEnd()
		v.MustAppendString("banana").InTheEnd()
		v.MustAppendString("cherry").InTheEnd()
		v.MustAppendString("date").InTheEnd()

		walkCount := 0
		foundBanana := false
		v.Walk(func(path Path, val *V) bool {
			walkCount++
			if val.String() == "banana" {
				foundBanana = true
				return false // Stop when we find banana
			}
			return true
		})

		so(foundBanana, isTrue)    // Should have found banana
		so(walkCount >= 1, isTrue) // Should have visited at least 1 element
		so(walkCount <= 2, isTrue) // Should have stopped by element 2
	})
}
