<font size=6>イテレーション</font>

[前のページ](./06_import_export.md) | [目次](./README.md) | [次のページ](./08_caseless.md)

---

- [概要](#概要)
- [配列値のイテレーション](#配列値のイテレーション)
- [オブジェクト値のイテレーション](#オブジェクト値のイテレーション)
- [オブジェクトの元のキー順序の取得](#オブジェクトの元のキー順序の取得)
- [すべての子値のウォーク](#すべての子値のウォーク)

---

## 概要

jsonvalue では、配列型またはオブジェクト型の JSON 値をイテレートすることができます。

イテレーションには2つのモードがあります：

1. `jsonparser` の [`ArrayEach`](https://pkg.go.dev/github.com/buger/jsonparser#ArrayEach) や [`ObjectEach`](https://pkg.go.dev/github.com/buger/jsonparser#ObjectEach) のように、コールバック関数を使用してイテレーションデータを受け取る方法
2. `for-range` を使用してイテレートする方法。これにより、`map` や `slice` の操作により近いコードが書けます

**重要**: jsonvalue のすべてのイテレーションメソッドはゴルーチンセーフではありません。マルチゴルーチン環境では、ロックやその他の保護機能を追加してください。ただし、イテレーション中に `Set` 系メソッドや `Caseless` メソッドが使用されない場合は、ゴルーチンセーフになります。つまり、これらの操作には読み取りロックのみが必要です。

## 配列値のイテレーション

配列型の JSON 値をイテレートするには、以下のメソッドを使用します。

```go
func (v *V) RangeArray(callback func(i int, v *V) bool)
func (v *V) ForRangeArr() []*V
```

コールバックパターンの場合は、`RangeArray` を使用します。コールバック内で `true` を返すとイテレーションを継続し（`continue` の意味）、`false` を返すとイテレーションを中断します（`break` の意味）。

`ForRangeArr` は `[]*jsonvalue.V` のスライスを返すため、これを `for-range` 文で使用できます。

詳細な例：

```go
anArr.RangeArray(func(i int, v *jsonvalue.V) bool {
    // ...... i と v を処理
    return true // 継続
})

for i, v := range anArr.ForRangeArr() {
    // ...... i と v を処理
}
```

## オブジェクト値のイテレーション

オブジェクト型の JSON 値をイテレートするには、以下のメソッドを使用します：

```go
func (v *V) RangeObjects(callback func(k string, v *V) bool)
func (v *V) ForRangeObj() map[string]*V
```

同様に、例：

```go
anObj.RangeObject(func(key string, v *jsonvalue.V) bool {
    // ...... key と v を処理
    return true // 継続
})

for key, v := range anObj.ForRangeObj() {
    // ...... key と v を処理
}
```

## オブジェクトの元のキー順序の取得

理論的には、JSON オブジェクトのキー順序は未定義で予期できないものであるべきです。しかし実際には、この機能が非常に人気があることは驚くべきことです。

v1.3.1 以降、この機能が追加され、アンマーシャル効率にほとんど影響を与えません。

呼び出し側は通常通り `Unmarshal` 操作を実行した後、`RangeObjectsBySetSequence` メソッドを使用できます。このメソッドは `RangeObjects` と同じコールバックを受け取ります。

例：

```go
const raw = `{"a":1,"b":2,"c":3}`
v := jsonvalue.MustUnmarshalString(raw)
keys := []string{}
v.RangeObjectsBySetSequence(func(key string, _ *V) bool {
    keys = append(keys, key)
    return true
})
fmt.Println(keys)
```

出力は `[a b c]` となり、常に保証されます。

## すべての子値のウォーク

`Walk` メソッドは、JSON 構造内のすべての子値を深さ優先でトラバースする簡単な方法を提供します。直接の子要素のみを反復する `RangeArray` や `RangeObjects` とは異なり、`Walk` メソッドは JSON ツリー内のすべてのリーフノードを再帰的に訪問します。このメソッドは `path/filepath` の `Walk` メソッドと同様のパターンに従います。

```go
func (v *V) Walk(fn WalkFunc)

type WalkFunc func(path []PathItem, v *V) bool

type PathItem struct {
    Idx int    // 配列インデックス、-1 はこの要素が配列要素でないことを示す
    Key string // オブジェクトキー名、"" はこの要素がオブジェクト要素でないことを示す
}
```

`Walk` メソッドは、各リーフ値に対して提供された `WalkFunc` コールバック関数を呼び出し、以下を提供します：
- `path`: ルートから現在の値への完全なパスを示す `PathItem` のスライス
- `v`: 現在訪問されている JSON 値

コールバック関数は、トラバーサルを続行するために `true` を返すか、早期にトラバーサルを停止するために `false` を返す必要があります。

