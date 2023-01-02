
<font size=6>Introduction</font>

[Contents](./README.md) | [Next Page](./02_quick_start.md)

---

- [Brief Introduction](#brief-introduction)
- [Why Design Jsonvalue](#why-design-jsonvalue)
- [Applications](#applications)
- [Concurrency](#concurrency)

---

## Brief Introduction

Welcome to [jsonvalue](https://github.com/Andrew-M-C/go.jsonvalue). This is a repository developed by [Go](https://go.dev/).

This repo is designed to provide a package to process non-structural JSON data, as a replacement of original `map[string]interface{}` patten type in `encoding/json`, which is quite in-convenient and inefficient. If you are interested in this, please refer to my benchmark testing [repo](https://github.com/Andrew-M-C/go.jsonvalue-test).

Also, I received some requirement and suggestions for this repo. Some are for practical works, some are bug reports. All of these scenarios will be explained later in this wiki.

[Author](https://github.com/Andrew-M-C/) os this repo is a programmer in [Tencent](https://www.tencent.com), so in fact this repo is used more in Tencent than Github.

---

## Why Design Jsonvalue

Why I decided to design this repo:

1. First of all, most works for JSON in Go is achieved by `struct` type. However, if you meet non-structural JSON data, you may have to use `map[string]interface{}` type. But as well known, it is quite complex to handle `interface{}` type.
   - As a former ANSI-C programmer, [cJSON](https://github.com/DaveGamble/cJSON) is the most popular choice. So I developed this jsonvalue, which may be similar to cJSON.
2. Some JSON API of cloud service is strange or complex. For example, you may need to check some object in object in object ... etc. 
   - So I add methods to access sub values in deep path. Also methods to create deep, complex paths automatically.
   - This may save quite A LOT if codes.
3. Some unexpected or unsupported situations for `encoding/json`. For example, `Inf` and `NaN` in IEEE float types, Caseless support for `map[string]interface{}`, etc.

---

## Applications

There are no silver bullets. For miscellaneous JSON scenarios, we should choose the fittest solution.

As for Jsonvalue, the most recommended scenarios are:

- Quickly construct a complex JSON data structure and then serialize it. This is the most comment usages as I known.
- Fully parsing and dump/convert/analyzing non-structural JSON keys and values. This is widely used in offline data cleansing by my colleagues.
- Handle some illegal JSON values like `Inf` and `NaN`.
- Other special scenarios. Please refer to Section [Additional Options](./12_option.md)

---

## Concurrency

**IMPORTANT**: Jsonvalue is NOT goroutine-safe. For concurrency use, please use additional locker.
