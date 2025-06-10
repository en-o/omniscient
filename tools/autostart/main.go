package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ToolName = "autostart"
	Version  = "0.0.1"
)

// 项目配置结构
type ProjectConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`        // jar, py, go, node, etc.
	Path        string            `json:"path"`        // 项目路径
	ExecPath    string            `json:"exec_path"`   // 可执行文件路径
	Args        []string          `json:"args"`        // 启动参数
	Env         map[string]string `json:"env"`         // 环境变量
	WorkDir     string            `json:"work_dir"`    // 工作目录
	User        string            `json:"user"`        // 运行用户
	Description string            `json:"description"` // 服务描述
}

func main() {
	// 检查操作系统支持
	fmt.Printf("system：%s \n", runtime.GOOS)
	if runtime.GOOS != "linux" {
		fmt.Printf("Error: This tool currently only supports Linux systems\n")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	// 检查是否有足够权限
	if needsRoot(command) && os.Geteuid() != 0 {
		fmt.Printf("Error: This operation requires root privileges. Please run with sudo:\n")
		fmt.Printf("  sudo %s %s\n", os.Args[0], strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}

	switch command {
	case "list":
		listAutostartServices()
	case "add":
		addAutostartService()
	case "remove":
		removeAutostartService()
	case "status":
		showServiceStatus()
	case "start":
		startService()
	case "stop":
		stopService()
	case "restart":
		restartService()
	case "version":
		fmt.Printf("%s v%s\n", ToolName, Version)
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}

// 检查命令是否需要root权限
func needsRoot(command string) bool {
	rootCommands := []string{"add", "remove", "start", "stop", "restart"}
	for _, cmd := range rootCommands {
		if command == cmd {
			return true
		}
	}
	return false
}

// 打印帮助信息
func printHelp() {
	fmt.Printf("%s v%s - Universal Project Autostart Management Tool\n\n", ToolName, Version)
	fmt.Println("Usage:")
	fmt.Printf("  %s <command> [arguments]\n\n", ToolName)
	fmt.Println("Commands:")
	fmt.Println("  list                          - List all autostart services")
	fmt.Println("  add <name> <type> <path>      - Add project to autostart")
	fmt.Println("  remove <name>                 - Remove project from autostart")
	fmt.Println("  status <name>                 - Show service status")
	fmt.Println("  start <name>                  - Start service")
	fmt.Println("  stop <name>                   - Stop service")
	fmt.Println("  restart <name>                - Restart service")
	fmt.Println("  version                       - Show version")
	fmt.Println("  help                          - Show this help")
	fmt.Println("")
	fmt.Println("Supported Project Types:")
	fmt.Println("  jar                           - Java JAR applications")
	fmt.Println("  py                            - Python applications")
	fmt.Println("  go                            - Go applications")
	fmt.Println("  node                          - Node.js applications")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Add a Java application")
	fmt.Printf("  sudo %s add myapp jar /path/to/app.jar\n", ToolName)
	fmt.Println("")
	fmt.Println("  # Add a Python application")
	fmt.Printf("  sudo %s add pyapp py /path/to/app.py\n", ToolName)
	fmt.Println("")
	fmt.Println("  # Add a Go application")
	fmt.Printf("  sudo %s add goapp go /path/to/app\n", ToolName)
	fmt.Println("")
	fmt.Println("  # List all services")
	fmt.Printf("  %s list\n", ToolName)
	fmt.Println("")
	fmt.Println("Note: Adding, removing, and controlling services requires root privileges.")
}

// 列出所有自启服务
func listAutostartServices() {
	fmt.Println("Autostart Services managed by this tool:")
	fmt.Println("----------------------------------------")

	// 查找所有以 autostart- 开头的服务
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: Failed to list services: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	found := false

	fmt.Printf("%-20s %-12s %-12s\n", "SERVICE", "AUTOSTART", "STATUS")
	fmt.Println("----------------------------------------------------")

	for _, line := range lines {
		if strings.Contains(line, "autostart-") && strings.Contains(line, ".service") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				serviceName := strings.TrimSuffix(parts[0], ".service")
				serviceName = strings.TrimPrefix(serviceName, "autostart-")
				autostartStatus := parts[1]

				// 获取服务运行状态
				statusCmd := exec.Command("systemctl", "is-active", "autostart-"+serviceName)
				statusOutput, _ := statusCmd.Output()
				activeStatus := strings.TrimSpace(string(statusOutput))

				fmt.Printf("%-20s %-12s %-12s\n", serviceName, autostartStatus, activeStatus)
				found = true
			}
		}
	}

	if !found {
		fmt.Println("No autostart services found.")
	} else {
		fmt.Println("")
		fmt.Println("Legend:")
		fmt.Println("  enabled/disabled - Autostart on boot")
		fmt.Println("  active/inactive  - Current running status")
		fmt.Println("")
		fmt.Printf("Use 'systemctl enable/disable autostart-<name>' to control autostart\n")
	}
}

// 添加自启服务
func addAutostartService() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: add <service-name> <project-type> <project-path> [extra-args...]")
		fmt.Println("Example: add myapp jar /path/to/app.jar --server.port=8080")
		return
	}

	serviceName := os.Args[2]
	projectType := os.Args[3]
	projectPath := os.Args[4]
	extraArgs := os.Args[5:]

	// 验证项目路径
	if !filepath.IsAbs(projectPath) {
		abs, err := filepath.Abs(projectPath)
		if err != nil {
			fmt.Printf("Error: Failed to get absolute path: %v\n", err)
			return
		}
		projectPath = abs
	}

	// 检查文件是否存在
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		fmt.Printf("Error: Project path does not exist: %s\n", projectPath)
		return
	}

	// 构建项目配置
	config, err := buildProjectConfig(serviceName, projectType, projectPath, extraArgs)
	if err != nil {
		fmt.Printf("Error: Failed to build project config: %v\n", err)
		return
	}

	// 创建 systemd 服务
	err = createSystemdService(config)
	if err != nil {
		fmt.Printf("Error: Failed to create systemd service: %v\n", err)
		return
	}

	// 保存配置文件
	err = saveProjectConfig(config)
	if err != nil {
		fmt.Printf("Warning: Failed to save project config: %v\n", err)
	}

	fmt.Printf("✓ Service '%s' added successfully!\n", serviceName)
	fmt.Printf("  Type: %s\n", projectType)
	fmt.Printf("  Path: %s\n", projectPath)
	fmt.Printf("  Working Directory: %s\n", config.WorkDir)
	fmt.Println("")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  sudo systemctl enable autostart-%s   # Enable autostart on boot\n", serviceName)
	fmt.Printf("  sudo systemctl start autostart-%s    # Start the service now\n", serviceName)
	fmt.Printf("  %s status %s                        # Check service status\n", ToolName, serviceName)
}

// 构建项目配置
func buildProjectConfig(name, projectType, projectPath string, extraArgs []string) (*ProjectConfig, error) {
	config := &ProjectConfig{
		Name:    name,
		Type:    projectType,
		Path:    projectPath,
		WorkDir: filepath.Dir(projectPath),
		User:    "root",
		Env:     make(map[string]string),
	}

	switch projectType {
	case "jar":
		config.ExecPath = "java"
		config.Args = []string{"-jar", projectPath}
		config.Description = fmt.Sprintf("Java application: %s", name)

	case "py", "python":
		config.ExecPath = "python3"
		config.Args = []string{projectPath}
		config.Description = fmt.Sprintf("Python application: %s", name)

	case "go":
		config.ExecPath = projectPath
		config.Args = []string{}
		config.Description = fmt.Sprintf("Go application: %s", name)

	case "node", "nodejs":
		config.ExecPath = "node"
		config.Args = []string{projectPath}
		config.Description = fmt.Sprintf("Node.js application: %s", name)

	case "shell", "sh":
		config.ExecPath = "bash"
		config.Args = []string{projectPath}
		config.Description = fmt.Sprintf("Shell script: %s", name)

	default:
		return nil, fmt.Errorf("unsupported project type: %s (supported: jar, py, go, node, shell)", projectType)
	}

	// 添加额外参数
	config.Args = append(config.Args, extraArgs...)

	return config, nil
}

// 创建 systemd 服务
func createSystemdService(config *ProjectConfig) error {
	serviceName := fmt.Sprintf("autostart-%s", config.Name)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)

	// 构建执行命令
	execStart := config.ExecPath
	if len(config.Args) > 0 {
		execStart += " " + strings.Join(config.Args, " ")
	}

	// 构建环境变量
	envVars := ""
	for key, value := range config.Env {
		envVars += fmt.Sprintf("Environment=%s=%s\n", key, value)
	}

	// 创建服务文件内容
	serviceContent := fmt.Sprintf(`[Unit]
