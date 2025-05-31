<font size=6>値の変換</font>

[前のページ](./08_caseless.md) | [目次](./README.md) | [次のページ](./10_scenarios.md)

---

- [概要](#概要)
- [文字列から数値への変換](#文字列から数値への変換)
- [文字列と数値からブール値への変換](#文字列と数値からブール値への変換)
- [`GetString`と`MustGet(...).String()`の違い](#getstringとmustgetstringの違い)

---

## 概要

`encoding/json`には「string」と呼ばれる`tag`があり、これは現在の値（型に関係なく）をJSON文字列に変換することを意味します。例えば：

```go
type st struct {
    Num int `json:"num,string"`
}

func main() {
    st := st{Num: 12345}
    b, _ := json.Marshal(&st)
    fmt.Println(string(b))
	// 出力: {"num":"12345"}
}
```

`num`が数値ではなく文字列としてシリアライズされていることが分かります。

Jsonvalueも、このような型の値を文字列から読み取ることをサポートしています。詳細については以下の説明をご覧ください。

---

## 文字列から数値への変換

すべての数値型の`GetXxx`メソッドは、文字列値からの数値解析をサポートしています。

ただし、対象の値が数値型のJSON値でない場合、`error`はnilになりません。状況に応じて様々な種類のエラーが返される可能性があります：

- 対象のキーが見つからない場合：`ErrNotFound`を返します
- 対象の値が数値の場合：数値型が返され、`nil`のerrが返されます
- 対象の値が文字列でも数値でもない場合：`ErrTypeNotMatch`を返します
- 対象の値が文字列の場合：文字列の値を解析し、以下のようになります
  - 解析が成功した場合：エラー`ErrTypeNotMatch`と共に数値を返します
  - 解析が失敗した場合（不正な数値）：詳細な解析の説明を含む`ErrParseNumberFromString`エラーを返します

`errors.Is()`を使用して異なる種類のエラーをチェックできます。例：

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

出力：

```
01 - float: 0.125
01 - err: not match given type
02 - float: 0
02 - err: failed to parse number from string: parsing number at index 0: zero string
```

---

## 文字列と数値からブール値への変換

数値と同様に、文字列もブール値を保持できます。

文字列をブール値に変換する際に`true`が返されるのは、以下の場合のみです：

- 文字列の値が正確に小文字の`"true"`である場合

その他の場合、変換では`false`が返されます。

数値もブール値に変換できます。数値の値がゼロでない限り、`true`が返されます。そうでなければ`false`が返されます。

---

## `GetString`と`MustGet(...).String()`の違い

`GetXxx()`シリーズのメソッドと対応する`MustGet().Xxx()`関数のロジックはほぼ同じです。ただし、`Xxx()`メソッドはエラーを返しません。

例：

```go
	const raw = `{"legal":"1.25E-1","illegal":"ABCD"}`
	v := jsonvalue.MustUnmarshalString(raw)
	f := v.MustGet("legal").Float64()
	fmt.Println("float:", f)
	// 出力: float: 0.125
```

しかし、一つの例外があります：`GetString`と`MustGet().String()`です。

`GetString`は文字列型のJSON値の文字列を返します。対象が見つからない場合、または文字列でない場合、エラーが返されます。

一方、`String`は非常に特別です。文字列型のJSON値の文字列値を取得できるだけでなく、`fmt.Stringer`インターフェースを実装しており、`fmt`パッケージの`%v`キーワードで呼び出されます。

そのため、`String`のロジックは少し複雑です：

- 現在の`*jsonvalue.V`が文字列型の場合：その文字列値を返します
- 数値型の場合：数値の文字列表現を返します
- null型（JSON標準で定義された「null」型）の場合：`null`を返します
- ブール型の場合：`true`または`false`を返します
- オブジェクト型の場合：`{K1: V1, K2: V2}`の形式の文字列を返します（エスケープなし、JSON形式ではありません）
- 配列型の場合：`[v1, v2]`の形式の文字列を返します（同様にエスケープなし）
