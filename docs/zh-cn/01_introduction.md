# Jsonvalue 简介

[总目录](./README.md) | [下一页](./02_quick_start.md)

---

[TOC]

---

## 引言

欢迎您关注 [jsonvalue](https://github.com/Andrew-M-C/go.jsonvalue)。本项目使用 [Go 语言](https://golang.org/)进行开发。

本项目的设计目的，是提供一个能够简易地处理非结构化 JSON 数据的工具包，用以替代 Go 原生的 `map[string]interface{}` 方案。

Go 原生 `map[string]interface{}` 方案在使用的过程中，不仅不够简便，而且效率也不高，读者如果感兴趣，可以移步本人的一篇[博客](https://cloud.tencent.com/developer/article/1676060)查阅。

在迭代的过程中，也因应作者本人工作上的实际用途，新增其他一些功能。这些功能将在本 wiki 中详细解释。作者 [Andrew-M-C](https://github.com/Andrew-M-C/) 本人目前是[腾讯](https://www.tencent.com/)的一个后台小码农，因此这个 package 实际上在腾讯内部用得比外部多。

## 设计初衷

Jsonvalue 的设计初衷有下面几个:

1. 作者在实际编码中需要处理一些类型不确定的 JSON 数据，因此不能用 struct。用原生的话，只能用 `map[string]interface{}`，但是这种方案要自己处理 interface{}，太难受了。作为 C 语言出身的程序员，笔者因此借鉴了 [cJSON](https://github.com/DaveGamble/cJSON) 的模式，开发了这个库
1. 一些云 API 的 JSON 格式很奇葩，比如说需要判断是否包含某个对象，并且判断对象中的某个字段是否为某个值。不论是要构造，还是解析，都很难受。因此在 jsonvalue 里就加入了深层 path 解析和配置的功能（你可以参考 Get 和 Set().At() 函数的那一长串参数）
1. 在处理 JSON 的过程中，遇到过一些奇奇怪怪的、原生 `json` 无法解决的问题，比如说对浮点 `Inf` 和 `NaN` 不支持、对忽略大小写不完全支持等等，都促使作者在这个轮子一步一步地迭代一些奇奇怪怪的功能。

## 应用场景

天下没有哪一个代码是万金油，针对不同的应用场景，选取最合适的才是王道。

之前笔者门针对我常用的几个 JSON 库进行了比较分析:《[Go 语言原生的 json 包有什么问题？如何更好地处理 JSON 数据？](https://segmentfault.com/a/1190000039957766)》。

在文中我也没推荐在所有场景下使用 jsonvalue。Jsonvalue 最典型和不可替代的应用场景是下面几个:

1. 快速构建复杂的 JSON 并序列化——我观察了一下使用这个库的其他 repo，大部分应用场景是这个，这也是很标准的设计初衷
1. 全量解析并 dump/转换 KV——这个一般在离线清洗的应用场景中，在我们公司内部的应用中有这么搞的
1. 部分非法的 JSON 值处理——这个主要是针对 float 的，在 JSON 中不支持双精度数的 Inf 和 N/A，但是 Go 的 struct 不可避免地会出现这个数（特别是在算法领域）。我们遇到了，所以就借用 jsonvalue 这个壳来解决，特别是搭配 v1.3.0 之后新增的 `Import/Export` 以及 jsonvalue 原生的序列化函数
1. 其他特殊的场景——请读者可以具体参阅 [额外选项配置](./08_option.md) 小节

## 其他说明

- 本 package 不支持协程安全，如果要在多协程环境中使用，请额外实现锁操作
