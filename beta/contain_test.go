package beta

import (
	"testing"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

func testContains(_ *testing.T) {
	cv("general array", func() {
		v := jsonvalue.MustUnmarshalString(`[1,2,3,4]`)
		sub := jsonvalue.MustUnmarshalString(`[2,3]`)
		notSub := jsonvalue.MustUnmarshalString(`[1,3]`)

		so(Contains(v, sub), isTrue)
		so(Contains(v, notSub), isFalse)
	})

	cv("general object", func() {
		v := jsonvalue.MustUnmarshalString(`{"num":1234,"str":"Hello"}`)
		sub := jsonvalue.MustUnmarshalString(`{"num":1234.00}`)
		notSub := jsonvalue.MustUnmarshalString(`{"num":1234,"str":"hello"}`)

		so(Contains(v, sub), isTrue)
		so(Contains(v, notSub), isFalse)
	})

	cv("test case by author's colleague", func() {
		type P = []interface{}

		f := func(res bool, vStr string, path P, subStr string) {
			// t.Logf("%v - %v", vStr, subStr)

			v := jsonvalue.MustUnmarshalString(vStr)
			sub := jsonvalue.MustUnmarshalString(subStr)
			b := Contains(v, sub, path...)
			so(b, eq, res)
		}

		f(true, `{}`, nil, `{}`)                                                    // 两个空的json
		f(false, `{"obj":{}}`, nil, `{"Obj":{}}`)                                   // key不一样
		f(false, `{"a":2}`, nil, `{"a":3}`)                                         // 值不一样
		f(false, `{"a":2.00000000000}`, nil, `{"a":2.000000000000000000000001}`)    // 这里与 testcase 理解不同
		f(false, `{"a":1}`, nil, `{"a":2, "obj":null}`)                             // 没有obj
		f(false, `{"a":[1.false,2.0,false,3.0]}`, nil, `{"a":[true,2,3]}`)          // 这里与 testcase 理解不同
		f(true, `{"a":[2,2,3]}`, nil, `{"a":[2,2,3]}`)                              // 数组要个数一致
		f(true, `{"a":[2,2,3]}`, nil, `{"a":[2,3]}`)                                // 这里与 testcase 理解不同
		f(true, `{"a":[2,3], "obj":{"a":{"b":23}}}`, nil, `{"obj":{"a":{"b":23}}}`) // 支持嵌套

		f(true, `{}`, nil, `{}`)
		f(true, `{"a":[2,2,3]}`, P{"a", 2}, `3`)
		f(true, `{"a":[2,2,3]}`, P{"a"}, `[2,3]`)
		f(false, `{"a":[2,2,3]}`, P{"a"}, `2`)

		f(true, `{"a":[{}, {"b":22,"c":222}, 23, "hello"]}`, P{"a", 1, "b"}, `22`)
		f(false, `{"a":[{}, {"b":22,"c":222}, 23, "hello"]}`, P{"a", 4}, `22`) // {index out of range: len: 4, idx: 5"}
		f(true, `{"a":[{}, {"b":22,"c":222}, 23, "hello"]}`, P{"a", 1}, `{"c":222}`)
		// f(false, testjson, `$.store.book[*].category`, `{}`)
	})
}
