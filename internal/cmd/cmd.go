package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcfg"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"omniscient/internal/controller/jpid"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"omniscient/internal/controller/hello"
)

const (
	ServiceName       = "omniscient"
	DefaultConfigFile = "config.prod.yaml"
)

var (
	// 主命令
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return runServer(ctx)
		},
	}

	// run 命令 - 运行当前程序
	Run = gcmd.Command{
		Name:  "run",
		Usage: "run",
		Brief: "run the http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return runServer(ctx)
		},
	}

	// sh 命令 - 服务管理命令行功能
	Shell = gcmd.Command{
		Name:  "sh",
		Usage: "sh [sub-command]",
		Brief: "service management shell commands",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return handleShellCommand(ctx)
		},
	}
)

// 运行服务器
func runServer(ctx context.Context) error {
	// 检查当前目录下是否存在 config.prod.yaml
	workDir, _ := os.Getwd()
	configPath := filepath.Join(workDir, DefaultConfigFile)

	// 获取命令行配置的文件路径
	configFile := g.Cfg().GetAdapter().(*gcfg.AdapterFile).GetFileName()
	if configFile != "" && gfile.Exists(configFile) {
		g.Log().Infof(ctx, "Using command line config file: %s", configFile)
	} else if gfile.Exists(configPath) {
		g.Log().Infof(ctx, "Using external config file: %s", configPath)
		// 设置默认配置文件路径
		g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(configPath)
	} else {
		g.Log().Info(ctx, "Using built-in config")
	}

	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(
			hello.NewV1(),
			jpid.NewV1(),
		)
	})
	// 绑定静态资源
	s.SetServerRoot("resource/public")
	s.Run()
	return nil
}

// 处理 shell 命令
func handleShellCommand(ctx context.Context) error {
	// 直接从 os.Args 获取参数，这是最可靠的方法
	osArgs := os.Args
	var args []string

	// 找到 "sh" 命令的位置，获取后面的参数
	for i, arg := range osArgs {
		if arg == "sh" && i+1 < len(osArgs) {
			args = osArgs[i+1:]
			break
		}
	}

	if len(args) == 0 {
		printShellHelp()
		return nil
	}

	switch args[0] {
	case "status":
		return showServiceStatus()
	case "enable":
		return enableService()
	case "disable":
		return disableService()
	case "start":
		return startService()
	case "stop":
		return stopService()
	case "restart":
		return restartService()
	case "reload":
		return reloadService()
	case "install":
		return installService()
	case "uninstall":
		return uninstallService()
	case "config":
		if len(args) > 1 {
			return setDefaultConfig(args[1])
		}
		return showCurrentConfig()
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		printShellHelp()
		return nil
	}
}

// 打印帮助信息
func printShellHelp() {
	fmt.Println("Service Management Commands (requires root privileges):")
	fmt.Println("Usage: sudo omniscient sh <command>")
	fmt.Println("  status     - Show service status")
	fmt.Println("  enable     - Enable service auto-start")
	fmt.Println("  disable    - Disable service auto-start")
	fmt.Println("  start      - Start service")
	fmt.Println("  stop       - Stop service")
	fmt.Println("  restart    - Restart service")
	fmt.Println("  reload     - Reload service")
	fmt.Println("  install    - Install systemd service")
	fmt.Println("  uninstall  - Uninstall systemd service")
	fmt.Println("  config [file] - Set default config file or show current config")
}

// 显示服务状态
func showServiceStatus() error {
	cmd := exec.Command("systemctl", "status", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Service %s is not installed or not running\n", ServiceName)
	}
	return nil
}

// 启用开机自启
func enableService() error {
	cmd := exec.Command("sudo", "systemctl", "enable", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to enable service: %v", err)
	}
	fmt.Printf("Service %s enabled for auto-start\n", ServiceName)
	return nil
}

// 关闭开机自启
func disableService() error {
	cmd := exec.Command("sudo", "systemctl", "disable", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to disable service: %v", err)
	}
	fmt.Printf("Service %s disabled from auto-start\n", ServiceName)
	return nil
}

// 启动服务
func startService() error {
	cmd := exec.Command("sudo", "systemctl", "start", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}
	fmt.Printf("Service %s started\n", ServiceName)
	return nil
}

// 停止服务
func stopService() error {
	cmd := exec.Command("sudo", "systemctl", "stop", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stop service: %v", err)
	}
	fmt.Printf("Service %s stopped\n", ServiceName)
	return nil
}

// 重启服务
func restartService() error {
	cmd := exec.Command("sudo", "systemctl", "restart", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to restart service: %v", err)
	}
	fmt.Printf("Service %s restarted\n", ServiceName)
	return nil
}

// 重载服务
func reloadService() error {
	cmd := exec.Command("sudo", "systemctl", "reload-or-restart", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reload service: %v", err)
	}
	fmt.Printf("Service %s reloaded\n", ServiceName)
	return nil
}

