package jsonvalue

import (
	"sort"
)

// ---------------- array sorting ----------------

// ArrayLessFunc is used in SortArray()
type ArrayLessFunc func(v1, v2 *V) bool

// SortArray is used to re-arrange sequence of the array. Invokers should pass less function for sorting.
// Nothing would happens either lessFunc is nil or v is not an array.
func (v *V) SortArray(lessFunc ArrayLessFunc) {
	if nil == lessFunc {
		return
	}
	if false == v.IsArray() {
		return
	}

	sv := newSortV(v, lessFunc)
	sv.Sort()
	return
}

type sortV struct {
	v        *V
	lessFunc ArrayLessFunc
	children []*V
}

func newSortV(v *V, lessFunc ArrayLessFunc) *sortV {
	sv := sortV{
		v:        v,
		lessFunc: lessFunc,
		children: make([]*V, 0, v.Len()),
	}

	for e := v.arrayChildren.Front(); e != nil; e = e.Next() {
		sv.children = append(sv.children, e.Value.(*V))
	}
	return &sv
}

func (v *sortV) Sort() {
	// sort
	sort.Sort(v)

	// re-arrange children
	v.v.arrayChildren.Init()
	for _, child := range v.children {
		v.v.arrayChildren.PushBack(child)
	}
	return
}

func (v *sortV) Len() int {
	return len(v.children)
}

func (v *sortV) Less(i, j int) bool {
	v1 := v.children[i]
	v2 := v.children[j]
	return v.lessFunc(v1, v2)
}

func (v *sortV) Swap(i, j int) {
	v.children[i], v.children[j] = v.children[j], v.children[i]
}
