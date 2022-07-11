package beta

import (
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// IsSubsetOf is the invert version of HasSubset.
//
// IsSubsetOf 是 HasSubset 函数的反向操作。
func IsSubsetOf(v, macroSet *jsonvalue.V) bool {
	return HasSubset(macroSet, v)
}

// HasSubset identidies whether has a subset. This only takes effect to object
// and array types.
//
// HasSubset 表示是否包含某个子集。只对 object 和 array 类型有效, 其他类型则需完全相等时,
// 才返回 true
func HasSubset(v, sub *jsonvalue.V) bool {
	if v == nil || sub == nil {
		return false
	}
	if v.ValueType() != sub.ValueType() {
		return false
	}

	switch v.ValueType() {
	default:
		return v.Equal(sub)
	case jsonvalue.Object:
		return objectHasSubset(v, sub)
	case jsonvalue.Array:
		return arrayHasSubset(v, sub)
	}
}

func objectHasSubset(v, sub *jsonvalue.V) bool {
	res := true

	sub.RangeObjects(func(k string, subV *jsonvalue.V) bool {
		vv, err := v.Get(k)
		if err != nil {
			res = false
			return false
		}
		if !HasSubset(vv, subV) {
			res = false
			return false
		}
		return true
	})

	return res
}

func arrayHasSubset(v, sub *jsonvalue.V) bool {
	lenV := v.Len()
	lenSub := sub.Len()

	if lenSub == 0 {
		return true
	} else if lenV == lenSub {
		return sub.Equal(v)
	} else if lenSub > lenV {
		return false
	}

	subEqual := func(start int) bool {
		remain := lenV - start
		if remain < lenSub {
			return false
		}

		res := true

		sub.RangeArray(func(i int, subV *jsonvalue.V) bool {
			vv := v.MustGet(start + i)
			if !vv.Equal(subV) {
				res = false
				return false
			}
			return true
		})

		return res
	}

	for i := 0; i <= lenV-lenSub; i++ {
		if subEqual(i) {
			return true
		}
	}
	return false
}
