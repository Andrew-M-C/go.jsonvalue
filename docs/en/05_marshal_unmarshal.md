
<font size=6>Marshal and Unmarshal</font>

[Prev Page](./04_get.md) | [Contents](./README.md) | [Next Page](./06_import_export.md)

---

- [Unmarshal Functions](#unmarshal-functions)
  - [Basic Unmarshal](#basic-unmarshal)
  - [其他 Unmarshal 函数](#其他-unmarshal-函数)
- [Marshal Functions](#marshal-functions)
- [Official `encoding/json` Support](#official-encodingjson-support)

---

## Unmarshal Functions

### Basic Unmarshal

We use marshal / unmarshal to describe serialization and de-serialization process in jsonvalue.

Jsonvalue uses the following function to parse a raw JSON text:

```go
func Unmarshal(b []byte) (ret *V, err error)
```

No matter whether error occurs, an un-nil `*jsonvalue.V` object will be returned. However, when the raw text is illegal, en error object will be returned, describing what error is. In this case, the `Type()` of the returned jsonvalue value will be `jsonvalue.NotExist`

### 其他 Unmarshal 函数

Practically, a raw JSON text will be given in format of `string` instead of `[]byte`. It will take a little time to do `string(b)` conversion. To save this copying time, you can use `string` version unmarshal function:

```go
func UnmarshalString(s string) (ret *V, err error)
```

Besides, if the correctness of given JSON text need no care about, or it is sure to be legal, we can simply ignore error and use functions below:

```go
func MustUnmarshal(b []byte) *V
func MustUnmarshalString(s string) *V
```

As functions above, this two functions will definitely return a un-nil `jsonvalue.V`.

---

## Marshal Functions

Serialization in jsonvalue is "marshal". Like unmarshal, four functions below are provided:

```go
func (v *V) Marshal          (opts ...Option) (b []byte, err error)
func (v *V) MarshalString    (opts ...Option) (s string, err error)
func (v *V) MustMarshal      (opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

In current version of jsonvalue, error will occurred in situations below:

1. `*V` is `NotExist` type
2. Illegal floating numbers included like `+Inf`, `-Inf` or `NaN`, while no special operation to these floating values a specified.
   - Special options with these floating value will mentioned later in other sections.
3. Illegal configurations included in additional options.

---

## Official `encoding/json` Support

Type `*jsonvalue.V` also implements `json.Marshaler` and `json.Unmarshaler` interfaces. This enables marshaling and unmarshaling via `encoding/json`. For example:

```go
var v &jsonvalue.V{}
err := json.Unmarshal(data, v)
```

Or

```go
v := jsonvalue.NewObject()
v.MustSet("Hello, JSON!").At("greeting")
b, err := json.Marshal(v)
```

