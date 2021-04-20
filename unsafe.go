package jsonvalue

import (
	"reflect"
	"unsafe"
)

func unsafeStoB(s string) []byte {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sh.Cap = sh.Len
	return *(*[]byte)(unsafe.Pointer(sh))
}

// func unsafeBtoS(b []byte) string {
// 	return *(*string)(unsafe.Pointer(&b))
// }
