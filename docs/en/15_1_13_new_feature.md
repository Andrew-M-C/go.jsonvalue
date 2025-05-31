<font size=6>New Features in v1.3.x</font>

[Prev Page](./14_1_14_new_feature.md) | [Contents](./README.md) | [Next Page](./16_1_12_new_feature.md)

---

- [`%` Will Not be Escaped by Default](#-will-not-be-escaped-by-default)
- [Supporting Import from / Export to `encoding/json`](#supporting-import-from--export-to-encodingjson)
- [Supporting Generics-Like Operations](#supporting-generics-like-operations)
- [Comparing JSON Values](#comparing-json-values)
- [Support for Indentation when Marshaling](#support-for-indentation-when-marshaling)
- [Getting Original Sequence of Keys](#getting-original-sequence-of-keys)
- [Marshal by the Sequence in Which Keys Are Set](#marshal-by-the-sequence-in-which-keys-are-set)
- [Escaping Non-Visible ASCII Characters](#escaping-non-visible-ascii-characters)
- [MustXxx methods](#mustxxx-methods)
- [Support for Official Marshaler and Unmarshaler Interfaces](#support-for-official-marshaler-and-unmarshaler-interfaces)

---

## `%` Will Not be Escaped by Default

In the previous version of jsonvalue, the `%` character was escaped, but this is not standard. From v1.3.0, `%` will not be escaped by default. If you need this feature, please file an [issue](https://github.com/Andrew-M-C/go.jsonvalue/issues/new).

## Supporting Import from / Export to `encoding/json`

Please refer to Section [Import and Export with `encoding/json`](./06_import_export.md). It was a beta feature in v1.2.x and became official in v1.3.x.

## Supporting Generics-Like Operations

In the previous version, when you set a sub value into `*V`, you should specify types like:

```go
v.SetString("Hello, JSON!").At("msg")
v.SetInt64(time.Now().Unix()).At("time")
```

From v1.3.0, you do not need to do this any more. Just use `Set`, which accepts any type of legal JSON values:

```go
v.Set("Hello, JSON!").At("msg")
v.Set(time.Now().Unix()).At("time")
```

This feature also applies to other functions including:

- `Append`
- `Insert`
- `Add`

Besides, the `New` function was also introduced in v1.3.0, which is actually a simple wrapper for `Import`. Please refer to Section [Import and Export with `encoding/json`](./06_import_export.md).

## Comparing JSON Values

The `Equal` method is introduced to compare whether two JSON values are equal. The [Contains](./13_beta.md) function is also added in Package beta.

## Support for Indentation when Marshaling

Please refer to "Visible Indention" in Section [Additional Options](./12_option.md).

## Getting Original Sequence of Keys

Please refer to "Acquiring the Original Key Sequence of An Object" in Section [Iteration](./07_iteration.md).

## Marshal by the Sequence in Which Keys Are Set

Please refer to "Serialize a JSON Object with Sequence of When Keys Are Set" in Section [Special Application Scenarios](./10_scenarios.md).

## Escaping Non-Visible ASCII Characters

From v1.3.1, if non-visible ASCII characters appear in keys or string-typed values, they will be escaped in the format `\u00XX`, to prevent inappropriate display effects.

## MustXxx methods

From v1.3.4, `MustAdd`, `MustAppend`, `MustInsert`, `MustSet`, `MustDelete` methods are introduced. These methods will not return sub-values or errors.

## Support for Official Marshaler and Unmarshaler Interfaces

From v1.3.4, `*jsonvalue.V` implements the following official interfaces:

- `json.Marshaler`、`json.Unmarshaler`
- `encoding.BinaryMarshaler`、`encoding.BinaryUnmarshaler`

Please refer to [Marshal and Unmarshal](./05_marshal_unmarshal.md)
