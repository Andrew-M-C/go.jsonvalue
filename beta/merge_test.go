package beta_test

import (
	"testing"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/Andrew-M-C/go.jsonvalue/beta"
)

func testMerge(_ *testing.T) {
	cv("object merge", func() {
		cv("disjoint keys", func() {
			to, cs := mergeJSON(`{"a":1}`, `{"b":2}`)
			soJSON(to, `{"a":1,"b":2}`)
			so(len(cs), eq, 0)
		})
		cv("nested objects", func() {
			to, cs := mergeJSON(`{"a":{"x":1}}`, `{"a":{"y":2}}`)
			soJSON(to, `{"a":{"x":1,"y":2}}`)
			so(len(cs), eq, 0)
		})
		cv("empty object plus nested", func() {
			to, cs := mergeJSON(`{}`, `{"a":{"b":1}}`)
			soJSON(to, `{"a":{"b":1}}`)
			so(len(cs), eq, 0)
		})
		cv("equal nested number is not a conflict", func() {
			to, cs := mergeJSON(`{"a":1}`, `{"a":1.0}`)
			soJSON(to, `{"a":1}`)
			so(len(cs), eq, 0)
		})
		cv("unequal number conflicts and keeps to", func() {
			to, cs := mergeJSON(`{"a":1,"b":{"x":1,"y":2}}`, `{"a":2,"b":{"x":1,"y":3,"z":4}}`)
			soJSON(to, `{"a":1,"b":{"x":1,"y":2,"z":4}}`)
			so(pathStrings(cs), resemble, []string{"a", "b.y"})
		})
		cv("type mismatch keeps to", func() {
			to, cs := mergeJSON(`{"a":1}`, `{"a":"x"}`)
			soJSON(to, `{"a":1}`)
			so(pathStrings(cs), resemble, []string{"a"})
		})
	})

	cv("null fill", func() {
		cv("root null", func() {
			to := jsonvalue.NewNull()
			from := jsonvalue.MustUnmarshalString(`{"a":1}`)
			cs := beta.Merge(&to, from)
			soJSON(to, `{"a":1}`)
			so(len(cs), eq, 0)
		})
		cv("nested null", func() {
			to, cs := mergeJSON(`{"n":null}`, `{"n":1}`)
			soJSON(to, `{"n":1}`)
			so(len(cs), eq, 0)
		})
		cv("nil to pointer is filled", func() {
			var to *jsonvalue.V
			from := jsonvalue.MustUnmarshalString(`{"a":1}`)
			cs := beta.Merge(&to, from)
			soJSON(to, `{"a":1}`)
			so(len(cs), eq, 0)
		})
		cv("from null against object conflicts", func() {
			to, cs := mergeJSON(`{"a":1}`, `null`)
			soJSON(to, `{"a":1}`)
			so(pathStrings(cs), resemble, []string{""})
		})
	})

	cv("string", func() {
		cv("equal strings", func() {
			to := jsonvalue.NewString("hello")
			from := jsonvalue.NewString("hello")
			cs := beta.Merge(&to, from)
			so(to.String(), eq, "hello")
			so(len(cs), eq, 0)
		})
		cv("unequal strings conflict", func() {
			to := jsonvalue.NewString("hello")
			from := jsonvalue.NewString("world")
			cs := beta.Merge(&to, from)
			so(to.String(), eq, "hello")
			so(pathStrings(cs), resemble, []string{""})
		})
		cv("concat at root", func() {
			to := jsonvalue.NewString("hello")
			from := jsonvalue.NewString("world")
			cs := beta.Merge(&to, from, beta.WithMergeConcatStrings())
			so(to.String(), eq, "helloworld")
			so(len(cs), eq, 0)
		})
		cv("concat nested", func() {
			to, cs := mergeJSON(`{"s":"hello"}`, `{"s":"world"}`, beta.WithMergeConcatStrings())
			soJSON(to, `{"s":"helloworld"}`)
			so(len(cs), eq, 0)
		})
		cv("concat equal strings doubles them", func() {
			to := jsonvalue.NewString("hello")
			from := jsonvalue.NewString("hello")
			cs := beta.Merge(&to, from, beta.WithMergeConcatStrings())
			so(to.String(), eq, "hellohello")
			so(len(cs), eq, 0)
		})
		cv("nested without concat conflicts", func() {
			to, cs := mergeJSON(`{"k":"a"}`, `{"k":"b"}`)
			soJSON(to, `{"k":"a"}`)
			so(pathStrings(cs), resemble, []string{"k"})
		})
	})

	cv("boolean and number", func() {
		cv("equal bool", func() {
			to := jsonvalue.NewBool(true)
			from := jsonvalue.NewBool(true)
			cs := beta.Merge(&to, from)
			so(to.Bool(), isTrue)
			so(len(cs), eq, 0)
		})
		cv("unequal bool", func() {
			to := jsonvalue.NewBool(true)
			from := jsonvalue.NewBool(false)
			cs := beta.Merge(&to, from)
			so(to.Bool(), isTrue)
			so(pathStrings(cs), resemble, []string{""})
		})
		cv("unequal number", func() {
			to := jsonvalue.NewInt(1)
			from := jsonvalue.NewInt(2)
			cs := beta.Merge(&to, from)
			so(to.Int(), eq, 1)
			so(pathStrings(cs), resemble, []string{""})
		})
	})

	cv("array default is whole-array conflict", func() {
		cv("from longer does not hang or drop silently", func() {
			to, cs := mergeJSON(`[1,2]`, `[3,4,5]`)
			soJSON(to, `[1,2]`)
			so(pathStrings(cs), resemble, []string{""})
		})
		cv("to longer does not hang", func() {
			to, cs := mergeJSON(`[1,2,3]`, `[4]`)
			soJSON(to, `[1,2,3]`)
			so(pathStrings(cs), resemble, []string{""})
		})
		cv("nested array", func() {
			to, cs := mergeJSON(`{"a":[1]}`, `{"a":[2,3]}`)
			soJSON(to, `{"a":[1]}`)
			so(pathStrings(cs), resemble, []string{"a"})
		})
		cv("equal arrays", func() {
			to, cs := mergeJSON(`[1,2]`, `[1,2]`)
			soJSON(to, `[1,2]`)
			so(len(cs), eq, 0)
		})
	})

	cv("array concat", func() {
		cv("root", func() {
			to, cs := mergeJSON(`[1,2]`, `[3,4,5]`, beta.WithMergeConcatArrays())
			soJSON(to, `[1,2,3,4,5]`)
			so(len(cs), eq, 0)
		})
		cv("nested", func() {
			to, cs := mergeJSON(`{"a":[1]}`, `{"a":[2,3]}`, beta.WithMergeConcatArrays())
			soJSON(to, `{"a":[1,2,3]}`)
			so(len(cs), eq, 0)
		})
	})

	cv("array ignore", func() {
		cv("keeps to array and overrides concat", func() {
			to, cs := mergeJSON(`[1]`, `[2]`, beta.WithMergeIgnoreArrays(), beta.WithMergeConcatArrays())
			soJSON(to, `[1]`)
			so(len(cs), eq, 0)
		})
		cv("nested keep array add other keys", func() {
			to, cs := mergeJSON(`{"a":[1],"b":0}`, `{"a":[9],"b":0,"c":3}`, beta.WithMergeIgnoreArrays())
			soJSON(to, `{"a":[1],"b":0,"c":3}`)
			so(len(cs), eq, 0)
		})
	})

	cv("array by index", func() {
		cv("from longer appends extra", func() {
			to, cs := mergeJSON(`[1,2]`, `[3,4,5]`, beta.WithMergeArraysByIndex())
			soJSON(to, `[1,2,5]`)
			so(pathStrings(cs), resemble, []string{"[0]", "[1]"})
		})
		cv("to longer keeps extra and does not hang", func() {
			to, cs := mergeJSON(`[1,2,3]`, `[4]`, beta.WithMergeArraysByIndex())
			soJSON(to, `[1,2,3]`)
			so(pathStrings(cs), resemble, []string{"[0]"})
		})
		cv("merge objects inside array", func() {
			to, cs := mergeJSON(`[{"a":1},{"b":2}]`, `[{"c":3}]`, beta.WithMergeArraysByIndex())
			soJSON(to, `[{"a":1,"c":3},{"b":2}]`)
			so(len(cs), eq, 0)
		})
		cv("nested", func() {
			to, cs := mergeJSON(`{"a":[{"x":1}]}`, `{"a":[{"y":2},3]}`, beta.WithMergeArraysByIndex())
			soJSON(to, `{"a":[{"x":1,"y":2},3]}`)
			so(len(cs), eq, 0)
		})
		cv("concat overrides by index", func() {
			to, cs := mergeJSON(`[1]`, `[2]`, beta.WithMergeConcatArrays(), beta.WithMergeArraysByIndex())
			soJSON(to, `[1,2]`)
			so(len(cs), eq, 0)
		})
	})

	cv("nil parameters", func() {
		cv("nil to double-pointer", func() {
			cs := beta.Merge(nil, jsonvalue.NewInt(1))
			so(pathStrings(cs), resemble, []string{""})
		})
		cv("nil from is no-op", func() {
			to := jsonvalue.MustUnmarshalString(`{"a":1}`)
			cs := beta.Merge(&to, nil)
			soJSON(to, `{"a":1}`)
			so(len(cs), eq, 0)
		})
	})

	cv("deep copy isolation", func() {
		from := jsonvalue.MustUnmarshalString(`{"shared":[1]}`)
		to := jsonvalue.NewNull()
		_ = beta.Merge(&to, from)
		from.MustSetInt(99).At("shared", 0)
		so(from.MustGet("shared", 0).Int(), eq, 99)
		so(to.MustGet("shared", 0).Int(), eq, 1)

		toArr := jsonvalue.MustUnmarshalString(`[1]`)
		fromArr := jsonvalue.MustUnmarshalString(`[{"z":1}]`)
		_ = beta.Merge(&toArr, fromArr, beta.WithMergeConcatArrays())
		fromArr.MustSetInt(99).At(0, "z")
		so(fromArr.MustGet(0, "z").Int(), eq, 99)
		so(toArr.MustGet(1, "z").Int(), eq, 1)
	})

	cv("self merge", func() {
		to := jsonvalue.MustUnmarshalString(`{"a":{"x":1},"b":[1,2]}`)
		cs := beta.Merge(&to, to)
		soJSON(to, `{"a":{"x":1},"b":[1,2]}`)
		so(len(cs), eq, 0)
	})

	cv("override when conflict", func() {
		override := beta.WithMergeOverrideWhenConflict()

		cv("basic types at root", func() {
			cv("string", func() {
				to := jsonvalue.NewString("hello")
				cs := beta.Merge(&to, jsonvalue.NewString("world"), override)
				so(to.String(), eq, "world")
				so(pathStrings(cs), resemble, []string{""})
			})
			cv("number", func() {
				to := jsonvalue.NewInt(1)
				cs := beta.Merge(&to, jsonvalue.NewInt(2), override)
				so(to.Int(), eq, 2)
				so(pathStrings(cs), resemble, []string{""})
			})
			cv("boolean", func() {
				to := jsonvalue.NewBool(true)
				cs := beta.Merge(&to, jsonvalue.NewBool(false), override)
				so(to.Bool(), isFalse)
				so(pathStrings(cs), resemble, []string{""})
			})
		})

		cv("basic types nested", func() {
			cv("number", func() {
				to, cs := mergeJSON(`{"a":1}`, `{"a":2}`, override)
				soJSON(to, `{"a":2}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("boolean", func() {
				to, cs := mergeJSON(`{"a":true}`, `{"a":false}`, override)
				soJSON(to, `{"a":false}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("string", func() {
				to, cs := mergeJSON(`{"a":"x"}`, `{"a":"y"}`, override)
				soJSON(to, `{"a":"y"}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("type mismatch between basic types", func() {
				to, cs := mergeJSON(`{"a":1}`, `{"a":"y"}`, override)
				soJSON(to, `{"a":"y"}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("mixed conflicts report every path", func() {
				to, cs := mergeJSON(
					`{"a":1,"b":{"x":1,"y":2}}`, `{"a":2,"b":{"x":1,"y":3,"z":4}}`, override)
				soJSON(to, `{"a":2,"b":{"x":1,"y":3,"z":4}}`)
				so(pathStrings(cs), resemble, []string{"a", "b.y"})
			})
		})

		cv("equal values are still not conflicts", func() {
			to, cs := mergeJSON(`{"a":1,"b":"x","c":true}`, `{"a":1,"b":"x","c":true}`, override)
			soJSON(to, `{"a":1,"b":"x","c":true}`)
			so(len(cs), eq, 0)
		})

		cv("null", func() {
			cv("from null overrides a basic type", func() {
				to, cs := mergeJSON(`{"a":1}`, `{"a":null}`, override)
				soJSON(to, `{"a":null}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("from null overrides an object", func() {
				to, cs := mergeJSON(`{"a":{"b":1}}`, `{"a":null}`, override)
				soJSON(to, `{"a":null}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("to null is still filled without conflict", func() {
				to, cs := mergeJSON(`{"n":null}`, `{"n":1}`, override)
				soJSON(to, `{"n":1}`)
				so(len(cs), eq, 0)
			})
		})

		cv("object and array in from are never overridden", func() {
			cv("object over basic type", func() {
				to, cs := mergeJSON(`{"a":1}`, `{"a":{"b":2}}`, override)
				soJSON(to, `{"a":1}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("array over basic type", func() {
				to, cs := mergeJSON(`{"a":1}`, `{"a":[2]}`, override)
				soJSON(to, `{"a":1}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
			cv("whole array conflict at root", func() {
				to, cs := mergeJSON(`[1,2]`, `[3]`, override)
				soJSON(to, `[1,2]`)
				so(pathStrings(cs), resemble, []string{""})
			})
			cv("whole array conflict nested", func() {
				to, cs := mergeJSON(`{"a":[1,2]}`, `{"a":[3]}`, override)
				soJSON(to, `{"a":[1,2]}`)
				so(pathStrings(cs), resemble, []string{"a"})
			})
		})

		cv("combined with array options", func() {
			cv("by index overrides scalar elements", func() {
				to, cs := mergeJSON(`[1,2]`, `[3,4,5]`, override, beta.WithMergeArraysByIndex())
				soJSON(to, `[3,4,5]`)
				so(pathStrings(cs), resemble, []string{"[0]", "[1]"})
			})
			cv("by index keeps extra elements of to", func() {
				to, cs := mergeJSON(`[1,2,3]`, `[4]`, override, beta.WithMergeArraysByIndex())
				soJSON(to, `[4,2,3]`)
				so(pathStrings(cs), resemble, []string{"[0]"})
			})
			cv("concat arrays wins over override", func() {
				to, cs := mergeJSON(`[1]`, `[2]`, override, beta.WithMergeConcatArrays())
				soJSON(to, `[1,2]`)
				so(len(cs), eq, 0)
			})
			cv("ignore arrays wins over override", func() {
				to, cs := mergeJSON(`{"a":[1],"b":1}`, `{"a":[9],"b":2}`,
					override, beta.WithMergeIgnoreArrays())
				soJSON(to, `{"a":[1],"b":2}`)
				so(pathStrings(cs), resemble, []string{"b"})
			})
		})

		cv("concat strings wins over override", func() {
			to, cs := mergeJSON(`{"s":"a"}`, `{"s":"b"}`, override, beta.WithMergeConcatStrings())
			soJSON(to, `{"s":"ab"}`)
			so(len(cs), eq, 0)
		})

		cv("overridden value is deep copied", func() {
			to := jsonvalue.MustUnmarshalString(`{"a":1}`)
			from := jsonvalue.MustUnmarshalString(`{"a":2}`)
			cs := beta.Merge(&to, from, override)
			so(pathStrings(cs), resemble, []string{"a"})
			from.MustSetInt(99).At("a")
			so(from.MustGet("a").Int(), eq, 99)
			so(to.MustGet("a").Int(), eq, 2)
		})

		cv("not specified means to always wins", func() {
			to, cs := mergeJSON(`{"a":1,"b":true,"c":"x"}`, `{"a":2,"b":false,"c":"y"}`)
			soJSON(to, `{"a":1,"b":true,"c":"x"}`)
			so(pathStrings(cs), resemble, []string{"a", "b", "c"})
		})
	})
}

func mergeJSON(toStr, fromStr string, opts ...beta.MergeOption) (*jsonvalue.V, []jsonvalue.Path) {
	to := jsonvalue.MustUnmarshalString(toStr)
	from := jsonvalue.MustUnmarshalString(fromStr)
	cs := beta.Merge(&to, from, opts...)
	return to, cs
}

func soJSON(got *jsonvalue.V, want string) {
	wantV := jsonvalue.MustUnmarshalString(want)
	so(got.Equal(wantV), isTrue)
}

func pathStrings(ps []jsonvalue.Path) []string {
	res := make([]string, 0, len(ps))
	for _, p := range ps {
		res = append(res, p.String())
	}
	return res
}
