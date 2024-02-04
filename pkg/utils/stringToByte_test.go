package utils

import (
	"testing"
)

// go test -v -run Test_StringToBytes stringToByte_test.go stringToByte.go
func Test_StringToBytes(t *testing.T) {
	var str = "hello"
	// []byte(str) // Although simple, this method has low performance in large string scenarios
	var bytes = StringToBytes(str)
	t.Log(bytes)
}

// go test -v -run Test_BytesToString stringToByte_test.go stringToByte.go
func Test_BytesToString(t *testing.T) {
	var bytes = []byte{104, 101, 108, 108, 111}
	//string(bytes) // Although simple, this method has low performance in large string scenarios
	var str = BytesToString(bytes)
	t.Log(str)
}