// 安装 systemd 服务
func installService() error {
	fmt.Printf("installService starting ================\n")

	// 检查是否以 root 身份运行
	if os.Geteuid() != 0 {
		return fmt.Errorf("installing service requires root privileges. Please run with sudo:\n  sudo ./omniscient sh install")
	}

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	execPath, _ = filepath.Abs(execPath)
	workDir := filepath.Dir(execPath)

	// 检查配置文件是否存在
	configPath := filepath.Join(workDir, DefaultConfigFile)
	if !gfile.Exists(configPath) {
		fmt.Printf("Warning: Default config file %s not found in %s\n", DefaultConfigFile, workDir)
		fmt.Println("Service will be installed but may not start properly without config file")
	}

	// 创建 systemd 服务文件内容
	serviceContent := fmt.Sprintf(`[Unit]
Description=Omniscient HTTP Server
After=network.target
Wants=network.target

[Service]
Type=simple
User=root
WorkingDirectory=%s
ExecStart=%s run --gf.gcfg.file=%s
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=%s

[Install]
WantedBy=multi-user.target
`, workDir, execPath, configPath, ServiceName)

	// 写入服务文件
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", ServiceName)

	// 使用标准库写入文件，更好的错误处理
	file, err := os.Create(servicePath)
	if err != nil {
		return fmt.Errorf("failed to create service file %s: %v\nMake sure you have write permissions to /etc/systemd/system/", servicePath, err)
	}
	defer file.Close()

	_, err = file.WriteString(serviceContent)
	if err != nil {
		return fmt.Errorf("failed to write service file content: %v", err)
	}

	// 设置文件权限
	err = os.Chmod(servicePath, 0644)
	if err != nil {
		fmt.Printf("Warning: failed to set file permissions: %v\n", err)
	}

	// 重载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	fmt.Printf("Service %s installed successfully\n", ServiceName)
	fmt.Printf("Service file: %s\n", servicePath)
	fmt.Printf("Executable: %s\n", execPath)
	fmt.Printf("Config file: %s\n", configPath)
	fmt.Printf("Use 'omniscient sh enable' to enable auto-start\n")
	fmt.Printf("Use 'omniscient sh start' to start the service\n")

	return nil
}

// 卸载 systemd 服务
func uninstallService() error {
	// 检查是否以 root 身份运行
	if os.Geteuid() != 0 {
		return fmt.Errorf("uninstalling service requires root privileges. Please run with sudo:\n  sudo ./omniscient sh uninstall")
	}

	// 先停止服务
	stopService()

	// 禁用服务
	disableService()

	// 删除服务文件
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", ServiceName)
	err := os.Remove(servicePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service file %s: %v", servicePath, err)
	}

	// 重载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	if isRedHatBased() {
		fmt.Println("Note: On CentOS/RHEL systems, you need to explicitly enable the service:")
		fmt.Printf("Use 'sudo systemctl enable %s' to enable auto-start\n", ServiceName)
	}
	fmt.Printf("Service %s uninstalled successfully\n", ServiceName)
	return nil
}

// 判断是否是 RedHat 系列系统
func isRedHatBased() bool {
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return true
	}
	return false
}

// 设置默认配置文件
func setDefaultConfig(configFile string) error {
	if !gfile.Exists(configFile) {
		return fmt.Errorf("config file %s does not exist", configFile)
	}

	// 这里可以实现配置文件路径的持久化存储
	// 例如写入到 ~/.omniscient/config 或者其他位置
	configDir := filepath.Join(os.Getenv("HOME"), ".omniscient")
	if !gfile.Exists(configDir) {
		err := gfile.Mkdir(configDir)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	configPath := filepath.Join(configDir, "default_config")
	err := gfile.PutContents(configPath, configFile)
	if err != nil {
		return fmt.Errorf("failed to save default config: %v", err)
	}

	fmt.Printf("Default config file set to: %s\n", configFile)
	return nil
}

// 显示当前配置
func showCurrentConfig() error {
	configDir := filepath.Join(os.Getenv("HOME"), ".omniscient")
	configPath := filepath.Join(configDir, "default_config")

	if gfile.Exists(configPath) {
		content := gfile.GetContents(configPath)
		fmt.Printf("Current default config file: %s\n", strings.TrimSpace(content))
	} else {
		fmt.Printf("No default config file set, using: %s\n", DefaultConfigFile)
	}

	return nil
}

// 获取进程PID（基于端口）
func getProcessByPort(port string) (int, error) {
	cmd := exec.Command("ss", "-tlnp")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ":"+port) {
			parts := strings.Fields(line)
			if len(parts) >= 6 {
				pidInfo := parts[5]
				if strings.Contains(pidInfo, "pid=") {
					pidStr := gstr.SubStr(pidInfo, gstr.Pos(pidInfo, "pid=")+4)
					if gstr.Pos(pidStr, ",") > 0 {
						pidStr = gstr.SubStr(pidStr, 0, gstr.Pos(pidStr, ","))
					}
					pid, err := strconv.Atoi(pidStr)
					if err != nil {
						return 0, err
					}
					return pid, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("no process found on port %s", port)
}

// 终止进程
func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	// 先尝试优雅关闭
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		// 如果优雅关闭失败，强制杀死
		return process.Signal(syscall.SIGKILL)
	}

	// 等待进程退出
	time.Sleep(2 * time.Second)

	// 检查进程是否还存在
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		// 进程仍然存在，强制杀死
		return process.Signal(syscall.SIGKILL)
	}

	return nil
}
