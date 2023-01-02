
<font size=6>与原生 json 的导入/导出</font>

[上一页](./05_marshal_unmarshal.md) | [总目录](./README.md) | [下一页](./07_iteration.md)

---

- [Import / Export](#import--export)
- [New 函数](#new-函数)
- [Set, Append, Insert, Add 函数](#set-append-insert-add-函数)

---

## Import / Export

Import 和 Export 最开始的作用，是在原生 `encoding/json` 和 `jsonvalue` 之间进行互转。

此外，作者在开发 `Import` 函数过程中，也顺便构建了不少功能，也就成就了 v1.3.0 版本新增的很多功能，这些功能主要体现在以下的几个内容：

---

## New 函数

从 v1.3.0 开始，jsonvalue 支持 `New` 函数，该函数接收任意类型的参数，并且将其解析为 JSON 并转换为一个 `*jsonvalue.V` 类型。如果入参不是合法的 JSON 值, 那么返回的对象类型为 `NotExist`。

实际上 `New` 函数是针对 `Import` 函数的封装，差别只是不返回错误而已。

---

## Set, Append, Insert, Add 函数

可以注意到，标题中提及的函数参数类型均为 `any`，或者应该说是 `any`。也就是说这几个函数也支持配置任意类型的子类型。比如:

```go
v := jsonvalue.NewObject()
```

我们可以往里添加一个子对象。

```go
child := map[string]string{
    "text": "Hello, jsonvalue!",
}
v.Set(child).At("child")
fmt.Println(v.MustMarshalString())
```

输出为: `{"data":{"message":"Hello, JSON!"}}`

也可以直接配置一个合法的 JSON 值：

```go
v := jsonvalue.NewObject()
v.Set("Hello, JSON!").At("msg")
fmt.Println(v.MustMarshalString())
// Output: {"msg":"Hello, JSON!"}
```


