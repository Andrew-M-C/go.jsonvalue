<font size=6>Iteration</font>

[Prev Page](./06_import_export.md) | [Contents](./README.md) | [Next Page](./08_caseless.md)

---

- [Overview](#overview)
- [Iterate array values](#iterate-array-values)
- [Iterate object values](#iterate-object-values)
- [Acquiring the Original Key Sequence of An Object](#acquiring-the-original-key-sequence-of-an-object)
- [Walk through all child values](#walk-through-all-child-values)

---

## Overview

In jsonvalue, you can iterate over an array or object typed JSON value. 

There are two modes of iteration:

1. Use a callback function to receive iteration data, like [`ArrayEach`](https://pkg.go.dev/github.com/buger/jsonparser#ArrayEach) and [`ObjectEach`](https://pkg.go.dev/github.com/buger/jsonparser#ObjectEach) in `jsonparser`.
2. Use `for-range` to iterate, which makes code much more like operating on `map` and `slice`.

**IMPORTANT**: All jsonvalue iterating methods are goroutine-unsafe. Please add locks or other protection in multi-goroutine operations. However, if none of the `Set` series methods and `Caseless` method are used during iteration, it will be goroutine-safe. On the other hand, only a read-lock is required for these operations.

## Iterate array values

Use the following methods to iterate over an array typed JSON value.

```go
func (v *V) RangeArray(callback func(i int, v *V) bool)
func (v *V) ForRangeArr() []*V
```

For the callback pattern, use `RangeArray`. Return `true` in the callback to continue iteration (which means `continue`), while `false` acts as `break`, breaking the iteration.

`ForRangeArr` returns a slice with `[]*jsonvalue.V`, which you can use in a `for-range` statement.

For a detailed example:

```go
anArr.RangeArray(func(i int, v *jsonvalue.V) bool {
    // ...... handle with i and v
    return true // continue
})

for i, v := range anArr.ForRangeArr() {
    // ...... handle with i and v
}
```

## Iterate object values

Use the following methods to iterate over an object typed JSON value:

```go
func (v *V) RangeObjects(callback func(k string, v *V) bool)
func (v *V) ForRangeObj() map[string]*V
```

Similarly, for example:

```go
anObj.RangeObject(func(key string, v *jsonvalue.V) bool {
    // ...... handle with key and v
    return true // continue
})

for key, v := range anObj.ForRangeObj() {
    // ...... handle with key and v
}
```

## Acquiring the Original Key Sequence of An Object

Theoretically, the key sequence of a JSON object should be undefined and unexpected. But in practice, it is quite surprising that this feature is quite popular.

After v1.3.1, this feature was added, and it has almost no effect on the unmarshal efficiency.

The caller may still execute the `Unmarshal` operation, and then use the `RangeObjectsBySetSequence` method, which accepts the same callback as `RangeObjects`.

For example:

```go
const raw = `{"a":1,"b":2,"c":3}`
v := jsonvalue.MustUnmarshalString(raw)
keys := []string{}
v.RangeObjectsBySetSequence(func(key string, _ *V) bool {
    keys = append(keys, key)
})
fmt.Println(keys)
```

The output is `[a, b, c]` and is always guaranteed.

## Walk through all child values

The `Walk` method provides a simple way to depth-first traverse all child values in a JSON structure. Unlike `RangeArray` and `RangeObjects` which only iterate over direct children, the `Walk` method recursively visits all leaf nodes in the JSON tree. This method follows a pattern similar to the `Walk` method in `path/filepath`.

```go
func (v *V) Walk(fn WalkFunc)

type WalkFunc func(path []PathItem, v *V) bool

type PathItem struct {
    Idx int    // Array index, -1 indicates this element is not an array element
    Key string // Object key name, "" indicates this element is not an object element
}
```

The `Walk` method calls the provided `WalkFunc` callback function for each leaf value, providing:
- `path`: A slice of `PathItem` showing the complete path from root to the current value
- `v`: The current JSON value being visited

The callback function should return `true` to continue traversal, or `false` to stop traversal early.

