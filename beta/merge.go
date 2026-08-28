package beta

import (
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// Merge merges the 'from' value into the 'to' value. It mutates to in place; when
// the root itself needs to be replaced (for example when to is null), *to is
// rewritten.
//
// Semantics:
//   - object: merge by key; missing keys in to are copied from from
//   - array: whole-array conflict; see array options below
//   - string / number / boolean: equal values are kept; unequal values conflict
//   - to is nil, null or NotExist: filled with a deep copy of from
//   - type mismatch: keep to and report the path
//   - from is nil: no-op
//
// Conflicting positions are returned as paths and are not overwritten. Nodes
// copied from from are deep-copied so later mutation of from does not affect to.
//
// Merge 将 from 合并进 to。该函数会原地修改 to；当根节点需要被替换时（例如 to 为
// null），会改写 *to。
//
// 语义:
//   - object: 按 key 合并, to 中缺失的 key 从 from 拷入
//   - array: 整段数组视为冲突, 详见下方的数组相关选项
//   - string / number / boolean: 值相等则保留, 不相等则视为冲突
//   - to 为 nil、null 或 NotExist: 用 from 的深拷贝填充
//   - 类型不匹配: 保留 to 并上报路径
//   - from 为 nil: 无操作
//
// 发生冲突的路径会返回, 且不会覆盖 to 中已有的值。从 from 拷入的节点会深拷贝,
// 之后修改 from 不会影响 to。
func Merge(to **jsonvalue.V, from *jsonvalue.V, opts ...MergeOption) (conflictPaths []jsonvalue.Path) {
	o := &mergeOptions{}
	for _, opt := range opts {
		opt(o)
	}
	// empty parameters are ignored
	if to == nil {
		return []jsonvalue.Path{nil} // contains a empty path, identifying root element
	}
	if from == nil {
		return nil
	}
	from = cloneV(from)
	merged, conflicts := merge(*to, from, nil, o)
	*to = merged
	return conflicts
}

func merge(
	to, from *jsonvalue.V,
	currentPath jsonvalue.Path,
	opts *mergeOptions,
) (merged *jsonvalue.V, conflictPaths []jsonvalue.Path) {
	// empty parameters are ignored
	if from == nil {
		return to, nil
	}
	// if 'to' is null, it will be replace to 'from' value
	if to == nil || to.ValueType() == jsonvalue.Null || to.ValueType() == jsonvalue.NotExist {
		return from, nil
	}
	// type mismatch
	if to.ValueType() != from.ValueType() {
		return resolveConflict(to, from, currentPath, opts) // an empty path identifies root element
	}

	concatString := opts.ConcatStrings && to.ValueType() == jsonvalue.String
	concatArray := opts.ConcatArrays && !opts.IgnoreArrays && to.ValueType() == jsonvalue.Array
	if !concatString && !concatArray && to.Equal(from) {
		return to, nil
	}

	// merge by type
	switch to.ValueType() {
	// string: concat
	case jsonvalue.String:
		if opts.ConcatStrings {
			return jsonvalue.NewString(to.String() + from.String()), nil
		}
		return resolveConflict(to, from, currentPath, opts)

	// Array: iteration, concat or ignore
	case jsonvalue.Array:
		return mergeArray(to, from, currentPath, opts)

	// object: merge
	case jsonvalue.Object:
		return mergeObject(to, from, currentPath, opts)

	// scalars: unequal values conflict
	case jsonvalue.Number, jsonvalue.Boolean:
		return resolveConflict(to, from, currentPath, opts)

	// other types will not be merged
	case jsonvalue.NotExist:
		fallthrough
	default:
		return to, []jsonvalue.Path{currentPath} // this should not occurs, which means the 'to' value is not initialized
	}
}

// resolveConflict decides which side wins a conflicting position. The conflict
// path is reported no matter which side wins.
//
// WithMergeOverrideWhenConflict() only applies to basic types. Object and array
// values are never taken from 'from' here, they are governed by the array
// options and by the key-by-key object merging instead.
func resolveConflict(
	to, from *jsonvalue.V, currentPath jsonvalue.Path, opts *mergeOptions,
) (merged *jsonvalue.V, conflictPaths []jsonvalue.Path) {
	if opts.OverrideWhenConflict && isBasicType(from.ValueType()) {
		return from, []jsonvalue.Path{currentPath}
	}
	return to, []jsonvalue.Path{currentPath}
}

func isBasicType(t jsonvalue.ValueType) bool {
	switch t {
	case jsonvalue.String, jsonvalue.Number, jsonvalue.Boolean, jsonvalue.Null:
		return true
	default:
		return false
	}
}

func mergeObject(
	to, from *jsonvalue.V, currentPath jsonvalue.Path, opts *mergeOptions,
) (merged *jsonvalue.V, conflictPaths []jsonvalue.Path) {
	from.RangeObjectsBySetSequence(func(k string, v *jsonvalue.V) bool {
		path := appendKey(currentPath, k)
		child, err := to.Get(k)
		if err != nil {
			to.MustSet(v).At(k)
			return true
		}
		newChild, res := merge(child, v, path, opts)
		if newChild != child {
			to.MustSet(newChild).At(k)
		}
		conflictPaths = append(conflictPaths, res...)
		return true
	})
	return to, conflictPaths
}

func mergeArray(
	to, from *jsonvalue.V, currentPath jsonvalue.Path, opts *mergeOptions,
) (merged *jsonvalue.V, conflictPaths []jsonvalue.Path) {
	if opts.IgnoreArrays {
		return to, nil
	}
	if opts.ConcatArrays {
		from.RangeArray(func(_ int, v *jsonvalue.V) bool {
			to.MustAppend(v).InTheEnd()
			return true
		})
		return to, nil
	}
	if opts.ArraysByIndex {
		return mergeArrayByIndex(to, from, currentPath, opts)
	}
	// other arrays will not be merged
	return to, []jsonvalue.Path{currentPath}
}

func mergeArrayByIndex(
	to, from *jsonvalue.V, currentPath jsonvalue.Path, opts *mergeOptions,
) (merged *jsonvalue.V, conflictPaths []jsonvalue.Path) {
	toLen := to.Len()
	fromLen := from.Len()
	n := toLen
	if fromLen < n {
		n = fromLen
	}

	// merge
	for i := 0; i < n; i++ {
		child := to.MustGet(i)
		fromChild := from.MustGet(i)
		path := appendIndex(currentPath, i)
		newChild, res := merge(child, fromChild, path, opts)
		if newChild != child {
			to.MustSet(newChild).At(i)
		}
		conflictPaths = append(conflictPaths, res...)
	}

	// simply append
	for i := toLen; i < fromLen; i++ {
		to.MustAppend(from.MustGet(i)).InTheEnd()
	}
	return to, conflictPaths
}

func appendKey(parent jsonvalue.Path, key string) jsonvalue.Path {
	res := make(jsonvalue.Path, len(parent)+1)
	copy(res, parent)
	res[len(parent)] = jsonvalue.PathItem{Idx: -1, Key: key}
	return res
}

func appendIndex(parent jsonvalue.Path, i int) jsonvalue.Path {
	res := make(jsonvalue.Path, len(parent)+1)
	copy(res, parent)
	res[len(parent)] = jsonvalue.PathItem{Idx: i}
	return res
}

func cloneV(v *jsonvalue.V) *jsonvalue.V {
	if v == nil {
		return nil
	}
	cloned, err := jsonvalue.Import(v)
	if err != nil || cloned == nil {
		return jsonvalue.MustUnmarshal(v.MustMarshal())
	}
	return cloned
}

type mergeOptions struct {
	ConcatStrings        bool
	ConcatArrays         bool
	IgnoreArrays         bool
	ArraysByIndex        bool
	OverrideWhenConflict bool // 发生冲突时是否采用 from 的值, 仅对基础类型生效
}

// MergeOption is an option for Merge() operation.
type MergeOption func(*mergeOptions)

// WithMergeConcatStrings specifies that the strings should be concatenated.
// Otherwise, string value will be ignored when 'to' has the same path key.
func WithMergeConcatStrings() MergeOption {
	return func(opts *mergeOptions) {
		opts.ConcatStrings = true
	}
}

// WithMergeConcatArrays specifies that the arrays should be concatenated instead
// if merge by each array element.
//
// This option will not take effect if WithMergeIgnoreArrays() is also specified.
func WithMergeConcatArrays() MergeOption {
	return func(opts *mergeOptions) {
		opts.ConcatArrays = true
	}
}

// WithMergeArraysByIndex merges arrays element by element. Extra elements in
// from are appended. Extra elements in to are kept.
//
// This option will not take effect if WithMergeIgnoreArrays() or
// WithMergeConcatArrays() is also specified.
//
// WithMergeArraysByIndex 按索引逐个合并数组元素。from 多出的元素会追加，to 多出
// 的元素会保留。若同时指定了 WithMergeIgnoreArrays() 或 WithMergeConcatArrays(),
// 则本选项不生效。
func WithMergeArraysByIndex() MergeOption {
	return func(opts *mergeOptions) {
		opts.ArraysByIndex = true
	}
}

// WithMergeIgnoreArrays specifies that the arrays should be ignored.
func WithMergeIgnoreArrays() MergeOption {
	return func(opts *mergeOptions) {
		opts.IgnoreArrays = true
	}
}

// WithMergeOverrideWhenConflict specifies that the value in from should
// override the one in to when they conflict. Please note that the conflict
// paths are still returned even though the values are overridden.
//
// This option only applies to basic types, which are string, number, boolean
// and null. Conflicts on objects and arrays are never overridden: objects are
// always merged key by key, while arrays are governed by
// WithMergeConcatArrays(), WithMergeIgnoreArrays() and
// WithMergeArraysByIndex(). WithMergeConcatStrings() also takes precedence over
// this option.
//
// WithMergeOverrideWhenConflict 指定发生冲突时采用 from 中的值覆盖 to。请注意,
// 即便值被覆盖, 冲突路径依然会被返回。
//
// 本选项仅对基础类型生效, 也即 string、number、boolean 和 null。object 和 array
// 的冲突不会被覆盖: object 始终按 key 逐个合并, 而 array 由
// WithMergeConcatArrays()、WithMergeIgnoreArrays() 和 WithMergeArraysByIndex()
// 控制。此外 WithMergeConcatStrings() 的优先级高于本选项。
func WithMergeOverrideWhenConflict() MergeOption {
	return func(opts *mergeOptions) {
		opts.OverrideWhenConflict = true
	}
}
