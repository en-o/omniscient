package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	gi := NewGlobalInstaller()

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
		"version":          func() { printVersion() },
		"-v":               func() { printVersion() },
		"--version":        func() { printVersion() },
		"help":             func() { utils.PrintHelp() },
		"-h":               func() { utils.PrintHelp() },
		"--help":           func() { utils.PrintHelp() },
		"install-global":   func() { handleInstallGlobal(gi) },
		"uninstall-global": func() { handleUninstallGlobal(gi) }, // 新增
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

// handleInstallGlobal 处理全局安装命令
func handleInstallGlobal(gi *GlobalInstaller) {
	fmt.Println("Installing autostart to global environment...")

	// 检查是否已经从全局位置运行
	if isGlobal, globalPath := gi.isRunningFromGlobal(); isGlobal {
		fmt.Printf("✓ Already installed globally at: %s\n", globalPath)
		fmt.Println("\nYou can now use 'autostart' command from anywhere:")
		fmt.Println("  autostart list          # List all services")
		fmt.Println("  sudo autostart add ...  # Add new service")
		return
	}

	// 检查是否已经安装在其他位置
	if existingPath := gi.findInstalledBinary(); existingPath != "" {
		fmt.Printf("Found existing installation at: %s\n", existingPath)
		if !confirmAction("Replace existing installation") {
			fmt.Println("Installation cancelled.")
			return
		}
	}

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		handleError(fmt.Errorf("failed to get current executable path: %v", err))
		return
	}

	// 确定安装路径
	targetPath := gi.getPreferredInstallPath()

	// 执行安装
	if err := gi.installBinary(execPath, targetPath); err != nil {
		handleError(fmt.Errorf("failed to install globally: %v", err))
		return
	}

	fmt.Printf("✓ Successfully installed to: %s\n", targetPath)
	fmt.Println("\nVerifying installation...")

	// 验证安装
	if err := gi.verifyInstallation(targetPath); err != nil {
		fmt.Printf("Warning: Installation verification failed: %v\n", err)
	} else {
		fmt.Println("✓ Installation verified successfully")
	}

	fmt.Println("\nYou can now use 'autostart' command from anywhere:")
	fmt.Println("  autostart list          # List all services")
	fmt.Println("  sudo autostart add ...  # Add new service")
	fmt.Println("  autostart help          # Show help")
}

// 全局卸载处理函数
func handleUninstallGlobal(gi *GlobalInstaller) {
	fmt.Println("Uninstalling autostart from global environment...")

	// 查找已安装的二进制文件
	installedPath := gi.findInstalledBinary()
	if installedPath == "" {
		fmt.Println("No global installation found.")
		fmt.Println("\nPossible locations checked:")
		for _, path := range gi.possiblePaths {
			fmt.Printf("  - %s\n", path)
		}
		return
	}

	fmt.Printf("Found installation at: %s\n", installedPath)

	// 确认卸载
	if !confirmAction("Remove global installation") {
		fmt.Println("Uninstall cancelled.")
		return
	}

	// 执行卸载
	if err := os.Remove(installedPath); err != nil {
		handleError(fmt.Errorf("failed to remove global installation: %v", err))
		return
	}

	fmt.Printf("✓ Successfully removed: %s\n", installedPath)
	fmt.Println("\nGlobal installation has been removed.")
	fmt.Println("You can still use autostart from the local binary if available.")
}

// installBinary 执行二进制文件安装
func (gi *GlobalInstaller) installBinary(srcPath, targetPath string) error {
	// 创建目标目录
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	// 读取源文件
	srcData, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	// 如果目标文件已存在，先备份
	if _, err := os.Stat(targetPath); err == nil {
		backupPath := targetPath + ".backup"
		if err := os.Rename(targetPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing file: %v", err)
		}
		defer func() {
			// 如果安装失败，恢复备份
			if _, err := os.Stat(targetPath); err != nil {
				os.Rename(backupPath, targetPath)
			} else {
				os.Remove(backupPath)
			}
		}()
	}

	// 写入目标文件
	if err := os.WriteFile(targetPath, srcData, 0755); err != nil {
		return fmt.Errorf("failed to write target file: %v", err)
	}

	return nil
}

// 验证安装
func (gi *GlobalInstaller) verifyInstallation(targetPath string) error {
	// 检查文件是否存在且可执行
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("binary not found: %v", err)
	}

	if info.Mode()&0111 == 0 {
		return fmt.Errorf("binary is not executable")
	}

	// 尝试执行版本命令
	cmd := exec.Command(targetPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("binary execution failed: %v", err)
	}

	if !strings.Contains(string(output), ToolName) {
		return fmt.Errorf("binary output validation failed")
	}

	return nil
}

// 确认操作的辅助函数
func confirmAction(action string) bool {
	fmt.Printf("%s? (y/N): ", action)

	var response string
	fmt.Scanln(&response)

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// GlobalInstaller 全局安装管理器结构
type GlobalInstaller struct {
	binaryName    string
	possiblePaths []string
}

// NewGlobalInstaller 创建全局安装管理器
func NewGlobalInstaller() *GlobalInstaller {
	return &GlobalInstaller{
		binaryName: ToolName,
		possiblePaths: []string{
			"/usr/local/bin/" + ToolName,
			"/usr/bin/" + ToolName,
			"/opt/bin/" + ToolName,
		},
	}
}

// 查找已安装的二进制文件
func (gi *GlobalInstaller) findInstalledBinary() string {
	for _, path := range gi.possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// 检查当前是否从全局安装位置运行
func (gi *GlobalInstaller) isRunningFromGlobal() (bool, string) {
	execPath, err := os.Executable()
	if err != nil {
		return false, ""
	}

	// 解析符号链接
	realPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		realPath = execPath
	}

	for _, globalPath := range gi.possiblePaths {
		if realPath == globalPath {
			return true, globalPath
		}
	}
	return false, ""
}

// 获取合适的安装路径
func (gi *GlobalInstaller) getPreferredInstallPath() string {
	// 检查 /usr/local/bin 是否可写
	if err := gi.testWritePermission("/usr/local/bin"); err == nil {
		return "/usr/local/bin/" + gi.binaryName
	}

	// 检查 /usr/bin 是否可写
	if err := gi.testWritePermission("/usr/bin"); err == nil {
		return "/usr/bin/" + gi.binaryName
	}

	// 回退到第一个路径
	return gi.possiblePaths[0]
}

// 测试目录写权限
func (gi *GlobalInstaller) testWritePermission(dir string) error {
	testFile := filepath.Join(dir, "."+gi.binaryName+"_test")

	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	file.Close()
	os.Remove(testFile)
	return nil
}
