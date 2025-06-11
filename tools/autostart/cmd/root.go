package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"autostart/internal/service"
	"autostart/internal/utils"
)

const (
	ToolName = "autostart"
	Version  = "0.0.2"
)

// commandHandler 定义命令处理函数类型
type commandHandler func()

// Execute 是应用程序的主入口点
func Execute() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run 执行主要的应用程序逻辑
func run() error {
	// 检查操作系统支持
	if err := checkSystemSupport(); err != nil {
		return err
	}

	// 解析命令行参数
	if len(os.Args) < 2 {
		utils.PrintHelp()
		return nil
	}

	command := strings.ToLower(os.Args[1])

	// 检查权限要求
	if err := checkPermissions(command); err != nil {
		return err
	}

	// 执行命令
	return executeCommand(command)
}

// checkSystemSupport 检查操作系统支持
func checkSystemSupport() error {
	currentOS := runtime.GOOS
	fmt.Printf("System: %s\n", currentOS)

	if currentOS != "linux" {
		return fmt.Errorf("this tool currently only supports Linux systems (current: %s)", currentOS)
	}

	return nil
}

// checkPermissions 检查命令执行所需的权限
func checkPermissions(command string) error {
	if !utils.NeedsRoot(command) {
		return nil
	}

	if os.Geteuid() != 0 {
		fmt.Printf("This operation requires root privileges. Please run with sudo:\n")
		fmt.Printf("  sudo %s %s\n", os.Args[0], strings.Join(os.Args[1:], " "))
		return fmt.Errorf("insufficient privileges")
	}

	return nil
}

// executeCommand 根据命令执行相应的操作
func executeCommand(command string) error {
	sm := service.NewServiceManager()

	// 命令映射表
	commandMap := map[string]commandHandler{
		// 查询命令
		"list": func() { handleError(sm.ListAutostartServices()) },
		"ls":   func() { handleError(sm.ListAutostartServices()) },

		// 服务管理命令
		"add":       func() { handleServiceAdd(sm) },
		"create":    func() { handleServiceAdd(sm) },
		"install":   func() { handleServiceAdd(sm) },
		"remove":    func() { handleServiceRemove(sm) },
		"rm":        func() { handleServiceRemove(sm) },
		"delete":    func() { handleServiceRemove(sm) },
		"uninstall": func() { handleServiceRemove(sm) },
		"edit":      func() { handleServiceEdit(sm) },

		// 服务控制命令
		"enable":  func() { handleServiceEnable(sm) },
		"disable": func() { handleServiceDisable(sm) },
		"start":   func() { handleServiceStart(sm) },
		"stop":    func() { handleServiceStop(sm) },
		"restart": func() { handleServiceRestart(sm) },
		"status":  func() { handleServiceStatus(sm) },
		"logs":    func() { handleServiceLogs(sm) },

		// 工具命令
		"version":   func() { printVersion() },
		"-v":        func() { printVersion() },
		"--version": func() { printVersion() },
		"help":      func() { utils.PrintHelp() },
		"-h":        func() { utils.PrintHelp() },
		"--help":    func() { utils.PrintHelp() },
	}

	handler, exists := commandMap[command]
	if !exists {
		fmt.Printf("Unknown command: %s\n\n", command)
		utils.PrintHelp()
		return fmt.Errorf("unknown command: %s", command)
	}

	handler()
	return nil
}

// handleServiceAdd 处理服务添加命令
func handleServiceAdd(sm *service.ServiceManager) {
	if len(os.Args) < 4 {
		fmt.Println("Usage: autostart add <name> <exec-start> [options...]")
		fmt.Println("Example: autostart add myapp \"java -jar /path/to/app.jar\" --workdir=/path/to")
		handleError(fmt.Errorf("insufficient arguments for add command"))
		return
	}

	serviceName := os.Args[2]
	execStart := os.Args[3]
	options := os.Args[4:]

	handleError(sm.AddAutostartService(serviceName, execStart, options))
}

// handleServiceRemove 处理服务移除命令
func handleServiceRemove(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("remove")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.RemoveAutostartService(serviceName))
}

// handleServiceEdit 处理服务编辑命令
func handleServiceEdit(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("edit")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.EditService(serviceName))
}

// handleServiceEnable 处理服务启用命令
func handleServiceEnable(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("enable")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.EnableService(serviceName))
}

// handleServiceDisable 处理服务禁用命令
func handleServiceDisable(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("disable")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.DisableService(serviceName))
}

// handleServiceStart 处理服务启动命令
func handleServiceStart(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("start")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.StartService(serviceName))
}

// handleServiceStop 处理服务停止命令
func handleServiceStop(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("stop")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.StopService(serviceName))
}

// handleServiceRestart 处理服务重启命令
func handleServiceRestart(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("restart")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.RestartService(serviceName))
}

// handleServiceStatus 处理服务状态查询命令
func handleServiceStatus(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("status")
	if err != nil {
		handleError(err)
		return
	}
	handleError(sm.ShowServiceStatus(serviceName))
}

// handleServiceLogs 处理服务日志查询命令
func handleServiceLogs(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("logs")
	if err != nil {
		handleError(err)
		return
	}

	lines := ""
	if len(os.Args) > 3 {
		lines = os.Args[3]
	}

	handleError(sm.ShowServiceLogs(serviceName, lines))
}

// getServiceNameArg 获取服务名称参数
func getServiceNameArg(command string) (string, error) {
	if len(os.Args) < 3 {
		return "", fmt.Errorf("missing service name for %s command\nUsage: autostart %s <service-name>", command, command)
	}
	return os.Args[2], nil
}

// printVersion 打印版本信息
func printVersion() {
	fmt.Printf("%s v%s\n", ToolName, Version)
	fmt.Printf("Built for %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go version: %s\n", runtime.Version())
}

// handleError 统一的错误处理
func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
