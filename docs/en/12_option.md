<font size=6>Additional Options</font>

[Prev Page](./11_comparation.md) | [Contents](./README.md) | [Next Page](./13_beta.md)

---

This section describes the usage of additional options in detail. Or you can refer to Section [Special Application Scenarios](./10_scenarios.md).

---

- [Overview](#overview)
- [Ignoring null values](#ignoring-null-values)
- [Visible Indentation](#visible-indentation)
- [Specify the Sequence of Keys in Object Value](#specify-the-sequence-of-keys-in-object-value)
  - [Sequence of When the Keys Are Set](#sequence-of-when-the-keys-are-set)
  - [Use Callback to Specify Sequence](#use-callback-to-specify-sequence)
  - [Use Alphabetical Sequence](#use-alphabetical-sequence)
  - [Use Pre-defined \[\]string Identifying Key Sequence](#use-pre-defined-string-identifying-key-sequence)
- [Handling NaN](#handling-nan)
  - [NaN to Another Floating Value](#nan-to-another-floating-value)
  - [NaN to Null Type](#nan-to-null-type)
  - [NaN to String](#nan-to-string)
- [Handling +/-Inf](#handling--inf)
  - [+/-Inf to Another Floating Value](#-inf-to-another-floating-value)
  - [+/-Inf to Null Type](#-inf-to-null-type)
  - [+/-Inf to String](#-inf-to-string)
- [Escaping Non-essential Characters](#escaping-non-essential-characters)
  - [SetEscapeHTML Options](#setescapehtml-options)
  - [The Slash `/`](#the-slash-)
  - [Escaping Unicode Greater than 0x7F](#escaping-unicode-greater-than-0x7f)
- [Ignoring Tag omitempty of A Struct](#ignoring-tag-omitempty-of-a-struct)

---

## Overview

Let us take a look back at the `Marshal` methods:

```go
func (v *V) Marshal          (opts ...Option) (b []byte, err error)
func (v *V) MarshalString    (opts ...Option) (s string, err error)
func (v *V) MustMarshal      (opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

You can see that every method supports additional `opts ...Option`, specifying additional options for serializing jsonvalue data.

For a simple example: you can ignore all null values with the following option:

```go
v := jsonvalue.NewObject()
v.SetNull().At("null")
fmt.Println(v.MustMarshalString())
fmt.Println(v.MustMarshalString(jsonvalue.OptOmitNull(true)))
```

Outputs:

```json
{"null":null}
{}
```

Currently, only the `OptIgnoreOmitempty()` option is designed for `Import()`, while all others are for marshaling.

---

## Ignoring null values

Already mentioned above.

---

## Visible Indentation

This is like the `json.MarshalIndent` function in `encoding/json`. Taking the previous example, you can add the following option:

```go
s := v.MustMarshalString(jsonvalue.OptIndent("", "  "))
fmt.Println(s)
```

Which will output:

```json
{
  "null": null
}
```

---

## Specify the Sequence of Keys in Object Value

In general, specifying the sequence of key-values is unnecessary and a waste of CPU time. But there are some special situations that make key sequence important:

- Hash checksum mentioned previously.
- Quickly finding some specific key-value pairs when debugging.
- Non-standard use of JSON, which depends on the order of keys.

This subsection will tell you how to specify the sequence of keys when marshaling.

### Sequence of When the Keys Are Set

Please refer to [Special Application Scenarios](./10_scenarios.md).

### Use Callback to Specify Sequence

```go
func OptKeySequenceWithLessFunc(f MarshalLessFunc) Option
```

`MarshalLessFunc` is a callback function with the following prototype:

```go
type MarshalLessFunc func(nilableParent *ParentInfo, key1, key2 string, v1, v2 *V) bool
```

This function acts like `sort.Sort`. The definitions of the input parameters are:

- `nilableParent` - parent keys of the current value.
- `key1`, `v1` - first K-V pair to rearrange.
- `key2`, `v2` - second K-V pair to rearrange.

It returns whether v1 should be ahead of v2. It acts like the `Less` function in package `sort`.

### Use Alphabetical Sequence

```go
func OptDefaultStringSequence() Option
```

### Use Pre-defined []string Identifying Key Sequence

It is quite complicated to use the "less" callback. You can simply pass a `[]string`, then jsonvalue will arrange the key sequence according to the given string slice. If there are any keys not specified in it, they will be put after all specified ones.

```go
func OptKeySequence(seq []string) Option
```

If both `OptKeySequence` and `OptKeySequenceWithLessFunc` are specified, `OptKeySequenceWithLessFunc` will be used with priority.

--- 

## Handling NaN

In standard JSON, NaN (not a number) and +/-Inf (plus / minus infinity) are illegal. But in some cases, we have to handle those values.

By default, jsonvalue will raise an error when handling those values. But some non-error conversion options are also provided. 

Let us see NaN first:

### NaN to Another Floating Value

```go
func OptFloatNaNToFloat(f float64) Option
```

By specifying another floating number, all NaN values will be replaced with it. You cannot specify the replacement as NaN or +/-Inf.

### NaN to Null Type

```go
func OptFloatNaNToNull() Option
```

Replace NaN with JSON null. Please be advised that this option is NOT affected by `OptOmitNull`. It will ALWAYS convert NaN to `null` with this option.

### NaN to String

```go
func OptFloatNaNToString   (s string) Option
func OptFloatNaNToStringNaN() Option
```

`OptFloatNaNToStringNaN()` is equivalent to `OptFloatNaNToString("NaN")`.

---

## Handling +/-Inf

The processing of +/-Inf is similar to NaN:

### +/-Inf to Another Floating Value

```go
func OptFloatInfToFloat(f float64) Option
```

+Inf will be replaced with `f`, while -Inf will be replaced with `-f`. You cannot specify the replacement as NaN or +/-Inf.

### +/-Inf to Null Type

```go
func OptFloatInfToNull() Option
```

Similarly, this option is NOT affected by `OptOmitNull`.

### +/-Inf to String

```go
func OptFloatInfToString   (positiveInf, negativeInf string) Option
func OptFloatInfToStringInf() Option
```

Specify two strings as replacements for +/-Inf. If the given string is empty, replace with the following priority:

- For +Inf, if the given string is empty, then replace with `"+Inf"`
- For -Inf, if the given string is empty, first find the +Inf configuration. If present, remove the `+` prefix and add a `-`.
- If the configuration for both +/-Inf is empty, -Inf will be replaced with `"-Inf"`

`OptFloatInfToStringInf()` is equivalent to `OptFloatInfToString("+Inf", "-Inf")`

## Escaping Non-essential Characters

In the JSON standard, there are several characters that need to be escaped:

1. Some important formatting or reserved characters
2. Unicode characters greater than 127

However, not all JSON encoders follow the full escaping rules. This section tells you how to specify special escaping rules.

### SetEscapeHTML Options

According to the JSON standard, `&`, `<` and `>` should be escaped to `\u00XX`. However, it is safe not to escape these three characters in practice. By default, jsonvalue will escape them all. However, you can disable the escaping of these characters by using this option:

```go
func OptEscapeHTML(on bool) Option
```

Pass `true` to enable the escaping while `false` disables escaping them.

This option is quite like the [`SetEscapeHTML`](https://pkg.go.dev/encoding/json#Encoder.SetEscapeHTML) method in `encoding/json`.

### The Slash `/`

According to the JSON standard, slash `/` should be escaped. But actually it is OK not to escape it. By default, jsonvalue escapes the slash, but you can use this switch option:

```go
func OptEscapeSlash(on bool) Option
```

Pass `false` to not escape slash. By default, it is `true`.

### Escaping Unicode Greater than 0x7F

```go
func OptUTF8() Option
```

By default, jsonvalue will escape all unicode values greater than 0x7F to `\uXXXX` format (UTF-16). This will avoid almost all encoding problems. However, it may be a waste of network traffic if most of your payload is unicode. In this case, you may consider using UTF-8 encoding. This is what `encoding/json` does.

## Ignoring Tag omitempty of A Struct

This is a rare feature. Please refer to the subsection "Ignoring `omitempty` JSON Tag of A Struct" in [Special Application Scenarios](./10_scenarios.md).
