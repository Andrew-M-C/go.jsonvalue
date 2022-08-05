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
	if !v.IsObject() {
		return
	}
	if nil == callback {
		return
	}

	for k, c := range v.children.object {
		ok := callback(k, c.v)
		if !ok {
			break
		}
	}
}

// Deprecated: IterObjects is deprecated, please Use ForRangeObj() instead.
func (v *V) IterObjects() <-chan *ObjectIter {
	ch := make(chan *ObjectIter, len(v.children.object))

	go func() {
		for k, c := range v.children.object {
			ch <- &ObjectIter{
				K: k,
				V: c.v,
			}
		}
		close(ch)
	}()
	return ch
}

// ForRangeObj returns a map which can be used in for - range block to iteration KVs in a JSON object value.
//
// ForRangeObj 返回一个 map 类型，用于使用 for - range 块迭代 JSON 对象类型的子成员。
func (v *V) ForRangeObj() map[string]*V {
	res := make(map[string]*V, len(v.children.object))
	for k, c := range v.children.object {
		res[k] = c.v
	}
	return res
}

// RangeArray goes through each children when this is an array value
//
// Return true in callback to continue range iteration, while false to break.
//
// 若当前 JSON 值是一个 array 类型时，RangeArray 遍历所有的数组成员。
//
// 在回调函数中返回 true 表示继续迭代，返回 false 表示退出迭代
func (v *V) RangeArray(callback func(i int, v *V) bool) {
	if !v.IsArray() {
		return
	}
	if nil == callback {
		return
	}

	for i, child := range v.children.arr {
		if ok := callback(i, child); !ok {
			break
		}
	}
}

// Deprecated: IterArray is deprecated, please Use ForRangeArr() instead.
func (v *V) IterArray() <-chan *ArrayIter {
	c := make(chan *ArrayIter, len(v.children.arr))

	go func() {
		for i, child := range v.children.arr {
			c <- &ArrayIter{
				I: i,
				V: child,
			}
		}
		close(c)
	}()
	return c
}

// ForRangeArr returns a slice which can be used in for - range block to iteration KVs in a JSON array value.
//
// ForRangeObj 返回一个切片，用于使用 for - range 块迭代 JSON 数组类型的子成员。
func (v *V) ForRangeArr() []*V {
	res := make([]*V, 0, len(v.children.arr))
	return append(res, v.children.arr...)
}
