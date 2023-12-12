package utils

import (
	"io"
	"net"
	"net/http"
	"strings"
)

func getPublicIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}

func getExternalIp() (string, error) {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetExternalIp() (string, error) {
	externalIp, err1 := getExternalIp()
	if err1 != nil {
		publicIP, err2 := getPublicIP()
		if err2 != nil {
			return "", err2
		}
		return publicIP, nil
	}
	return externalIp, nil
}

func GetIntranetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// Check the IP address to determine whether to loop back the address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println("ip:", ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}
	return "127.0.0.1", nil
}
