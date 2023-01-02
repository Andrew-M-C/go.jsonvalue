
<font size=6>Experimental Feature</font>

[Prev Page](./12_option.md) | [Contents](./README.md) | [Next Page](./14_1_13_new_feature.md)

---

There are some experimental features in jsonvalue. Import them by:

```go
import (
    "github.com/Andrew-M-C/go.jsonvalue/beta"
)
```

## Contains

From v1.3.0, I provided an experimental function `Contains`, to check whether a JSON value is the subset of another:

```go
func Contains(v *jsonvalue.V, sub any, inPath ...any) bool
```

This function returns whether `v` contains subset `sub`, in path `inPath`. sub can be any type, just like the parameters of `New` do.

This function will find sub values by `inPath` parameter (if no inPath, then check itself), then check by following sequence:

- If the types of two values are different, returns `false`.
- If not array or object, returns whether the two values equal to each other.
- If array typed, returns whether `sub` is the sub-array of `v`.
- If object typed, then iterate every sub values of both `v` and `sub` and check:
  - If `sub` contains keys those `v` does not have, returns `false`
  - For those keys both `v` and `sub` have but not equal to each other, then invoke `Contains` to them recursively. Only if all recursion return `true`, the final `true` will be returned.

## Import/Export

These two functions are provided from v1.2.x, but they are now official in v1.3.x. Please use them directly in jsonvalue package.
