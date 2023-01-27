package beta

import (
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// Contains tells whether a value has a subset. This only takes effect to object
// and array types, otherwise, this function acts like Equal().
//
// Contains 表示是否包含某个子集。只对 object 和 array 类型有效, 其他类型则需完全相等时,
// 才返回 true
func Contains(v *jsonvalue.V, sub interface{}, inPath ...interface{}) bool {
	if v == nil {
		return false
	}

	var err error
	subV, ok := sub.(*jsonvalue.V)
	if !ok {
		subV, err = jsonvalue.Import(sub)
		if err != nil {
			// fmt.Println("Import failed")
			return false
		}
	}
	if len(inPath) > 0 {
		v, err = v.Get(inPath[0], inPath[1:]...)
		if err != nil {
			// fmt.Println("Get failed:", err)
			return false
		}
	}

	if v.ValueType() != subV.ValueType() {
		// fmt.Println("type mismatch - ", v.ValueType(), subV.ValueType())
		return false
	}

	switch v.ValueType() {
	default:
		return v.Equal(subV)
	case jsonvalue.Object:
		return objectHasSubset(v, subV)
	case jsonvalue.Array:
		return arrayHasSubset(v, subV)
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
		if !Contains(vv, subV) {
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
