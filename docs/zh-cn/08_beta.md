# 实验性功能

[上一页](./07_conversion.md) | [总目录](./README.md) | [下一页](./09_new_feature.md)

---

- [struct 转 jsonvalue](./08_beta.md#struct-转-jsonvalue)

---

## struct 转 jsonvalue

在 jsonvalue 正篇中，提供了两个函数:

```go
func Import(src interface{}) (*V, error)
func (v *V) Export(dst interface{}) error
```

这两个函数实现的功能是将 struct（或者是其他 `encoding/json` 支持的变量类型）与 jsonvalue 实现互相转换。

在底层中，这两个函数实际上是通过 `encoding/json` 转换出 `[]byte` 序列之后再进行中转，可以说只是一个简单的流程封装而已，并没有专门优化性能。

为此，jsonvalue 提供了一个 [beta](../../beta) 包，也提供了一个形式上完全一致的 `Import` 函数。使用该函数，有以下明显的优势：

- 从 `interface{}` 直接转为 `*V` ，少了 `encoding/json` 的一次中转，因此性能大大提高。
- 由于转成了 jsonvalue，因此能够对 struct 进行进一步的扩展，比如添加附加参数，或者是将 JSON 进行继承、嵌套等等操作。
- 可以做到 `encoding/json` 所做不到的一些功能，比如忽略大小写、支持 NaN、+/-Inf 等。

该功能目前笔者正在使用，暂时还没移植到正篇中。也欢迎开发者们使用和反馈。
