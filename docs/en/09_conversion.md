<font size=6>Value Conversions</font>

[Prev Page](./08_caseless.md) | [Contents](./README.md) | [Next Page](./10_scenarios.md)

---

- [Overview](#overview)
- [String to Number](#string-to-number)
- [String and Number to Boolean](#string-and-number-to-boolean)
- [Difference between `GetString` and `MustGet(...).String()`](#difference-between-getstring-and-mustgetstring)

---

## Overview

In `encoding/json`, there is a `tag` called "string", which means converting the current value (no matter what type it is) to a JSON string. For example:

```go
type st struct {
    Num int `json:"num,string"`
}

func main() {
    st := st{Num: 12345}
    b, _ := json.Marshal(&st)
    fmt.Println(string(b))
	// Output: {"num":"12345"}
}
```

It is clear that `num` is serialized as a string instead of a number.

Jsonvalue also supports reading this type of value from a string. Please read the following descriptions for details.

---

## String to Number

All number-typed `GetXxx` methods support parsing a number from a string value.

However, if the target value is not a number-typed JSON value, the `error` will not be nil. Various types of errors may be returned due to different situations:

- If the target key is not found, returns `ErrNotFound`
- If the target value is a number, the number type will be returned with a `nil` err.
- If the target value is neither a string nor a number, returns `ErrTypeNotMatch`
- If the target value is a string, then parse the value of the string, and:
  - if the parsing succeeds, returns the number with error `ErrTypeNotMatch`
  - if the parsing fails (illegal number), returns error `ErrParseNumberFromString`, including the detailed parsing description.

You can use `errors.Is()` to check the different types of errors. For example:

```go
	const raw = `{"legal":"1.25E-1","illegal":"ABCD"}`
	v := jsonvalue.MustUnmarshalString(raw)

	f, err := v.GetFloat64("legal")
	fmt.Println("01 - float:", f)
	fmt.Println("01 - err:", err)

	f, err = v.GetFloat64("illegal")
	fmt.Println("02 - float:", f)
	fmt.Println("02 - err:", err)
```

Output:

```
01 - float: 0.125
01 - err: not match given type
02 - float: 0
02 - err: failed to parse number from string: parsing number at index 0: zero string
```

---

## String and Number to Boolean

Similar to numbers, strings can hold a boolean value.

A `true` will be returned in only one situation when converting a string to boolean:

- The value of the string is exactly the lower-cased string `"true"`.

In other cases, `false` will be returned for the conversion.

Numbers can also be converted into boolean values. As long as the number value does NOT equal zero, `true` will be returned. Otherwise, `false`.

---

## Difference between `GetString` and `MustGet(...).String()`

The logic between `GetXxx()` series methods and their corresponding `MustGet().Xxx()` functions is mostly the same. However, the `Xxx()` methods do not return an error.

For example:

```go
	const raw = `{"legal":"1.25E-1","illegal":"ABCD"}`
	v := jsonvalue.MustUnmarshalString(raw)
	f := v.MustGet("legal").Float64()
	fmt.Println("float:", f)
	// Output: float: 0.125
```

But there is one exception: `GetString` and `MustGet().String()`

`GetString` returns a string of a string-typed JSON value. If the target is not found, or is not a string, an error will be returned.

But `String` is quite special. On one hand, it can get the string value of a string-typed JSON value. On the other hand, it implements the `fmt.Stringer` interface, which is invoked by the `%v` keyword of the `fmt` package.

Therefore, the logic of `String` is a bit complicated:

- If the current `*jsonvalue.V` is string-typed, returns its string value.
- If number-typed, returns the string representation of the number.
- If null-typed ("null" type defined in JSON standard), returns `null`.
- If boolean-typed, returns `true` or `false`.
- If object-typed, returns a string with the format `{K1: V1, K2: V2}`, with no escaping, not in JSON format.
- If array-typed, returns a string with the format `[v1, v2]`. Also, no escaping.
