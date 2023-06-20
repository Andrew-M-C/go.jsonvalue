
<font size=6>序列化和反序列化</font>

[上一页](./04_get.md) | [总目录](./README.md) | [下一页](./06_import_export.md)

---

- [Unmarshal 系列函数](#unmarshal-系列函数)
  - [基础的 Unmarshal](#基础的-unmarshal)
  - [其他 Unmarshal 函数](#其他-unmarshal-函数)
- [Marshal 系列函数](#marshal-系列函数)
- [原生 `encoding/json` 支持](#原生-encodingjson-支持)

---

## Unmarshal 系列函数

### 基础的 Unmarshal

在 jsonvalue 中，采用 Go 的 marshal / unmarshal 语义描述序列化和反序列化的过程。

在原生 `encoding/go` 中，`json.Unmarshal` 的入参是一个 `[]byte` 类型。在 jsonvalue 中类似，可以使用以下函数来解析 JSON 字节串

```go
func Unmarshal(b []byte) (ret *V, err error)
```

不论是否 error，该函数都会返回一个 `*jsonvalue.V` 对象

- 当 JSON 字节串非法时，会返回 error 信息，此时返回的 json 对象的类型等于 `jsonvalue.NotExist`

### 其他 Unmarshal 函数

在实际操作中，JSON 字节串经常会以 `string` 而不是 `[]byte` 的格式出现。如果进行 `string(b)` 进行转换的话，实际上会进行一次内存拷贝。为了节省这个拷贝的开销，jsonvalue 也提供了入参为 `string` 的 unmarshal 函数:

```go
func UnmarshalString(s string) (ret *V, err error)
```

此外，当程序不需要关心 JSON 字节串格式是否正确的时候，也可以使用 `Must...` 系列的 unmarshal 函数：

```go
func MustUnmarshal(b []byte) *V
func MustUnmarshalString(s string) *V
```

这两个函数与前面一样，必然会返回一个非空的 `jsonvalue.V` 对象，但是不返回 `error` 类型，便于开发者编写一些极为简短的逻辑代码。这种简短代码的技巧，在本页面的最后会进行介绍。

---

## Marshal 系列函数

与 Unmarshal 对应，jsonvalue 的序列化函数也采用其相对的 marshal 语义。提供了以下四个方法：

```go
func (v *V) Marshal          (opts ...Option) (b []byte, err error)
func (v *V) MarshalString    (opts ...Option) (s string, err error)
func (v *V) MustMarshal      (opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

在当前版本下，marshal 只有两种情况会报错：

- `*V` 是 `NotExist` 类型
- 值中包含了不合法的浮点数值 `+Inf`, `-Inf` 或 `NaN`，并且没有明确说明如何处理这些数值。后文还会提到，jsonvalue 针对这些非法的 number 数值，还提供了额外的处理能力。
- 选项参数 opts 中包含非法配置

因此如果开发者能够确定规避掉上述错误的话，完全可以使用 `MustMarshal` 系列函数。

## 原生 `encoding/json` 支持

`*jsonvalue.V` 类型支持 `json.Marshaler` 和 `json.Unmarshaler`。这就意味着我们可以在原生的 `encoding/json` 的序列化和反序列化操作中使用本 package, 如:

```go
var v &jsonvalue.V{}
err := json.Unmarshal(data, v)
```

或者是

```go
v := jsonvalue.NewObject()
v.MustSet("Hello, JSON!").At("greeting")
b, err := json.Marshal(v)
```
