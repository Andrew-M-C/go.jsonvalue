package jsonvalue

// ObjectIter is used in IterObjects function.
//
// ObjectIter 用于 IterObjects 函数。
type ObjectIter struct {
	K string
	V *V
}

// ArrayIter is used in IterArray function.
//
// ArrayIter 用于 IterArray 函数。
type ArrayIter struct {
	I int
	V *V
}

// RangeObjects goes through each children when this is an object value
//
// Return true in callback to continue range iteration, while false to break.
//
// 当当前 JSON 值是一个 object 类型时，RangeObjects 遍历所有的键值对。
//
// 在回调函数中返回 true 表示继续迭代，返回 false 表示退出迭代
func (v *V) RangeObjects(callback func(k string, v *V) bool) {
	if !v.IsObject() {
		return
	}
	if nil == callback {
		return
	}

	for k, v := range v.children.object {
		ok := callback(k, v)
		if !ok {
			break
		}
	}
}

// IterObjects returns a channel for range statement of object type JSON.
//
// 当当前 JSON 值是一个 object 类型时，IterObjects 返回一个可用于 range 操作符的 channel。
func (v *V) IterObjects() <-chan *ObjectIter {
	c := make(chan *ObjectIter, len(v.children.object))

	go func() {
		for k, v := range v.children.object {
			c <- &ObjectIter{
				K: k,
				V: v,
			}
		}
		close(c)
	}()
	return c
}

// RangeArray goes through each children when this is an array value
//
// Return true in callback to continue range iteration, while false to break.
//
// 当当前 JSON 值是一个 array 类型时，RangeArray 遍历所有的数组成员。
//
// 在回调函数中返回 true 表示继续迭代，返回 false 表示退出迭代
func (v *V) RangeArray(callback func(i int, v *V) bool) {
	if !v.IsArray() {
		return
	}
	if nil == callback {
		return
	}

	for i, child := range v.children.array {
		if ok := callback(i, child); !ok {
			break
		}
	}
}

// IterArray returns a channel for range statement of array type JSON.
//
// 当当前 JSON 值是一个 array 类型时，IterArray 返回一个可用于 range 操作符的 channel。
func (v *V) IterArray() <-chan *ArrayIter {
	c := make(chan *ArrayIter, len(v.children.array))

	go func() {
		for i, child := range v.children.array {
			c <- &ArrayIter{
				I: i,
				V: child,
			}
		}
		close(c)
	}()
	return c
}
