
<font size=6>Get Values from JSON Structure</font>

[Prev Page](./03_set.md) | [Contents](./README.md) | [Next Page](./05_marshal_unmarshal.md)

---

- [Get Functions](#get-functions)
  - [Parameters](#parameters)
  - [GetXxx Series](#getxxx-series)
- [`MustGet` and Other Related Methods](#mustget-and-other-related-methods)
- [jsonvalue.V 对象的属性](#jsonvaluev-对象的属性)
  - [官方定义](#官方定义)
  - [jsonvalue 基础属性](#jsonvalue-基础属性)

---

## Get Functions

### Parameters

Get function is the core of reading information of jsonvalue. Here is the prototype:

```go
func (v *V) Get(param1 any, params ...any) (*V, error)
```

For a practical example:

```go
const raw = `{"someObject": {"someObject": {"someObject": {"message": "Hello, JSON!"}}}}`
child, _ := jsonvalue.MustUnmarshalString(s).Get("someObject", "someObject", "someObject", "message")
fmt.Println(child.String())
```

The meaning of the `Get` parameters above are:

-  Locate and get the sub value with key `someObject` from the `*V` instance, and then get another value with key `someObject` from the previous located `*V` instance, and then go on...
  - This operation could also be described as domain format, like: `child = v.someObject.someObject.someObject.message`

The type of parameters of `Get` is `any`. In fact, only ones with string or integer (both signed and unsigned are OK) kind are allowed. `Get` will check the parameter type and decide whether to treat the next value an object or array for the next iteration.

If the [Kind](https://pkg.go.dev/reflect#Kind) of current path node's parameter is string, then:

- If value type of current `jsonvalue` value is "Object", then locate the sub value specified by the string key.
  - If the sub value exists and it is "Object" typed, then continue iteration with this value and the parameter of next path node.
  - If the value with specified key does not exist, a value with type `NotExist` will be returned, with error [ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).
- If current value is not an "Object", a `NotExist` typed value will be returned, with another error [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).

If the [Kind](https://pkg.go.dev/reflect#Kind) of current path node's parameter is integer, then:

- If value type of current `jsonvalue` value is "Array", then find the sub value from specified index by the integer key. At this moment, the meaning of various value of this integer may be:
  - If index >= 0, it will be a normal index value, and locate sub value just like a ordinary Go slice. If the index is out of range, `NotExist` value and [ErrOutOfRange](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants) will be returned.
  - If index < 0, it will be treated as "XXth to the last", counting backwards. But also, should be within range of the JSON array.
    - For example, if the length of a JSON array is 5, then -5 locates the element in Index 0, while -6 leads to error returned.
  - If the searching success in current JSON node, iterations will continue if there are more parameters remaining.
- If current value of the path node is not an "Array", a `NotExist` typed value will be returned, with another error [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).

You may curious that why do I cut parameters of `Get` function into two parts, instead of a simple `...any`? Just like I answered in [this issue](https://github.com/Andrew-M-C/go.jsonvalue/issues/4), I designed it in purpose:

- This is to avoid the programming error like `v.Get` (lacking parameter). By making this method with at least one parameter, an error will thrown in compiling phase instead of runtime.
- If you are 100% sure that there is at least one parameter for the input `[]any`, you may call this method like this: `subValue, _ := Get(para[0], para[1:]...)`

### GetXxx Series

In practical codes the `Get` itself is rarely used, we use its "siblings" to access basic typed values instead:

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

There are some commons within this methods:

- If the sub value exists, and with correct type specified by the method, the returned `error` is nil. And all methods will return the corresponding value besides `GetNull`.
- If the sub value exists but the type does not match, the error will be [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).
- If the sub value does not exist, the error will be [ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants).

Also, some of these methods do not simply match types and return, they also provide some additional features, which will be described later in later sections. Here I will show you an example:

In many cases, we need to extract number from a string typed value, such as: `{"number":"12345"}`. In this case, `GetInt` method would return the corresponding integer value from this string correctly:

```go
raw := `{"number":"12345"}`
n, err := jsonvalue.MustUnmarshalString(raw).GetInt("number")
fmt.Println("n =", n)
fmt.Println("err =", err)
```

Output：

```
n = 12345
err = not match given type
```

As shown by the example, both `n` and `err` returns meaningful value. Now only a valid number is returned, but also the [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants) error thrown.

---

## `MustGet` and Other Related Methods

上文中提到了 `Get` 和 `GetXxx` 系列函数。除了 `GetNull` 之外，各个函数的返回值均为两个。而针对 Get 函数，jsonvalue 也提供了一个 `MustGet` 函数，仅返回一个参数，从而便于实现即为简单的逻辑。

为了便于理解，我们举个场景作为例子——

比如我们开发一个论坛功能，论坛支持将若干个帖子进行置顶。置顶功能的配置是一段 JSON 字符串的 "top" 字段，举例如下：


```json
{
    "//":"other configs",
    "top":[
        {
            "UID": "12345",
            "title": "发帖规范"
        }, {
            "UID": "67890",
            "title": "论坛精华"
        }
    ]
}
```

在实际，可能由于各种原因，获取到的配置字符串会有以下几种异常情况：

- 整个字符串都是一个空字符串 ""
- 字符串由于错误编辑，不合法，或者是格式错误
- "top" 字段可能是 `null`，或者是空字符串

如果按照传统的逻辑，需要对这些异常情况一一处理。但如果开发者不需要关心这些异常，只关心合法的配置。那么我们完全可以利用 `MustXxx` 函数必然返回一个 `*V` 对象的特点，将逻辑简化如下：

```go
    c := jsonvalue.MustUnmarshalString(confString) // 假设 confString 是获取到的配置字符串
    for _, v := range c.Get("top").ForRangeArr() {
        feeds = append(feeds, &Feed{               // 将帖子主题追加到返回列表中，假设帖子的结构体为 Feed
            ID:    v.MustGet("UID").String(),
            Title: v.MustGet("title").String(),
        }) 
    }
```

---

## jsonvalue.V 对象的属性

首先我们要了解一下 JSON 官方定义的一些属性，然后再说明这些属性在 `jsonvalue` 中是如何体现的。

### 官方定义

在标准的 [JSON 规范](https://www.json.org/json-en.html)中，规定了以下的几个概念：

- 一个有效的 JSON 值，称为一个 JSON 的 `value`。在本工具包中，则使用一个 `*V` 来表示一个 JSON value
- JSON 值的类型有以下几种：

|   类型    | 说明                                                                                                          |
| :-------: | :------------------------------------------------------------------------------------------------------------ |
| `object`  | 也就是一个对象，对应着一个 K-V 格式的值。其中 K 必然是一个 string，而 V 则是有效的 JSON `value`               |
|  `array`  | 一个数组，对应着一系列 `value` 的有序组合                                                                     |
| `string`  | 字符串类型，这很好理解                                                                                        |
| `number`  | 数字型，准确地说，是双精度浮点数                                                                              |
|           | 由于 JSON 是基于 JavaScript 定义的，而 JS 中只有 double 这一种数字，所以 number 实际上就是 double。这是个小坑 |
| `"true"`  | 表示布尔 “真”                                                                                                 |
| `"false"` | 表示布尔 “假”                                                                                                 |
| `"null"`  | 表示空值                                                                                                      |

### jsonvalue 基础属性

在 `*jsonvalue.V` 对象中，参照绝大多数 JSON 工具包的做法，将 `"true"` 和 `"false"` 合并为一个 `Boolean` 类型。此外，将 `"null"` 也映射为一个 `Null` 类型。

此外，还定义了一个 `NotExist` 类型，表示当前不是一个合法的 JSON 对象。此外还有一个 `Unknown`，开发者可以不用关心，使用中不会出现这个值。

使用以下函数，可以获得 value 的类型属性：

```go
func (v *V) ValueType() ValueType
func (v *V) IsObject()  bool
func (v *V) IsArray()   bool
func (v *V) IsString()  bool
func (v *V) IsNumber()  bool
func (v *V) IsBoolean() bool
func (v *V) IsNull()    bool
```
