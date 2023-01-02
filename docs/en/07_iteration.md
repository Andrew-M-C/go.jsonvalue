
<font size=6>Iteration</font>

[Prev Page](./06_import_export.md) | [Contents](./README.md) | [Next Page](./08_caseless.md)

---

- [Overview](#overview)
- [Iterate array values](#iterate-array-values)
- [Iterate object values](#iterate-object-values)
- [Acquiring the Original Key Sequence of An Object](#acquiring-the-original-key-sequence-of-an-object)

---

## Overview

In jsonvalue, you can iterate an array or object typed JSON value. 

There are two mode of iteration:

1. Use a callback function to receive iteration data, like [`ArrayEach`](https://pkg.go.dev/github.com/buger/jsonparser#ArrayEach) and [`ObjectEach`](https://pkg.go.dev/github.com/buger/jsonparser#ObjectEach) in `jsonparser`.
2. Use `for-range` to iterate, which make codes much more like operating `map` and `slice`.

**IMPORTANT**: All jsonvalue iterating methods are goroutine-unsafe. Please add lock or other protection in multi-goroutine operating. However, if none of `Set` series methods and `Caseless` method are used during iteration, it will be goroutine-safe. In the other hand, just read-lock required for these operation.

## Iterate array values

Use following methods to iterate an array typed JSON value.

```go
func (v *V) RangeArray(callback func(i int, v *V) bool)
func (v *V) ForRangeArr() []*V
```

For callback pattern, use `RangeArray`. Return `true` in callback to continue iteration (which means `continue`), while `false` as `break`, breaking iteration.

`ForRangeArr` returns a slice with `[]*jsonvalue.V`, which you can append to a `for-range` text.

For detailed example:

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

Use following methods to iterate an object typed JSON value:

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

for key, v := range anArr.ForRangeObj() {
    // ...... handle with key and v
}
```

## Acquiring the Original Key Sequence of An Object

Theoretically, the key sequence of an JSON object should be undefined, unexpected. But in practical, it is quite surprising that this feature is quite popular.

After v1.3.1, this feature is added, and it takes almost no affecting the the unmarshal efficiency.

The caller may still execute the `Unmarshal` operation, and then use `RangeObjectsBySetSequence` method, which accept the same callback like `RangeObjects`.

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

The output is `[a, b, c]` and always be guaranteed.

