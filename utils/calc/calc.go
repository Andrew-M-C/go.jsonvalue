// Package calc provides some utility functions for calculation.
//
// Calc 包提供一些计算相关的工具函数。
package calc

import (
	"fmt"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

var emptyValue = &jsonvalue.V{}

// ConcatStrings finds every string-typed sub value of a and b, and concatenates them into a new JSON value.
//
// ConcatStrings 找到 a 和 b 的所有 string 类型的子值，并将它们拼接成一个新的 JSON 值。
func ConcatStrings(a, b *jsonvalue.V) (*jsonvalue.V, error) {
	if a.ValueType() != b.ValueType() {
		return emptyValue, fmt.Errorf("%w: input parameters have different types", jsonvalue.ErrParameterError)
	}

	var res *jsonvalue.V

	switch a.ValueType() {
	case jsonvalue.String:
		// for string type, simply concatenate them
		return jsonvalue.NewString(a.String() + b.String()), nil

	case jsonvalue.Object:
		// for object type, iterate and concatenate each string-typed sub value
		res = jsonvalue.NewObject()
		// and go on

	case jsonvalue.Array:
		// for array type, iterate and concatenate each string-typed sub value
		res = jsonvalue.NewArray()
		// and go on

	default:
		return emptyValue, fmt.Errorf(
			"%w: only string, object and array type are supported", jsonvalue.ErrParameterError,
		)
	}

	// Walk a
	a.Walk(func(path jsonvalue.Path, v *jsonvalue.V) bool {
		// if not string type, simply set it
		if v.ValueType() != jsonvalue.String {
			res.MustSet(v).At(path)
			return true
		}

		// if string type, just concatenate if
		strB := b.MustGet(path).String()
		res.MustSet(v.String() + strB).At(path)
		return true
	})

	// Walk b
	b.Walk(func(path jsonvalue.Path, v *jsonvalue.V) bool {
		if _, err := res.Get(path); err != nil {
			res.MustSet(v).At(path) // value not exists, set it
		}
		return true
	})

	return res, nil
}
