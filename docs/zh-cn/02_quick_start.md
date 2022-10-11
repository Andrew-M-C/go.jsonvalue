<font size=6>快速上手</font>

[上一页](./01_introduction.md) | [总目录](./README.md) | [下一页](./03_set.md)

---

创建一个比如如下的复杂 JSON 对象：

```json
{
	"obj": {
		"obj": {
			"obj": {
				"str": "Hello, JSON!"
			}
		}
	}
}
```

使用 jsonvalue 只需要三行:

```go
	v := jsonvalue.NewObject()
	v.Set("Hello, JSON").At("obj", "obj", "obj", "str")
	fmt.Println(v.MustMarshalString())
```

输出结果为: `{"obj":{"obj":{"obj":{"str":"Hello, JSON!"}}}`

反过来，我们如果要直接读取上面的 json 数据，也可以这么用 jsonvalue: 

```go
const raw = `{"obj": {"obj": {"obj": {"str": "Hello, JSON!"}}}}`
s := jsonvalue.MustUnmarshalString(s).MustGet("obj", "obj", "obj", "str").String()
fmt.Println(s)
```

输出结果为: `Hello, JSON!`


