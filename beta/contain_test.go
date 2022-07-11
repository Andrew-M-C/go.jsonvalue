package beta

import (
	"testing"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

func testHasSubsetIsSubsetOf(t *testing.T) {
	cv("general array", func() {
		v := jsonvalue.MustUnmarshalString(`[1,2,3,4]`)
		sub := jsonvalue.MustUnmarshalString(`[2,3]`)
		notSub := jsonvalue.MustUnmarshalString(`[1,3]`)

		so(HasSubset(v, sub), isTrue)
		so(HasSubset(v, notSub), isFalse)
		so(IsSubsetOf(sub, v), isTrue)
		so(IsSubsetOf(notSub, v), isFalse)
	})

	cv("general object", func() {
		v := jsonvalue.MustUnmarshalString(`{"num":1234,"str":"Hello"}`)
		sub := jsonvalue.MustUnmarshalString(`{"num":1234.00}`)
		notSub := jsonvalue.MustUnmarshalString(`{"num":1234,"str":"hello"}`)

		so(HasSubset(v, sub), isTrue)
		so(HasSubset(v, notSub), isFalse)
		so(IsSubsetOf(sub, v), isTrue)
		so(IsSubsetOf(notSub, v), isFalse)
	})
}
