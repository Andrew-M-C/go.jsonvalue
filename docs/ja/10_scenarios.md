<font size=6>特殊な適用シナリオ</font>

[前のページ](./09_conversion.md) | [目次](./README.md) | [次のページ](./11_comparation.md)

---

- [JSONオブジェクトにおけるキーの順序](#jsonオブジェクトにおけるキーの順序)
	- [オブジェクトの元のキー順序の取得](#オブジェクトの元のキー順序の取得)
	- [キーが設定された順序でのJSONオブジェクトのシリアライズ](#キーが設定された順序でのjsonオブジェクトのシリアライズ)
	- [JSON値の冪等性](#json値の冪等性)
- [特殊な浮動小数点数 +/-Inf と NaN](#特殊な浮動小数点数--inf-と-nan)
- [構造体の `omitempty` JSONタグの無視](#構造体の-omitempty-jsonタグの無視)
- [数値型JSON値の元のテキストの取得](#数値型json値の元のテキストの取得)

---

[はじめに](./01_introduction.md)で前述したように、jsonvalueでは様々な特殊なシナリオがサポートされています。

このセクションでは、これらの特殊な適用について説明し、これらのシナリオの背景にある論理と経緯を説明します。

---

## JSONオブジェクトにおけるキーの順序

標準的なJSON実装によると、JSONオブジェクトは一連のキー・バリューペアです。しかし実際には、異なるチームの多くの同僚がこの一連のキー・バリューペアを順序付きとして扱っていることがわかりました...

そこで彼らから以下の要望が寄せられました：

1. JSONテキストを解析する際、キーの元の順序を知りたい。
2. jsonvalueをシリアライズする際、キーの順序を指定したい。

`encoding/json`を使用する場合、（`map`型では）サポートされていません。他のJSONパッケージでは可能かもしれませんが、困難です。

最初は、これは不正だと思って拒否していました。しかし、このような奇妙な（または愚かな）作業をする人がますます多くなっていることがわかったため、最終的に実装することにしました。

### オブジェクトの元のキー順序の取得

これは最初の要望に対するものです。例を見てみましょう：

```go
const raw = `{"a":1,"b":2,"c":3}`
```

通常通りデシリアライズできます：

```go
v := jsonvalue.MustUnmarshalString(raw)
```

元の順序でキーを取得したい場合は、`RangeObjectsBySetSequence`メソッドを使用します：

```go
keys := []string{}
v.RangeObjectsBySetSequence(func(key string, _ *V) bool {
    keys = append(keys, key)
    return true
})
fmt.Println(keys)
```

出力は`[a, b, c]`で、常に保証されます。これは生データの元の順序です。

この機能はv1.3.1から最初にサポートされました。

---

2番目の要望については、いくつかの実装方法があります。以下で説明します。

### キーが設定された順序でのJSONオブジェクトのシリアライズ

これは`Marshal`メソッドを呼び出す際のオプションです。このセクションのタイトルを説明しましょう：

1. コードによって段階的に生成されたjsonvalue値の場合、キー・バリューペアはオブジェクトに設定された順序でシリアライズされます。
2. このjsonvalue値が生テキストからデシリアライズされた場合、マーシャリング時には生JSONテキストの元の順序でシリアライズされます。
3. 状況1と2を組み合わせると、まず生JSONテキストの元の順序、次にその後設定された順序になります。
4. キーが最初に削除されてから設定された場合、または後で上書きされた場合は、最後のものを使用します。

これを実現するために、`OptSetSequence()`のような追加オプションを追加できます：

```go
const raw = `{"a":1,"b":2,"c":3}`
v := jsonvalue.MustUnmarshalString(raw) // キー順序: a, b, c
v.Delete("b")                           // bが削除される、順序: a, c
v.Set(4).At("d")                        // 新しいdを追加、キー順序: a, c, d
v.Set(1).At("a")                        // aを再設定、最新のものを使用、キー順序: c, d, a
s := v.MustMarshalString(OptSetSequence())
fmt.Println(s)
```

出力：

```go
{"c":3,"d":4,"a":1}
```

この機能は`struct`データの解析（`Import()`関数）でも有効です。

---

### JSON値の冪等性

キーの順序が重要になる状況があります。JSON生テキストを生成した後、そのチェックサム値を計算することがあります。これはHTTP検証操作では非常に一般的です。

しかし、JSONオブジェクトの順序が予期しないものである場合、この操作は不可能です。そのため、同じオブジェクトは常に同じJSONテキストを生成することを確実にする必要があります。

この問題は、前述の「設定順序」とは少し異なります。

これを実現する最も簡単な方法は、追加オプション`OptDefaultStringSequence()`を使用することです。これにより、キーが文字列比較順序でシリアライズされることが保証されます。

これは`strings.Compare`のシンプルなラッパーです。

例：

```go
// mapの範囲指定のランダム化により、aとbの順序を予測することはできません。
v := jsonvalue.New(map[string]any{
    "a": 1,
    "b": 2,
    "c": 3,
})
s := v.MustMarshalString(OptDefaultStringSequence())
fmt.Println(s)
// 出力:
// {"a":1,"b":2,"c":3}
```

---

## 特殊な浮動小数点数 +/-Inf と NaN

`encoding/json`を使用する場合、+/-無限大やNaNの浮動小数点数をシリアライズすると、次のようになります：

```go
func main() {
	_, err := json.Marshal(map[string]float64{
		"score": math.Inf(-1),
	})
	fmt.Println(err)
}
```

エラーが発生する可能性があります：`json: unsupported value: -Inf`。

解決策は簡単だと思うかもしれません：これらの特殊な値を避けるだけです。

しかし、推薦システムのプログラマーとして、マテリアルスコアリングシステムから-Inf値を取得することは非常に一般的であることがわかりました。元のエンコーディングプロトコルはprotobufで、IEEE-754を完全にサポートしています。しかし、JSONでは機能しません。多くの場合、JSONを使用する必要があります。

この場合、jsonvalueを使用して、マーシャリング時に+/-InfとNaNをどう処理するかを指定できます。3つのタイプまたはオプションがあります：

- JSON nullで置換
- 別の有効な浮動小数点値で置換
- JSON文字列で置換

前の例も使用できます：

```go
func main() {
	v, _ := jsonvalue.Import(map[string]float64{
		"score": math.Inf(-1),
	})
	s := v.MustMarshalString(jsonvalue.OptFloatInfToFloat(23333))
	fmt.Println(s)
	// 出力:
	// {"score":-23333}
}
```

上記のオプションについては、セクション「[追加オプション](./12_option.md)」を参照してください。

---

## 構造体の `omitempty` JSONタグの無視

このオプションは`Import`で使用されます。簡単に言うと、`struct`を`*jsonvalue.V`に変換する際、`omitempty`タグが無視され、`struct`の構造が完全にエクスポートされます。

次の構造体を例に取ります：

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

追加オプションが指定されていない場合、出力は`{}`になります。しかし、上記の問題によると、`{"a":0}`であるべきです。

この問題を受けた時（TencentでGithubではありません）、混乱しました：`omitempty`が不要なら、なぜ定義するのでしょうか？

![？？？](https://bkimg.cdn.bcebos.com/pic/8cb1cb1349540923e1860cc29958d109b2de499a?x-bce-process=image/resize,m_lfit,w_536,limit_1/format,f_jpg)

この問題を研究した後、エクスポートされる構造体がprotobufによって生成されることがわかりました。goツール`protoc`によって生成されたすべての構造体のフィールドには、`omitempty`タグが追加されます。しかし、問題提起者はすべてのprotobufデータを分析したかったため、タグを無視する必要がありました。

これは合理的な要求でした。StackOverflowでも同様の質問を見つけました：[golang protobuf remove omitempty tag from generated json tags](https://stackoverflow.com/questions/34716238/)

これを解決するために`jsonpb`を使用することは可能でした。しかし、他の機能（例えば、上記で言及したキー順序）と一致しませんでした。そこで同僚はjsonvalueを使用することを決定し、この機能を求めました。

これを実装するのは難しくありませんでした。追加オプション`OptIgnoreOmitempty()`を追加しました。これは`Import`関数でサポートされる唯一のオプションでもあります。コードは次のように書き直すことができます：

```go
type st struct {
	A int `json:"a,omitempty"`
}

func main() {
	st := st{}
	v, _ := jsonvalue.Import(&st, jsonvalue.OptIgnoreOmitempty())
	s := v.MustMarshalString()
	fmt.Println(s)
	// 出力:
	// {"a":0}
}
```

---

## 数値型JSON値の元のテキストの取得

数値はコンピュータでバイナリ桁として保存されます。しかし、JSON標準は10進数を使用し、バイナリと10進数の間の変換は定義されていません。

`map[string]any`を使用して数値をアンマーシャルする場合、`encoding/json`ではすべての数値が`float64`型の変数として解析されます。誤って`int64`や`uint64`をJSONテキストに設定した場合、精度損失の可能性により、一部の桁が失われる可能性があります。

これは64ビットハッシュ数を転送する際の重要な問題です。`encoding/json`を使用している場合は、`json.RawMessage`を使用してこれを解決できます。jsonvalueについては、`GetUint64`または`Get(...).String()`を使用するだけで、生の64ビットハッシュ数を取得できます。