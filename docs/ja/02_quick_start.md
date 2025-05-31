<font size=6>クイックスタート</font>

[前のページ](./01_introduction.md) | [目次](./README.md) | [次のページ](./03_set.md)

---

以下のような複雑なJSONオブジェクトを作成する場合：

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

jsonvalueを使えばたった3行で実現できます：

```go
	v := jsonvalue.NewObject()
	v.Set("Hello, JSON").At("obj", "obj", "obj", "str")
	fmt.Println(v.MustMarshalString())
```

出力: `{"obj":{"obj":{"obj":{"str":"Hello, JSON!"}}}`

一方、上記のJSONバイトから単一のデータを取得する場合も簡単です：

```go
const raw = `{"obj": {"obj": {"obj": {"str": "Hello, JSON!"}}}}`
s := jsonvalue.MustUnmarshalString(raw).MustGet("obj", "obj", "obj", "str").String()
fmt.Println(s)
```

出力: `Hello, JSON!`


