
<font size=6>值的比较</font>

[上一页](./10_scenarios.md) | [总目录](./README.md) | [下一页](./12_option.md)

---

- [Equal](#equal)
- [数值比较](#数值比较)

---

## Equal

从 v1.3.0 开始，支持使用 `Equal` 函数判断两个 JSON 值是否相等。判断的规则如下：

首先，如果两个 JSON 值的类型不同，则返回 `false`。

当两个 JSON 类型相同时，则按照以下逻辑进行判断:

- `string` 类型: 检查两个 string 值是否相等
- `number` 类型: 检查两个十进制数值是否相等
    - 这里请注意，是比对十进制数值，实现上使用了 [decimal](https://pkg.go.dev/github.com/shopspring/decimal) 库。
- `boolean` 类型: 检查两个布尔值是否相等
- `null` 类型: 两个 null 值永远相等
- `array` 类型: 两个数组相等的充分必要条件是: 长度相等并且在每一个索引位置上的 JSON 值相等
- `object` 类型: 两个对象相等的充分必要条件是: key 列表完全相同，并且对应的每一个值都相等

也即针对 array 和 object 类型，内部进行了递归判断。

---

## 数值比较

从 v1.4.0 开始，支持 `GreaterThan`, `LessThan`, `GreaterThanOrEqual`, `LessThanOrEqual` 等方法，用于比对两个数值的大小。如果任一比较值不是数值类型，则前述方法无论如何都会返回 `false`。
