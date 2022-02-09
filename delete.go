package jsonvalue

// Delete deletes specified JSON value. Forexample, parameters ("data", "list") identifies deleting value in data.list.
// While ("list", 1) means deleting 2nd (count from one) element from the "list" array.
//
// Delete 从 JSON 中删除参数指定的对象。比如参数 ("data", "list") 表示删除 data.list 值；参数 ("list", 1) 则表示删除 list
// 数组的第2（从1算起）个值。
func (v *V) Delete(firstParam interface{}, otherParams ...interface{}) error {
	if v.impl == nil {
		return ErrValueUninitialized
	}
	return v.impl.delete(false, firstParam, otherParams...)
}
