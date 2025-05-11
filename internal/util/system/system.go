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
	interfaces, err := net.Interfaces()
	if err != nil {
		g.Log().Error(context.TODO(), "获取网络接口失败:", err)
		return "0"
	}

	for _, iface := range interfaces {
		// 过滤掉down和loopback接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址
		addrs, err := iface.Addrs()
		if err != nil {
			g.Log().Error(context.TODO(), "获取接口地址失败:", err)
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			ip := ipnet.IP.To4()
			if ip == nil {
				continue // 跳过IPv6地址
			}

			ipStr := ip.String()
			g.Log().Info(context.TODO(), "接口名称:", iface.Name)
			g.Log().Info(context.TODO(), "IP地址:", ipStr)

			// 排除特殊IP范围（如docker网桥等）
			if !strings.HasPrefix(ipStr, "172.") && !strings.HasPrefix(ipStr, "10.") && !strings.HasPrefix(ipStr, "192.168.") {
				continue
			}

			parts := strings.Split(ipStr, ".")
			if len(parts) == 4 {
				return parts[3]
			}
		}
	}

	g.Log().Warning(context.TODO(), "未找到有效的IPv4地址")
	return "0"
}
