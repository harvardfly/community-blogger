package netutil

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/peer"
	"net"
	"net/http"
	"strings"
)

// GetLocalIP4 获取本地的ipv4地址
func GetLocalIP4() (ip string) {
	interfaces, err := net.Interfaces()
	_, _ = net.InterfaceAddrs()
	if err != nil {
		return
	}
	if len(interfaces) == 2 {
		for _, face := range interfaces {
			if strings.Contains(face.Name, "lo") {
				continue
			}
			as, err := face.Addrs()
			if err != nil {
				return
			}
			for _, addr := range as {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						currIP := ipNet.IP.String()
						if !strings.Contains(currIP, ":") && currIP != "127.0.0.1" {
							ip = currIP
						}
					}
				}
			}
		}
	}
	for _, face := range interfaces {
		if strings.Contains(face.Name, "lo") {
			continue
		}
		as, err := face.Addrs()
		if err != nil {
			return
		}
		for _, addr := range as {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					currIP := ipNet.IP.String()
					if !strings.Contains(currIP, ":") && currIP != "127.0.0.1" && isIntranetIpv4(currIP) {
						ip = currIP
					}
				}
			}
		}
	}
	return
}

// isIntranetIpv4 判断是否内网ipv4地址
func isIntranetIpv4(ip string) bool {
	if strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "169.254.") ||
		strings.HasPrefix(ip, "172.") ||
		strings.HasPrefix(ip, "10.30.") ||
		strings.HasPrefix(ip, "10.31.") {
		return true
	}
	return false
}

// GetClientIP4 获取client的ipv4地址
func GetClientIP4(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}

// RemoteIP 返回远程客户端的 IP，如 192.168.1.1
func RemoteIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// ExternalIP 外部IP
func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := GetIPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

// GetIPFromAddr 获取地址的IP
func GetIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}

// GetIP 获取IP
func GetIP(w http.ResponseWriter, r *http.Request) string {
	// 尝试从 X-Forwarded-For 中获取
	xForwardedFor := r.Header.Get(`X-Forwarded-For`)
	ip := strings.TrimSpace(strings.Split(xForwardedFor, `,`)[0])
	if ip == `` {
		// 尝试从 X-Real-Ip 中获取
		ip = strings.TrimSpace(r.Header.Get(`X-Real-Ip`))
		if ip == `` {
			// 直接从 Remote Addr 中获取
			_ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
			if err != nil {
				panic(err)
			} else {
				ip = _ip
			}
		}
	}
	return ip
}
