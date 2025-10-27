package autostart

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/internal/util/system"
	"os/exec"
	"strings"
	"time"
)

// 调用 autostart命令行工具进行自启动服务管理

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

// VerifyAutostartService 验证自启动服务状态
func VerifyAutostartService(ctx context.Context, autoName string) error {
	cmd := exec.Command("autostart", "ls")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return gerror.Wrap(err, "获取自启动服务列表失败")
	}

	g.Log().Info(ctx, "当前自启动服务列表", "output", string(output))

	// 检查服务是否在列表中
	if !strings.Contains(string(output), autoName) {
		return gerror.New("服务未出现在自启动列表中")
	}

	return nil
}

// RemoveAutostartServiceNonInteractive 非交互式移除自启服务
func RemoveAutostartServiceNonInteractive(ctx context.Context, autoName string) error {
	// 方法1: 使用echo 'y'通过管道
	cmdCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "bash", "-c", "echo 'y' | autostart rm "+autoName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		g.Log().Warning(ctx, "方法1移除自启服务失败，尝试方法2",
			"error", err, "output", string(output))

		// 方法2: 直接禁用服务
		err3 := system.ExecCommand(ctx, "autostart", "disable", autoName)
		if err3 != nil {
			return gerror.Wrapf(err, "所有移除方法都失败了，原始错误: %v", err)
		}
	}

	return nil
}
