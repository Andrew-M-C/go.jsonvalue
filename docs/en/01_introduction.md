<font size=6>Introduction</font>

[Contents](./README.md) | [Next Page](./02_quick_start.md)

---

- [Brief Introduction](#brief-introduction)
- [Why Design Jsonvalue](#why-design-jsonvalue)
- [Applications](#applications)
- [Concurrency](#concurrency)

---

## Brief Introduction

Welcome to [jsonvalue](https://github.com/Andrew-M-C/go.jsonvalue). This is a repository developed in [Go](https://go.dev/).

This repo is designed to provide a package to process non-structural JSON data, as a replacement for the original `map[string]any` pattern type in `encoding/json`, which is quite inconvenient and inefficient. If you are interested in this, please refer to my benchmark testing [repo](https://github.com/Andrew-M-C/go.jsonvalue-test).

Also, I have received some requirements and suggestions for this repo. Some are for practical work, some are bug reports. All of these scenarios will be explained later in this wiki.

The [author](https://github.com/Andrew-M-C/) of this repo is a programmer at [Tencent](https://www.tencent.com), so in fact this repo is used more at Tencent than on Github.

---

## Why Design Jsonvalue

Why I decided to design this repo:

1. First of all, most work for JSON in Go is achieved by `struct` type. However, if you encounter non-structural JSON data, you may have to use `map[string]any` type. But as is well known, it is quite complex to handle `interface{}` type.
   - As a former ANSI-C programmer, [cJSON](https://github.com/DaveGamble/cJSON) is the most popular choice. So I developed this jsonvalue, which may be similar to cJSON.
2. Some JSON APIs of cloud services are strange or complex. For example, you may need to check some object in object in object ... etc.
   - So I added methods to access sub values in deep paths. Also methods to create deep, complex paths automatically.
   - This may save quite A LOT of code.
3. Some unexpected or unsupported situations for `encoding/json`. For example, `Inf` and `NaN` in IEEE float types, caseless (case-insensitive) support for `map[string]any`, etc.

---

## Applications

There are no silver bullets. For miscellaneous JSON scenarios, we should choose the most fitting solution.

As for Jsonvalue, the most recommended scenarios are:

- Quickly construct a complex JSON data structure and then serialize it. This is the most common usage as I know.
- Fully parsing and dumping/converting/analyzing non-structural JSON keys and values. This is widely used in offline data cleansing by my colleagues.
- Handle some illegal JSON values like `Inf` and `NaN`.
- Other special scenarios. Please refer to Section [Additional Options](./12_option.md)

---

## Concurrency

**IMPORTANT**: Jsonvalue is NOT goroutine-safe. For concurrent use, please use additional locking.
