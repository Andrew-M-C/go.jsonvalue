package jsonvalue

import (
	"fmt"
	"strconv"
	"testing"
)

func testSort(t *testing.T) {
	cv("sort array", func() { testSortArray(t) })
	cv("sort array errors", func() { testSortArrayError(t) })
	cv("sort marshal", func() { testSortMarshal(t) })
	cv("sort by string slice", func() { testSortByStringSlice(t) })
}

func testSortArray(t *testing.T) {
	arr := NewArray()

	arr.MustAppendInt(0).InTheEnd()
	arr.MustAppendInt(1).InTheEnd()
	arr.MustAppendInt(2).InTheEnd()
	arr.MustAppendInt(3).InTheEnd()
	arr.MustAppendInt(4).InTheEnd()
	arr.MustAppendInt(5).InTheEnd()
	arr.MustAppendInt(6).InTheEnd()
	arr.MustAppendInt(7).InTheEnd()
	arr.MustAppendInt(8).InTheEnd()
	arr.MustAppendInt(9).InTheEnd()

	t.Logf("pre-sorted: '%s'", arr.MustMarshalString())

	lessFunc := func(v1, v2 *V) bool {
		return v1.Int() > v2.Int()
	}
	arr.SortArray(lessFunc)

	res := arr.MustMarshalString()
	t.Logf("sorted res: '%s'", res)

	so(res, eq, `[9,8,7,6,5,4,3,2,1,0]`)
}

func testSortArrayError(t *testing.T) {
	// simple test, should not panic
	v := NewInt(1)
	v.SortArray(func(v1, v2 *V) bool { return false })

	v = NewArray()
	v.SortArray(nil)
}

func testSortMarshal(t *testing.T) {
	// default sequence
	expected := `{"0":0,"1":"1","2":2,"3":"3","4":4,"5":"5","6":6,"7":"7","8":8,"9":"9"}`
	t.Logf("expected string: %s", expected)

	for count := 0; count < 10; count++ {
		v := NewObject()
		for i := 0; i < 10; i++ {
			iStr := strconv.Itoa(i)
			if i&1 == 0 {
				v.MustSetInt(i).At(iStr)
			} else {
				v.MustSetString(iStr).At(iStr)
			}
		}

		s := v.MustMarshalString(Opt{MarshalLessFunc: DefaultStringSequence})
		so(s, eq, expected)

		s = v.MustMarshalString(OptDefaultStringSequence())
		so(s, eq, expected)
	}

	// key path
	orig := `{
		"object!":{
			"string!!!": "a string",
			"object!!":{
				"array!!!!":[
					1234,
					{
						"stringBB":"aa string",
						"stringA":"a string",
					}
				]
			},
			"null":null
		}
	}`

	v, err := UnmarshalString(orig)
	so(err, isNil)

	less := func(parentInfo *ParentInfo, keyA, keyB string, _, _ *V) bool {
		so(parentInfo, notNil)
		t.Logf("parentInfo: %v", parentInfo.KeyPath)
		s := ""
		for _, k := range parentInfo.KeyPath {
			s += fmt.Sprintf(`"%s"<%d><%v|%v>  `, k.String(), k.Int(), k.IsString(), k.IsInt())
		}
		t.Logf("Key path: %v, keyA: '%s', keyB: '%s'", s, keyA, keyB)

		return len(keyA) <= len(keyB)
	}

	s := v.MustMarshalString(Opt{
		OmitNull:        true,
		MarshalLessFunc: less,
	})

	expected = `{"object!":{"object!!":{"array!!!!":[1234,{"stringA":"a string","stringBB":"aa string"}]},"string!!!":"a string"}}`
	t.Logf("marshaled string: %v", s)
	so(s, eq, expected)

	s = v.MustMarshalString(OptOmitNull(true), OptKeySequenceWithLessFunc(less))
	so(s, eq, expected)
}

func testSortByStringSlice(t *testing.T) {
	seq := []string{
		"grandpa",
		"grandma",
		"father",
		"mother",
		"son",
		"daughter",
	}

	v := NewObject()
	v.MustSetString("Beef").At("friendB")
	v.MustSetString("Fish").At("friendA")
	v.MustSetString("Mayonnaise").At("daughter")
	v.MustSetString("Ketchup").At("son")
	v.MustSetString("Kentucky").At("grandpa")
	v.MustSetString("McDonald").At("grandma")
	v.MustSetString("Hanberger").At("father")
	v.MustSetString("Chips").At("mother")
	v.MustSetNull().At("relative")

	s := v.MustMarshalString(Opt{
		OmitNull:           true,
		MarshalKeySequence: seq,
	})

	expected := `{"grandpa":"Kentucky","grandma":"McDonald","father":"Hanberger","mother":"Chips","son":"Ketchup","daughter":"Mayonnaise","friendA":"Fish","friendB":"Beef"}`
	t.Logf("marshaled: '%s'", s)
	so(s, eq, expected)

	s = v.MustMarshalString(OptOmitNull(true), OptKeySequence(seq))
	so(s, eq, expected)
}
