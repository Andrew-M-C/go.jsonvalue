
<font size=6>创建并序列化 JSON</font>

[上一页](./02_quick_start.md) | [总目录](./README.md) | [下一页](./04_get.md)

---

本小节说明如何设置和生成一个 jsonvalue 值。

---

- [创建 JSON 值](#创建-json-值)
- [往 jsonvalue 中设置值](#往-jsonvalue-中设置值)
  - [基础用法](#基础用法)
  - [At 参数语义](#at-参数语义)
- [往 JSON 数组中添加值 —— Append 和 Insert 系列函数](#往-json-数组中添加值--append-和-insert-系列函数)

---

## 创建 JSON 值

在绝大部分情况下，我们要创建的最外层 JSON 值是一个 object 或者是 array 类型值。此时我们可以使用以下的两个函数：

```go
o := jsonvalue.NewObject()
a := jsonvalue.NewArray()
```

也可以指定任意可以合法地转换成 JSON 的 Go 类型，使用 `New` 函数直接创建 JSON 值。比如上面的对象和数组类型值，也可以用这种方式创建:

```go
o := jsonvalue.New(struct{}{})  // 生成一个 JSON object
a := jsonvalue.New([]int{})     // 生成一个 JSON array
```

如果你想要新建的是简单的 JSON 元素，也可以创建其他的合法类型，如：

```go
i := jsonvalue.New(100)             // 生成一个 JSON number
f := jsonvalue.New(188.88)          // 生成一个 JSON number
s := jsonvalue.New("Hello, JSON!")  // 生成一个 JSON string
b := jsonvalue.New(true)            // 生成一个 JSON boolean
n := jsonvalue.New(nil)             // 返回一个 JSON null
```

---

## 往 jsonvalue 中设置值

在创建了最外层的 object 或者是 array 之后，下一步就是构建 JSON 的内部结构。相对于上一小节的 `Get` 系列函数，jsonvalue 则提供了 `Set` 和 `MustSet` 系列函数来处理 JSON 子结构的创建。

`Set` 和 `MustSet` 方法的差别是: 前者会返回设置后的子 `*jsonvalue.V` 对象和 `error` 类型值, 而后者则不。如果调用方不关心是否设置成功 (或者有把握设置成功), 那么可以使用 `MustSet` 系列函数, 这也可以避免 golangci-lint 的告警提示。

### 基础用法

Set 系列函数，一般使用以下的模式进行调用：

```go
v.MustSet(child).At(path...)
```

对应英语中的语法：`SET value AT some position.`，请注意，value 在前，path 在后

目前 jsonvalue 的函数使用 `any`, 因此获得了一个类似于泛型的体验，如：

```go
v := jsonvalue.NewObject()
v.MustSet("Hello, JSON!").At("data", "message")
fmt.Println(v.MustMarshalString())
```

输出: `{"data":{"message":"Hello, JSON!"}}`

### At 参数语义

可以看到，通过 `Set` 系列函数后，还需要紧跟 `At` 函数来将欲设置的值落地到真正的 JSON 结构中。因此 `At` 函数的参数自然是重点。`At` 函数的原型如下：

```go
type Setter interface {
	At(firstParam interface{}, otherParams ...interface{}) (*V, error)
}
```

At 函数的参数语义，与前文提及的 `Get` 函数语义基本一致。同样地，为了防止编程错误，这个函数至少需要传一个参数。

不过函数更为重要的是自动创建目标结构的能力：

- `At` 函数在在指定位置上设置子值时，首先会采用与 `Get` 函数类似的逻辑，层层迭代找到目标结构，然后在指定层级的指定 key 或 index 中设置子值。
- 如果目标结构不存在，则按照参数中指定的参数类型创建相应的结构。同样是 string 类型对应 object，整型对应 array。

以下例子中，自动创建了数据结构：

```go
v := jsonvalue.NewObject()                       // {}
v.MustSet("Hello, object!").At("obj", "message") // {"obj":{"message":"Hello, object!"}}
v.MustSet("Hello, array!").At("arr", 0)          // {"obj":{"message":"Hello, object!"},"arr":["Hello, array!"]}
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
        v.MustSet(words[i]).At("array", i, "word")
        v.MustSet(lessons[i]).At("array", i, "lesson")
    }
    fmt.Println(c.MustMarshalString())
```

最终输出为:

```json
{"array":[{"word":"apple","lesson":1},{"word":"banana","lesson":2},{"word":"cat","lesson":3},{"word":"dog","lesson":4}]}
```

---

## 往 JSON 数组中添加值 —— Append 和 Insert 系列函数

函数 `Append` 和 `Insert` 专门针对数组操作使用。其中 `Append` 函数需搭配 `InTheBeginning` 和 `InTheEnd` 函数，而 `Insert` 则搭配 `After` 和 `Before`

对应着以下几个英语语法：

- Append some value in the beginning of ...
- Append some value to the end of ...
- Insert some value after ...
- Insert some value before ...

与 `Set` 函数一样，请注意路径参数是后置的。此外, `Append` 和 `Insert` 也有其对应的 `MustAppend` 和 `MustInsert` 方法, 原因相同。

这几个函数的原型如下：

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

基本语义与前文的 `Set` 和配套函数基本一致，但有以下几点小差异：

- `InTheBeginning` 和 `InTheEnd` 允许空参数，此时表示当前的 value 就已经是一个数组，语义是在当前数组的开头或末尾追加子值。
- `After` 和 `Before` 的最后一个参数（如果只有一个参数，则最后一个即为第一个）必须是一个整型数字，代表在数组中的下标位。与 `Set(...).At(...)` 类似，允许负下标。


