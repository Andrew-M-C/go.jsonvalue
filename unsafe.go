package jsonvalue

import (
	"reflect"
	"unsafe"
)

// unsafeBtoS []byte to string
func unsafeBtoS(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// unsafeStoB string to []byte
func unsafeStoB(s string) []byte {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sh.Cap = sh.Len
	return *(*[]byte)(unsafe.Pointer(sh))
}
