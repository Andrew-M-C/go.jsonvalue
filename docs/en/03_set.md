<font size=6>Create and Serialize JSON</font>

[Prev Page](./02_quick_start.md) | [Contents](./README.md) | [Next Page](./04_get.md)

---

This section shows how to generate a jsonvalue value and set it

---

- [Create JSON Value](#create-json-value)
- [Set Sub Value into JSON](#set-sub-value-into-json)
  - [Basic Usage](#basic-usage)
  - [Semantics of `At` Method](#semantics-of-at-method)
- [Append Values to JSON Array](#append-values-to-json-array)

---

## Create JSON Value

In most situations, we need to create an outer JSON value with object or array type. In jsonvalue, we use the following functions:

```go
o := jsonvalue.NewObject()
a := jsonvalue.NewArray()
```

We can also convert (import) any JSON-legal data types to jsonvalue. Just use the `New` function. Take the object and array above for example:

```go
o := jsonvalue.New(struct{}{})  // construct a JSON object
a := jsonvalue.New([]int{})     // construct a JSON array
```

Other simple JSON elements are also supported:

```go
i := jsonvalue.New(100)             // construct a JSON number
f := jsonvalue.New(188.88)          // construct a JSON number
s := jsonvalue.New("Hello, JSON!")  // construct a JSON string
b := jsonvalue.New(true)            // construct a JSON boolean
n := jsonvalue.New(nil)             // construct a JSON null
```

---

## Set Sub Value into JSON

After generating the outer object or array, the next step is to create the inside structures. Like the `Get` method shown in the previous section, we can use `Set` or `MustSet` to achieve this.

The `Set(xxx).At(yyy)` methods will return sub-value and error. While `MustSet(xxx).At(yyy)` will not. If you do not care about the return values, please use `MustSet` methods, which will avoid golangci-lint's "return value unused" warning.

### Basic Usage

Generally, we can use `Set` to construct child values:

```go
v.MustSet(child).At(path...)
```

The semantics is "SET something AT some position". Please be advised that the value comes ahead of the key.

As the parameter type of the `Set` method is `any`, therefore you can set any supported type (even complex object or array data) into a jsonvalue.

Complete example:

```go
v := jsonvalue.NewObject()
v.MustSet("Hello, JSON!").At("data", "message")
v.MustSet(221101).At("data", "date")
fmt.Println(v.MustMarshalString())
```

Output: `{"data":{"message":"Hello, JSON!","date":221101}}`

### Semantics of `At` Method

After calling `Set`, `At` should be called afterward to set child value into JSON. The prototype of `At` is:

```go
type Setter interface {
	At(firstParam any, otherParams ...any) (*V, error)
}
```

The basic semantics of this method is consistent with `Get`. To prevent programming errors, at least one parameter should be given, which is the meaning of `firstParam`.

The more important feature of `At` is that it can generate the target JSON structure automatically. It processes with the following logic:

- First, locate the sub position by the given parameter, just like `Get`. If the target path already exists previously, simply set the sub value in the specified path.
- If the target position does not exist, the structure will be created automatically. Either `string` or `int`-like type parameters are supported in this method. String type identifies an object while integer identifies an array.

Here is an example with automatic path generation:

```go
v := jsonvalue.NewObject()                       // {}
v.MustSet("Hello, object!").At("obj", "message") // {"obj":{"message":"Hello, object!"}}
v.MustSet("Hello, array!").At("arr", 0)          // {"obj":{"message":"Hello, object!"},"arr":["Hello, array!"]}
```

As for array auto-creation, the procedure is a bit complicated:

- If the array specified in the given parameters does not exist, the index value SHOULD be zero to make the operation successful.
- If the array already exists, either of the two cases will be OK:
  - The corresponding child value specified by the given index parameter value exists. In this case, the value in that slot may be replaced.
  - The given index value equals the length of the array. In this case, the value will be appended to the end of the array.

This feature is so complicated that we will not use it in most cases. But there is one situation in which it is useful:

```go
    var words = []string{"apple", "banana", "cat", "dog"}
    var lessons = []int{1, 2, 3, 4}
    v := jsonvalue.NewObject()
    for i := range words {
        v.MustSet(words[i]).At("array", i, "word")
        v.MustSet(lessons[i]).At("array", i, "lesson")
    }
    fmt.Println(v.MustMarshalString())
```

If you like placing keys ahead, you can use the `v.At(...).Set(...)` pattern:

```go
    // ...
        v.At("array", i, "word").Set(words[i])
        v.At("array", i, "lesson").Set(lessons[i])
    // ...
```

Final output:

```json
{"array":[{"word":"apple","lesson":1},{"word":"banana","lesson":2},{"word":"cat","lesson":3},{"word":"dog","lesson":4}]}
```

You can also pass a slice or array to identify the parameter chain, for example, the following code:

```go
v.MustSet("Hello, object!").At("obj", "message")
```

is equivalent to this:

```go
v.MustSet("Hello, object!").At([]any{"obj", "message"})
```

or:

```go
v.MustSet("Hello, object!").At([]string{"obj", "message"})
```

This feature makes it easy to pass parameters from outer sources or configurations.

---

## Append Values to JSON Array

Methods `Append` and `Insert` are designed for array type JSON. `Append` should work with `InTheBeginning` and `InTheEnd`, while `Insert` method works with `After` and `Before`.

They work with the semantics below:

- Append some value at the beginning of ...
- Append some value to the end of ...
- Insert some value after ...
- Insert some value before ...

Please be advised of the parameter sequence.

Like `Set` methods, there are also `MustAppend` and `MustInsert` methods for the same reason.

Prototypes of these methods are as below:

```go
func (v *V) Append(child any) Appender
type Appender interface {
	InTheBeginning(params ...any) (*V, error)
	InTheEnd(params ...any) (*V, error)
}

func (v *V) Insert(child any) Inserter
type Inserter interface {
	After(firstParam any, otherParams ...any) (*V, error)
	Before(firstParam any, otherParams ...any) (*V, error)
}

func (v *V) MustAppend(child any) MustAppender
type MustAppender interface {
	InTheBeginning(params ...any)
	InTheEnd(params ...any)
}

func (v *V) MustInsert(child any) MustInserter
type MustInserter interface {
	After(firstParam any, otherParams ...any)
	Before(firstParam any, otherParams ...any)
}
```

Basic semantics are like `Set` methods. But there are a few differences:

- Empty parameters are allowed in `InTheBeginning` and `InTheEnd`, which identifies that the current JSON is already an array, and sub values will be appended to the beginning or end of it.
- The last parameter SHOULD be a number for `After` and `Before`, identifying the index of the array. Negative index is allowed.
