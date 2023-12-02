package utils

import (
	"fmt"
	"net"
	"strconv"
)

func AddressToHostPost(address string) (host string, port int, err error) {
	host, stringPost, err := net.SplitHostPort(address)
	if err != nil {
		return host, port, err
	}
	post, err := strconv.Atoi(stringPost)
	if err != nil {
		return host, port, err
	}

	return host, post, nil
}

func HostPostToAddress(host string, post int) string {
	return fmt.Sprintf("%s:%d", host, post)
}
