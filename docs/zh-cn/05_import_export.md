# 与原生 json 的导入/导出

[上一页](./04_set.md) | [总目录](./README.md) | [下一页](./06_iteration.md)

---

## 功能说明

在 jsonvalue 中主要操作的类型是 `*jsonvalue.V`，但在逻辑过程中其实也支持其它类型，就像原生的 `json.Unmarshal` 函数的第二个参数一样，是 `interface{}` 类型。

## Set, Append, Insert, Add 函数

可以注意到，标题中提及的函数参数类型均为 `interface{}`。也就是说这几个函数也支持配置任意类型的子类型。比如:

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

## Import / Export



