package jsonvalue

import "unsafe"

func unsafeBtoS(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
