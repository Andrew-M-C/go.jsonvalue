<font size=6>New Features in v1.2.x</font>

[Prev Page](./16_1_12_new_feature.md) | [Contents](./README.md)

---

- [Change to NewFloat64 and NewFloat32](#change-to-newfloat64-and-newfloat32)
- [Passing options with Function Typed Options](#passing-options-with-function-typed-options)
- [Handling NaN and +/-Inf](#handling-nan-and--inf)
- [Supporting Overriding Default Options](#supporting-overriding-default-options)

---

## Change to NewFloat64 and NewFloat32

This version introduces the first non-backward-compatible feature: `NewFloat64` and `NewFloat32` functions.

Before v1.2.0, these functions were like:

```go
func NewFloat64(f float64, prec int) *V
func NewFloat32(f float32, prec int) *V
```

The parameter `prec` is used in calling `strconv.FormatFloat`, while the corresponding `fmt` parameter is set to `'f'`.

However, [Issue #8](https://github.com/Andrew-M-C/go.jsonvalue/issues/8) made me realize that the `'f'` format is not the best format to describe every floating-point number. In some cases, it is better to use scientific notation. Therefore, it is necessary to expose the parameters of `strconv.FormatFloat`.

After careful consideration, I decided to release a non-backward-compatible version. The changes are:

- Remove `prec` from `NewFloat64` and `NewFloat32`. The formatting behavior remains the same as what `encoding/json` does.
- Add `NewFloat64f` and `NewFloat32f`, with the same option parameters as `strconv.FormatFloat`, including floating value, format, and precision.
  - **NOTE**: According to the JSON standard, only `f`, `E`, `e`, `G`, `g` formats are supported. If illegal values are received, jsonvalue will use `g` instead.

---

## Passing options with Function Typed Options

Before v1.2.0, additional options were defined in a struct. From v1.2.0, please use `OptXxx` functions instead.

---

## Handling NaN and +/-Inf

Please refer to "Handling +/-Inf" in Section [Additional Options](./12_option.md).

---

## Supporting Overriding Default Options

```go
func SetDefaultMarshalOptions(opts ...Option)
func ResetDefaultMarshalOptions()
```

In jsonvalue, the default serialization options are:

- Escape ALL characters explicitly defined by the JSON standard including `", /, \, <, >, &`, vertical/horizontal tabs, carriage return, and backspace.
- Escape ALL characters greater than 127.

But according to feedback, those strict rules are not needed in most cases. Therefore, if you are sure about what serialization rules you need, you can invoke `SetDefaultMarshalOptions` to override jsonvalue's default options after your process is initialized.
