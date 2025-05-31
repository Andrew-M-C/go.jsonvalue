<font size=6>Import and Export with `encoding/json`</font>

[Prev Page](./05_marshal_unmarshal.md) | [Contents](./README.md) | [Next Page](./07_iteration.md)

---

- [Import / Export](#import--export)
- [Function `New()`](#function-new)
- [Methods `Set`, `Append`, `Insert`, `Add`](#methods-set-append-insert-add)

---

## Import / Export

The initial purpose of designing `Import` and `Export` is to convert data between `encoding/json` and `jsonvalue`.

However, the development of `Import` resulted in many additional features as described below:

---

## Function `New()`

From v1.3.0, a new function called `New` is provided. This function receives any kind of parameter (`any` or `interface{}` in older versions), then parses it to a `*jsonvalue.V` value. If the input parameter is invalid, the type of the returned value will be `NotExist`.

Actually, `New` is a simple wrapper of the function `Import`; however, it does not return an `error` type.

---

## Methods `Set`, `Append`, `Insert`, `Add`

Before version v1.3.0, when you wanted to add a sub-value into a jsonvalue node, you had to specify the parameter type of the input parameter. For example:

```go
v.SetString("Hello, world").At("greeting")
```

After v1.3.0, methods like `Set`, `Append`, `Insert`, `Add` will accept `any` type of parameter, parse it, and then set the corresponding sub-value.

For example, after creating an empty JSON object:

```go
v := jsonvalue.NewObject()
```

We can add a sub-object into it:

```go
child := map[string]string{
    "text": "Hello, jsonvalue!",
}
v.Set(child).At("child")
fmt.Println(v.MustMarshalString())
```

Outputs: `{"child":{"text":"Hello, jsonvalue!"}}`

Or we can set a normal JSON value:

```go
v := jsonvalue.NewObject()
v.Set("Hello, JSON!").At("msg")
fmt.Println(v.MustMarshalString())
// Output: {"msg":"Hello, JSON!"}
```


