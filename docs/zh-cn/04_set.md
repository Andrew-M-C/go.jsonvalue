# 创建并序列化 JSON

[上一页](./03_get.md) | [总目录](./README.md) | [下一页](./99_TODO.md)

---

[TOC]

---

## 创建 JSON 值

TODO:

## Marshal 系列函数

与 Unmarshal 对应，jsonvalue 的序列化函数也采用其相对的 marshal 语义。提供了以下四个方法：

```go
func (v *V) Marshal(opts ...Option) (b []byte, err error)
func (v *V) MarshalString(opts ...Option) (s string, err error)
func (v *V) MustMarshal(opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

在当前版本下，marshal 只有两种情况会报错：

- `*V` 是 `NotExist` 类型
- 值中包含了不合法的浮点数值 `+Inf`, `-Inf` 或 `NaN`，并且没有明确说明如何处理这些数值。

因此如果开发者能够确定规避掉这两种错误的话，完全可以使用 `MustMarshal` 系列函数。
