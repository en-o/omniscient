package autostart

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"os/exec"
	"strings"
	"time"
)

// IsAutostartInstalled 检查autostart命令是否安装
func IsAutostartInstalled() bool {
	cmd := exec.Command("which", "autostart")
	err := cmd.Run()
	if err != nil {
		// 也尝试检查常见安装路径
		paths := []string{
			"/usr/local/bin/autostart",
			"/usr/bin/autostart",
			"/opt/autostart/bin/autostart",
		}

		for _, path := range paths {
			if cmd := exec.Command("test", "-x", path); cmd.Run() == nil {
				return true
			}
		}
		return false
	}
	return true
}

// CheckAutostartServiceExists 检查autostart服务是否存在
func CheckAutostartServiceExists(ctx context.Context, autoName string) bool {
	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "autostart", "exists", autoName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		g.Log().Debug(ctx, "检查自启服务状态",
			"autoName", autoName, "error", err, "output", string(output))
		return false
	}

	// 检查输出中是否包含确认信息
	return strings.Contains(string(output), "Service '"+autoName+"' exists") ||
		strings.Contains(string(output), autoName)
}
