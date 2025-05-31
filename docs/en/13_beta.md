<font size=6>Experimental Features</font>

[Prev Page](./12_option.md) | [Contents](./README.md) | [Next Page](./14_1_14_new_feature.md)

---

There are some experimental features in jsonvalue. Import them by:

```go
import (
    "github.com/Andrew-M-C/go.jsonvalue/beta"
)
```

## Contains

From v1.3.0, I have provided an experimental function `Contains` to check whether a JSON value is a subset of another:

```go
func Contains(v *jsonvalue.V, sub any, inPath ...any) bool
```

This function returns whether `v` contains the subset `sub` at the path `inPath`. The `sub` parameter can be any type, just like the parameters of `New`.

This function will find sub-values using the `inPath` parameter (if no inPath is provided, it checks the value itself), then performs checks in the following sequence:

- If the types of the two values are different, it returns `false`.
- If neither is an array nor an object, it returns whether the two values are equal to each other.
- If array-typed, it returns whether `sub` is a sub-array of `v`.
- If object-typed, it iterates through every sub-value of both `v` and `sub` and checks:
  - If `sub` contains keys that `v` does not have, it returns `false`
  - For keys that both `v` and `sub` have but are not equal to each other, it invokes `Contains` on them recursively. Only if all recursive calls return `true` will the final result be `true`.

## Import/Export

These two functions were provided from v1.2.x, but they are now official in v1.3.x. Please use them directly in the jsonvalue package.
