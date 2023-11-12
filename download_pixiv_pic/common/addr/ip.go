package addr

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func GetInterIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "GetInterIP err"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

func GetExternalIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "GetExternalIP err"
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return string(content)
}

func GetMacAddr() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "GetMacAddr err"
	}
	inter := interfaces[0]
	mac := inter.HardwareAddr.String() //获取本机MAC地址
	fmt.Println("MAC = ", mac)
	return mac
}
