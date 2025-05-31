<font size=6>Caseless</font>

[Prev Page](./07_iteration.md) | [Contents](./README.md) | [Next Page](./09_conversion.md)

---

- [Issue with `encoding/json`](#issue-with-encodingjson)
- [Doing Caseless in Jsonvalue](#doing-caseless-in-jsonvalue)

---

## Issue with `encoding/json`

With `encoding/json`, it handles object keys case-insensitively when using JSON tags in a `struct`. For example:

```go
type st struct {
    Name string `json:"name"` // lower case
}

func main() {
    raw := []byte(`{"NAME":"json"}`) // upper case
    s := st{}
    json.Unmarshal(raw, &s)
    fmt.Println("name:", s.Name)
    // Output:
    // name: json
}
```

Although the JSON text uses the full-uppercase `NAME` as the key, while the tag definition in the `struct` is lowercase `name`, the value is parsed into the `struct` correctly.

But it is different with `map`. With `encoding/json`, access to a map is case-sensitive:

```go
func main() {
    raw := []byte(`{"NAME":"json"}`)
    var m map[string]any
    json.Unmarshal(raw, &m)
    fmt.Println("name:", m["name"])
    // Output:
    // name: <nil>
}
```

Can we solve this issue in jsonvalue? Yes!

---

## Doing Caseless in Jsonvalue

In jsonvalue, I use a map to store object-typed JSON values. By default, the `Get` method is case-sensitive.

If you want to perform case-insensitive operations when getting sub-values in a JSON value, just insert a `Caseless()` before the `Get` method. For example:

```go
func main() {
    raw := []byte(`{"NAME":"json"}`)
    v := jsonvalue.MustUnmarshal(raw)
    fmt.Println("name =", v.MustGet("name").String())
    fmt.Println("NAME =", v.Caseless().MustGet("name").String())
    // Output:
    // name =
    // NAME = json
}
```

The second `Println` outputs the value.

When getting sub-values with `Caseless`, jsonvalue will first try to match the key exactly. If not found, other keys will be searched until one is found. Therefore, the keys `name` and `NAME` may coexist in the same jsonvalue object.

**IMPORTANT:** Please be aware of the following characteristics of the `Caseless()` method:

- `Caseless` may change the internal structure of `*jsonvalue.V`. Since `jsonvalue` is not goroutine-safe, please add a write-lock when using `Caseless` in multi-goroutine operations.
- `Caseless` takes additional CPU time. If not necessary, please do not use it.
- The `Equal` method is NOT supported after `Caseless`. Because it is quite difficult to define "equal" when comparing two caseless values, I decided not to support `Equal` for `Caseless`.
