# 创建并序列化 JSON

[上一页](./03_get.md) | [总目录](./README.md) | [下一页](./05_import_export.md)

---

[TOC]

---

## 创建 JSON 值

在 jsonvalue 中提供了一系列的 `NewXxx` 函数，用于创建指定类型的 JSON 值。在绝大部分情况下，我们要创建的最外层 JSON 值是一个 object 或者是 array 类型值。此时我们可以使用以下的两个函数：

```go
anObj := jsonvalue.NewObject()
anArr := jsonvalue.NewArray()
```

也可以创建其他的基础类型值，如：

```go
func NewInt    (i int)     *V
func NewString (s string)  *V
func NewFloat64(f float64) *V
func NewBool   (b bool)    *V
func NewNull() *V
```

---

## Set 系列函数

在创建了最外层的 object 或者是 array 之后，下一步就是构建 JSON 的内部结构。相对于上一小节的 `Get` 系列函数，jsonvalue 则提供了 `Set` 系列函数来处理（目标类型为 object 为主）的 JSON 子结构的创建。

### 基础用法

Set 系列函数，一般使用以下的模式进行调用：

```go
v.Set(child).At(path...)
```

对应英语中的语法：`SET some sub value AT some position.`

实际操作中，我们经常直接设置指定的基础类型值，如：

```go
v := jsonvalue.NewObject()
v.Set("Hello, JSON!").At("data", "message")
fmt.Println(v.MustMarshalString())
```

输出: `{"data":{"message":"Hello, JSON!"}}`

### At 参数语义

可以看到，通过 `Set` 系列函数后，还需要紧跟 `At` 函数来将欲设置的值落地到真正的 JSON 结构中。因此 `At` 函数的参数自然是重点。`At` 函数的原型如下：

```go
func (s *Set) At(param1 interface{}, params ...interface{}) (*V, error)
```

At 函数的参数语义，与前文提及的 `Get` 函数语义基本一致。同样地，为了防止编程错误，这个函数至少需要传一个参数。

不过函数更为重要的是自动创建目标结构的能力：

- `At` 函数在在指定位置上设置子值时，首先会采用与 `Get` 函数类似的逻辑，层层迭代找到目标结构，然后在指定层级的指定 key 或 index 中设置子值。
- 如果目标结构不存在，则按照参数中指定的参数类型创建相应的结构。同样是 string 类型对应 object，整型对应 array。

以下例子中，自动创建了数据结构：

```go
v := jsonvalue.NewObject()                   // {}
v.Set("Hello, object!").At("obj", "message") // {"obj":{"message":"Hello, object!"}}
v.Set("Hello, array!").At("arr", 0)          // {"obj":{"message":"Hello, object!"},"arr":["Hello, array!"]}
```

在 At() 自动创建数组的逻辑其实稍微有点复杂，需要解释一下：

- 当调用方在参数中指定在某个尚未存在的数组中设置一个值的时候，那么 `At` 指定的下标（position）数字， 应当为0，操作才能成功
- 当数组已经存在，那么 `At` 指定的位置数，要么在数组中已存在，要么正好等于数组的长度，当后者的情况下，会在数组的最后追加值。
- 如果在 At() 中指定了数组的某个已存在值的下标，那么那个位置上的值会被替换掉，请注意。

这个特性在使用 for-range 块时会非常有用，比如：

```go
    const words = []string{"apple", "banana", "cat", "dog"}
    const lessons = []int{1, 2, 3, 4}
    v := jsonvalue.NewObject()
    for i := range words {
        v.Set(words[i]).At("array", i, "word")
        v.Set(lessons[i]).At("array", i, "lesson")
    }
    fmt.Println(c.MustMarshalString())
```

最终输出为:

```json
{"array":[{"word":"apple","lesson":1},{"word":"banana","lesson":2},{"word":"cat","lesson":3},{"word":"dog","lesson":4}]}
```

---

## Append 和 Insert 系列函数

函数 `Append` 和 `Insert` 专门针对数组操作使用。其中 `Append` 函数需搭配 `InTheBeginning` 和 `InTheEnd` 函数，而 `Insert` 则搭配 `After` 和 `Before`

对应着以下几个英语语法：

- Append some value in the beginning of ...
- Append some value to the end of ...
- Insert some value after ...
- Insert some value before ...

这几个函数的原型如下：

```go
func (v *V) Append(child interface{}) *Append
func (apd *Append) InTheBeginning(params ...interface{}) (*V, error)
func (apd *Append) InTheEnd      (params ...interface{}) (*V, error)

func (v *V) Insert(child interface{}) *Insert
func (ins *Insert) After (firstParam interface{}, otherParams ...interface{}) (*V, error)
func (ins *Insert) Before(firstParam interface{}, otherParams ...interface{}) (*V, error)
```

基本语义与前问的 `Set` 和配套函数基本一致，但有以下几点小差异：

- `InTheBeginning` 和 `InTheEnd` 允许空参数，此时表示当前的 value 就已经是一个数组，语义是在当前数组的开头或末尾追加子值。
- `After` 和 `Before` 的最后一个参数（如果只有一个参数，则最后一个即为第一个）必须是一个整型数字，代表在数组中的下标位。与 `Set(...).At(...)` 类似，允许负下标。

---

## Marshal 系列函数

与 Unmarshal 对应，jsonvalue 的序列化函数也采用其相对的 marshal 语义。提供了以下四个方法：

```go
func (v *V) Marshal          (opts ...Option) (b []byte, err error)
func (v *V) MarshalString    (opts ...Option) (s string, err error)
func (v *V) MustMarshal      (opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

在当前版本下，marshal 只有两种情况会报错：

- `*V` 是 `NotExist` 类型
- 值中包含了不合法的浮点数值 `+Inf`, `-Inf` 或 `NaN`，并且没有明确说明如何处理这些数值。后文还会提到，jsonvalue 针对这些非法的 number 数值，还提供了额外的处理能力。
- 选项参数 opts 中包含非法配置

因此如果开发者能够确定规避掉上述错误的话，完全可以使用 `MustMarshal` 系列函数。
