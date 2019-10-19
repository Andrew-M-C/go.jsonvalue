package jsonvalue_test

import (
	"fmt"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

func ExampleV_String() {
	v := jsonvalue.NewObject()
	v.SetString("Hello, string").At("object", "message")
	fmt.Println(v)

	child, _ := v.Get("object")
	fmt.Println(child)

	child, _ = v.Get("object", "message")
	fmt.Println(child)
	// Output:
	// {object: {message: Hello, string}}
	// {message: Hello, string}
	// Hello, string
}

func ExampleNewFloat64() {
	f := 123.123456789
	var v *jsonvalue.V

	v = jsonvalue.NewFloat64(f, 9)
	fmt.Println(v)

	v = jsonvalue.NewFloat64(f, 6)
	fmt.Println(v)

	v = jsonvalue.NewFloat64(f, 10)
	fmt.Println(v)
	// Output:
	// 123.123456789
	// 123.123457
	// 123.1234567890
}

func ExampleOpt() {
	raw := `{"string":"hello, world", "null":null}`
	v, _ := jsonvalue.UnmarshalString(raw)

	s := v.MustMarshalString()
	fmt.Println(s)
	s = v.MustMarshalString(jsonvalue.Opt{OmitNull: true})
	fmt.Println(s)
	// Output:
	// {"string":"hello, world","null":null}
	// {"string":"hello, world"}
}

func ExampleAppend_InTheBeginning() {
	s := `{"obj":{"arr":[1,2,3,4,5]}}`
	v, _ := jsonvalue.UnmarshalString(s)

	// append a zero in the bebinning of v.obj.arr
	v.AppendInt(0).InTheBeginning("obj", "arr")
	s = v.MustMarshalString()

	fmt.Println(s)
	// Output:
	// {"obj":{"arr":[0,1,2,3,4,5]}}
}

func ExampleAppend_InTheEnd() {
	s := `{"obj":{"arr":[1,2,3,4,5]}}`
	v, _ := jsonvalue.UnmarshalString(s)

	// append a zero in the end of v.obj.arr
	v.AppendInt(0).InTheEnd("obj", "arr")
	s = v.MustMarshalString()

	fmt.Println(s)
	// Output:
	// {"obj":{"arr":[1,2,3,4,5,0]}}
}

func ExampleInsert_After() {
	s := `{"obj":{"arr":["hello","world"]}}`
	v, _ := jsonvalue.UnmarshalString(s)

	// insert a word in the middle, which is after the first word of the array
	v.InsertString("my").After("obj", "arr", 0)

	fmt.Println(v.MustMarshalString())
	// Output:
	// {"obj":{"arr":["hello","my","world"]}}
}

func ExampleInsert_Before() {
	s := `{"obj":{"arr":["hello","world"]}}`
	v, _ := jsonvalue.UnmarshalString(s)

	// insert a word in the middle, which is before the second word of the array
	v.InsertString("my").Before("obj", "arr", 1)

	fmt.Println(v.MustMarshalString())
	// Output:
	// {"obj":{"arr":["hello","my","world"]}}
}

// For a simplest example:
func ExampleSet_At_1() {
	v := jsonvalue.NewObject()                        // {}
	v.SetObject().At("obj")                           // {"obj":{}}
	v.SetString("Hello, world!").At("obj", "message") // {"obj":{"message":"Hello, world!"}}
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"obj":{"message":"Hello, world!"}}
}

// Or you can make it even more simpler, as At() function will automatically create objects those do not exist
func ExampleSet_At_2() {
	v := jsonvalue.NewObject()                        // {}
	v.SetString("Hello, world!").At("obj", "message") // {"obj":{"message":"Hello, world!"}}
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"obj":{"message":"Hello, world!"}}
}

// As for array, At() also works
func ExampleSet_At_3() {
	v := jsonvalue.NewObject()                // {}
	v.SetString("Hello, world!").At("arr", 0) // {"arr":[Hello, world!]}
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"arr":["Hello, world!"]}
}

// Auto-array-creating in At() function is actually a bit complicated. It fails when specifying an position that the array does not have yet. But with one exception: the index value is equal to the length of an array, in this case, a new value will be append to the end of the array. This is quite convient when setting array elements in a for-range block.
func ExampleSet_At_4() {
	v := jsonvalue.NewObject()                         // {}
	_, err := v.SetString("Hello, world").At("arr", 1) // failed because there are no children of v.arr
	if err != nil {
		fmt.Println("got error:", err)
	}

	fmt.Println(v.MustMarshalString()) // as error occurred, the "arr" array would not be set

	integers := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for i, n := range integers {
		// this will succeed because i is equal to len(v.arr) every time
		v.SetInt(n).At("arr", i)
	}

	fmt.Println(v.MustMarshalString())
	// Output:
	// got error: out of range
	// {}
	// {"arr":[10,20,30,40,50,60,70,80,90,100]}
}

// As for elements those in positions that the array already has, At() will REPLACE it.
func ExampleSet_At_5() {
	v := jsonvalue.NewObject()
	for i := 0; i < 10; i++ {
		v.SetInt(i).At("arr", i)
	}

	fmt.Println(v.MustMarshalString())

	v.SetFloat64(123.12345, -1).At("arr", 3)
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"arr":[0,1,2,3,4,5,6,7,8,9]}
	// {"arr":[0,1,2,123.12345,4,5,6,7,8,9]}
}

func ExampleV_Get() {
	s := `{"objA":{"objB":{"message":"Hello, world!"}}}`
	v, _ := jsonvalue.UnmarshalString(s)
	msg, _ := v.Get("objA", "objB", "message")
	fmt.Println(msg.String())
	// Output:
	// Hello, world!
}

func ExampleV_GreaterThanInt64Max() {
	v1 := jsonvalue.NewUint64(uint64(9223372036854775808)) // 0x8000000000000000
	v2 := jsonvalue.NewUint64(uint64(9223372036854775807)) // 0x7FFFFFFFFFFFFFFF
	v3 := jsonvalue.NewInt64(int64(-9223372036854775807))
	fmt.Println(v1.GreaterThanInt64Max())
	fmt.Println(v2.GreaterThanInt64Max())
	fmt.Println(v3.GreaterThanInt64Max())
	// Output:
	// true
	// false
	// false
}

func ExampleV_RangeArray() {
	s := `[1,2,3,4,5,6,7,8,9,10]`
	v, _ := jsonvalue.UnmarshalString(s)

	v.RangeArray(func(i int, v *jsonvalue.V) bool {
		fmt.Println(v)
		return i < 5
	})
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}
