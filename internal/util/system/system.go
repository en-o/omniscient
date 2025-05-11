package system

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/context"
	"net"
	"os"
	"strings"
)

// GetWorkerName 获取worker名称（服务器名-IP最后一段）
func GetWorkerName() string {
	// 获取服务器主机名
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// 获取本机IP地址
	lastSegment := getLastIPSegment()

	// 构造worker标识：主机名-IP最后一段
	return fmt.Sprintf("%s-%s", hostname, lastSegment)
}

// getLastIPSegment 获取IP地址的最后一段
func getLastIPSegment() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		g.Log().Error(context.TODO(), "获取网络接口地址失败:", err)
		return "0"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				g.Log().Info(context.TODO(), "当前服务器IP:", ip)
				parts := strings.Split(ip, ".")
				if len(parts) == 4 {
					return parts[3]
				}
			}
		}
	}
	return "0"
}
