package utils

import "testing"

// go test -v -run Test_GetExternalIp ip_test.go ip.go
func Test_GetExternalIp(t *testing.T) {
	ip, err := GetExternalIp()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(ip)
}

// go test -v -run Test_GetIntranetIp ip_test.go ip.go
func Test_GetIntranetIp(t *testing.T) {
	intranetIp, err := GetIntranetIp()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(intranetIp)
}
