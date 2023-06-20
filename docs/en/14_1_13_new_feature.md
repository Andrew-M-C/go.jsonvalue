
<font size=6>New Features in v1.3.x</font>

[Prev Page](./13_beta.md) | [Contents](./README.md) | [Next Page](./15_1_12_new_feature.md)

---

- [`%` Will Not be Escaped by Default](#-will-not-be-escaped-by-default)
- [Supporting Import from / Export to `encoding/json`](#supporting-import-from--export-to-encodingjson)
- [Supporting Generics-Like Operations](#supporting-generics-like-operations)
- [Comparing JSON Values](#comparing-json-values)
- [Support Indent when Marshal](#support-indent-when-marshal)
- [Getting Original Sequence of Keys](#getting-original-sequence-of-keys)
- [Marshal by Sequence of When Keys Are Set](#marshal-by-sequence-of-when-keys-are-set)
- [Escaping Non Visible ASCII Characters](#escaping-non-visible-ascii-characters)
- [MustXxx methods](#mustxxx-methods)
- [Support some official marshaler and unmarshaler interfaces](#support-some-official-marshaler-and-unmarshaler-interfaces)

---

## `%` Will Not be Escaped by Default

In previous version of jsonvalue, the `%` character will be escaped, but this is not standard. From v1.3.0, `%` will not be escaped by default. If you need this feature, please raise me an [issue](https://github.com/Andrew-M-C/go.jsonvalue/issues/new)。

## Supporting Import from / Export to `encoding/json`

Please refer to Section [Import and Export with `encoding/json`](./06_import_export.md) 小节。It is beta feature in v1.2.x and official in v1.3.x.

## Supporting Generics-Like Operations

In previous version, when you set a sub value into `*V`, you should specify types like:

```go
v.SetString("Hello, JSON!").At("msg")
v.SetInt64(time.Now().Unix()).At("time")
```

From v1.3.0, you do not need to do this any more. Just use `Set`, which accepts any type of legal JSON values:

```go
v.Set("Hello, JSON!").At("msg")
v.Set(time.Now().Unix()).At("time")
```

This feature is also OK to other functions including:

- `Append`
- `Insert`
- `Add`

Besides, `New` function is also introduced from v1.3.0, which actually a simple packaging to `Import`. Please refer to Section [Import and Export with `encoding/json`](./06_import_export.md).

## Comparing JSON Values

`Equal` method is introduced to compare whether two JSON values equal to each other. Also [Contains](./13_beta.md) function is added in Package beta.

## Support Indent when Marshal

Please refer to "Visible Indention" in Section [Additional Options](./12_option.md).

## Getting Original Sequence of Keys

Please refer to "Acquiring the Original Key Sequence of An Object" in Section [Iteration](./07_iteration.md).

## Marshal by Sequence of When Keys Are Set

Please refer to "Serialize a JSON Object with Sequence of When Keys Are Set" in Section [Special Application Scenarios](./10_scenarios.md).

## Escaping Non Visible ASCII Characters

From v1.3.1, if non-visible ASCII characters appears in keys of string typed values, they will be escaped by format `\u00XX`, to prevent inappropriate display effect.

## MustXxx methods

From v1.3.4, `MustAdd`, `MustAppend`, `MustInsert`, `MustSet`, `MustDelete` methods are introduced. These methods will not return sub-values or errors.

## Support some official marshaler and unmarshaler interfaces

From v1.3.4, `*jsonvalue.V` implements following official interfaces:

- `json.Marshaler`、`json.Unmarshaler`
- `encoding.BinaryMarshaler`、`encoding.BinaryUnmarshaler`

Please refer to [Marshal and Unmarshal](./05_marshal_unmarshal.md)
