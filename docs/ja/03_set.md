<font size=6>JSONの作成とシリアライゼーション</font>

[前のページ](./02_quick_start.md) | [目次](./README.md) | [次のページ](./04_get.md)

---

このセクションでは、jsonvalue値を生成し、設定する方法について説明します。

---

- [JSON値の作成](#json値の作成)
- [JSONにサブ値を設定](#jsonにサブ値を設定)
  - [基本的な使用方法](#基本的な使用方法)
  - [`At`メソッドのセマンティクス](#atメソッドのセマンティクス)
- [JSON配列に値を追加](#json配列に値を追加)

---

## JSON値の作成

ほとんどの場合、オブジェクト型または配列型の外側のJSON値を作成する必要があります。jsonvalueでは、以下の関数を使用します：

```go
o := jsonvalue.NewObject()
a := jsonvalue.NewArray()
```

また、任意のJSON対応データ型をjsonvalueに変換（インポート）することもできます。`New`関数を使用するだけです。上記のオブジェクトと配列を例にとると：

```go
o := jsonvalue.New(struct{}{})  // JSONオブジェクトを構築
a := jsonvalue.New([]int{})     // JSON配列を構築
```

その他の単純なJSON要素もサポートされています：

```go
i := jsonvalue.New(100)             // JSON数値を構築
f := jsonvalue.New(188.88)          // JSON数値を構築
s := jsonvalue.New("Hello, JSON!")  // JSON文字列を構築
b := jsonvalue.New(true)            // JSONブール値を構築
n := jsonvalue.New(nil)             // JSONのnullを構築
```

---

## JSONにサブ値を設定

外側のオブジェクトまたは配列を生成した後、次のステップはその内部構造を作成することです。前のセクションで示した`Get`メソッドと同様に、`Set`または`MustSet`を使用してこれを実現できます。

`Set(xxx).At(yyy)`メソッドはサブ値とエラーを返します。一方、`MustSet(xxx).At(yyy)`は返しません。戻り値を気にしない場合は、golangci-lintの「戻り値が使用されていません」警告を回避するために`MustSet`メソッドを使用してください。

### 基本的な使用方法

一般的に、`Set`を使用して子の値を構築できます：

```go
v.MustSet(child).At(path...)
```

セマンティクスは「何かをある位置に設定する」です。値がキーより前に来ることに注意してください。

`Set`メソッドのパラメータ型は`any`であるため、サポートされている任意の型（複雑なオブジェクトや配列データでも）をjsonvalueに設定できます。

完全な例：

```go
v := jsonvalue.NewObject()
v.MustSet("Hello, JSON!").At("data", "message")
v.MustSet(221101).At("data", "date")
fmt.Println(v.MustMarshalString())
```

出力：`{"data":{"message":"Hello, JSON!","date":221101}}`

### `At`メソッドのセマンティクス

`Set`を呼び出した後、JSONに子の値を設定するために`At`を続ける必要があります。`At`のプロトタイプは：

```go
type Setter interface {
	At(firstParam any, otherParams ...any) (*V, error)
}
```

このメソッドの基本的なセマンティクスは`Get`と一致しています。プログラミングエラーを防ぐため、少なくとも1つのパラメータを指定する必要があります。これが`firstParam`の意味です。

`At`のより重要な機能は、対象のJSON構造を自動的に生成できることです。以下のロジックで処理されます：

- まず、指定されたパラメータによってサブ位置を特定します。これは`Get`と同様です。対象パスが既に存在する場合は、指定されたパスのサブ値を単純に設定します。
- 対象位置が存在しない場合、構造が自動的に作成されます。このメソッドでは`string`または`int`型のパラメータがサポートされています。文字列型はオブジェクトを識別し、整数は配列を識別します。

自動パス生成の例：

```go
v := jsonvalue.NewObject()                       // {}
v.MustSet("Hello, object!").At("obj", "message") // {"obj":{"message":"Hello, object!"}}
v.MustSet("Hello, array!").At("arr", 0)          // {"obj":{"message":"Hello, object!"},"arr":["Hello, array!"]}
```

配列の自動作成については、手順が少し複雑です：

- 指定されたパラメータで指定された配列が存在しない場合、操作を成功させるためにインデックス値はゼロである必要があります。
- 配列が既に存在する場合、以下の2つのケースのいずれかが適用されます：
  - 指定されたインデックスパラメータ値で指定された対応する子の値が存在する場合。この場合、そのスロットの値が置き換えられる可能性があります。
  - 指定されたインデックス値が配列の長さと等しい場合。この場合、値は配列の末尾に追加されます。

この機能は非常に複雑なため、ほとんどの場合は使用しません。しかし、有用な状況が1つあります：

```go
    var words = []string{"apple", "banana", "cat", "dog"}
    var lessons = []int{1, 2, 3, 4}
    v := jsonvalue.NewObject()
    for i := range words {
        v.MustSet(words[i]).At("array", i, "word")
        v.MustSet(lessons[i]).At("array", i, "lesson")
    }
    fmt.Println(v.MustMarshalString())
```

キーを前に置くことを好む場合は、`v.At(...).Set(...)`パターンを使用できます：

```go
    // ...
        v.At("array", i, "word").Set(words[i])
        v.At("array", i, "lesson").Set(lessons[i])
    // ...
```

最終出力：

```json
{"array":[{"word":"apple","lesson":1},{"word":"banana","lesson":2},{"word":"cat","lesson":3},{"word":"dog","lesson":4}]}
```

パラメータチェーンを識別するためにスライスまたは配列を渡すこともできます。例えば、以下のコード：

```go
v.MustSet("Hello, object!").At("obj", "message")
```

は以下と同等です：

```go
v.MustSet("Hello, object!").At([]any{"obj", "message"})
```

または：

```go
v.MustSet("Hello, object!").At([]string{"obj", "message"})
```

この機能により、外部ソースや設定からパラメータを簡単に渡すことができます。

---

## JSON配列に値を追加

メソッド`Append`と`Insert`は配列型JSONのために設計されています。`Append`は`InTheBeginning`と`InTheEnd`と組み合わせて動作し、`Insert`メソッドは`After`と`Before`と組み合わせて動作します。

これらは以下のセマンティクスで動作します：

- ...の始めに何らかの値を追加
- ...の終わりに何らかの値を追加
- ...の後に何らかの値を挿入
- ...の前に何らかの値を挿入

パラメータの順序に注意してください。

`Set`メソッドと同様に、同じ理由で`MustAppend`と`MustInsert`メソッドもあります。

これらのメソッドのプロトタイプは以下の通りです：

```go
func (v *V) Append(child any) Appender
type Appender interface {
	InTheBeginning(params ...any) (*V, error)
	InTheEnd(params ...any) (*V, error)
}

func (v *V) Insert(child any) Inserter
type Inserter interface {
	After(firstParam any, otherParams ...any) (*V, error)
	Before(firstParam any, otherParams ...any) (*V, error)
}

func (v *V) MustAppend(child any) MustAppender
type MustAppender interface {
	InTheBeginning(params ...any)
	InTheEnd(params ...any)
}

func (v *V) MustInsert(child any) MustInserter
type MustInserter interface {
	After(firstParam any, otherParams ...any)
	Before(firstParam any, otherParams ...any)
}
```

基本的なセマンティクスは`Set`メソッドと似ています。しかし、少し違いがあります：

- `InTheBeginning`と`InTheEnd`では空のパラメータが許可されており、これは現在のJSONが既に配列であることを識別し、サブ値がその始めまたは終わりに追加されることを示します。
- `After`と`Before`では、最後のパラメータは配列のインデックスを識別する数値である必要があります。負のインデックスも許可されています。