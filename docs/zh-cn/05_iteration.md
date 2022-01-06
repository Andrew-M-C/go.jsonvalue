# 迭代遍历

[上一页](./04_set.md) | [总目录](./README.md) | [下一页](./06_caseless.md)

---

- [概述](./05_iteration.md#概述)
- [遍历 array 类型值](./05_iteration.md#遍历-array-类型值)
- [遍历 object 类型值](./05_iteration.md#遍历-object-类型值)
- [已弃用的函数](./05_iteration.md#已弃用的函数)

---

## 概述

在 jsonvalue 中，开发者也可以遍历一个 array 或 object 类型值，前者是按顺序迭代数组中的每一个值；后者则是迭代对象类型中的每一个键值对。

支持两种模式的迭代，第一种是使用回调函数的方式，这类似于 `jsonparser` 的 [`ArrayEach`](https://pkg.go.dev/github.com/buger/jsonparser#ArrayEach) 和 [`ObjectEach`](https://pkg.go.dev/github.com/buger/jsonparser#ObjectEach) 函数。

另一种模式，则允许开发者使用 Go 的 `for-range` 语法，更加接近于对 `map` 和 `slice` 的操作。

需要特别**注意**的是：jsonvalue 提供的迭代函数，均为线程不安全的！在多协程环境下操作时，请开发者注意加锁保护；如果在过程中不需要调用 `Set` 系列函数和 `Caseless` 函数的话，可以只加读锁。

在使用 `for-range` 的模式中，jsonvalue 需要先遍历一遍 object 或 array 值，组装成返回值之后再给业务函数遍历一次，因此实际上是遍历了两次。因此如果业务代码对效率极为敏感，或者是只需要遍历极少数子成员，那么建议使用回调函数模式。

## 遍历 array 类型值

遍历数组类型值，可以使用以下两个函数：

```go
func (v *V) RangeArray(callback func(i int, v *V) bool)
func (v *V) ForRangeArr() []*V
```

`RangeArray` 采用回调函数模式。在回调函数中，需要返回 `true` 以继续迭代，相当于 `continue`；返回 `false` 则相当于 `break`，表示终止迭代。

`ForRangeArr` 函数则返回一个 `[]*jsonvalue.V` 切片，开发者可以直接把这个函数放在 `for-range` 代码块后面。

具体的举例如下：

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
anObj.RangeObject(func(key string, v *jsonvalue.V) bool {
    // ...... handle with key and v
    return true // 表示 continue
})

for key, v := range anArr.ForRangeObj() {
    // ...... handle with key and v
}
```

## 已弃用的函数

从 v1.2.0 版开始，弃用的遍历函数为：

```go
func (v *V) IterArray() <-chan *ArrayIter
func (v *V) IterObjects() <-chan *ObjectIter
```

这两个函数是通过返回 `channel` 的方式来实现 `for-range`，效率低下，因此弃用。请勿再使用这两个函数。
