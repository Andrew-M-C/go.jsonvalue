package jsonvalue

// ObjectIter is used in IterObjects function
type ObjectIter struct {
	K string
	V *V
}

// ArrayIter is used in IterArray function
type ArrayIter struct {
	I int
	V *V
}

// RangeObjects goes through each children when this is an object value
//
// Return true in callback to continue range iteration, while false to break.
func (v *V) RangeObjects(callback func(k string, v *V) bool) {
	if false == v.IsObject() {
		return
	}
	if nil == callback {
		return
	}

	for k, v := range v.children.object {
		ok := callback(k, v)
		if false == ok {
			break
		}
	}
	return
}

// IterObjects returns a channel for range statement of object type JSON.
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
func (v *V) RangeArray(callback func(i int, v *V) bool) {
	if false == v.IsArray() {
		return
	}
	if nil == callback {
		return
	}

	i := 0
	for e := v.children.array.Front(); e != nil; e = e.Next() {
		v := e.Value.(*V)
		ok := callback(i, v)
		if false == ok {
			break
		}
		i++
	}
}

// IterArray returns a channel for range statement of array type JSON.
func (v *V) IterArray() <-chan *ArrayIter {
	c := make(chan *ArrayIter, v.children.array.Len())

	go func() {
		i := 0
		for e := v.children.array.Front(); e != nil; e = e.Next() {
			v := e.Value.(*V)
			c <- &ArrayIter{
				I: i,
				V: v,
			}
			i++
		}
		close(c)
	}()
	return c
}
