
<font size=6>Value Conversions</font>

[Prev Page](./08_caseless.md) | [Contents](./README.md) | [Next Page](./10_scenarios.md)

---

- [Overview](#overview)
- [String to Number](#string-to-number)
- [String and Number to Boolean](#string-and-number-to-boolean)
- [Difference between `GetString` and `MustGet(...).String()`](#difference-between-getstring-and-mustgetstring)

---

## Overview

In `encoding/json`, there is a `tag` called "string", which means to convert current value (no matter what type it is) to a JSON string. For example:

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

It is clear that `num` is serialized to a string instead of number

Jsonvalue also supports reading such type of value from a string. Please read the following descriptions for detail.

---

## String to Number

All number typed `GetXxx` methods support parsing a number from string value.

But at the same time, if the target value is not a number typed JSON value, the `error` will not be nil. But various types of error may returned due to different situation:

- If target key is not found, returns `ErrNotFound`
- If target value is a number, the number type will be returned and with a `nil` err.
- If target value is neither a string or number, returns `ErrTypeNotMatch`
- If target value is a string, then parse the value of the string, and:
  - if the parsing succeed, returns the number with error `ErrTypeNotMatch`
  - if the parsing failed (illegal number), returns error `ErrParseNumberFromString`, including the detailed parsing description.

You can use `errors.Is()` to check the different type of errors. For example:

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

Similar with number, strings can hold a boolean.

A `true` will be returned in only one situation then converting string to boolean:

- The value of the string value is exactly the lower-cased string `"true"`.

In other cases, `false` will be returned for the conversion.

Number could also converted into boolean. As long as the number value does NOT equal to zero, `true` will be returned. Otherwise, `false`

---

## Difference between `GetString` and `MustGet(...).String()`

The logic between `GetXxx()` series methods and its corresponding `MustGet().Xxx()` function are mostly the same. Instead, the `Xxx()` methods does not return error.

For example:

```go
	const raw = `{"legal":"1.25E-1","illegal":"ABCD"}`
	v := jsonvalue.MustUnmarshalString(raw)
	f := v.MustGet("legal").Float64()
	fmt.Println("float:", f)
	// Output: float: 0.125
```

But there is one exception: `GetString` and `MustGet().String()`

`GetString` returns a string of a string-typed JSON value. If target is not found, or not a string, error will be returned.

But `String` is quite special. In on hand, it can get the string value of a string-typed JSON value. In the other hand, it implements the `fmt.Stringer` interface, which is invoked in `%v` keyword of `fmt` package.

Therefore the logic of `String` is a bit complicated:

- If current `*jsovalue.V` is string-typed, returns its string value.
- If number-typed, returns the string of the number.
- If null-typed ("null" type defined in JSON standard), returns `null`.
- If boolean-typed, returns `true` or `false`.
- If object-typed, returns a string with format of `{K1: V1, K2: V2}`, and no escaping, not JSON format.
- If array-typed, returns a string with format of `[v1, v2]`. Also, no escaping.
