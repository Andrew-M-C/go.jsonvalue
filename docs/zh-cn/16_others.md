
<font size=6>其他</font>

[上一页](./12_new_feature.md) | [总目录](./README.md)

---

- [关于泛型的思考](#关于泛型的思考)
- [不使用泛型的原因](#不使用泛型的原因)
- [泛型支持计划](#泛型支持计划)

---

## 关于泛型的思考

在 jsonvalue 的 [issue](https://github.com/Andrew-M-C/go.jsonvalue/issues?q=) 列表中，[harryhan1989](https://github.com/Andrew-M-C/go.jsonvalue/issues/14) 批评本库不支持泛型。但是这个 issue 提得很奇怪，我也不知道 issuer 希望我以何种方式支持，也没有提出任何建设性意见。

实际上在 v1.3.0 中实现了类泛型（Generics-like）的[操作](./12_new_feature.md)，因此这个 issue 不再存在。实际上，我并不是通过泛型、而是使用 `any` 来解决的。

## 不使用泛型的原因

在 Go 1.20 推出之前，笔者暂不打算支持泛型，原因如下：

- 泛型的核心是 DRY，也就是 Don't Repeat Yourself。这主要是不同类型、但是逻辑写起来一模一样的场景。但是实际上，jsonvalue 来说，不同类型的处理逻辑差异是很大的，并不存在很多 DRY 的需求，仅有一些需求是 `SetXxx` 方法。但是 Go 的泛型目前并不支持方法，所以无法实现。
- 在使用体验上，jsonvalue 已经尽量降低对类型的依赖了。最典型的就是 `Set().At()` 和 `Get()` 的参数，可以让使用者无需区分 string 和 int 参数
- 从 v1.3.0 开始，jsonvalue 魔改了 `Set()` 函数以及其他[相关函数](./14_1_13_new_feature.md)，将这些函数的参数从 `*V` 换成了 `any`，提供与泛型几乎无异的编程体验。
- 从 [go.mod](../../go.mod) 文件中也可以看到，jsonvalue 的支持版本是 Go 1.13+，Go 对泛型的正式支持在 1.18+，一旦支持了泛型，那么依赖 jsonvalue 的代码将可能无法编译，这种过河拆桥的行为是为 Go 社区所不齿的。

## 泛型支持计划

目前的 Go 泛型是不支持方法的，因此笔者还在密切跟进泛型的开发进展，根据未来变化再做决定。
