<font size=6>Value Comparing</font>

[Prev Page](./10_scenarios.md) | [Contents](./README.md) | [Next Page](./12_option.md)

---

- [Equal](#equal)
- [Number Comparison](#number-comparison)

---

## Equal

From v1.3.0, the `Equal` method is added to tell whether two JSON values are the same. The comparison rules are:

Firstly, if the types of two values are different, it returns `false`.

If the types are the same, check them by types as follows:

- `string`: check if two string values are the same.
- `number`: check if the DECIMAL numbers of two values are equal. I use the [decimal](https://pkg.go.dev/github.com/shopspring/decimal) package to do this.
- `boolean`: check if the two boolean values are the same.
- `null`: all null values are equal to each other.
- `array`: the necessary and sufficient condition for two arrays to be equal to each other is: they share the same array length, and each element equals the other one at the same index.
- `object`: the necessary and sufficient condition for two objects to be equal to each other is: they share the exact same keys, and each element equals the other one with the same key.

---

## Number Comparison

From v1.4.0, `GreaterThan`, `LessThan`, `GreaterThanOrEqual`, `LessThanOrEqual` are added for comparing two numeric values. If either value is not a number, these methods will always return `false`.
