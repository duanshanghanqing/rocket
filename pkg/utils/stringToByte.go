package utils

import (
	"reflect"
	"unsafe"
)

func StringToBytes(s string) []byte {
	// 1.Convert string header information to reflect StringHeader type to obtain the address and length of the underlying data in the string
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	// 2.Create a reflection SliceHeader type structure used to construct byte slice header information
	sliceHeader := &reflect.SliceHeader{
		Data: stringHeader.Data, // The address of the underlying data in the string
		Len:  stringHeader.Len,  // The length of a string
		Cap:  stringHeader.Len,  // The capacity of a byte slice, equal to the length of a string
	}

	// 3.Utilize unsafe Pointer will reflect SliceHeader is converted to a pointer of type [] byte,
	// Then use the * operator to dereference and obtain byte slices
	return *(*[]byte)(unsafe.Pointer(sliceHeader))
}

func BytesToString(b []byte) string {
	// 1.Convert byte slice header information to reflect SliceHeader type to obtain the address and length of byte slice underlying data
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	// 2.Create a reflection A structure of type StringHeader, used to construct string header information
	stringHeader := &reflect.StringHeader{
		Data: sliceHeader.Data, // Address of byte slice underlying data
		Len:  sliceHeader.Len,  // The length of byte slices
	}

	// 3.Utilize unsafe Pointer will reflect Convert StringHeader to a pointer of type string,
	// Then use the * operator to dereference and obtain the string
	return *(*string)(unsafe.Pointer(stringHeader))
}
