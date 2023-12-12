package utils

import "testing"

// go test -v -run Test_AddressToHostPost address_test.go address.go
func Test_AddressToHostPost(t *testing.T) {
	host, port, err := AddressToHostPost("192.168.0.50:8090")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(host, port)
}

// go test -v -run Test_HostPostToAddress address_test.go address.go
func Test_HostPostToAddress(t *testing.T) {
	address := HostPostToAddress("192.168.0.50", 8090)
	t.Log(address)
}

// go test -v -run Test_GetFreePort address_test.go address.go
func Test_GetFreePort(t *testing.T) {
	port, err := GetFreePort()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(port)
}
