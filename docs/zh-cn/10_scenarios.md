<font size=6>特殊应用场景</font>

[上一页](./09_conversion.md) | [总目录](./README.md) | [下一页](./11_comparation.md)

---

- [JSON object 中 key 的顺序问题](#json-object-中-key-的顺序问题)
	- [获取原始 JSON 串的 key-value 顺序](#获取原始-json-串的-key-value-顺序)
	- [按照 object 的 key 被设置的顺序进行序列化](#按照-object-的-key-被设置的顺序进行序列化)
	- [JSON 值的幂等](#json-值的幂等)
- [特殊浮点数 +/-Inf 和 NaN](#特殊浮点数--inf-和-nan)
- [忽略 Go 结构体的 omitempty 标记](#忽略-go-结构体的-omitempty-标记)
- [获取数字类型的原始串](#获取数字类型的原始串)

---

正如[概述](./01_introduction.md)中所说，除了 jsonvalue 最典型的应用之外，我也支持了一些奇奇怪怪的场景。

本小节尽可能罗列笔者添加的一些特殊功能背后的逻辑，并给出具体的说明。

---

## JSON object 中 key 的顺序问题

按照标准 JSON 的实现来说，object 类型的 JSON 就是一组 K-V，这是没有顺序属性的。但是在实际应用中，我发现有好几个团队的同事，把 object 类型玩成了一个有序的 KV……

这就带来了下面的强需求：

1. 解析 JSON 串的时候，希望能够获取到 object 中 key 的顺序。
2. 反过来，对 JSON 进行序列化的时候，也希望能够指定 key 的顺序。

哎，用原生 `encoding/json` 的话呢，肯定是不支持的～～用其他的 JSON 库呢，倒也不是不能支持，但是会有点绕。

其实最开始我对这种场景感到有点无语（非标准的 JSON 场景），不过嘛既然是个需求，而且好多人遇到了，于是我就研究了一波，并且做出了实现。

### 获取原始 JSON 串的 key-value 顺序

咱们先来看第一个需求：解析 JSON 串的时候，希望能够获取到 object 中 key 的顺序。

首先有一个 JSON 串，比如：

```go
const raw = `{"a":1,"b":2,"c":3}`
```

咱们就先按照常规的反序列化逻辑去解析:

```go
v := jsonvalue.MustUnmarshalString(raw)
```

如果想要按顺序获取 key，可以调用 `RangeObjectsBySetSequence` 函数:

```go
keys := []string{}
v.RangeObjectsBySetSequence(func(key string, _ *V) bool {
    keys = append(keys, key)
    return true
})
fmt.Println(keys)
```

保证每次你都可以按顺序获得 `[a, b, c]`，也就是原始数据的顺序。

这个功能是从 1.3.1 版开始支持的。

至于第二个需求，则有好几种实现方法，下面也会一一提及。

---

### 按照 object 的 key 被设置的顺序进行序列化

从字面含义上，正如标题所示。笔者在这里给大家解释一下这是什么意思：

1. 如果这是一个程序生成的新的 jsonvalue 值，那么这个选项则表示按照每一个 key 被设置的先后顺序，进行 object 类型的序列化
2. 如果这段数据是从一个原始字节流反序列化出来的 jsonvalue 值的话，那么在序列化的时候，依然按照原始字节流中 key 的顺序进行序列化
3. 第三点就是前面两点的结合，该选项在序列化的时候呢，先序列化原始字节流中的 key，再按照新 key 被设置的顺序进行序列化

如果一个 key 先被 delete，然后再被 set 的话，以最新（最后）一次的顺序为准。

要实现这个功能，需要在序列化的时候传入额外参数 `OptSetSequence()`，如：

```go
const raw = `{"a":1,"b":2,"c":3}`
v := jsonvalue.MustUnmarshalString(raw) // 此时 key 的顺序就是原始的 a, b, c
v.Delete("b")                           // b 被删了，顺序变成 a, c
v.Set(4).At("d")                        // 新增一个 d，顺序是 a, c, d
v.Set(1).At("a")                        // 重新设置 a，最新位置优先，顺序是 c, d, a
s := v.MustMarshalString(OptSetSequence()) // 添加序列化选项
fmt.Println(s)
```

按照注释中的说明，可以确保得到以下顺序的字节流:

```go
{"c":3,"d":4,"a":1}
```

这个功能实际上在解析结构体的时候也会生效（`Import` 函数）。

---

### JSON 值的幂等

有一种情况是生成了 JSON 之后，对这段 JSON 计算 checksum 值。这经常用于 HTTP JSON API 中对请求包体进行签名或校验的场景，这个时候我们希望同样结构的 JSON，不论是在任何时候输出，都能够保证一模一样的序列。

这个需求与前述的 set 顺序还有点不一样，因为一系列 key 被设置的顺序可能不同。

一个最简单的方法是使用 `OptDefaultStringSequence()`，这个函数确保在序列化时，按照 string 的比较顺序来排序。实现上是简单封装了一下原生的 `strings.Compare` 函数。

如：

```go
// 由于 map 的 key 的随机性，a 和 b 谁先谁后其实是无法保证的
v := jsonvalue.New(map[string]any{
    "a": 1,
    "b": 2,
    "c": 3,
})
s := v.MustMarshalString(OptDefaultStringSequence())
fmt.Println(s)
```

上文中使用 `OptDefaultStringSequence` 可以确保每次输出都是 `{"a":1,"b":2,"c":3}`

---

## 特殊浮点数 +/-Inf 和 NaN

使用原生 `encoding/json`，当你把一个浮点数设置为正/负无限，或者是 NaN 的时候，比如以下的简单函数: 

```go
func main() {
	_, err := json.Marshal(map[string]float64{
		"score": math.Inf(-1),
	})
	fmt.Println(err)
}

```

你会简单干脆地拿到一个错误: `json: unsupported value: -Inf`。

可能读者觉得这有什么，避开就好了呗，没事干嘛设这个值。但是，笔者是做推荐系统的，在算法评分中得出 -Inf 值并不奇怪。而我们的原始协议采用的是 protobuf，完美支持 IEEE-754 所定义的所有浮点数格式。但是在 JSON 中，行不通……

最好的解决方式是与对端协商，采用 protobuf 传递数据。但是在实际情况下，很多 API 实际上是使用 JSON 来通信的，改成 protobuf 的话改动量很大。

这种情况下，我们依然可以使用 jsonvalue，指定在序列化时，针对 Inf 和 NaN 进行不同的替换逻辑。大体上有三类：

- 替换成 JSON null
- 替换成另一个浮点数
- 替换成 JSON string

比如上面的例子，我们就可以改写为:

```go
func main() {
	v, _ := jsonvalue.Import(map[string]float64{
		"score": math.Inf(-1),
	})
	s := v.MustMarshalString(jsonvalue.OptFloatInfToFloat(23333))
	fmt.Println(s)
}
```

得到的输出为: `{"score":-23333}`

其他的几种模式请参见 "[额外选项配置](./12_option.md)" 章节。

---

## 忽略 Go 结构体的 omitempty 标记

这个需求是用在 `Import` 函数中的。简单而言，就是在将结构体转为 `jsonvalue` 的时候，希望能够无视结构体的字段后面的 `omitempty` 标记，而将结构完整导出。

比如说下面的这个 struct:

```go
type st struct {
	A int `json:"a,omitempty"`
}

func main() {
	st := st{}
	b, _ := json.Marshal(&st)
	fmt.Println(string(b))
}
```

输出是 `{}`，但是按照这个需求，希望输出的是 `{"a":0}`

当时我接到这个 issue 的时候（公司内部的，所以在 Github 上看不到），我是一脸黑人问号的——如果确实不想 `omitempty`，那为什么还要定义这个呢？

![黑人问号？](https://bkimg.cdn.bcebos.com/pic/8cb1cb1349540923e1860cc29958d109b2de499a?x-bce-process=image/resize,m_lfit,w_536,limit_1/format,f_jpg)

进一步了解之后我才知道，issuer 应用的场景是对 protobuf 生成的 Go 代码中的 struct 进行转换，而 protoc 的 Go 工具生成的 Go 结构体字段中，都会加上 `omitempty` 标记。但是我们就是有需求获取完整的数据结构，那咋办呢？

这还真不是一个无厘头的需求，甚至在 stackoverflow 上都有专门的提问：[golang protobuf remove omitempty tag from generated json tags](https://stackoverflow.com/questions/34716238/)

当然用 `jsonpb` 可以解决这个问题，但是 `jsonpb` 也有不满足一些应用场景的情况（其实就是前文的指定 key 顺序功能），所以当时同事选用了我的 `jsonvalue`，并且也同时向我提出了这个需求。

实现这个 feature 不难，我添加了一个 `OptIgnoreOmitempty()` 选项，这是第一个也是目前唯一一个用于 `Import` 函数的 option。上面的代码可以改造为:

```go
type st struct {
	A int `json:"a,omitempty"`
}

func main() {
	st := st{}
	v, _ := jsonvalue.Import(&st, jsonvalue.OptIgnoreOmitempty())
	s := v.MustMarshalString()
	fmt.Println(s) // 输出: {"a":0}
}
```

---

## 获取数字类型的原始串

很简单，对于 JSON 来说，并没有明确定义数字类型的位宽和定义。如果用 `map[string]any` 来承载的话，`encoding/json` 会将数字统一解析为 `float64` 类型。但是如果我们不小心把 `int64` 或 `uint64` 写入 JSON 中，由于浮点数的精度损失，当数值过大时，最低几位将会丢失。

这种问题，在使用 `uint64` 来传递哈希值的时候，就成为硬伤。对 `encoding/json` 熟悉的同学看了标题之后可能会说：这题我会，用 `json.RawMessage` 啊。

当然是因为用 `encoding/json` 不方便才用 jsonvalue 的呀！这种场景往往就是层级深、或者是结构不确定的情况下，原生 json 实在是不方便。

只要确保对端传过来的数字是小于或等于 `uint64 MAX` 的值，那么使用 jsonvalue 可以精确地解析并获取到原始数值，使用 `v.Uint64()` 函数即可。当然，也可以使用 `v.String()` 直接获得原始字符串。



