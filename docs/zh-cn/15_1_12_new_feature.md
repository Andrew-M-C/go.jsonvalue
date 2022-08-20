
<font size=6>1.2.x 新特性</font>

[上一页](./14_1_13_new_feature.md) | [总目录](./README.md) | [下一页](./16_others.md)

---

- [NewFloat64 和 NewFloat32 功能变更](#newfloat64-和-newfloat32-功能变更)
- [使用函数形式配置序列化的额外参数](#使用函数形式配置序列化的额外参数)
- [支持处理 NaN 和 Inf 浮点值](#支持处理-nan-和-inf-浮点值)
- [支持覆盖默认序列化配置](#支持覆盖默认序列化配置)

---

## NewFloat64 和 NewFloat32 功能变更

本版本中，jsonvalue 将迎来第一个不向后兼容（backward compatible）的特性，那就是 `NewFloat64` 和 `NewFloat32` 函数。

在 v1.2.0 之前，这两个函数的形式是这样的：

```go
func NewFloat64(f float64, prec int) *V
func NewFloat32(f float32, prec int) *V
```

在底层中，prec 字段是用在 `strconv.FormatFloat` 函数的 prec 字段中的，而对应的 `fmt` 参数，则填入 `'f'`。

但是经过 [Issue #8](https://github.com/Andrew-M-C/go.jsonvalue/issues/8) 的提醒，作者意识到 `'f'` 格式并不是对所有的浮点数是最优的表达形式，在部分情况下，使用科学计数法会更好，因此 `strconv.FormatFloat` 的 `fmt` 参数有必要开放给开发者。

如果考虑向后兼容，原本可以简单粗暴地在 NewFloat64 函数的最后面加一个 `fmt byte` 参数，但是作者思考再三，否决了这个方案，理由如下：

- fmt 和 prec 参数的顺序与 strconv.FormatFloat 的参数顺序不一样，让人费解；而这参数，本质上都是数字，而由于顺序与 strconv.FormatFloat 不同，很容易传错，且编译器不会报任何错误，导致 bug 的发生。

因此作者最终决定，发布一版不向后兼容的版本，变更如下：

- `NewFloat64` 和 `NewFloat32` 函数，取消第二个 `prec` 参数。开发者直接传入浮点数即可，格式化浮点数的逻辑与 `encoding/json` 保持一致；
- 增加 `NewFloat64f` 和 `NewFloat32f` 函数，参数与 `strconv.FormatFloat` 相同，均为浮点数、格式、精度三个参数。
    - 但这里需要注意的是，由于 JSON number 的规定，`fmt` 参数只支持 `f`, `E`, `e`, `G`, `g` 五个。如果输入了非法值，则默认采用 `g`。

---

## 使用函数形式配置序列化的额外参数

这一点请参见[额外选项配置](./12_option.md)小节，正如其最后一部分所说的，之前的配置模式是通过传入一个 struct 实现的，现在改为使用 `OptXxx` 系列函数创建配置值。

---

## 支持处理 NaN 和 Inf 浮点值

参见[值的自动转换](./10_conversion.md)小节。

---

## 支持覆盖默认序列化配置

```go
func SetDefaultMarshalOptions(opts ...Option)
func ResetDefaultMarshalOptions()
```

在 jsonvalue 中，序列化的默认逻辑为:

- 按照 JSON 标准序列化所有需要转义的字符，包括 `", /, \, <, >, &`, 水平/垂直制表符, 换行符, 退格符等，均进行转义, 此外 Go 格式化输出的标志字符 `%` 也进行转义。
- 非 ASCII 字符，均进行转义

但是从反馈上来看，很多具体的需求中，不需要这么严格的转义，普通的即可。在这种情况下，当开发者能够明确全局的序列化格式时，可以在程序启动后调用一次 `SetDefaultMarshalOptions` 函数来覆盖掉 jsonvalue 的配置逻辑。当然，也可以调用 `ResetDefaultMarshalOptions` 复位。

