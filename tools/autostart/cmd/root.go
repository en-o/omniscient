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

		// 服务状态查询命令
		"exists": func() { handleServiceExists(sm) },
		"check":  func() { handleServiceExists(sm) },

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

// handleServiceExists 处理服务存在性检查命令
func handleServiceExists(sm *service.ServiceManager) {
	serviceName, err := getServiceNameArg("exists")
	if err != nil {
		handleError(err)
		return
	}

	exists := sm.ServiceExists(serviceName)
	if exists {
		fmt.Printf("✓ Service '%s' exists\n", serviceName)
		sm.ShowServiceBriefStatus(serviceName)
		os.Exit(0)
	} else {
		fmt.Printf("✗ Service '%s' does not exist\n", serviceName)
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
	// 检查并修复 CentOS sudo 环境变量问题
	gi.checkAndFixSudoPath(targetPath)
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
			"/usr/bin/" + ToolName, // CentOS 首选路径
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
	// 在 CentOS 系统上优先使用 /usr/bin
	if gi.isCentOSFamily() {
		if err := gi.testWritePermission("/usr/bin"); err == nil {
			return "/usr/bin/" + gi.binaryName
		}
	}

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

// 修复 CentOS sudo 环境变量问题的检查和处理
func (gi *GlobalInstaller) checkAndFixSudoPath(targetPath string) {
	fmt.Println("\nChecking sudo environment...")

	// 检测操作系统是否为 CentOS/RHEL
	if !gi.isCentOSFamily() {
		fmt.Println("✓ Non-CentOS system, sudo path should work correctly")
		return
	}

	// 测试 sudo 是否能找到命令
	cmd := exec.Command("sudo", "-n", "which", gi.binaryName)
	if err := cmd.Run(); err == nil {
		fmt.Println("✓ sudo can find autostart command")
		return
	}

	fmt.Println("⚠ Detected sudo path issue (common on CentOS/RHEL)")
	fmt.Println("\nTrying to fix sudo secure_path...")

	// 尝试修复 sudoers 配置
	if err := gi.fixSudoSecurePath(targetPath); err != nil {
		fmt.Printf("Warning: Could not automatically fix sudo path: %v\n", err)
		gi.printManualSudoFix(targetPath)
	} else {
		fmt.Println("✓ Successfully updated sudo secure_path")
		fmt.Println("\nVerifying fix...")

		// 验证修复结果
		verifyCmd := exec.Command("sudo", "-n", "which", gi.binaryName)
		if err := verifyCmd.Run(); err == nil {
			fmt.Println("✓ sudo can now find autostart command")
		} else {
			fmt.Println("⚠ sudo may still have issues, please check manually")
		}
	}
}

// 检测是否为 CentOS/RHEL 系列
func (gi *GlobalInstaller) isCentOSFamily() bool {
	// 检查 /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		content := string(data)
		return strings.Contains(content, "centos") ||
			strings.Contains(content, "rhel") ||
			strings.Contains(content, "Red Hat")
	}

	// 检查 /etc/redhat-release
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return true
	}

	return false
}

// 修复 sudo secure_path
func (gi *GlobalInstaller) fixSudoSecurePath(targetPath string) error {
	targetDir := filepath.Dir(targetPath)

	// 读取当前 sudoers 配置
	sudoersPath := "/etc/sudoers"
	content, err := os.ReadFile(sudoersPath)
	if err != nil {
		return fmt.Errorf("failed to read sudoers file: %v", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// 查找 secure_path 行
	var securePathIndex = -1
	var currentPath string

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Defaults") &&
			strings.Contains(line, "secure_path") {
			securePathIndex = i
			// 提取当前路径
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				currentPath = strings.Trim(parts[1], "\" \t")
			}
			break
		}
	}

	if securePathIndex == -1 {
		return fmt.Errorf("secure_path not found in sudoers")
	}

	// 检查目标目录是否已在 secure_path 中
	if strings.Contains(currentPath, targetDir) {
		return nil // 已经包含，无需修改
	}

	// 添加目标目录到 secure_path
	newPath := currentPath + ":" + targetDir
	newLine := fmt.Sprintf("Defaults    secure_path = \"%s\"", newPath)
	lines[securePathIndex] = newLine

	// 备份原文件
	backupPath := sudoersPath + ".autostart.backup"
	if err := os.WriteFile(backupPath, content, 0440); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	// 写入修改后的内容
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(sudoersPath, []byte(newContent), 0440); err != nil {
		// 恢复备份
		os.WriteFile(sudoersPath, content, 0440)
		return fmt.Errorf("failed to update sudoers: %v", err)
	}

	// 验证 sudoers 语法
	cmd := exec.Command("visudo", "-c")
	if err := cmd.Run(); err != nil {
		// 恢复备份
		os.WriteFile(sudoersPath, content, 0440)
		return fmt.Errorf("sudoers syntax error, restored backup: %v", err)
	}

	fmt.Printf("✓ Added %s to sudo secure_path\n", targetDir)
	fmt.Printf("✓ Backup saved to %s\n", backupPath)

	return nil
}

// 打印手动修复说明
func (gi *GlobalInstaller) printManualSudoFix(targetPath string) {
	targetDir := filepath.Dir(targetPath)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("MANUAL SUDO FIX REQUIRED")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nTo fix the sudo issue, you have two options:")
	fmt.Println("\n1. Update sudoers secure_path (Recommended):")
	fmt.Println("   sudo visudo")
	fmt.Println("   Find the line: Defaults    secure_path = \"/usr/local/sbin:...\"")
	fmt.Printf("   Add '%s' to the path, like:\n", targetDir)
	fmt.Printf("   Defaults    secure_path = \"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:%s\"\n", targetDir)
	fmt.Println("\n2. Use full path with sudo:")
	fmt.Printf("   sudo %s add myservice \"command\"\n", targetPath)
	fmt.Println("\n3. Alternative: Install to /usr/bin instead:")
	fmt.Printf("   sudo cp %s /usr/bin/\n", targetPath)
	fmt.Println(strings.Repeat("=", 60))
}
