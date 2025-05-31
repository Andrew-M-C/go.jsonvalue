<font size=6>`encoding/json` との Import と Export</font>

[前のページ](./05_marshal_unmarshal.md) | [目次](./README.md) | [次のページ](./07_iteration.md)

---

- [Import / Export](#import--export)
- [関数 `New()`](#関数-new)
- [メソッド `Set`、`Append`、`Insert`、`Add`](#メソッド-setappendinsertadd)

---

## Import / Export

`Import` と `Export` を設計した当初の目的は、`encoding/json` と `jsonvalue` の間でデータを変換することでした。

しかし、`Import` の開発過程で、以下のような多くの追加機能が実装されました：

---

## 関数 `New()`

v1.3.0 から、`New` という新しい関数が提供されています。この関数は任意の型のパラメータ（古い形式では `any` または `interface{}`）を受け取り、それを `*jsonvalue.V` 値に解析します。入力パラメータが不正な場合、返される値の型は `NotExist` になります。

実際には、`New` は関数 `Import` の簡単なラッパーであり、`error` 型を返さない点が異なります。

---

## メソッド `Set`、`Append`、`Insert`、`Add`

バージョン v1.3.0 以前では、jsonvalue ノードにサブ値を追加する際、入力パラメータの型を明示的に指定する必要がありました。例えば：

```go
v.SetString("Hello, world").At("greeting")
```

v1.3.0 以降、`Set`、`Append`、`Insert`、`Add` などのメソッドは `any` 型のパラメータを受け入れ、それを解析してから対応するサブ値を設定するようになりました。

例えば、空の JSON オブジェクトを作成した後：

```go
v := jsonvalue.NewObject()
```

その中にサブオブジェクトを追加することができます：

```go
child := map[string]string{
    "text": "Hello, jsonvalue!",
}
v.Set(child).At("child")
fmt.Println(v.MustMarshalString())
```

出力：`{"child":{"text":"Hello, jsonvalue!"}}`

または、通常の JSON 値を設定することもできます：

```go
v := jsonvalue.NewObject()
v.Set("Hello, JSON!").At("msg")
fmt.Println(v.MustMarshalString())
// 出力：{"msg":"Hello, JSON!"}
```