Description=%s
After=network.target
Wants=network.target

[Service]
Type=simple
User=%s
WorkingDirectory=%s
ExecStart=%s
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=%s
%s
[Install]
WantedBy=multi-user.target
`, config.Description, config.User, config.WorkDir, execStart, serviceName, envVars)

	// 写入服务文件
	err := ioutil.WriteFile(servicePath, []byte(serviceContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create service file: %v", err)
	}

	// 重载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	return nil
}

// 保存项目配置
func saveProjectConfig(config *ProjectConfig) error {
	configDir := "/etc/autostart-manager"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return err
		}
	}

	configFile := filepath.Join(configDir, config.Name+".json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, data, 0644)
}

// 移除自启服务
func removeAutostartService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: remove <service-name>")
		return
	}

	serviceName := os.Args[2]
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", fullServiceName)

	// 检查服务是否存在
	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		fmt.Printf("Error: Service '%s' does not exist\n", serviceName)
		return
	}

	// 停止服务
	fmt.Printf("Stopping service '%s'...\n", serviceName)
	stopCmd := exec.Command("systemctl", "stop", fullServiceName)
	stopCmd.Run() // 忽略错误，可能服务已经停止

	// 禁用服务
	fmt.Printf("Disabling service '%s'...\n", serviceName)
	disableCmd := exec.Command("systemctl", "disable", fullServiceName)
	disableCmd.Run() // 忽略错误，可能服务已经禁用

	// 删除服务文件
	err := os.Remove(servicePath)
	if err != nil {
		fmt.Printf("Error: Failed to remove service file: %v\n", err)
		return
	}

	// 删除配置文件
	configFile := fmt.Sprintf("/etc/autostart-manager/%s.json", serviceName)
	os.Remove(configFile) // 忽略错误，配置文件可能不存在

	// 重载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Warning: Failed to reload systemd: %v\n", err)
	}

	fmt.Printf("✓ Service '%s' removed successfully!\n", serviceName)
}

// 显示服务状态
func showServiceStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: status <service-name>")
		return
	}

	serviceName := fmt.Sprintf("autostart-%s", os.Args[2])
	cmd := exec.Command("systemctl", "status", serviceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// 启动服务
func startService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: start <service-name>")
		return
	}

	serviceName := fmt.Sprintf("autostart-%s", os.Args[2])
	cmd := exec.Command("systemctl", "start", serviceName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: Failed to start service: %v\n", err)
		return
	}

	fmt.Printf("✓ Service '%s' started successfully!\n", os.Args[2])

	// 显示启动后的状态
	fmt.Println("\nCurrent status:")
	statusCmd := exec.Command("systemctl", "is-active", serviceName)
	output, _ := statusCmd.Output()
	status := strings.TrimSpace(string(output))
	fmt.Printf("Status: %s\n", status)
}

// 停止服务
func stopService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: stop <service-name>")
		return
	}

	serviceName := fmt.Sprintf("autostart-%s", os.Args[2])
	cmd := exec.Command("systemctl", "stop", serviceName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: Failed to stop service: %v\n", err)
		return
	}

	fmt.Printf("✓ Service '%s' stopped successfully!\n", os.Args[2])
}

// 重启服务
func restartService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: restart <service-name>")
		return
	}

	serviceName := fmt.Sprintf("autostart-%s", os.Args[2])
	cmd := exec.Command("systemctl", "restart", serviceName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: Failed to restart service: %v\n", err)
		return
	}

	fmt.Printf("✓ Service '%s' restarted successfully!\n", os.Args[2])

	// 显示重启后的状态
	fmt.Println("\nCurrent status:")
	statusCmd := exec.Command("systemctl", "is-active", serviceName)
	output, _ := statusCmd.Output()
	status := strings.TrimSpace(string(output))
	fmt.Printf("Status: %s\n", status)
}
