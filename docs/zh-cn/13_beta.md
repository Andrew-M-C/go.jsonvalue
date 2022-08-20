
<font size=6>实验性功能</font>

[上一页](./12_option.md) | [总目录](./README.md) | [下一页](./14_1_13_new_feature.md)

---

在 jsonvalue 中提供了一些 beta 功能，使用以下方法导入:

```go
import (
    "github.com/Andrew-M-C/go.jsonvalue/beta"
)
```

在 beta 包中, 提供实验性的功能。从 v1.3.0 版本则新增了另一个函数：`Contains`，用于判断一个 JSON 值是否包含另一个子集。函数原型如下：

```go
func Contains(v *jsonvalue.V, sub any, inPath ...any) bool
```

首先我们看 `inPath` 参数，如果传入了路径参数的话，那么首先会从 v 中 `Get(inPath...)` 对应的子值，然后再执行 Contains 逻辑。逻辑如下：

- 当两个 JSON 的类型不同，则返回 `false`
- 如果是对象和数组之外的类型，则返回两个值是否相等
- 如果是数组类型，则判断是否拥有子数组
- 如果是对象类型，则迭代每一个值：
    - 当指定的 “子集” 拥有 “父集” 以外的 key 时，则返回 `false`
    - 针对指定的 “子集” 所拥有的所有 key 下面，均递归执行 `Contains` 函数，全部递归均为 `true` 时，则返回 `true`

此外，自从 v1.2 开始提供的 `Import` 函数，现在已转入正式，如果开发者用到了，请直接到正式包里调用即可。
