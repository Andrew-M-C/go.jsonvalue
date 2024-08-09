
<font size=6>New Features in v1.4.x</font>

[Prev Page](./13_beta.md) | [Contents](./README.md) | [Next Page](./15_1_13_new_feature.md)

---

- [Allow Parsing Slice as Parameter](#allow-parsing-slice-as-parameter)
- [Number Comparing](#number-comparing)

---

## Allow Parsing Slice as Parameter

From v1.4.0, a single slice or array for Get, Set, Append, Insert, Delete methods. It takes the same effect with passing each parameters in the slice. Please refer to [Create and Serialize JSON](./03_set.md).

Please allow me to explain why I decided to support this feature. In some cases, programmers may wanted to use a `[]string` slice or `[]any` slice identifying the key chain. But types the `At` series methods is like `(any, ...any)` (in fact `(any, []any)`). Therefore programmers had to write some codes to convert the slice into an `any` and a `[]any`. So I decided to simplify this, now we can simply pass the slice as parameter.

---

## Number Comparing

From v1.4.0, comparing methods are supported.

