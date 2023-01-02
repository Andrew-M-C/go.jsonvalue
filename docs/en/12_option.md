
<font size=6>Additional Options</font>

[Prev Page](./11_comparation.md) | [Contents](./README.md) | [Next Page](./13_beta.md)

---

This sections describe the usage of additional options in detail. Or you can refer to Section [Special Application Scenarios](./10_scenarios.md).

---

- [Overview](#overview)
- [Ignoring null values](#ignoring-null-values)
- [Visible Indention](#visible-indention)
- [Specify the Sequence of Keys in Object Value.](#specify-the-sequence-of-keys-in-object-value)
  - [Sequence of When the Keys Are Set](#sequence-of-when-the-keys-are-set)
  - [Use Callback to Specify Sequence](#use-callback-to-specify-sequence)
  - [Use Alphabet Sequence](#use-alphabet-sequence)
  - [Use Pre-defined \[\]string Identifying Key Sequence](#use-pre-defined-string-identifying-key-sequence)
- [Handing NaN](#handing-nan)
  - [NaN to Another Floating Value](#nan-to-another-floating-value)
  - [NaN to Null Type](#nan-to-null-type)
  - [NaN to String](#nan-to-string)
- [Handing +/-Inf](#handing--inf)
  - [+/-Inf to Another Floating Value](#-inf-to-another-floating-value)
  - [+/-Inf to Null Type](#-inf-to-null-type)
  - [+/-Inf to String](#-inf-to-string)
- [Escaping Un-essential Characters](#escaping-un-essential-characters)
  - [SetEscapeHTML Options](#setescapehtml-options)
  - [The Slash `/`](#the-slash-)
  - [Escaping Unicode Greater than 0x7F](#escaping-unicode-greater-than-0x7f)
- [Ignoring Tag omitempty of A Struct](#ignoring-tag-omitempty-of-a-struct)

---

## Overview

Let us take a look back to `Marshal` methods:

```go
func (v *V) Marshal          (opts ...Option) (b []byte, err error)
func (v *V) MarshalString    (opts ...Option) (s string, err error)
func (v *V) MustMarshal      (opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

You can see that every methods supports additional `opts ...Option`, identifying additional options in serializing jsonvalue data.

For a simple example: you can ignore all null value with following option:

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

Currently, only `OptIgnoreOmitempty()` options os designed for `Import()`, while all others for marshal.

---

## Ignoring null values

Already mentioned above.

---

## Visible Indention

This is like the `json.MarshalIndent` function in `encoding/json`. Take the previous example, you can added the following option:

```go
s := v.MustMarshalString(jsonvalue.OptIndent("", "  "))
fmt.Println(s)
```

Which will outputs:

```json
{
  "null": null
}
```

---

## Specify the Sequence of Keys in Object Value.

In general, specifying the sequence of key-values are unnecessary and a waste of CPU time. But there are some special situations make key sequence important:

- Hash checksum mentioned previously.
- Quickly finding some specific key-value pairs when debugging.
- Nonstandard use of JSON, which depends on the order if keys.

This sub section will tell you how to specify the sequence of keys when marshaling.

### Sequence of When the Keys Are Set

Please refer to [Special Application Scenarios](./10_scenarios.md).

### Use Callback to Specify Sequence

```go
func OptKeySequenceWithLessFunc(f MarshalLessFunc) Option
```

`MarshalLessFunc` is a callback function, prototype:

```go
type MarshalLessFunc func(nilableParent *ParentInfo, key1, key2 string, v1, v2 *V) bool
```

This function acts like `sort.Sort`. The definitions of input parameters are:

- `nilableParent` - parent keys of current value.
- `key1`, `v1` - first K-V to rearrange.
- `key2`, `v2` - second K-V to rearrange.

return whether v1 should be ahead of v2. It acts like `Less` function in Package `sort`.

### Use Alphabet Sequence

```go
func OptDefaultStringSequence() Option
```

### Use Pre-defined []string Identifying Key Sequence

It is quite complicated to use "less" callback. You can simple pass a `[]string`, then jsonvalue will arrange the key sequence according the given string slice. If there is any keys not specified in it, it will be put after all specified ones.

```go
func OptKeySequence(seq []string) Option
```

If both `OptKeySequence` and `OptKeySequenceWithLessFunc` are specified, `OptKeySequenceWithLessFunc` will be used in priority.

--- 

## Handing NaN

In standard JSON, NaN (not a number) and +/-Inf (plus / minus infinity) are illegal. But in some cases, we had to handle those values.

By default, jsonvalue will raise an error when handling those values. But some non-error conversion options are also provided. 

Let us see NaN first:

### NaN to Another Floating Value

```go
func OptFloatNaNToFloat(f float64) Option
```

Specifying another floating number, all NaN will be replaced with it. You cannot specify the replacement as NaN or +/-Inf.

### NaN to Null Type

```go
func OptFloatNaNToNull() Option
```

Replace NaN with JSON null. Please be advised that this option does NOT affected by `OptOmitNull`. It will ALWAYS convert NaN to `null` with option.

### NaN to String

```go
func OptFloatNaNToString   (s string) Option
func OptFloatNaNToStringNaN() Option
```

`OptFloatNaNToStringNaN()` is equivalent to `OptFloatNaNToString("NaN")`。

---

## Handing +/-Inf

The processing with +/-Inf are similar to NaN:

### +/-Inf to Another Floating Value

```go
func OptFloatInfToFloat(f float64) Option
```

+Inf will be replaced to `f`, while -Inf to `-f`. You cannot specify the replacement as NaN or +/-Inf.

### +/-Inf to Null Type

```go
func OptFloatInfToNull() Option
```

Similarly, this option does NOT affected by `OptOmitNull`.

### +/-Inf to String

```go
func OptFloatInfToString   (positiveInf, negativeInf string) Option
func OptFloatInfToStringInf() Option
```

Specifying two strings as replacement for +/-Inf. If the given string is empty, replace in following priority:

- As for +Inf, if given string is empty, then replace as `"+Inf"`
- As for -Inf, if given string is empty, first find the +Inf configuration. If present, remove the `+` prefix and add a `-`.
- If configuration of both +/-Inf are empty, -Inf will be replaced with `"-Inf"`

`OptFloatInfToStringInf()` is equivalent to `OptFloatInfToString("+Inf", "-Inf")`

## Escaping Un-essential Characters

In JSON standard, there are several characters need to be escaped:

1. Some important formatting or reserved characters
2. Unicodes greater than 127

However, not all JSON encoder follows the full escaping rules. This sections tells how to specify special escaping rules.

### SetEscapeHTML Options

According to JSON standard, `&`, `<` and `>` should be escaped to `\u00XX`. However it is safe not escaping these three characters in practical. By default, jsonvalue will escape them all. However, you can disable the escaping to these characters by using option:

```go
func OptEscapeHTML(on bool) Option
```

Passing `true` to enable the escaping while `false` not escaping them.

This options is quite like [`SetEscapeHTML`](https://pkg.go.dev/encoding/json#Encoder.SetEscapeHTML) method In `encoding/json`.

### The Slash `/`

According to JSON standard, slash `/` should be escaped. But actually it is OK to not escaping it. By default, jsonvalue escapes the slash, but you can use this switch option:

```go
func OptEscapeSlash(on bool) Option
```

Passing `false` to not escape slash. By default, it is `true`.

### Escaping Unicode Greater than 0x7F

```go
func OptUTF8() Option
```

By default, jsonvalue will escape all unicode values greater than 0x7F to `\uXXXX` format (UTF-16). This will avoid almost all encoding problem. However, it may be a waste of network traffic if most of your payload are unicode. In this case, you may consider using UTF-8 encoding. This is what `encoding/json` does.

## Ignoring Tag omitempty of A Struct

This is a rare feature. Please refer to sub section "Ignoring `omitempty` JSON Tag of A Struct" in [Special Application Scenarios](./10_scenarios.md).
