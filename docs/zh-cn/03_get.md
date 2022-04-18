# 解析并获取 JSON

[上一页](./02_quick_start.md) | [总目录](./README.md) | [下一页](./04_set.md)

---

- [反序列化: Unmarshal 系列函数](./03_get.md#unmarshal-系列函数)
- [jsonvalue.V 对象的属性](./03_get.md#jsonvaluev-对象的属性)
- [Get 系列函数](./03_get.md#get-系列函数)
- [迭代 Object 和 Array 的成员](./03_get.md#迭代-object-和-array-的成员)
- [MustGet 和相关函数](./03_get.md#mustget-和相关函数)

---

## Unmarshal 系列函数

### 基础的 Unmarshal

在 jsonvalue 中，采用 Go 的 marshal / unmarshal 语义描述序列化和反序列化的过程。

在原生 `encoding/go` 中，`json.Unmarshal` 的入参是一个 `[]byte` 类型。在 jsonvalue 中类似，可以使用以下函数来解析 JSON 字节串

```go
func Unmarshal(b []byte) (ret *V, err error)
```

不论是否 error，该函数都会返回一个 `*jsonvalue.V` 对象

- 当 JSON 字节串非法时，会返回 error 信息，此时返回的 json 对象的类型等于 `jsonvalue.NotExist`

### 其他 Unmarshal 函数

在实际操作中，JSON 字节串经常会以 `string` 而不是 `[]byte` 的格式出现。如果进行 `string(b)` 进行转换的话，实际上会进行一次内存拷贝。为了节省这个拷贝的开销，jsonvalue 也提供了入参为 `string` 的 unmarshal 函数:

```go
func UnmarshalString(s string) (ret *V, err error)
```

此外，当程序不需要关心 JSON 字节串格式是否正确的时候，也可以使用 `Must...` 系列的 unmarshal 函数：

```go
func MustUnmarshal(b []byte) *V
func MustUnmarshalString(s string) *V
```

这两个函数与前面一样，必然会返回一个非空的 `jsonvalue.V` 对象，但是不返回 `error` 类型，便于开发者编写一些极为简短的逻辑代码。这种简短代码的技巧，在本页面的最后会进行介绍。

---

## jsonvalue.V 对象的属性

首先我们要了解一下 JSON 官方定义的一些属性，然后再说明这些属性在 `jsonvalue` 中是如何体现的。

### 官方定义

在标准的 [JSON 规范](https://www.json.org/json-en.html)中，规定了以下的几个概念：

- 一个有效的 JSON 值，称为一个 JSON 的 `value`。在本工具包中，则使用一个 `*V` 来表示一个 JSON value
- JSON 值的类型有以下几种：

|类型|说明|
|:---:|:---|
|`object`|也就是一个对象，对应着一个 K-V 格式的值。其中 K 必然是一个 string，而 V 则是有效的 JSON `value`|
|`array`|一个数组，对应着一系列 `value` 的有序组合|
|`string`|字符串类型，这很好理解|
|`number`|数字型，准确地说，是双精度浮点数|
||由于 JSON 是基于 JavaScript 定义的，而 JS 中只有 double 这一种数字，所以 number 实际上就是 double。这是个小坑|
|`"true"`|表示布尔 “真”|
|`"false"`|表示布尔 “假”|
|`"null"`|表示空值|

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

---

## Get 系列函数

### 函数参数含义

Get 系列函数是 jsonvalue 中读取 JSON 信息的核心函数。函数格式如下：

```go
func (v *V) Get(param1 interface{}, params ...interface{}) (*V, error)
```

实际的使用示例为：

```go
const raw = `{"someObject": {"someObject": {"someObject": {"message": "Hello, JSON!"}}}}`
child, _ := jsonvalue.MustUnmarshalString(s).Get("someObject", "someObject", "someObject", "message")
fmt.Println(child.String())
```

在上面的例子中，`Get` 函数的参数的含义为：

- 获取 `*V` 对象中，key 为 `someObject` 的 object 值，再从这个值中，获取 key 为 `someObject` 的 object 值……
  - 如果写成域格式，则相当于 `child = v.someObject.someObject.someObject.message`

`Get` 函数的参数是一个 `interface{}`，但实际上，这个函数只接受两种类别的参数：一是字符串类型，二是整型数字（有符号无符号均可）。

`Get` 函数会解析入参，迭代检查每一个参数的类型，从而决定下一轮迭代的逻辑：

如果当前层级的参数是一个字符串时，则：

- 如果当前的 `jsonvalue` 对象是一个 `Object` 类型时，则查找当浅层级字符串所指定的 value
  - 如果找到，若有下一层参数，则使用当前 value 和下一层参数继续迭代查找
  - 如果无法找到参数所指定的错误，则返回类型为 `NotExist` 的对象，以及 [ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue@v1.1.1#pkg-constants) 错误。
  - 如果当前的对象不是 `Object` 类型，则返回 `NotExist` 对象以及 [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue@v1.1.1#pkg-constants) 错误。

如果当前层级的参数是一个整数时，则：

- 如果当前的对象是一个 `Array` 类型时，则将整数视为 index 参数，查找在指定 index 中是否包含 JSON value。此时 Index 的含义如下：
  - 当 Index >= 0，则按照正常的切片下标逻辑来查找。如果 JSON array 的长度不足，则返回 `NotExist` 对象以及 [ErrOutOfRange](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue@v1.1.1#pkg-constants) 错误。
  - 当 Index < 0，则视为 “倒数第几个” 的语义，但最多依然不大于 JSON array 的长度。比如说 array 长度为 5，那么 -5 会返回下标为 0 的子成员，而 -6 则会返回错误。
  - 如果找到，则根据后续参数情况继续迭代或返回。如果无法找到，则返回 `NotExist` 对象以及 [ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue@v1.1.1#pkg-constants) 错误。
- 如果当前的对象不是 `Array` 类型，则返回 `NotExist` 对象以及 [ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue@v1.1.1#pkg-constants) 错误。

相信开发者会有[这样的一个疑问](https://github.com/Andrew-M-C/go.jsonvalue/issues/4)：为什么输入参数要强行切为两个部分，而不是直接一个 `...interface{}` 就搞定呢？

- 理由是：这是为了是避免出现 `v.Get()` 这样的笔误。让函数至少需要一个参数，就可以在编译阶段就检查出类似的错误，而不会带到线上程序中。
- 如果开发者需要传入类似参数的话，那么开发者需要检查 `[]interface{}` 参数的长度是否大于一；如果能确保大于一的话，可以采用 `v, _ := Get(para[0], para[1:]...)` 的格式进行调用。

### GetXxx 系列函数

实际操作中，开发者完全不关心 `*V` 对象本身，而只关心它所承载的值。在开发者可以确定或限定某个字段只能是某个值的时候，可以使用以下函数：

```go
func (v *V) GetObject (param1 interface{}, params ...interface{}) (*V, error)
func (v *V) GetArray  (param1 interface{}, params ...interface{}) (*V, error)
func (v *V) GetBool   (param1 interface{}, params ...interface{}) (bool, error)
func (v *V) GetString (param1 interface{}, params ...interface{}) (string, error)
func (v *V) GetBytes  (param1 interface{}, params ...interface{}) ([]byte, error)
func (v *V) GetInt    (param1 interface{}, params ...interface{}) (int, error)
func (v *V) GetInt32  (param1 interface{}, params ...interface{}) (int32, error)
func (v *V) GetInt64  (param1 interface{}, params ...interface{}) (int64, error)
func (v *V) GetNull   (param1 interface{}, params ...interface{}) error
func (v *V) GetUint   (param1 interface{}, params ...interface{}) (uint, error)
func (v *V) GetUint32 (param1 interface{}, params ...interface{}) (uint32, error)
func (v *V) GetUint64 (param1 interface{}, params ...interface{}) (uint64, error)
func (v *V) GetFloat32(param1 interface{}, params ...interface{}) (float32, error)
func (v *V) GetFloat64(param1 interface{}, params ...interface{}) (float64, error)
```

这些函数都有以下共同点：

- 如果参数指定的子值存在，并且类型匹配上，那么 error 字段为 nil；除了 `GetNull` 函数之外，其他函数都会返回对应的值。
- 如果参数指定的子值不存在，或者值存在但是类型不匹配，则 error 必然非 nil；但在不同的情况下，返回的 error 值会有不同。

此外，这些函数并不是简简单单地只是匹配类型并返回，它们还拥有更加方便的功能，在后续小节中会着重说明，这里笔者先举一个小例子：

比如很多情况下，我们可能需要使用 JSON 的 string 类型，实际上承载数字值，比如: `{"number":"12345"}`。

按照 JSON 标准的定义，`number` 成员是一个 string 值。但使用 jsonvalue 的 GetInt 值，是能够正确获得数字值的：

```go
raw := `{"number":"12345"}`
n, err := jsonvalue.MustUnmarshalString(raw).GetInt("number")
fmt.Println("n =", n)
fmt.Println("err =", err)
```

输出内容为：

```
n = 12345
err = not match given type
```

可见，n 和 err 都返回了值，这算是打了巴掌又给糖吃（笑）——一方面尽责地帮开发者解析 JSON 内容，另一方面还是提示开发者数据可能存在的错误。当然如果这是预期范围内的正常错误的话，在程序中完全可以忽略。

---

## 迭代 Object 和 Array 的成员

对于基础类型（number, string, boolean, null），我们只关心它的一个值。但对于复杂类型（object, array），我们有必要关心其中的各种结构。除了使用 `Get` 系列函数之外，jsonvalue 还提供了 iter 函数，有以下两种风格：

### 回调函数风格

```go
func (v *V) RangeArray  (callback func(i int, v *V) bool)
func (v *V) RangeObjects(callback func(k string, v *V) bool)
```

使用回调函数的风格来迭代 object 和 array 中的每一个成员。如果回调函数返回 true，则继续迭代；返回 false 则会中止迭代，退出回调。

### for-range 风格

```go
func (v *V) ForRangeArr() []*V
func (v *V) ForRangeObj() map[string]*V
```

这种模式返回了一个预先存好了 kv 信息的 channel，并且已经 close 了，因此开发者可以使用 `for` 语法，进行更加直观的开发：

```go
    v := jsonvalue.MustUnmarshalString(`["A","B","C","D"]`)
    for i, v := range v.ForRangeArr() {
        fmt.Println(i, "-", v)
    }
```

在 `ForRangeObj` 函数中，由于 jsonvalue 是使用 map 来实现 object 的 kv 存储，因此 key 的顺序不予保证。

在命名的角度上，由于历史原因，笔者先开发了 RangeXxx 系列函数，所以导致 for-range 风格的函数反而不使用 range 命名，还请开发者们谅解。

## MustGet 和相关函数

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

