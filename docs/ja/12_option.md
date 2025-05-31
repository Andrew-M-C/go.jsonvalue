<font size=6>追加オプション</font>

[前のページ](./11_comparation.md) | [目次](./README.md) | [次のページ](./13_beta.md)

---

このセクションでは、追加オプションの詳細な使用方法について説明します。または、[特別なアプリケーションシナリオ](./10_scenarios.md)のセクションを参照してください。

---

- [概要](#概要)
- [null値の無視](#null値の無視)
- [可視インデント](#可視インデント)
- [オブジェクト値内のキーの順序を指定する](#オブジェクト値内のキーの順序を指定する)
  - [キーが設定された時の順序](#キーが設定された時の順序)
  - [コールバックを使用して順序を指定する](#コールバックを使用して順序を指定する)
  - [アルファベット順序を使用する](#アルファベット順序を使用する)
  - [事前定義された\[\]stringを使用してキー順序を識別する](#事前定義されたstringを使用してキー順序を識別する)
- [NaNの処理](#nanの処理)
  - [NaNを別の浮動小数点値に変換](#nanを別の浮動小数点値に変換)
  - [NaNをNull型に変換](#nanをnull型に変換)
  - [NaNを文字列に変換](#nanを文字列に変換)
- [+/-Infの処理](#-infの処理)
  - [+/-Infを別の浮動小数点値に変換](#-infを別の浮動小数点値に変換)
  - [+/-InfをNull型に変換](#-infをnull型に変換)
  - [+/-Infを文字列に変換](#-infを文字列に変換)
- [非必須文字のエスケープ](#非必須文字のエスケープ)
  - [SetEscapeHTMLオプション](#setescapehtmlオプション)
  - [スラッシュ `/`](#スラッシュ-)
  - [0x7Fより大きいUnicodeのエスケープ](#0x7fより大きいunicodeのエスケープ)
- [構造体のomitemptyタグの無視](#構造体のomitemptyタグの無視)

---

## 概要

`Marshal`メソッドを振り返ってみましょう：

```go
func (v *V) Marshal          (opts ...Option) (b []byte, err error)
func (v *V) MarshalString    (opts ...Option) (s string, err error)
func (v *V) MustMarshal      (opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

すべてのメソッドが追加の`opts ...Option`をサポートしており、jsonvalueデータのシリアライズ時に追加オプションを指定できることがわかります。

簡単な例として、以下のオプションですべてのnull値を無視できます：

```go
v := jsonvalue.NewObject()
v.SetNull().At("null")
fmt.Println(v.MustMarshalString())
fmt.Println(v.MustMarshalString(jsonvalue.OptOmitNull(true)))
```

出力：

```json
{"null":null}
{}
```

現在、`OptIgnoreOmitempty()`オプションのみが`Import()`用に設計されており、その他はすべてmarshal用です。

---

## null値の無視

上記で既に説明しました。

---

## 可視インデント

これは`encoding/json`の`json.MarshalIndent`関数に似ています。前の例を使用して、以下のオプションを追加できます：

```go
s := v.MustMarshalString(jsonvalue.OptIndent("", "  "))
fmt.Println(s)
```

これは以下を出力します：

```json
{
  "null": null
}
```

---

## オブジェクト値内のキーの順序を指定する

一般的に、キー値の順序を指定することは不要であり、CPU時間の無駄です。しかし、キーの順序が重要になる特別な状況があります：

- 前述したハッシュチェックサム
- デバッグ時に特定のキー値ペアを素早く見つける
- キーの順序に依存するJSONの非標準的な使用

このサブセクションでは、マーシャリング時にキーの順序を指定する方法を説明します。

### キーが設定された時の順序

[特別なアプリケーションシナリオ](./10_scenarios.md)を参照してください。

### コールバックを使用して順序を指定する

```go
func OptKeySequenceWithLessFunc(f MarshalLessFunc) Option
```

`MarshalLessFunc`はコールバック関数で、プロトタイプは以下の通りです：

```go
type MarshalLessFunc func(nilableParent *ParentInfo, key1, key2 string, v1, v2 *V) bool
```

この関数は`sort.Sort`のように動作します。入力パラメータの定義は以下の通りです：

- `nilableParent` - 現在の値の親キー
- `key1`, `v1` - 再配置する最初のキー・値ペア
- `key2`, `v2` - 再配置する2番目のキー・値ペア

v1がv2より前に来るべきかどうかを返します。パッケージ`sort`の`Less`関数のように動作します。

### アルファベット順序を使用する

```go
func OptDefaultStringSequence() Option
```

### 事前定義された[]stringを使用してキー順序を識別する

「less」コールバックを使用するのは非常に複雑です。単純に`[]string`を渡すことで、jsonvalueは指定された文字列スライスに従ってキーの順序を配置します。その中で指定されていないキーがある場合、指定されたすべてのキーの後に配置されます。

```go
func OptKeySequence(seq []string) Option
```

`OptKeySequence`と`OptKeySequenceWithLessFunc`の両方が指定された場合、`OptKeySequenceWithLessFunc`が優先的に使用されます。

--- 

## NaNの処理

標準JSONでは、NaN（非数値）と+/-Inf（正/負の無限大）は不正な値です。しかし、場合によってはこれらの値を処理する必要があります。

デフォルトでは、jsonvalueはこれらの値を処理する際にエラーを発生させます。しかし、エラーなしの変換オプションも提供されています。

まずNaNを見てみましょう：

### NaNを別の浮動小数点値に変換

```go
func OptFloatNaNToFloat(f float64) Option
```

別の浮動小数点数を指定すると、すべてのNaNがそれに置き換えられます。置換値としてNaNや+/-Infを指定することはできません。

### NaNをNull型に変換

```go
func OptFloatNaNToNull() Option
```

NaNをJSONのnullに置き換えます。このオプションは`OptOmitNull`の影響を受けないことに注意してください。このオプションでは常にNaNを`null`に変換します。

### NaNを文字列に変換

```go
func OptFloatNaNToString   (s string) Option
func OptFloatNaNToStringNaN() Option
```

`OptFloatNaNToStringNaN()`は`OptFloatNaNToString("NaN")`と同等です。

---

## +/-Infの処理

+/-Infの処理はNaNと似ています：

### +/-Infを別の浮動小数点値に変換

```go
func OptFloatInfToFloat(f float64) Option
```

+Infは`f`に置き換えられ、-Infは`-f`に置き換えられます。置換値としてNaNや+/-Infを指定することはできません。

### +/-InfをNull型に変換

```go
func OptFloatInfToNull() Option
```

同様に、このオプションは`OptOmitNull`の影響を受けません。

### +/-Infを文字列に変換

```go
func OptFloatInfToString   (positiveInf, negativeInf string) Option
func OptFloatInfToStringInf() Option
```

+/-Infの置換として2つの文字列を指定します。指定された文字列が空の場合、以下の優先順位で置換されます：

- +Infについて、指定された文字列が空の場合、`"+Inf"`として置換されます
- -Infについて、指定された文字列が空の場合、まず+Infの設定を探します。存在する場合、`+`プレフィックスを削除して`-`を追加します
- +/-Infの設定が両方とも空の場合、-Infは`"-Inf"`に置換されます

`OptFloatInfToStringInf()`は`OptFloatInfToString("+Inf", "-Inf")`と同等です。

## 非必須文字のエスケープ

JSON標準では、エスケープが必要な文字がいくつかあります：

1. 重要なフォーマットまたは予約文字
2. 127より大きいUnicode

しかし、すべてのJSONエンコーダーが完全なエスケープルールに従うわけではありません。このセクションでは、特別なエスケープルールを指定する方法を説明します。

### SetEscapeHTMLオプション

JSON標準によると、`&`、`<`、`>`は`\u00XX`にエスケープされるべきです。しかし、実際にはこれら3つの文字をエスケープしなくても安全です。デフォルトでは、jsonvalueはそれらすべてをエスケープします。しかし、以下のオプションを使用してこれらの文字のエスケープを無効にできます：

```go
func OptEscapeHTML(on bool) Option
```

エスケープを有効にするには`true`を、エスケープしない場合は`false`を渡します。

このオプションは`encoding/json`の[`SetEscapeHTML`](https://pkg.go.dev/encoding/json#Encoder.SetEscapeHTML)メソッドに非常に似ています。

### スラッシュ `/`

JSON標準によると、スラッシュ`/`はエスケープされるべきです。しかし、実際にはエスケープしなくても問題ありません。デフォルトでは、jsonvalueはスラッシュをエスケープしますが、このスイッチオプションを使用できます：

```go
func OptEscapeSlash(on bool) Option
```

スラッシュをエスケープしない場合は`false`を渡します。デフォルトでは`true`です。

### 0x7Fより大きいUnicodeのエスケープ

```go
func OptUTF8() Option
```

デフォルトでは、jsonvalueは0x7Fより大きいすべてのunicode値を`\uXXXX`形式（UTF-16）にエスケープします。これにより、ほぼすべてのエンコーディング問題を回避できます。しかし、ペイロードの大部分がunicodeの場合、ネットワークトラフィックの無駄になる可能性があります。この場合、UTF-8エンコーディングの使用を検討できます。これは`encoding/json`が行うことです。

## 構造体のomitemptyタグの無視

これは稀な機能です。[特別なアプリケーションシナリオ](./10_scenarios.md)の「構造体の`omitempty` JSONタグの無視」サブセクションを参照してください。
