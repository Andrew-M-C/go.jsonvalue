
<font size=6>Value Comparing</font>

[Prev Page](./10_scenarios.md) | [Contents](./README.md) | [Next Page](./12_option.md)

---

From v 1.3.0, `Equal` method is added to tell whether two JSON values are the same. The comparing rules are:

Firstly, if the types of two values are different, returns `false`

If types are the same, check them by types as follow:

- `string`: check if two string values are the same.
- `number`: check if the DECIMAL numbers of two values are equal. I use Package [decimal](https://pkg.go.dev/github.com/shopspring/decimal) to do this.
- `boolean`: check if the two boolean values are the same.
- `null`: all null values equal to each other.
- `array`: the necessary and sufficient condition of that two array equal to each other are: they shares the same array length, and each elements equal to the other one in the same index.
- `object`: the necessary and sufficient condition of that two array equal to each other are: they shares the exact same keys, and each elements equal to the other one with same key.
