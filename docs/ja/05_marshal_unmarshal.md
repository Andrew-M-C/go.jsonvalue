<font size=6>マーシャルとアンマーシャル</font>

[前のページ](./04_get.md) | [目次](./README.md) | [次のページ](./06_import_export.md)

---

- [アンマーシャル関数](#アンマーシャル関数)
  - [基本的なアンマーシャル](#基本的なアンマーシャル)
  - [その他のアンマーシャル関数](#その他のアンマーシャル関数)
- [マーシャル関数](#マーシャル関数)
- [公式 `encoding/json` サポート](#公式-encodingjson-サポート)

---

## アンマーシャル関数

### 基本的なアンマーシャル

jsonvalue では、シリアライゼーションとデシリアライゼーションのプロセスを表現するために marshal / unmarshal という用語を使用しています。

Jsonvalue は以下の関数を使用して生の JSON テキストを解析します：

```go
func Unmarshal(b []byte) (ret *V, err error)
```

エラーが発生するかどうかに関わらず、非 nil の `*jsonvalue.V` オブジェクトが返されます。ただし、生のテキストが不正な場合は、エラーオブジェクトが返され、どのようなエラーかが説明されます。この場合、返される jsonvalue 値の `Type()` は `jsonvalue.NotExist` になります。

### その他のアンマーシャル関数

実際には、生の JSON テキストは `[]byte` ではなく `string` 形式で与えられることが多いです。`string(b)` 変換を行うには少し時間がかかります。このコピー時間を節約するために、`string` 版のアンマーシャル関数を使用できます：

```go
func UnmarshalString(s string) (ret *V, err error)
```

さらに、与えられた JSON テキストの正確性を気にする必要がない場合、または確実に正当であることが分かっている場合は、単純にエラーを無視して以下の関数を使用できます：

```go
func MustUnmarshal(b []byte) *V
func MustUnmarshalString(s string) *V
```

上記の関数と同様に、これら2つの関数は確実に非 nil の `jsonvalue.V` を返します。

---

## マーシャル関数

jsonvalue におけるシリアライゼーションは「マーシャル」と呼ばれます。アンマーシャルと同様に、以下の4つの関数が提供されています：

```go
func (v *V) Marshal          (opts ...Option) (b []byte, err error)
func (v *V) MarshalString    (opts ...Option) (s string, err error)
func (v *V) MustMarshal      (opts ...Option) []byte
func (v *V) MustMarshalString(opts ...Option) string
```

現在のバージョンの jsonvalue では、以下の状況でエラーが発生します：

1. `*V` が `NotExist` タイプの場合
2. `+Inf`、`-Inf`、`NaN` などの不正な浮動小数点数が含まれており、これらの浮動小数点値に対する特別な操作が指定されていない場合
   - これらの浮動小数点値に対する特別なオプションについては、他のセクションで後述します。
3. 追加オプションに不正な設定が含まれている場合

---

## 公式 `encoding/json` サポート

`*jsonvalue.V` 型は `json.Marshaler` と `json.Unmarshaler` インターフェースも実装しています。これにより、`encoding/json` を介したマーシャルとアンマーシャルが可能になります。例えば：

```go
var v *jsonvalue.V
err := json.Unmarshal(data, &v)
```

または

```go
v := jsonvalue.NewObject()
v.MustSet("Hello, JSON!").At("greeting")
b, err := json.Marshal(v)
```

