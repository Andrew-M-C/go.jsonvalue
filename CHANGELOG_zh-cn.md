# 变更日志

[English](./CHANGELOG.md)

- [变更日志](#变更日志)
  - [v1.3.4](#v134)
  - [v1.3.3](#v133)
  - [v1.3.2](#v132)
  - [v1.3.1](#v131)
  - [v1.3.0](#v130)
  - [v1.2.1](#v121)
  - [v1.2.0](#v120)
  - [v1.1.1](#v111)


## v1.3.4

发布 jsonvalue v1.3.4，相比 v1.3.3:

- 添加 `*jsonvalue.V` 对以下接口的支持, 使之可以适配 `encoding/json`:
  - `json.Marshaler`、`json.Unmarshaler`
  - `encoding.BinaryMarshaler`、`encoding.BinaryUnmarshaler`
- `Import` 和 `New` 函数支持导入实现了 `json.Marshaler` 和 `encoding.TextMarshaler` 的数据类型, 就像 `encoding/json` 一样 ([#24](https://github.com/Andrew-M-C/go.jsonvalue/issues/24))
- 添加 `MustAdd`, `MustAppend`, `MustInsert`, `MustSet`, `MustDelete` 等方法, 这些方法不返回子值和 error 信息
  - 可以避免 golangci-lint 的警告提示
- Bug 修复: 当 `Append(xxx).InTheBeginning()` 时, 实际上追加到了末尾
- 在 Unmarshal 时, 通过预估对象 / 数组大小的方式预创建空间, 略微提升一点点速度, 以及节省平均约 40% 的 alloc 数

## v1.3.3

发布 jsonvalue v1.3.3，相比 v1.3.2:

- 发布英文 [wiki](https://github.com/Andrew-M-C/go.jsonvalue/blob/master/docs/en/README.md)。
- 修复了 #19 和 #22。

## v1.3.2

发布 jsonvalue 1.3.2，相比 1.3.1，主要是修复了 [#17](https://github.com/Andrew-M-C/go.jsonvalue/issues/17)。

## v1.3.1

发布 jsonvalue 1.3.1，相比 1.3.0，改动如下：

- 支持获取 object 类型的键（key）的顺序
- 按照 key 被设置的顺序进行 object 的序列化
- 当序列化时，如果遇到 ASCII 的非可打印字符（包括控制字符），均进行 `\u00XX` 转义

具体说明情参见 [wiki](https://github.com/Andrew-M-C/go.jsonvalue/blob/feature/v1.3.0/docs/zh-cn/12_new_feature.md)

## v1.3.0

发布 v1.3.0，相比 v1.2.x，改动如下：

- 不再默认转义 % 字符
- 支持与原生 json 的转换，也即支持从 `struct` 和其他类型中导入为 jsonvalue 类型，也支持将 jsonvalue 导出到 Go 标准类型中
- 支持类泛型操作，最直接的影响就是 `Set, Append, Insert, Add, New` 等函数可以传入任意类型
- 新增 `OptIndent` 函数以支持可视化锁进

具体说明情参见 [wiki](https://github.com/Andrew-M-C/go.jsonvalue/blob/feature/v1.3.0/docs/zh-cn/12_new_feature.md)

## v1.2.1

发布 v1.2.1，相比 v1.2.1，改动如下：

- Marshal 时，允许不转义斜杠符号 '/'，参见 `OptEscapeSlash` 函数。
- Marshal 时，允许不转义几个 HTML 关键符号 (issue #13)。参见 OptEscapeHTML 函数
- Marshal 时，对于大于 `\u00FF` 的字符，允许保留 UTF-8 编码，而不进行转义。参见 OptUTF8 函数
- `Append(...).InTheBeginning()` 和 `Append(...).InTheEnd()` 操作时，支持自动创建不存在的层级。但 insert 操作依然不会自动创建，这一点请留意
- 支持设置默认序列化选项，请参见 SetDefaultMarshalOptions 函数

其他变更:

- CI 弃用 Travis-CI，改用 Github Actions

## v1.2.0

发布 v1.2.0，相比 v1.1.1，改动如下：

- 修复不支持科学计数法浮点数的 bug（[issue #8](https://github.com/Andrew-M-C/go.jsonvalue/issues/8)）
- 添加详细的中文 [wiki](https://github.com/Andrew-M-C/go.jsonvalue/blob/master/docs/zh-cn/README.md)。
- Marshal 时的额外参数配置，原使用的是 `Opt` 结构体，现推荐改为使用一个或多个 `OptXxx()` 函数
- 发布第一个不向前兼容的变更: 修改 `NewFloat64` 和 `NewFloat32` 的入参
  - 原函数的 `prec` 变量及相关功能改由 `NewFloat64f` 和 `NewFloat32f` 承担。
- 添加 `ForRangeArr` 和 `ForRangeObj` 以取代 `IterArray` 和 `IterObject` 函数，参见 [wiki](https://github.com/Andrew-M-C/go.jsonvalue/blob/master/docs/zh-cn/05_iteration.md#%E6%A6%82%E8%BF%B0)
- 支持 `SetEscapeHTML`（[issue11](https://github.com/Andrew-M-C/go.jsonvalue/issues/11)）
- 当浮点值中出现 `+/-NaN` 和 `+/-Inf` 时，支持进行特殊转换以符合 JSON 标准

## v1.1.1

发布 v1.1.1，相比 v1.1.0，变更如下：

- 无论是否错误，所有返回 `*jsonvalue.V` 的函数都不会返回 nil。如果错误发生，那么至少会返回一个有效的 `*V` 实例，但是 `Type` 等于 `NotExist`
  - 这样一来，程序员可以在代码中放心地直接使用返回的 *V 对象，简化代码，而忽略不必要的错误检查。
- 支持 `MustGet` 函数，只返回 `*V` 而不返回 error 变量，便于一些级联的代码。
- `ValueType` 类型支持 `String()` 函数。

