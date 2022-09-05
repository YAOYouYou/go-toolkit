package string

import (
	"reflect"
	"unsafe"
)

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func String2Bytes(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}
