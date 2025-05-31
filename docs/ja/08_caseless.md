<font size=6>大文字小文字の区別</font>

[前のページ](./07_iteration.md) | [目次](./README.md) | [次のページ](./09_conversion.md)

---

- [Go標準jsonの問題](#go標準jsonの問題)
- [jsonvalueで大文字小文字を無視する](#jsonvalueで大文字小文字を無視する)

---

## Go標準jsonの問題

Go標準の`encoding/json`では、`struct`内のJSONタグに対する処理原則は大文字小文字を区別しません。例えば以下の例：

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

JSON本文で使用されているキーは全て大文字の`NAME`で、struct内で定義されているタグは全て小文字の`name`ですが、それでも正しく構造体に解析することができます。

しかし`map`の場合は話が違います。mapの特性により、標準jsonは大文字小文字を無視することができません。

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

この問題は、jsonvalueで解決できるでしょうか？答えは「はい」です。

---

## jsonvalueで大文字小文字を無視する

jsonvalueではmapを使用してobject型のK-V情報を格納します。デフォルトでは、jsonvalueのGet関数は大文字小文字を区別します。

しかし、開発者が`Get`操作時に大文字小文字を無視したい場合は、Get操作の前に`Caseless()`呼び出しを挿入するだけで済みます。以下の例をご覧ください：

```go
func main() {
    raw := []byte(`{"NAME":"json"}`)
    v := jsonvalue.MustUnmarshal(raw)
    fmt.Println("name =", v.MustGet("name").String())
    fmt.Println("NAME =", v.Caseless().MustGet("name").String())
}
```

出力結果では、最初の`Println`は文字列値を読み取ることができませんが、2番目の行では読み取ることができます。

`Caseless`関数の仕組みは、まず`*V`内部の`caseless`スイッチをオンにし、`*V`オブジェクト内のすべてのK-Vセグメントを走査して、大文字小文字のマッピングを構築することです。caselessとしてマークされた`*V`は、後続のset操作でも新しいkey情報をcaselessマッピングに追加します。

大文字小文字を区別せずに値を読み取る際、jsonvalueはまず大文字小文字が一致するキーを優先的に検索し、見つからない場合に他の大文字小文字を区別しないキーを検索します。例えば上記の場合、`name`と`NAME`は同じobject値内に共存することができます。

したがって仕組み上、`Caseless()`関数には以下の特徴があります：

- `Caseless`関数は`*V`内部の構造を変更します。`jsonvalue`はゴルーチンセーフではないため、マルチゴルーチン環境で同じjsonvalueの`caseless`機能を使用する場合は、書き込みロックを追加する必要があります
- `Caseless`はjsonvalueに追加のオーバーヘッドをもたらします。特に必要でない限り、大文字小文字を区別することを推奨します

注意すべき点は、`Caseless`は`Equal`関数をサポートしていないことです。これは主に、大文字小文字を区別しない場合、Equalかどうかに多くの曖昧さが生じる可能性があることを考慮し、作者が熟慮の末にこの機能を放棄することを決定したためです。実際の作業においても、開発者は`Caseless`関数を使用する必要がある状況をできるだけ避けるよう注意してください。
