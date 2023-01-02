
<font size=6>Caseless</font>

[Prev Page](./07_iteration.md) | [Contents](./README.md) | [Next Page](./09_conversion.md)

---

- [Issue with `encoding/json`](#issue-with-encodingjson)
- [Doing Caseless in Jsonvalue](#doing-caseless-in-jsonvalue)

---

## Issue with `encoding/json`

With `encoding/json`, it handles the key of an object caselessly with JSON tags in a `struct`. Such as:

```go
type st struct {
    Name string `json:"name"` // lower cases
}

func main() {
    raw := []byte(`{"NAME":"json"}`) // upper cases
    s := st{}
    json.Unmarshal(raw, &s)
    fmt.Println("name:", s.Name)
    // Output:
    // name: json
}
```

Although the text of JSON uses full-upper-case `NAME` as key, while the tag definition in the `struct` is lower-case `name`, the value is parsed into the `struct` correctly.

But it is different to `map`. With `encoding/json`, the access to a map is case-sensitive:

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

In jsonvalue, I use a map to store object typed JSON value. By default, the `Get` method is case-sensitive.

If you want to be caseless when getting sub values in an JSON value, just insert a `Caseless()` before the `Get` method. Such as:

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

When getting sub values with `Caseless`, jsonvalue will hit the strictly equaled key. If not found, other keys will be searched until one found. Therefore, the key `name` and `NAME` may be together in the same jsonvalue object.

**IMPORTANT:** Please be advised of following knowledge of `Caseless()` method:

- `Caseless` may change the inside structure of `*jsonvalue.V`. As `jsonvalue` is goroutine-unsafe, please add a write-lock when using `Caseless` for multi-goroutine operating.
- `Caseless` takes additional CPU time, if not necessary, please do not use it.
- `Equal` method is NOT supported after `Caseless`. Because it is quite difficult to define "equal" when comparing two caseless values. Therefore I gave up supporting `Equal` for `Caseless`.
