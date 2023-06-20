# Change Log

[中文](./CHANGELOG_zh-cn.md)

- [Change Log](#change-log)
  - [v1.3.4](#v134)
  - [v1.3.3](#v133)
  - [v1.3.2](#v132)
  - [v1.3.1](#v131)
  - [v1.3.0](#v130)
  - [v1.2.1](#v121)
  - [v1.2.0](#v120)
  - [v1.1.1](#v111)


## v1.3.4

Release jsonvalue v1.3.4. Comparing to v1.3.3:

- `*jsonvalue.V` type now implements following interfaces thus it can used directly in `encoding/json`:
  - `json.Marshaler`、`json.Unmarshaler`
  - `encoding.BinaryMarshaler`、`encoding.BinaryUnmarshaler`
- Supports importing data types those implement `json.Marshaler` or `encoding.TextMarshaler` interfaces, just like `encoding/json` does.
- Add `MustAdd`, `MustAppend`, `MustInsert`, `MustSet`, `MustDelete` methods. These methods will not return sub-values or errors.
  - This will prevent golangci-lint warning of "return value is not used"
- Bug fix: When invoking `Append(xxx).InTheBeginning()`, the sub-value was actually and faulty append to the end.
- Pre-allocate some buffer and space when marshaling and unmarshaling. It will speed-up a little bit and save 40% alloc count by average.

## v1.3.3

Release jsonvalue v1.3.3. Comparing to v1.3.2:

- Release Engligh [wiki](https://github.com/Andrew-M-C/go.jsonvalue/blob/master/docs/en/README.md).
- Fix #19 and #22.

## v1.3.2

Release jsonvalue 1.3.2. Comparing to v1.3.1, [#17](https://github.com/Andrew-M-C/go.jsonvalue/issues/17) is fixed.

## v1.3.1

Release jsonvalue 1.3.1. Comparing to v1.3.0:

- It is available to get key sequences of an unmarshaled object JSON
- Supporting marshal object by key sequence of when they are set.
- All invisible ASCII characters will be escaped in `\u00XX` format.

## v1.3.0

Release v1.3.0. Comparing to v1.2.x:

- No longer escape % symbol by default.
- Add `Import` and `Export` function to support conversion between Go standard types and jsonvalue.
- Support generics-like operations in `Set, Append, Insert, Add, New` functions, which allows programmers to set any legal types without specifying types explicitly.
- Add `OptIndent` to support indent in marshaling.

## v1.2.1

Release v1.2.1. Comparing to v1.2.0:

- Support not escaping slash symbol / in marshaling. Please refer to `OptEscapeSlash` function.
- Support not escaping several HTML symbols in marshaling (issue #13). Please refer to OptEscapeHTML function.
- Support not escaping unicodes greater than `\u00FF` in marshaling. Please refer to OptUTF8 function.
- Support children-auto-generating in `Append(...).InTheBeginning()` and `Append(...).InTheEnd()`. But keep in mind that Insert operations remains not doing that.
- You can override default marshaling options from now on. Please refer to `SetDefaultMarshalOptions` function.

Other changes:

- Instead of Travis-CI, I will use Github Actions in order to generate go test reports from now on.

## v1.2.0

Release v1.2.0, comparing to v1.1.1:

- Fix the bug that float-pointed numbers with scientific notation format are not supported. ([issue #8](https://github.com/Andrew-M-C/go.jsonvalue/issues/8))
- Add detailed wiki written in simplified Chinese. If you need English one, please leave me an issue.
- Deprecate `Opt{}` for additional options, and replace with `OptXxx()` functions.
- A first non-forward-compatible feature is released: Parameter formats for `NewFloat64` and `NewFloat32` is changed. 
  - If you want to specify format of float-pointed numbers, please use `NewFloat64f` and `NewFloat32f` instead.
- Add `ForRangeArr` and `ForRangeObj` functions, deprecating `IterArray` and `IterObject`.
- Support `SetEscapeHTML`. ([issue11](https://github.com/Andrew-M-C/go.jsonvalue/issues/11))
- Support converting number to legal JSON data when they are `+/-NaN` and `+/-Inf`.

## v1.1.1

Release v1.1.1, comparing to v1.1.0:

- All functions will NOT return nil `*jsonvalue.V` object instead of an instance with `NotExist` type when error occurs.
  - Therefore, programmers can use the returned instance for shorter operations.
- Supports `MustGet()` function, without error returning.
- Supports `String()` for `ValueType` type
