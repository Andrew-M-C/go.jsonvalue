<font size=6>值的自动转换</font>

[上一页](./08_caseless.md) | [总目录](./README.md) | [下一页](./10_scenarios.md)

---

- [自动转换简介](#自动转换简介)
- [String 转 number](#string-转-number)
- [String 和 number 转 boolean](#string-和-number-转-boolean)
- [GetString 和 MustGet(...).String() 的区别](#getstring-和-mustgetstring-的区别)

---

## 自动转换简介

在 Go 原生 `encoding/json` 中，有一个 `tag` 称为 "string"，它的作用是将当前值强制序列化为一个 JSON string。比如：

```go
type st struct {
    Num int `json:"num,string"`
}

func main() {
    st := st{Num: 12345}
    b, _ := json.Marshal(&st)
    fmt.Println(string(b))
}
```

输出 `{"num":"12345"}`，可见 `num` 字段被序列化成了一个字符串而不是数字类型。

在 jsonvalue 中，支持对这种类型的字符串进行直接转换读取。但是针对不同类型，读取模式有细微差别，具体说明如下：

---

## String 转 number

所有 number 类型的 `GetXxx` 系列函数，均支持从 string 类型值中提取 number 字段。

但同时，只要目标 value 类型不是 number，必然会返回 `error` 信息以标记。在不同的情况下会返回不同的 error（采用 `errors.Is()` 函数判断）：

- 如果目标查找不到，则返回 `ErrNotFound`
- 如果目标是一个 number，那么直接返回数字，err 为 nil
- 如果目标不是 number 也不是 string，err 为 `ErrTypeNotMatch`
- 如果目标是一个 string，则解析 string 中的数字：
    - 如果解析成功，则 err 为 `ErrTypeNotMatch`
    - 如果解析失败（不是一个合法的数字），则 err 为 `ErrParseNumberFromString` 并带有具体的错误原因

比如：

```go
	const raw = `{"legal":"1.25E-1","illegal":"ABCD"}`
	v := jsonvalue.MustUnmarshalString(raw)

	f, err := v.GetFloat64("legal")
	fmt.Println("01 - float:", f)
	fmt.Println("01 - err:", err)

	f, err = v.GetFloat64("illegal")
	fmt.Println("02 - float:", f)
	fmt.Println("02 - err:", err)
```

输出内容：

```
01 - float: 0.125
01 - err: not match given type
02 - float: 0
02 - err: failed to parse number from string: parsing number at index 0: zero string
```

---

## String 和 number 转 boolean

与 number 类似，string 也可以承载布尔值。同时，number 也可以转成 boolean 值。

String 转 boolean 只有一种情况会返回 `true`，那就是字符串字面值完全等同于小写的 `"true"`，其他情况均返回 false。

Number 转 boolean 则是判断数字值是否不等于 0。

无论在什么情况下，只要目标值存在但不为 Boolean 类型，err 均会返回 `ErrTypeNotMatch`。

---

## GetString 和 MustGet(...).String() 的区别

前文所述的 `GetXxx` 自动转换函数，在其对应的 `Xxx` 函数中，逻辑也是一致的，不同之处仅仅是 `Xxx` 函数并不返回 error 信息。比如：

```go
	const raw = `{"legal":"1.25E-1","illegal":"ABCD"}`
	v := jsonvalue.MustUnmarshalString(raw)
	f := v.MustGet("legal").Float64()
	fmt.Println("float:", f)
```

输出：

```
float: 0.125
```

但是有一种情况例外，那就是 `GetString` 和 `String` 函数。

`GetString` 的语义是获取一个 string 类型的 JSON 值。当目标值找不到时，返回 error；当目标值不是 string 类型时，返回 `ErrTypeNotMatch` 错误，不会进行其他转换。

`String` 函数比较特别，它既能够提取 `*V` 的 string 类型值，又需要符合 Go 的 `String()` 函数要求，也就是用在 `fmt` 包中的 `%v` 关键字中。

因此 `String` 函数的逻辑如下：

- 如果当前 `*V` 是一个 string 类型，则返回 string 值；
- 如果当前是一个 number 类型，则返回数字的字符串；
- 如果当前是一个 null，则返回 `null`；
- 如果当前是一个 boolean，则返回 `true` 或 `false`；
- 如果当前是一个 object，则返回 `{k1: v1, k2: v2}` 格式的字符串，并且不会刻意去进行转义，也不符合 JSON 格式；
- 如果当前是一个 array，则返回 `[v1, v2]` 格式的字符串，也同样不会刻意去转义。

