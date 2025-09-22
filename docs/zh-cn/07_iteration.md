<font size=6>遍历和迭代</font>

[上一页](./06_import_export.md) | [总目录](./README.md) | [下一页](./08_caseless.md)

---

- [概述](#概述)
- [遍历 array 类型值](#遍历-array-类型值)
- [遍历 object 类型值](#遍历-object-类型值)
- [遍历所有子值](#遍历所有子值)
- [获取 object 类型值的原始顺序](#获取-object-类型值的原始顺序)
- [已弃用的函数](#已弃用的函数)

---

## 概述

在 jsonvalue 中，开发者可以遍历一个 array 或 object 类型值，前者是按顺序迭代数组中的每一个值，后者则是迭代对象类型中的每一个键值对。

支持两种模式的迭代：第一种是使用回调函数的方式，这类似于 `jsonparser` 的 [`ArrayEach`](https://pkg.go.dev/github.com/buger/jsonparser#ArrayEach) 和 [`ObjectEach`](https://pkg.go.dev/github.com/buger/jsonparser#ObjectEach) 函数。

另一种模式则允许开发者使用 Go 的 `for-range` 语法，更加接近于对 `map` 和 `slice` 的操作。

需要特别**注意**的是：jsonvalue 提供的迭代函数均为线程不安全的！在多协程环境下操作时，请开发者注意加锁保护。如果在过程中不需要调用 `Set` 系列函数和 `Caseless` 函数的话，可以只加读锁。

在使用 `for-range` 的模式中，jsonvalue 需要先遍历一遍 object 或 array 值，组装成返回值之后再给业务函数遍历一次，因此实际上是遍历了两次。如果业务代码对效率极为敏感，或者只需要遍历极少数子成员，那么建议使用回调函数模式。

在命名方面，由于历史原因，笔者先开发了 RangeXxx 系列函数，所以导致 for-range 风格的函数反而不使用 range 命名，还请开发者们海涵。

## 遍历 array 类型值

遍历数组类型值，可以使用以下两个函数：

```go
func (v *V) RangeArray(callback func(i int, v *V) bool)
func (v *V) ForRangeArr() []*V
```

`RangeArray` 采用回调函数模式。在回调函数中，需要返回 `true` 以继续迭代，相当于 `continue`；返回 `false` 则相当于 `break`，表示终止迭代。

`ForRangeArr` 函数则返回一个 `[]*jsonvalue.V` 切片，开发者可以直接把这个函数放在 `for-range` 代码块后面。

具体举例如下：

```go
anArr.RangeArray(func(i int, v *jsonvalue.V) bool {
    // ...... handle with i and v
    return true // 表示 continue
})

for i, v := range anArr.ForRangeArr() {
    // ...... handle with i and v
}
```

## 遍历 object 类型值

遍历对象类型值，可以使用以下两个函数：

```go
func (v *V) RangeObjects(callback func(k string, v *V) bool)
func (v *V) ForRangeObj() map[string]*V
```

类似地，举例如下：

```go
anObj.RangeObjects(func(key string, v *jsonvalue.V) bool {
    // ...... handle with key and v
    return true // 表示 continue
})

for key, v := range anObj.ForRangeObj() {
    // ...... handle with key and v
}
```

## 遍历所有子值

`Walk` 方法提供了一种简单的方式来深度优先遍历 JSON 结构中的所有子值。与只遍历直接子元素的 `RangeArray` 和 `RangeObjects` 不同，`Walk` 方法会递归访问 JSON 树中的所有叶子节点。这个方法的模式类似于 `path/filepath` 的 `Walk` 方法。

```go
func (v *V) Walk(fn WalkFunc)

type WalkFunc func(path []PathItem, v *V) bool

type PathItem struct {
    Idx int    // 数组索引，-1 表示此元素不是数组元素
    Key string // 对象键名，"" 表示此元素不是对象元素
}
```

`Walk` 方法会为每个叶子值调用提供的 `WalkFunc` 回调函数，并提供：
- `path`：一个 `PathItem` 切片，显示从根节点到当前值的完整路径
- `v`：当前被访问的 JSON 值

回调函数应返回 `true` 以继续遍历，或返回 `false` 以提前停止遍历。

## 获取 object 类型值的原始顺序

该功能的呼声其实还不小，但毕竟这个功能相对小众，笔者担心影响正常的 unmarshal 功能的性能。不过从 1.3.1 版本之后，笔者采用了一个简单的手段实现了它，对原有 unmarshal 性能几乎没有带来影响。

在实际使用上，调用方可以先正常执行 `Unmarshal` 操作，然后使用 `RangeObjectsBySetSequence` 函数。这个函数的参数与 `RangeObjects` 完全相同，也是使用一个回调函数，通过返回 `true` 来继续下一次迭代。

比如下面的代码段：

```go
const raw = `{"a":1,"b":2,"c":3}`
v := jsonvalue.MustUnmarshalString(raw)
keys := []string{}
v.RangeObjectsBySetSequence(func(key string, _ *V) bool {
    keys = append(keys, key)
    return true
})
fmt.Println(keys)
```

可以稳定地保证获得 `[a, b, c]`。

## 已弃用的函数

从 v1.2.0 版本开始，弃用的遍历函数为：

```go
func (v *V) IterArray() <-chan *ArrayIter
func (v *V) IterObjects() <-chan *ObjectIter
```

这两个函数是通过返回 `channel` 的方式来实现 `for-range`，效率低下，因此弃用。请勿再使用这两个函数。
