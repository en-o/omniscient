package system

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/context"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
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
			//g.Log().Info(context.TODO(), "接口名称:", iface.Name)
			//g.Log().Info(context.TODO(), "IP地址:", ipStr)

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

// CheckSudoNoPassword 检查是否可以无密码执行sudo命令
func CheckSudoNoPassword(command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", "-n", command, "--help")
	if err := cmd.Run(); err != nil {
		return gerror.New("无法无密码执行sudo命令，请配置sudo免密:\n" +
			"1. 执行: sudo visudo\n" +
			"2. 添加: " + GetCurrentUser() + " ALL=(ALL) NOPASSWD: /*/autostart\n" +
			"3. 保存并退出")
	}
	return nil
}

// GetCurrentUser 获取当前用户名
func GetCurrentUser() string {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return "username" // 默认值
	}
	return strings.TrimSpace(string(output))
}

// ExecCommand 执行命令，带超时和错误处理
func ExecCommand(ctx context.Context, args ...string) error {
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return gerror.Wrapf(err, "执行命令失败: %s, 输出: %s",
			strings.Join(args, " "), string(output))
	}
	return nil
}

// ExecSudoCommand 执行sudo命令，带超时和错误处理
func ExecSudoCommand(ctx context.Context, args ...string) error {
	// 设置命令超时
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 构建完整的sudo命令
	fullArgs := append([]string{"sudo", "-n"}, args...)
	cmd := exec.CommandContext(cmdCtx, fullArgs[0], fullArgs[1:]...)

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return gerror.Wrapf(err, "执行sudo命令失败: %s, 输出: %s", strings.Join(args, " "), string(output))
	}

	return nil
}
