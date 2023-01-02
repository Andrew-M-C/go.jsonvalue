
<font size=6>New Features in v1.2.x</font>

[Prev Page](./14_1_13_new_feature.md) | [Contents](./README.md)

---

- [Change to NewFloat64 and NewFloat32](#change-to-newfloat64-and-newfloat32)
- [Passing options with Function Typed Options](#passing-options-with-function-typed-options)
- [Handling NaN and +/-Inf](#handling-nan-and--inf)
- [Supporting Overriding Default Option](#supporting-overriding-default-option)

---

## Change to NewFloat64 and NewFloat32

This version introduces first non-backward-compatible feature: `NewFloat64` and `NewFloat32` functions。

Before v1.2.0, these functions are like:

```go
func NewFloat64(f float64, prec int) *V
func NewFloat32(f float32, prec int) *V
```

The parameter `prec` in used in calling `strconv.FormatFloat`, while the corresponding `fmt` parameter is set to `'f'`.

However, [Issue #8](https://github.com/Andrew-M-C/go.jsonvalue/issues/8) made me realized that `'f'` format is not the best format to describe every floating number. In some cases, it is better to use scientific notation. Therefore, it is necessary to open parameters of `strconv.FormatFloat`

After considering carefully, I decided to release a non-backward-compatible version. The changes are:

- Remove `prec` in `NewFloat64` and `NewFloat32`. The formatting behavior keeps the same as what `encoding/json` does.
- Add `NewFloat64f` and `NewFloat32f`, with the same option parameters to `strconv.FormatFloat`, including floating value, format, precision.
  - **NOTE**: according to JSON standard, only `f`, `E`, `e`, `G`, `g` formats are supported. If receiving illegal values, jsonvalue will use `g` instead.

---

## Passing options with Function Typed Options

Before v1.2.0, additional options are defined in a struct. From v1.2.0, please use `OptXxx` functions instead.

---

## Handling NaN and +/-Inf

Please refer to "Handing +/-Inf" in Section [Additional Options](./12_option.md).

---

## Supporting Overriding Default Option

```go
func SetDefaultMarshalOptions(opts ...Option)
func ResetDefaultMarshalOptions()
```

In jsonvalue, the default serialize options are:

- Escape ALL characters explicitly defined by JSON standard including `", /, \, <, >, &`, vertical / horizontal tabs, returning character, backspace.
- Escape ALL characters greater than 127.

But according to feed back, those strict rules are not needed in most cases. Therefore, if you are sure of what serialize rules you need, you can invoke `SetDefaultMarshalOptions` to override jsonvalue's default options after your process initialized.
