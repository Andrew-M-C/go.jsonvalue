package jsonvalue

// Deprecated: ObjectIter is a deprecated type.
type ObjectIter struct {
	K string
	V *V
}

// Deprecated: ArrayIter is a deprecated type.
type ArrayIter struct {
	I int
	V *V
}

// RangeObjects goes through each children when this is an object value
//
// Return true in callback to continue range iteration, while false to break.
//
// 若当前 JSON 值是一个 object 类型时，RangeObjects 遍历所有的键值对。
//
// 在回调函数中返回 true 表示继续迭代，返回 false 表示退出迭代
func (v *V) RangeObjects(callback func(k string, v *V) bool) {
	if v.impl == nil {
		return
	}
	v.impl.RangeObjects(callback)
}

// Deprecated: IterObjects is deprecated, please Use ForRangeObj() instead.
func (v *V) IterObjects() <-chan *ObjectIter {
	if v.impl == nil {
		ch := make(chan *ObjectIter)
		close(ch)
		return ch
	}
	return v.impl.IterObjects()
}

// ForRangeObj returns a map which can be used in for - range block to iteration KVs in a JSON object value.
//
// ForRangeObj 返回一个 map 类型，用于使用 for - range 块迭代 JSON 对象类型的子成员。
func (v *V) ForRangeObj() map[string]*V {
	if v.impl == nil {
		return map[string]*V{}
	}
	return v.impl.ForRangeObj()
}

// RangeArray goes through each children when this is an array value
//
// Return true in callback to continue range iteration, while false to break.
//
// 若当前 JSON 值是一个 array 类型时，RangeArray 遍历所有的数组成员。
//
// 在回调函数中返回 true 表示继续迭代，返回 false 表示退出迭代
func (v *V) RangeArray(callback func(i int, v *V) bool) {
	if v.impl == nil {
		return
	}
	v.impl.RangeArray(callback)
}

// Deprecated: IterArray is deprecated, please Use ForRangeArr() instead.
func (v *V) IterArray() <-chan *ArrayIter {
	if v.impl == nil {
		ch := make(chan *ArrayIter)
		close(ch)
		return ch
	}
	return v.impl.IterArray()
}

// ForRangeArr returns a slice which can be used in for - range block to iteration KVs in a JSON array value.
//
// ForRangeObj 返回一个切片，用于使用 for - range 块迭代 JSON 数组类型的子成员。
func (v *V) ForRangeArr() []*V {
	if v.impl == nil {
		return nil
	}
	return v.impl.ForRangeArr()
}
