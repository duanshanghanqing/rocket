package utils

import (
	"fmt"
	"net"
	"strconv"
)

func AddressToHostPost(address string) (string, int, error) {
	host, stringPost, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, err
	}
	post, err := strconv.Atoi(stringPost)
	if err != nil {
		return "", 0, err
	}

	return host, post, nil
}

func HostPostToAddress(host string, post int) string {
	return fmt.Sprintf("%s:%d", host, post)
}

// GetFreePort get free port
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
