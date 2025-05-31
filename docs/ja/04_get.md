<font size=6>JSON構造から値を取得する</font>

[前のページ](./03_set.md) | [目次](./README.md) | [次のページ](./05_marshal_unmarshal.md)

---

- [Get関数](#get関数)
  - [パラメータ](#パラメータ)
  - [GetXxxシリーズ](#getxxxシリーズ)
- [`MustGet`とその他の関連メソッド](#mustgetとその他の関連メソッド)
- [jsonvalue.Vオブジェクトの属性](#jsonvaluevオブジェクトの属性)
  - [公式定義](#公式定義)
  - [jsonvalue基本属性](#jsonvalue基本属性)

---

## Get関数

### パラメータ

Get関数はjsonvalueの情報読み取りの中核となる機能です。以下がプロトタイプです：

```go
func (v *V) Get(param1 any, params ...any) (*V, error)
```

実用的な例：

```go
const raw = `{"someObject": {"someObject": {"someObject": {"message": "Hello, JSON!"}}}}`
child, _ := jsonvalue.MustUnmarshalString(s).Get("someObject", "someObject", "someObject", "message")
fmt.Println(child.String())
```

上記の`Get`パラメータの意味は以下の通りです：

- `*V`インスタンスからキー`someObject`を持つサブ値を特定して取得し、次に前回特定した`*V`インスタンスから別のキー`someObject`を持つ値を取得し、以下同様に続けます...
  - この操作はドメイン形式でも記述できます：`child = v.someObject.someObject.someObject.message`

`Get`のパラメータの型は`any`です。実際には、文字列または整数（符号付きと符号なしの両方が可能）の種類のもののみが許可されます。`Get`はパラメータの型をチェックし、次の値をオブジェクトまたは配列として扱うかを決定します。

現在のパスノードのパラメータの[Kind](https://pkg.go.dev/reflect#Kind)が文字列の場合：

- 現在の`jsonvalue`値の値型が「Object」の場合、文字列キーで指定されたサブ値を特定します。
  - サブ値が存在し、それが「Object」型の場合、この値と次のパスノードのパラメータで反復を続行します。
  - 指定されたキーを持つ値が存在しない場合、`NotExist`型の値がエラー[ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants)と共に返されます。
- 現在の値が「Object」でない場合、`NotExist`型の値が別のエラー[ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants)と共に返されます。

現在のパスノードのパラメータの[Kind](https://pkg.go.dev/reflect#Kind)が整数の場合：

- 現在の`jsonvalue`値の値型が「Array」の場合、整数キーで指定されたインデックスからサブ値を見つけます。この時、この整数の様々な値の意味は以下の通りです：
  - インデックス >= 0の場合、通常のインデックス値となり、通常のGoスライスのようにサブ値を特定します。インデックスが範囲外の場合、`NotExist`値と[ErrOutOfRange](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants)が返されます。
  - インデックス < 0の場合、「最後からXX番目」として扱われ、逆方向にカウントします。ただし、JSON配列の範囲内である必要があります。
    - 例えば、JSON配列の長さが5の場合、-5はインデックス0の要素を特定し、-6はエラーを返します。
  - 現在のJSONノードで検索が成功した場合、残りのパラメータがあれば反復が続行されます。
- パスノードの現在の値が「Array」でない場合、`NotExist`型の値が別のエラー[ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants)と共に返されます。

なぜ`Get`関数のパラメータを2つの部分に分けたのか、単純な`...any`ではなく、疑問に思うかもしれません。[この問題](https://github.com/Andrew-M-C/go.jsonvalue/issues/4)で答えたように、意図的に設計しました：

- これは`v.Get`（パラメータ不足）のようなプログラミングエラーを避けるためです。このメソッドを少なくとも1つのパラメータを持つようにすることで、実行時ではなくコンパイル段階でエラーが発生します。
- 入力`[]any`に少なくとも1つのパラメータがあることを100%確信している場合、このメソッドを次のように呼び出すことができます：`subValue, _ := Get(para[0], para[1:]...)`

### GetXxxシリーズ

実際のコードでは`Get`自体はほとんど使用されず、代わりに基本型の値にアクセスするために「兄弟」メソッドを使用します：

```go
func (v *V) GetObject (param1 any, params ...any) (*V, error)
func (v *V) GetArray  (param1 any, params ...any) (*V, error)
func (v *V) GetBool   (param1 any, params ...any) (bool, error)
func (v *V) GetString (param1 any, params ...any) (string, error)
func (v *V) GetBytes  (param1 any, params ...any) ([]byte, error)
func (v *V) GetInt    (param1 any, params ...any) (int, error)
func (v *V) GetInt32  (param1 any, params ...any) (int32, error)
func (v *V) GetInt64  (param1 any, params ...any) (int64, error)
func (v *V) GetNull   (param1 any, params ...any) error
func (v *V) GetUint   (param1 any, params ...any) (uint, error)
func (v *V) GetUint32 (param1 any, params ...any) (uint32, error)
func (v *V) GetUint64 (param1 any, params ...any) (uint64, error)
func (v *V) GetFloat32(param1 any, params ...any) (float32, error)
func (v *V) GetFloat64(param1 any, params ...any) (float64, error)
```

これらのメソッドには共通点があります：

- サブ値が存在し、メソッドで指定された正しい型である場合、返される`error`はnilです。そして`GetNull`以外のすべてのメソッドは対応する値を返します。
- サブ値は存在するが型が一致しない場合、エラーは[ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants)になります。
- サブ値が存在しない場合、エラーは[ErrNotFound](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants)になります。

また、これらのメソッドの一部は単純に型をマッチして返すだけでなく、追加機能も提供します。これについては後のセクションで説明します。ここでは例を示します：

多くの場合、文字列型の値から数値を抽出する必要があります。例：`{"number":"12345"}`。この場合、`GetInt`メソッドはこの文字列から対応する整数値を正しく返します：

```go
raw := `{"number":"12345"}`
n, err := jsonvalue.MustUnmarshalString(raw).GetInt("number")
fmt.Println("n =", n)
fmt.Println("err =", err)
```

出力：

```
n = 12345
err = not match given type
```

例で示されているように、`n`と`err`の両方が意味のある値を返します。有効な数値が返されるだけでなく、[ErrTypeNotMatch](https://pkg.go.dev/github.com/Andrew-M-C/go.jsonvalue#pkg-constants)エラーも発生します。

---

## `MustGet`とその他の関連メソッド

上記で`Get`と`GetXxx`シリーズ関数について言及しました。`GetNull`を除いて、各関数の戻り値は2つです。Get関数に対して、jsonvalueは`MustGet`関数も提供しており、1つのパラメータのみを返すため、非常にシンプルなロジックの実装を容易にします。

理解しやすくするために、シナリオを例として挙げましょう——

例えば、フォーラム機能を開発しているとします。フォーラムは複数の投稿をピン留めする機能をサポートしています。ピン留め機能の設定は、JSON文字列の「top」フィールドです。例：

```json
{
    "//":"other configs",
    "top":[
        {
            "UID": "12345",
            "title": "投稿規範"
        }, {
            "UID": "67890",
            "title": "フォーラム精選"
        }
    ]
}
```

実際には、様々な理由により、取得した設定文字列には以下のような異常な状況が発生する可能性があります：

- 文字列全体が空文字列""
- 文字列が誤った編集により不正、または形式エラー
- 「top」フィールドが`null`、または空文字列の可能性

従来のロジックに従えば、これらの異常な状況を一つずつ処理する必要があります。しかし、開発者がこれらの異常を気にせず、有効な設定のみに関心がある場合、`MustXxx`関数が必ず`*V`オブジェクトを返すという特性を利用して、ロジックを以下のように簡素化できます：

```go
    c := jsonvalue.MustUnmarshalString(confString) // confStringは取得した設定文字列と仮定
    for _, v := range c.Get("top").ForRangeArr() {
        feeds = append(feeds, &Feed{               // 投稿テーマを返却リストに追加、投稿の構造体をFeedと仮定
            ID:    v.MustGet("UID").String(),
            Title: v.MustGet("title").String(),
        })
    }
```

---

## jsonvalue.Vオブジェクトの属性

まず、JSON公式定義のいくつかの属性を理解し、次にこれらの属性が`jsonvalue`でどのように表現されるかを説明します。

### 公式定義

標準の[JSON仕様](https://www.json.org/json-en.html)では、以下の概念が規定されています：

- 有効なJSON値は、JSONの`value`と呼ばれます。このツールパッケージでは、`*V`を使用してJSON valueを表現します
- JSON値の型には以下の種類があります：

|    型     | 説明                                                                                                                    |
| :-------: | :---------------------------------------------------------------------------------------------------------------------- |
| `object`  | オブジェクト、K-V形式の値に対応します。Kは必ず文字列で、Vは有効なJSON`value`です                                        |
|  `array`  | 配列、一連の`value`の順序付きの組み合わせに対応します                                                                   |
| `string`  | 文字列型、これは理解しやすいです                                                                                        |
| `number`  | 数値型、正確には倍精度浮動小数点数です                                                                                  |
|           | JSONはJavaScriptに基づいて定義されており、JSにはdoubleという1つの数値型しかないため、numberは実際にはdoubleです。これは小さな落とし穴です |
| `"true"`  | ブール値「真」を表します                                                                                                |
| `"false"` | ブール値「偽」を表します                                                                                                |
| `"null"`  | null値を表します                                                                                                        |

### jsonvalue基本属性

`*jsonvalue.V`オブジェクトでは、ほとんどのJSONツールパッケージの手法を参考にして、`"true"`と`"false"`を1つの`Boolean`型に統合しています。さらに、`"null"`も`Null`型にマッピングしています。

さらに、現在が有効なJSONオブジェクトではないことを表す`NotExist`型も定義されています。また、`Unknown`もありますが、開発者は気にする必要はなく、使用中にこの値が現れることはありません。

以下の関数を使用して、valueの型属性を取得できます：

```go
func (v *V) ValueType() ValueType
func (v *V) IsObject()  bool
func (v *V) IsArray()   bool
func (v *V) IsString()  bool
func (v *V) IsNumber()  bool
func (v *V) IsBoolean() bool
func (v *V) IsNull()    bool
```
