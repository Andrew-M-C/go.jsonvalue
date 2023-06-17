package jsonvalue_test

import (
	"fmt"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

func ExampleV_String() {
	v := jsonvalue.NewObject()
	v.MustSetString("Hello, string").At("object", "message")
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

	v = jsonvalue.NewFloat64(f)
	fmt.Println(v)

	v = jsonvalue.NewFloat64f(f, 'f', 6)
	fmt.Println(v)

	v = jsonvalue.NewFloat64f(f, 'e', 10)
	fmt.Println(v)
	// Output:
	// 123.123456789
	// 123.123457
	// 1.2312345679e+02
}

func ExampleOpt() {
	raw := `{"null":null}`
	v, _ := jsonvalue.UnmarshalString(raw)

	s := v.MustMarshalString()
	fmt.Println(s)
	s = v.MustMarshalString(jsonvalue.OptOmitNull(true))
	fmt.Println(s)
	// Output:
	// {"null":null}
	// {}
}

func ExampleV_Append() {
	s := `{"obj":{"arr":[1,2,3,4,5]}}`
	v, _ := jsonvalue.UnmarshalString(s)

	// append a zero in the bebinning of v.obj.arr
	v.MustAppendInt(0).InTheBeginning("obj", "arr")
	fmt.Println(v.MustMarshalString())

	// append a zero in the end of v.obj.arr
	v.MustAppendInt(0).InTheEnd("obj", "arr")
	fmt.Println(v.MustMarshalString())

	// Output:
	// {"obj":{"arr":[0,1,2,3,4,5]}}
	// {"obj":{"arr":[0,1,2,3,4,5,0]}}
}

func ExampleV_Insert() {
	s := `{"obj":{"arr":["hello","world"]}}`
	v, _ := jsonvalue.UnmarshalString(s)

	// insert a word in the middle, which is after the first word of the array
	v.MustInsertString("my").After("obj", "arr", 0)
	fmt.Println(v.MustMarshalString())

	// insert a word in the middle, which is before the second word of the array
	v.MustInsertString("beautiful").Before("obj", "arr", 2)
	fmt.Println(v.MustMarshalString())

	// Output:
	// {"obj":{"arr":["hello","my","world"]}}
	// {"obj":{"arr":["hello","my","beautiful","world"]}}
}

// For a simplest example:
//
// 这是最简单的例子：
func ExampleV_Set() {
	v := jsonvalue.NewObject()                      // {}
	v.MustSetObject().At("obj")                     // {"obj":{}}
	v.MustSet("Hello, world!").At("obj", "message") // {"obj":{"message":"Hello, world!"}}
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"obj":{"message":"Hello, world!"}}
}

// Or you can make it even more simpler, as At() function will automatically create objects those do not exist
//
// 或者你还可以更加简洁，因为 At() 函数会自动创建在值链中所需要但未创建的对象
func ExampleV_Set_another() {
	v := jsonvalue.NewObject()                      // {}
	v.MustSet("Hello, world!").At("obj", "message") // {"obj":{"message":"Hello, world!"}}
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"obj":{"message":"Hello, world!"}}
}

// As for array, At() also works
//
// 对于数组类型，At() 也是能够自动生成的
func ExampleV_Set_another2() {
	v := jsonvalue.NewObject()              // {}
	v.MustSet("Hello, world!").At("arr", 0) // {"arr":[Hello, world!]}
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"arr":["Hello, world!"]}
}

// Auto-array-creating in At() function is actually a bit complicated. It fails when specifying an position that the
// array does not have yet. But with one exception: the index value is equal to the length of an array, in this case,
// a new value will be append to the end of the array. This is quite convient when setting array elements in a for-range
// block.
//
// 在 At() 自动创建数组的逻辑其实稍微有点复杂，需要解释一下。当调用方在参数中指定在某个尚未存在的数组中设置一个值的时候，那么 At() 指定的位置（position）数字，
// 应当为0，操作才能成功；而当数组已经存在，那么 At() 指定的位置数，要么在数组中已存在，要么正好等于数组的长度，当后者的情况下，会在数组的最后追加值。
// 这个特性在使用 for-range 块时会非常有用。
func ExampleV_Set_another3() {
	v := jsonvalue.NewObject()                   // {}
	_, err := v.Set("Hello, world").At("arr", 1) // failed because there are no children of v.arr
	if err != nil {
		fmt.Println("got error:", err)
	}

	fmt.Println(v.MustMarshalString()) // as error occurred, the "arr" array would not be set

	integers := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for i, n := range integers {
		// this will succeed because i is equal to len(v.arr) every time
		v.MustSet(n).At("arr", i)
	}

	fmt.Println(v.MustMarshalString())
	// Output:
	// got error: out of range
	// {}
	// {"arr":[10,20,30,40,50,60,70,80,90,100]}
}

