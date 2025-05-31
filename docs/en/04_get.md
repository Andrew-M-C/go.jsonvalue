<font size=6>Get Values from JSON Structure</font>

[Prev Page](./03_set.md) | [Contents](./README.md) | [Next Page](./05_marshal_unmarshal.md)

---

- [Get Functions](#get-functions)
  - [Parameters](#parameters)
  - [GetXxx Series](#getxxx-series)
- [`MustGet` and Other Related Methods](#mustget-and-other-related-methods)
- [Properties of jsonvalue.V Object](#properties-of-jsonvaluev-object)
  - [Official Definition](#official-definition)
  - [jsonvalue Basic Properties](#jsonvalue-basic-properties)

---

## Get Functions

### Parameters

The Get function is the core of reading information from jsonvalue. Here is the prototype:

```go
func (v *V) Get(param1 any, params ...any) (*V, error)
```

For a practical example:

```go
const raw = `{"someObject": {"someObject": {"someObject": {"message": "Hello, JSON!"}}}}`
child, _ := jsonvalue.MustUnmarshalString(raw).Get("someObject", "someObject", "someObject", "message")
fmt.Println(child.String())
```

The meaning of the `Get` parameters above is:

- Locate and get the sub-value with key `someObject` from the `*V` instance, then get another value with key `someObject` from the previously located `*V` instance, and continue...
  - This operation could also be described in dot notation format, like: `child = v.someObject.someObject.someObject.message`

The type of parameters for `Get` is `any`. In fact, only those with string or integer (both signed and unsigned are OK) kinds are allowed. `Get` will check the parameter type and decide whether to treat the next value as an object or array for the next iteration.

If the [Kind](https://pkg.go.dev/reflect#Kind) of the current path node's parameter is string, then:

- If the value type of the current `jsonvalue` value is "Object", then locate the sub-value specified by the string key.
  - If the sub-value exists and it is "Object" typed, then continue iteration with this value and the parameter of the next path node.
  - If the value with the specified key does not exist, a value with type `NotExist` will be returned, along with error [ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).
- If the current value is not an "Object", a `NotExist` typed value will be returned, along with error [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).

If the [Kind](https://pkg.go.dev/reflect#Kind) of the current path node's parameter is integer, then:

- If the value type of the current `jsonvalue` value is "Array", then find the sub-value at the specified index using the integer key. At this moment, the meaning of various values of this integer may be:
  - If index >= 0, it will be a normal index value, and locate the sub-value just like an ordinary Go slice. If the index is out of range, a `NotExist` value and [ErrOutOfRange](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants) will be returned.
  - If index < 0, it will be treated as "XXth from the end", counting backwards. However, it should still be within the range of the JSON array.
    - For example, if the length of a JSON array is 5, then -5 locates the element at index 0, while -6 leads to an error being returned.
  - If the search succeeds in the current JSON node, iterations will continue if there are more parameters remaining.
- If the current value at the path node is not an "Array", a `NotExist` typed value will be returned, along with error [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).

You may be curious about why I split the parameters of the `Get` function into two parts, instead of using a simple `...any`? Just like I answered in [this issue](https://github.com/Andrew-M-C/go.jsonvalue/issues/4), I designed it this way on purpose:

- This is to avoid programming errors like `v.Get()` (lacking parameters). By making this method require at least one parameter, an error will be thrown at compile time instead of runtime.
- If you are 100% sure that there is at least one parameter in the input `[]any`, you may call this method like this: `subValue, _ := v.Get(para[0], para[1:]...)`

### GetXxx Series

In practical code, `Get` itself is rarely used; we use its "siblings" to access basic typed values instead:

```go
func (v *V) GetObject (param1 any, params ...any) (*V, error)
func (v *V) GetArray  (param1 any, params ...any) (*V, error)
func (v *V) GetBool   (param1 any, params ...any) (bool, error)
func (v *V) GetString (param1 any, params ...any) (string, error)
func (v *V) GetBytes  (param1 any, params ...any) ([]byte, error)
func (v *V) GetInt    (param1 any, params ...any) (int, error)
func (v *V) GetInt32  (param1 any, params ...any) (int32, error)
func (v *V) GetInt64  (param1 any, params ...any) (int64, error)
func (v *V) GetNull   (param1 any, params ...any) error
func (v *V) GetUint   (param1 any, params ...any) (uint, error)
func (v *V) GetUint32 (param1 any, params ...any) (uint32, error)
func (v *V) GetUint64 (param1 any, params ...any) (uint64, error)
func (v *V) GetFloat32(param1 any, params ...any) (float32, error)
func (v *V) GetFloat64(param1 any, params ...any) (float64, error)
```

There are some commonalities among these methods:

- If the sub-value exists and has the correct type specified by the method, the returned `error` is nil. All methods will return the corresponding value except `GetNull`.
- If the sub-value exists but the type does not match, the error will be [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).
- If the sub-value does not exist, the error will be [ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).

Also, some of these methods do not simply match types and return; they also provide some additional features, which will be described later in subsequent sections. Here I will show you an example:

In many cases, we need to extract a number from a string-typed value, such as: `{"number":"12345"}`. In this case, the `GetInt` method will return the corresponding integer value from this string correctly:

```go
raw := `{"number":"12345"}`
n, err := jsonvalue.MustUnmarshalString(raw).GetInt("number")
fmt.Println("n =", n)
fmt.Println("err =", err)
```

Output:

```
n = 12345
err = type does not match
```

As shown in the example, both `n` and `err` return meaningful values. Not only is a valid number returned, but also the [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants) error is thrown.

---

## `MustGet` and Other Related Methods

As mentioned above, we have the `Get` and `GetXxx` series functions. Except for `GetNull`, each function returns two values. For the Get function, jsonvalue also provides a `MustGet` function that returns only one value, making it convenient to implement very simple logic.

To facilitate understanding, let's use a scenario as an example:

Suppose we are developing a forum feature, and the forum supports pinning several posts to the top. The pinning feature configuration is stored in a "top" field in a JSON string, as shown in the following example:

```json
{
    "//":"other configs",
    "top":[
        {
            "UID": "12345",
            "title": "Posting Guidelines"
        }, {
            "UID": "67890",
            "title": "Forum Highlights"
        }
    ]
}
```

In practice, due to various reasons, the obtained configuration string may have the following exceptional situations:

- The entire string is an empty string ""
- The string is invalid due to incorrect editing, or has format errors
- The "top" field might be `null`, or an empty string

If following traditional logic, we would need to handle these exceptional situations one by one. But if developers don't need to care about these exceptions and only care about valid configurations, we can completely utilize the characteristic that `MustXxx` functions always return a `*V` object to simplify the logic as follows:

```go
    c := jsonvalue.MustUnmarshalString(confString) // Assume confString is the obtained configuration string
    for _, v := range c.MustGet("top").ForRangeArr() {
        feeds = append(feeds, &Feed{               // Append post topics to the return list, assuming the post structure is Feed
            ID:    v.MustGet("UID").String(),
            Title: v.MustGet("title").String(),
        })
    }
```

---

## Properties of jsonvalue.V Object

First, we need to understand some properties defined by the official JSON specification, and then explain how these properties are reflected in `jsonvalue`.

### Official Definition

In the standard [JSON specification](https://www.json.org/json-en.html), the following concepts are defined:

- A valid JSON value is called a JSON `value`. In this toolkit, a `*V` is used to represent a JSON value
- JSON value types include the following:

|   Type    | Description                                                                                                                                    |
| :-------: | :--------------------------------------------------------------------------------------------------------------------------------------------- |
| `object`  | An object, corresponding to a key-value format value. The key must be a string, while the value is a valid JSON `value`                      |
|  `array`  | An array, corresponding to an ordered combination of a series of `values`                                                                      |
| `string`  | String type, which is easy to understand                                                                                                       |
| `number`  | Numeric type, more precisely speaking, a double-precision floating-point number                                                               |
|           | Since JSON is defined based on JavaScript, and JS only has double as the numeric type, number is actually double. This is a small pitfall    |
| `"true"`  | Represents boolean "true"                                                                                                                      |
| `"false"` | Represents boolean "false"                                                                                                                     |
| `"null"`  | Represents null value                                                                                                                          |

### jsonvalue Basic Properties

In the `*jsonvalue.V` object, following the approach of most JSON toolkits, `"true"` and `"false"` are merged into a `Boolean` type. Additionally, `"null"` is also mapped to a `Null` type.

Furthermore, a `NotExist` type is defined to indicate that the current value is not a valid JSON object. There is also an `Unknown` type, which developers don't need to worry about, as this value will not appear during normal usage.

The following functions can be used to get the type properties of a value:

```go
func (v *V) ValueType() ValueType
func (v *V) IsObject()  bool
func (v *V) IsArray()   bool
func (v *V) IsString()  bool
func (v *V) IsNumber()  bool
func (v *V) IsBoolean() bool
func (v *V) IsNull()    bool
```
