# Jsonvalue - 方便而快速的 Go JSON 包

[![Travis](https://travis-ci.org/Andrew-M-C/go.jsonvalue.svg?branch=master)](https://travis-ci.org/Andrew-M-C/go.jsonvalue)
[![Coveralls](https://coveralls.io/repos/github/Andrew-M-C/go.jsonvalue/badge.svg?branch=master)](https://coveralls.io/github/Andrew-M-C/go.jsonvalue)
[![Go report](https://goreportcard.com/badge/github.com/Andrew-M-C/go.jsonvalue)](https://goreportcard.com/report/github.com/Andrew-M-C/go.jsonvalue)
[![Codebeat](https://codebeat.co/badges/ecf87760-2987-48a7-a6dd-4d9fcad57256)](https://codebeat.co/projects/github-com-andrew-m-c-go-jsonvalue-master)<br>
[![GoDoc](https://godoc.org/github.com/Andrew-M-C/go.jsonvalue?status.svg)](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue)
[![Latest](https://img.shields.io/badge/latest-v1.0.4-blue.svg)](https://github.com/Andrew-M-C/go.jsonvalue/tree/v1.0.4)
[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

[English](./README.md)

**Jsonvalue** 是一个用于处理 JSON 序列化/反序列化的 Go 语言包。当处理普通的 `struct` 不适合处理的 JSON 场景时，以前我们会使用 `map[string]interface{}` 来处理。但是这个方法其实有不少问题和困难。参见[这篇文章](https://cloud.tencent.com/developer/article/1676060)。这个包就是为了解决这个问题而开发的。

## 快速入门

有时候我们需要创建一个如下的复杂 JSON 对象：

```json
{
    "someObject": {
        "someObject": {
            "someObject": {
                "message": "Hello, JSON!"
            }
        }
    }
}
```

使用 `jsonvalue`，这个过程是很快的（[Playground](https://play.golang.org/p/u5846Wk6mq2)）：

```go
    v := jsonvalue.NewObject()
    v.SetString("Hello, JSON").At("someObject", "someObject", "someObject", "message")
    fmt.Println(v.MustMarshalString())
    // Output:
    // {"someObject":{"someObject":{"someObject":{"message":"Hello, JSON!"}}}
```

创建层级比较深的数组对象也类似（[Playground](https://play.golang.org/p/iTxnJDNdny3)）：

```json
[
    {
        "someArray": [
            "Hello, JSON!"
        ]
    }
]
```

```go
    v := jsonvalue.NewArray()
    v.SetString("Hello, JSON").At(0, "someObject", 0)
    fmt.Println(v.MustMarshalString())
    // Output:
    // [{"someObject":["Hello, JSON"]}]
```

对于更多信息，请参见 [godoc](https://godoc.org/github.com/Andrew-M-C/go.jsonvalue).

## 相比 map[string]interface{} 的优势

本人的文章《[还在用 map[string]interface{} 处理 JSON？告诉你一个更高效的方法——jsonvalue](https://cloud.tencent.com/developer/article/1676060)》中说明了。

主要是两个：

- 获取处于层级较深的数据很麻烦，需要进行多级的类型断言，这在代码中非常麻烦，还不如定义一个 struct 了。
- Marshal 效率低——因为在 marshal 时，需要对 `map[string]interface{}` 中的每一个值进行 reflect 类型判断。这就导致了 marshal 效率低下。

此外，使用 jsonvalue 还有一些额外的优势

### 支持 caseless

从 `v1.0.4` 版本开始，在进行 Get 操作中支持 caseless，也即是不区分大小写。该逻辑在 go 默认的 struct 中是默认逻辑，比如:

```go
package main

import (
    "fmt"
    "encoding/json"
)

type example struct {
    ID string `json:"id"`
}

func main(){
    s := `{"Id":"ID001"}`
    e := example{}

    json.Unmarshal([]byte(s), &e)
    fmt.Println(e.ID)

    // Output:
    // ID001
}
```

其中在 tag 中，指定了 ID 所对应的 JSON 字段名为 `id`，但是在实际传入中，使用的是 `Id`。Go 在进行 `Unmarshal` 时，是能够正确解析出来的。但是如果我们使用 `map[string]interface{}`，就无法识别了：

```go
package main

import (
    "fmt"
    "encoding/json"
)

func main(){
    s := `{"Id":"ID001"}`
    e := map[string]interface{}{}

    json.Unmarshal([]byte(s), &e)
    fmt.Println(e["id"])

    // Output:
    // <nil>
}
```

但是使用 `jsonvalue`，这个问题迎刃而解，而且还简化了代码：

```go
package main

import (
    "fmt"
    "github.com/Andrew-M-C/go.jsonvalue"
)

func main(){
    s := `{"Id":"ID001"}`
    v, _ := jsonvalue.UnmarshalString(s)
    id, _ := b.GetString("id")
    fmt.Println(id)

    // Output:
    // ID001
}
```

### 原样保留数字类型值文本

其实 JSON 中对于数字是没有整型和浮点型的区别的，原因是 Javascript 中没有整型，所有的数字统一采用 IEEE 754 双精度浮点型来表示。
在 `map[string]interface{}` 中对于数字类型，也是采用 `float64` 来保存数字的。这就导致当在 JSON 中使用 64 位整型时，如果采用 map，那么可能会导致整型精度的损失。

这个问题的标准解决方案其实是改用 string 来传递相关的类型值。但是对于历史问题的参数，或者是其他特殊需求而必须采用数值型的话，那么可以采用 jsonvalue 来对原始值进行解析。
针对这个场景，jsonvalue 提供了以下功能：

- 可以使用 `v.IsFloat()`, `v.GreaterThanInt64Max()` 等函数，判断该 JSON 数值的原始信息。
- 对于数值类型的 JSON，使用 `v.String()` 函数，获取原始的 JSON 文本内容，一字不差，不会因为一次反序列化+序列化的过程而损失信息。

```go
package main

import (
	"fmt"
	"github.com/Andrew-M-C/go.jsonvalue"
)

func main() {
	s := `{"num":12.001000}`
	v, _ := jsonvalue.UnmarshalString(s)
	n, _ := v.Get("num")
	fmt.Println(n.String())
	fmt.Println(v.MustMarshalString())
	// Output:
	// 12.001000
	// {"num":12.001000}
}
```

此外，从 v1.0.5 开始，jsonvalue 将会支持获取保存在 string 类型 JSON 值中的数字，举例如下：

```go
package main

import (
	"fmt"
	"github.com/Andrew-M-C/go.jsonvalue"
)

func main() {
	s := `{"num":"12"}`
	v, _ := jsonvalue.UnmarshalString(s)
	n, _ := v.Get("num")
	fmt.Println(n.String())
    fmt.Println(n.Int())
    fmt.Println(n.MustMarshalString())
	// Output:
	// 12
	// 12
    // "12"
}
```
