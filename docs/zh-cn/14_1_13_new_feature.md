
<font size=6>1.3.x 新特性</font>

[上一页](./13_beta.md) | [总目录](./README.md) | [下一页](./15_1_12_new_feature.md)

---

- [不再默认转义 % 字符](#不再默认转义--字符)
- [支持与原生 json 的导入/导出](#支持与原生-json-的导入导出)
- [支持类泛型操作](#支持类泛型操作)
- [支持 JSON 比较](#支持-json-比较)
- [支持 JSON 序列化锁进](#支持-json-序列化锁进)

---

## 不再默认转义 % 字符

在 v1.2.x 及以前版本中，jsonvalue 当遇到 % 字符时会进行转义，但这并不是标准 JSON 的做法。从 v1.3.0 开始，jsonvalue 不再默认转义 % 字符。

如有需要继续转义 %，可以给作者提 [issue](https://github.com/Andrew-M-C/go.jsonvalue/issues/new)。

## 支持与原生 json 的导入/导出

参见 [与原生 json 的导入/导出](./06_import_export.md) 小节。这个功能在 v1.2 中是 beta 功能，在本版本中提升为正式功能提供。

## 支持类泛型操作

在 v1.2 及之前版本中，往 `*V` 类型中添加子成员时，需要明确指定类型如：

```go
v.SetString("Hello, JSON!").At("msg")
v.SetInt64(time.Now().Unix()).At("time")
```

从 v1.3.0 开始，不再需要这么繁琐了，`Set` 函数支持传入任意可以序列化为 JSON 的值:

```go
v.Set("Hello, JSON!").At("msg")
v.Set(time.Now().Unix()).At("time")
```

这个特性也可以扩展到其他相关函数中，包括:

- `Append`
- `Insert`
- `Add`

此外，基于该特性，在 v1.3.0 中也新增了 `New` 函数，其实相当于是 `Import` 的一个封装。参见 [与原生 json 的导入/导出](./06_import_export.md) 小节。

## 支持 JSON 比较

主要是新增了 [Equal](./11_comparation.md) 函数，用于判断两个 JSON 是否相等。此外在 beta 包中新增了 [Contains](./13_beta.md) 函数，也可以用于判断子集。

## 支持 JSON 序列化锁进

请参见 “额外选项配置” 的[可视化锁进](./09_option.md#可视化锁进)小节