// As for elements those in positions that the array already has, At() will REPLACE it.
//
// 正如上文所述，如果在 At() 中指定了已存在的数组的某个位置，那么那个位置上的值会被替换掉，请注意。
func ExampleV_Set_another4() {
	v := jsonvalue.NewObject()
	for i := 0; i < 10; i++ {
		v.MustSetInt(i).At("arr", i)
	}

	fmt.Println(v.MustMarshalString())

	v.MustSet(123.12345).At("arr", 3)
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"arr":[0,1,2,3,4,5,6,7,8,9]}
	// {"arr":[0,1,2,123.12345,4,5,6,7,8,9]}
}

// In addition, any legal json type parameters are supported in Set(...).At(...).
// For example, we can set a struct as following:
//
// 此外，Set(...).At(...) 支持任意合法的 json 类型变量参数。比如我可以传入一个结构体:
func ExampleV_Set_another5() {
	type st struct {
		Text string `json:"text"`
	}
	child := st{
		Text: "Hello, jsonvalue!",
	}
	v := jsonvalue.NewObject()
	v.MustSet(child).At("child")
	fmt.Println(v.MustMarshalString())
	// Output:
	// {"child":{"text":"Hello, jsonvalue!"}}
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

func ExampleV_ForRangeArr() {
	s := `[1,2,3,4,5,6,7,8,9,10]`
	v, _ := jsonvalue.UnmarshalString(s)

	for i, v := range v.ForRangeArr() {
		fmt.Println(v)
		if i < 5 {
			// continue
		} else {
			break
		}
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}

func ExampleV_RangeObjects() {
	s := `{"message":"Hello, JSON!"}`
	v, _ := jsonvalue.UnmarshalString(s)

	v.RangeObjects(func(k string, v *jsonvalue.V) bool {
		fmt.Println(k, "-", v)
		return true
	})
	// Output:
	// message - Hello, JSON!
}

func ExampleV_ForRangeObj() {
	s := `{"message":"Hello, JSON!"}`
	v, _ := jsonvalue.UnmarshalString(s)

	for k, v := range v.ForRangeObj() {
		fmt.Println(k, "-", v)
	}
	// Output:
	// message - Hello, JSON!
}

func ExampleOptUTF8() {
	v := jsonvalue.NewObject()
	v.MustSetString("🇺🇸🇨🇳🇷🇺🇬🇧🇫🇷").At("UN_leaderships")

	asciiString := v.MustMarshalString()
	utf8String := v.MustMarshalString(jsonvalue.OptUTF8())
	fmt.Println("ASCII -", asciiString)
	fmt.Println("UTF-8 -", utf8String)
	// Output:
	// ASCII - {"UN_leaderships":"\uD83C\uDDFA\uD83C\uDDF8\uD83C\uDDE8\uD83C\uDDF3\uD83C\uDDF7\uD83C\uDDFA\uD83C\uDDEC\uD83C\uDDE7\uD83C\uDDEB\uD83C\uDDF7"}
	// UTF-8 - {"UN_leaderships":"🇺🇸🇨🇳🇷🇺🇬🇧🇫🇷"}
}

func ExmapleOptEscapeHTML() {
	v := jsonvalue.NewObject()
	v.MustSetString("https://hahaha.com?para1=<&para2=>").At("url")

	defaultStr := v.MustMarshalString()
	htmlOn := v.MustMarshalString(jsonvalue.OptEscapeHTML(true))
	htmlOff := v.MustMarshalString(jsonvalue.OptEscapeHTML(false))

	fmt.Println("default  -", defaultStr)
	fmt.Println("HTML ON  -", htmlOn)
	fmt.Println("HTML OFF -", htmlOff)
	// Output:
	// default  - {"url":"https:\/\/hahaha.com?para1=\u003C\u0026para2=\u0025"}
	// HTML ON  - {"url":"https:\/\/hahaha.com?para1=\u003C\u0026para2=\u0025"}
	// HTML OFF - {"url":"https:\/\/hahaha.com?para1=<&para2=>"}
}

func ExampleOptEscapeSlash() {
	v := jsonvalue.NewObject()
	v.MustSetString("https://google.com").At("google")

	defaultStr := v.MustMarshalString()
	escapeStr := v.MustMarshalString(jsonvalue.OptEscapeSlash(true))
	nonEscape := v.MustMarshalString(jsonvalue.OptEscapeSlash(false))

	fmt.Println("default -", defaultStr)
	fmt.Println("escape  -", escapeStr)
	fmt.Println("non-esc -", nonEscape)
	// Output:
	// default - {"google":"https:\/\/google.com"}
	// escape  - {"google":"https:\/\/google.com"}
	// non-esc - {"google":"https://google.com"}
}
