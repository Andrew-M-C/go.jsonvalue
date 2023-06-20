
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

In most situations, we need to create an outer JSON value with object-or-array-typed. In jsonvalue, we use following functions:

```go
o := jsonvalue.NewObject()
a := jsonvalue.NewArray()
```

Also we can convert (import) any JSON-legal data types to jsonvalue. Just use `New` function. Take the object and array above for example:

```go
o := jsonvalue.New(struct{}{})  // construct a JSON object
a := jsonvalue.New([]int{})     // construct a JSON array
```

Also other simple JSON elements are supported:

```go
i := jsonvalue.New(100)             // construct a JSON number
f := jsonvalue.New(188.88)          // construct a JSON number
s := jsonvalue.New("Hello, JSON!")  // construct a JSON string
b := jsonvalue.New(true)            // construct a JSON boolean
n := jsonvalue.New(nil)             // construct a JSON null
```

---

## Set Sub Value into JSON

After generating the outer object or array, the next step is to create there inside structures. Like `Get` method shown in previous section, we can use `Set` or `MustSet` to achieve this.

The `Set(xxx).At(yyy)` methods will return sub-value and error. While `MustSet(xxx).At(yyy)` not. If you does not case the return values please use `MustSet` methods, which will avoid golangci-lint's "return value unused" warning.

### Basic Usage

Generally, we can Use `Set` to construct child value:

```go
v.MustSet(child).At(path...)
```

The semantics is "SET something AT some position". Please be advised that value the ahead of key.

As the parameter type of `Set` method is `any` (`interface{}`), therefore you can set any supported type (even complex object or array data) into a jsonvalue.

Complete example:

```go
v := jsonvalue.NewObject()
v.MustSet("Hello, JSON!").At("data", "message")
v.MustSet(221101).At("data", "date")
fmt.Println(v.MustMarshalString())
```

Output: `{"data":{"message":"Hello, JSON!","date":221101}}`

### Semantics of `At` Method

After calling `Set`, `At` should be followed afterward to set child value into JSON. The prototype of `At` is:

```go
type Setter interface {
	At(firstParam interface{}, otherParams ...interface{}) (*V, error)
}
```

Basic semantics of this method is consistent with `Get`. To prevent programming error, at least one parameter should be given, this the meaning of `param1`.

The more important feature of `At` is that it can generate target JSON structure automatically. It process with following logic:

- Firstly locate sub position by given parameter, just like `Get`. If the target path already exists previously, simple set the sub value in specified path.
- If target position does not exist, the structure wil be created automatically. Either `string` or `int`-like type parameters supported in this method. String type identifies an object while integer an array.

Here is an example with automatic path generating:

```go
v := jsonvalue.NewObject()                       // {}
v.MustSet("Hello, object!").At("obj", "message") // {"obj":{"message":"Hello, object!"}}
v.MustSet("Hello, array!").At("arr", 0)          // {"obj":{"message":"Hello, object!"},"arr":["Hello, array!"]}
```

As for array auto-creating, the procedure is a bit complicated:

- If the array specified in given parameters does not exist, the index value SHOULD be zero to make the operation success.
- If the array already exists, either two cases will be OK:
  - The corresponding child value specified by given index parameter value exists. In this case, the value in that slot my be replaced.
  - The given index value equals to length of the array. In this case, the value will be append to the end of array.

This feature is so complicated that we will not use in most cases. But there is one situation which is useful: 

```go
    const words = []string{"apple", "banana", "cat", "dog"}
    const lessons = []int{1, 2, 3, 4}
    v := jsonvalue.NewObject()
    for i := range words {
        v.MustSet(words[i]).At("array", i, "word")
        v.MustSet(lessons[i]).At("array", i, "lesson")
    }
    fmt.Println(c.MustMarshalString())
```

Final output:

```json
{"array":[{"word":"apple","lesson":1},{"word":"banana","lesson":2},{"word":"cat","lesson":3},{"word":"dog","lesson":4}]}
```

---

## Append Values to JSON Array

Method `Append` and `Insert` are designed for array type JSON. While `Append` should works with `InTheBeginning` and `InTheEnd`, while `Insert` method `After` and `Before`.

They work with semantics below:

- Append some value in the beginning of ...
- Append some value to the end of ...
- Insert some value after ...
- Insert some value before ...

Please be advised of the parameter sequence.

Like `Set` methods, there are also `MustAppend` and `MustInsert` methods by same reason.

Prototypes of these methods as below:

```go
func (v *V) Append(child any) Appender
type Appender interface {
	InTheBeginning(params ...interface{}) (*V, error)
	InTheEnd(params ...interface{}) (*V, error)
}

func (v *V) Insert(child any) Inserter
type Inserter interface {
	After(firstParam interface{}, otherParams ...interface{}) (*V, error)
	Before(firstParam interface{}, otherParams ...interface{}) (*V, error)
}

func (v *V) MustAppend(child any) MustAppender
type MustAppender interface {
	InTheBeginning(params ...interface{})
	InTheEnd(params ...interface{})
}

func (v *V) MustInsert(child any) MustInserter
type MustInserter interface {
	After(firstParam interface{}, otherParams ...interface{})
	Before(firstParam interface{}, otherParams ...interface{})
}
```

Basic semantics are like `Set` methods. But there are a bit differences:

- Empty parameter is allowed in `InTheBeginning` and `InTheEnd`, which identifies that current JSON is already an array, and sub values will append to the beginning or end of it.
- The last parameter SHOULD be a number for `After` and `Before`, identifying the index of the array. Negative index is allowed.
