
<font size=6>大小写</font>

[上一页](./07_iteration.md) | [总目录](./README.md) | [下一页](./09_option.md)

---

- [Go 原生 json 的问题](#go-原生-json-的问题)
- [在 jsonvalue 中忽略大小写](#在-jsonvalue-中忽略大小写)

---

## Go 原生 json 的问题

在 Go 原生的 `encoding/json` 中，针对 `struct` 中的 JSON tag，处理原则是不区分大小写的。比如下面的例子：

```go
type st struct {
    Name string `json:"name"`
}

func main() {
    raw := []byte(`{"NAME":"json"}`)
    s := st{}
    json.Unmarshal(raw, &s)
    fmt.Println("name:", s.Name)
    // Output:
    // name: json
}
```

虽然 JSON 正文中使用的 key 是全大写的 `NAME`，而 struct 中定义的 tag 是全小写的 `name`，但是依然能够正确地解析到结构体中。

但是到了 `map`，就不一样了，由于 map 的特性，原生 json 无法做到忽略大小写。

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

这种问题，在 jsonvalue 中可以解决吗？答案是肯定的。

---

## 在 jsonvalue 中忽略大小写

在 jsonvalue 中使用 map 来存储 object 类型的 K-V 信息，默认情况下，jsonvalue 的 Get 函数是区分大小写的。

但是如果开发者在做 `Get` 操作时需要忽略大小写，那么只需要在 Get 操作之前插入一个 `Caseless()` 调用即可，如以下例子：

```go
func main() {
    raw := []byte(`{"NAME":"json"}`)
    v := jsonvalue.MustUnmarshal(raw)
    fmt.Println("name =", v.MustGet("name").String())
    fmt.Println("NAME =", v.Caseless().MustGet("name").String())
}
```

输出结果，第一行 `Println` 将无法读取到字符串值，而第二行则能够读到。

`Caseless` 函数的原理，首先是打开 `*V` 内部的 `caseless` 开关，是遍历 `*V` 对象中的所有 K-V 段，建立大小写的映射。一个被标记为 caseless 的 `*V`，在后续的 set 操作中，也会将新 key 信息加入到 caseless 映射中。

当不区分大小写来读取值时，jsonvalue 会优先命中大小写一致的 key，如果找不到，才查找其他不区分大小写的 key。比如上文，`name` 和 `NAME` 可以共存在同一个 object 值中。


因此从原理上，`Caseless()` 函数有以下的特点：

- `Caseless` 函数会改变 `*V` 内部的结构。由于 `jsonvalue` 不是协程安全的，在多协程环境下使用同一个 jsonvalue 的 `caseless` 特性，需要加写锁
- `Caseless` 会给 jsonvalue 带来额外的开销，如果不是特别有必要，建议还是区分大小写

需要注意的是，`Caseless` 并不支持 `Equal` 函数，主要是考虑到当不区分大小写时，是否 Equal 可能会有很多歧义，因此作者斟酌再三决定放弃这一功能。其实也请开发者在实际工作中注意，尽量避免需要使用 `Caseless` 函数的情况
