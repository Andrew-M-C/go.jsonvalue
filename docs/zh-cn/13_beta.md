<font size=6>实验性功能</font>

[上一页](./12_option.md) | [总目录](./README.md) | [下一页](./14_1_14_new_feature.md)

---

在 jsonvalue 中提供了一些 beta 功能，使用以下方法导入：

```go
import (
    "github.com/Andrew-M-C/go.jsonvalue/beta"
)
```

在 beta 包中提供了实验性的功能。从 v1.3.0 版本开始新增了一个函数：`Contains`，用于判断一个 JSON 值是否包含另一个子集。函数原型如下：

```go
func Contains(v *jsonvalue.V, sub any, inPath ...any) bool
```

首先我们来看 `inPath` 参数，如果传入了路径参数，那么首先会从 v 中通过 `Get(inPath...)` 获取对应的子值，然后再执行 Contains 逻辑。具体逻辑如下：

- 当两个 JSON 的类型不同时，返回 `false`
- 如果是对象和数组之外的类型，则返回两个值是否相等
- 如果是数组类型，则判断是否包含子数组
- 如果是对象类型，则迭代每一个值：
    - 当指定的"子集"拥有"父集"以外的 key 时，返回 `false`
    - 针对指定的"子集"所拥有的所有 key，均递归执行 `Contains` 函数，当全部递归均为 `true` 时，返回 `true`

此外，自 v1.2 版本开始提供的 `Import` 函数现在已转为正式功能，如果开发者用到了该函数，请直接在正式包中调用即可。
