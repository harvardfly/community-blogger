package netutil

import "net"

// GetAvailablePort 获取有效的端口
func GetAvailablePort() int {
	l, _ := net.Listen("tcp", ":0")
	defer func() {
		_ = l.Close()
	}()
	port := l.Addr().(*net.TCPAddr).Port
	return port
}
