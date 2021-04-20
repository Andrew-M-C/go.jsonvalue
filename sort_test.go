package jsonvalue

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSortArray(t *testing.T) {
	arr := NewArray()

	arr.AppendInt(0).InTheEnd()
	arr.AppendInt(1).InTheEnd()
	arr.AppendInt(2).InTheEnd()
	arr.AppendInt(3).InTheEnd()
	arr.AppendInt(4).InTheEnd()
	arr.AppendInt(5).InTheEnd()
	arr.AppendInt(6).InTheEnd()
	arr.AppendInt(7).InTheEnd()
	arr.AppendInt(8).InTheEnd()
	arr.AppendInt(9).InTheEnd()

	t.Logf("pre-sorted: '%s'", arr.MustMarshalString())

	lessFunc := func(v1, v2 *V) bool {
		return v1.Int() > v2.Int()
	}
	arr.SortArray(lessFunc)

	res := arr.MustMarshalString()
	t.Logf("sorted res: '%s'", res)

	if res != `[9,8,7,6,5,4,3,2,1,0]` {
		t.Errorf("array sort failed")
		return
	}
}

func TestSortArrayError(t *testing.T) {
	// simple test, should not panic
	v := NewInt(1)
	v.SortArray(func(v1, v2 *V) bool { return false })

	v = NewArray()
	v.SortArray(nil)
}

func TestSortMarshal(t *testing.T) {
	// default sequence
	expected := `{"0":0,"1":"1","2":2,"3":"3","4":4,"5":"5","6":6,"7":"7","8":8,"9":"9"}`
	t.Logf("expected string: %s", expected)

	for count := 0; count < 10; count++ {
		v := NewObject()
		for i := 0; i < 10; i++ {
			iStr := strconv.Itoa(i)
			if i&1 == 0 {
				v.SetInt(i).At(iStr)
			} else {
				v.SetString(iStr).At(iStr)
			}
		}

		s := v.MustMarshalString(Opt{MarshalLessFunc: DefaultStringSequence})
		if s != expected {
			t.Errorf("unexpected string: %s", s)
			return
		}
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
						"stringA":"a string"
					}
				]
			},
			"null":null
		}
	}`

	v, err := UnmarshalString(orig)
	if err != nil {
		t.Errorf("unmarshal error: %v", err)
		return
	}

	s := v.MustMarshalString(Opt{
		OmitNull: true,
		MarshalLessFunc: func(parentInfo *ParentInfo, keyA, keyB string, vA, vB *V) bool {
			t.Logf("parentInfo: %v", parentInfo.KeyPath)
			s := ""
			for _, k := range parentInfo.KeyPath {
				s += fmt.Sprintf(`"%s"<%d><%v|%v>  `, k.String(), k.Int(), k.IsString(), k.IsInt())
			}
			t.Logf("Key path: %v", s)

			return len(keyA) <= len(keyB)
		},
	})
	t.Logf("marshaled string: %v", s)

	expected = `{"object!":{"object!!":{"array!!!!":[1234,{"stringA":"a string","stringBB":"aa string"}]},"string!!!":"a string"}}`
	if s != expected {
		t.Errorf("unpxpected marshaled string")
		return
	}
}

func TestSortByStringSlice(t *testing.T) {
	seq := []string{
		"grandpa",
		"grandma",
		"father",
		"mother",
		"son",
		"daughter",
	}

	v := NewObject()
	v.SetString("Beef").At("friendB")
	v.SetString("Fish").At("friendA")
	v.SetString("Mayonnaise").At("daughter")
	v.SetString("Ketchup").At("son")
	v.SetString("Kentucky").At("grandpa")
	v.SetString("McDonald").At("grandma")
	v.SetString("Hanberger").At("father")
	v.SetString("Chips").At("mother")
	v.SetNull().At("relative")

	s := v.MustMarshalString(Opt{
		OmitNull:           true,
		MarshalKeySequence: seq,
	})
	t.Logf("marshaled: '%s'", s)

	expected := `{"grandpa":"Kentucky","grandma":"McDonald","father":"Hanberger","mother":"Chips","son":"Ketchup","daughter":"Mayonnaise","friendA":"Fish","friendB":"Beef"}`
	if s != expected {
		t.Errorf("unexpected marshal result")
		return
	}
}
