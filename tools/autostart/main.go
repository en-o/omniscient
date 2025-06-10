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
	Version  = "0.0.2"
)

// 服务配置结构
type ServiceConfig struct {
	Name         string            `json:"name"`
	ExecStart    string            `json:"exec_start"`    // 完整的启动命令
	WorkDir      string            `json:"work_dir"`      // 工作目录
	User         string            `json:"user"`          // 运行用户
	Group        string            `json:"group"`         // 运行组
	Description  string            `json:"description"`   // 服务描述
	Env          map[string]string `json:"env"`           // 环境变量
	Restart      string            `json:"restart"`       // 重启策略
	RestartSec   int               `json:"restart_sec"`   // 重启间隔
	KillMode     string            `json:"kill_mode"`     // 终止模式
	KillSignal   string            `json:"kill_signal"`   // 终止信号
	TimeoutStart int               `json:"timeout_start"` // 启动超时
	TimeoutStop  int               `json:"timeout_stop"`  // 停止超时
	After        []string          `json:"after"`         // 依赖服务
	Wants        []string          `json:"wants"`         // 期望服务
	Requires     []string          `json:"requires"`      // 必需服务
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
	case "logs":
		showServiceLogs()
	case "edit":
		editService()
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
	rootCommands := []string{"add", "remove", "start", "stop", "restart", "edit"}
	for _, cmd := range rootCommands {
		if command == cmd {
			return true
		}
	}
	return false
}

// 打印帮助信息
func printHelp() {
	fmt.Printf("%s v%s - Universal Service Autostart Management Tool\n\n", ToolName, Version)
	fmt.Println("Usage:")
	fmt.Printf("  %s <command> [arguments]\n\n", ToolName)
	fmt.Println("Commands:")
	fmt.Println("  list                                    - List all autostart services")
	fmt.Println("  add <name> <exec-start> [options]      - Add service to autostart")
	fmt.Println("  remove <name>                           - Remove service from autostart")
	fmt.Println("  status <name>                           - Show service status")
	fmt.Println("  start <name>                            - Start service")
	fmt.Println("  stop <name>                             - Stop service")
	fmt.Println("  restart <name>                          - Restart service")
	fmt.Println("  logs <name> [lines]                     - Show service logs")
	fmt.Println("  edit <name>                             - Edit service configuration")
	fmt.Println("  version                                 - Show version")
	fmt.Println("  help                                    - Show this help")
	fmt.Println("")
	fmt.Println("Add Options:")
	fmt.Println("  --workdir=<path>                        - Working directory (default: exec file directory)")
	fmt.Println("  --user=<username>                       - Run as user (default: root)")
	fmt.Println("  --group=<groupname>                     - Run as group (default: user's primary group)")
	fmt.Println("  --description=<text>                    - Service description")
	fmt.Println("  --env=<KEY=VALUE>                       - Environment variable (can be used multiple times)")
	fmt.Println("  --restart=<policy>                      - Restart policy (always|on-failure|no, default: always)")
	fmt.Println("  --restart-sec=<seconds>                 - Restart delay in seconds (default: 5)")
	fmt.Println("  --kill-mode=<mode>                      - Kill mode (control-group|process|mixed, default: control-group)")
	fmt.Println("  --kill-signal=<signal>                  - Kill signal (default: SIGTERM)")
	fmt.Println("  --timeout-start=<seconds>               - Start timeout (default: 90)")
	fmt.Println("  --timeout-stop=<seconds>                - Stop timeout (default: 90)")
	fmt.Println("  --after=<service>                       - Start after service (can be used multiple times)")
	fmt.Println("  --wants=<service>                       - Wants service (can be used multiple times)")
	fmt.Println("  --requires=<service>                    - Requires service (can be used multiple times)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Add a Java application")
	fmt.Printf("  sudo %s add myapp \"java -jar /path/to/app.jar\" --workdir=/path/to\n", ToolName)
	fmt.Println("")
	fmt.Println("  # Add a Python application with environment variables")
	fmt.Printf("  sudo %s add pyapp \"python3 /path/to/app.py\" --user=www-data --env=PYTHONPATH=/path/to --env=DEBUG=true\n", ToolName)
	fmt.Println("")
	fmt.Println("  # Add a Node.js application with custom restart policy")
	fmt.Printf("  sudo %s add nodeapp \"node /path/to/app.js\" --restart=on-failure --restart-sec=10\n", ToolName)
	fmt.Println("")
	fmt.Println("  # Add a service with dependencies")
	fmt.Printf("  sudo %s add webapp \"./webapp\" --after=network.target --after=mysql.service --requires=mysql.service\n", ToolName)
	fmt.Println("")
	fmt.Println("  # View service logs")
	fmt.Printf("  %s logs myapp 50\n", ToolName)
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

	fmt.Printf("%-20s %-12s %-12s %-30s\n", "SERVICE", "AUTOSTART", "STATUS", "DESCRIPTION")
	fmt.Println("--------------------------------------------------------------------------------")

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

				// 获取服务描述
				description := getServiceDescription(serviceName)

				fmt.Printf("%-20s %-12s %-12s %-30s\n", serviceName, autostartStatus, activeStatus, description)
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
		fmt.Printf("Use '%s status <name>' for detailed status\n", ToolName)
		fmt.Printf("Use '%s logs <name>' to view service logs\n", ToolName)
	}
}

// 获取服务描述
func getServiceDescription(serviceName string) string {
	configFile := fmt.Sprintf("/etc/autostart-manager/%s.json", serviceName)
	if data, err := ioutil.ReadFile(configFile); err == nil {
		var config ServiceConfig
		if json.Unmarshal(data, &config) == nil {
			if config.Description != "" {
				return config.Description
			}
		}
	}
	return "-"
}

// 添加自启服务
func addAutostartService() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: add <service-name> <exec-start> [options]")
		fmt.Println("Example: add myapp \"java -jar /path/to/app.jar\" --workdir=/path/to --user=myuser")
		return
	}

	serviceName := os.Args[2]
	execStart := os.Args[3]
	options := os.Args[4:]

	// 解析选项
	config, err := parseAddOptions(serviceName, execStart, options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 创建 systemd 服务
	err = createSystemdService(config)
	if err != nil {
		fmt.Printf("Error: Failed to create systemd service: %v\n", err)
		return
	}

	// 保存配置文件
	err = saveServiceConfig(config)
	if err != nil {
		fmt.Printf("Warning: Failed to save service config: %v\n", err)
	}

	fmt.Printf("✓ Service '%s' added successfully!\n", serviceName)
	fmt.Printf("  Exec: %s\n", config.ExecStart)
	fmt.Printf("  User: %s\n", config.User)
	fmt.Printf("  Working Directory: %s\n", config.WorkDir)
	fmt.Printf("  Restart Policy: %s\n", config.Restart)
	fmt.Println("")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  sudo systemctl enable autostart-%s   # Enable autostart on boot\n", serviceName)
	fmt.Printf("  sudo systemctl start autostart-%s    # Start the service now\n", serviceName)
	fmt.Printf("  %s status %s                        # Check service status\n", ToolName, serviceName)
}

// 解析添加选项
func parseAddOptions(name, execStart string, options []string) (*ServiceConfig, error) {
	config := &ServiceConfig{
		Name:         name,
		ExecStart:    execStart,
		User:         "root",
		Group:        "",
		Description:  fmt.Sprintf("Autostart service: %s", name),
		Env:          make(map[string]string),
		Restart:      "always",
		RestartSec:   5,
		KillMode:     "control-group",
		KillSignal:   "SIGTERM",
		TimeoutStart: 90,
		TimeoutStop:  90,
		After:        []string{"network.target"},
		Wants:        []string{"network.target"},
		Requires:     []string{},
	}

	// 尝试从执行命令中推断工作目录
	parts := strings.Fields(execStart)
	if len(parts) > 0 {
		// 查找可能的文件路径
		for _, part := range parts {
			if strings.Contains(part, "/") && !strings.HasPrefix(part, "-") {
				if filepath.IsAbs(part) {
					if _, err := os.Stat(part); err == nil {
						config.WorkDir = filepath.Dir(part)
						break
					}
				}
			}
		}
	}

	// 如果没有推断出工作目录，使用当前目录
	if config.WorkDir == "" {
		if wd, err := os.Getwd(); err == nil {
			config.WorkDir = wd
		} else {
			config.WorkDir = "/tmp"
		}
	}

	// 解析选项
	for _, option := range options {
		if strings.HasPrefix(option, "--workdir=") {
			config.WorkDir = strings.TrimPrefix(option, "--workdir=")
		} else if strings.HasPrefix(option, "--user=") {
			config.User = strings.TrimPrefix(option, "--user=")
		} else if strings.HasPrefix(option, "--group=") {
			config.Group = strings.TrimPrefix(option, "--group=")
		} else if strings.HasPrefix(option, "--description=") {
			config.Description = strings.TrimPrefix(option, "--description=")
		} else if strings.HasPrefix(option, "--env=") {
			envStr := strings.TrimPrefix(option, "--env=")
			if parts := strings.SplitN(envStr, "=", 2); len(parts) == 2 {
				config.Env[parts[0]] = parts[1]
			}
		} else if strings.HasPrefix(option, "--restart=") {
			config.Restart = strings.TrimPrefix(option, "--restart=")
		} else if strings.HasPrefix(option, "--restart-sec=") {
			if sec := strings.TrimPrefix(option, "--restart-sec="); sec != "" {
				fmt.Sscanf(sec, "%d", &config.RestartSec)
			}
		} else if strings.HasPrefix(option, "--kill-mode=") {
			config.KillMode = strings.TrimPrefix(option, "--kill-mode=")
		} else if strings.HasPrefix(option, "--kill-signal=") {
			config.KillSignal = strings.TrimPrefix(option, "--kill-signal=")
		} else if strings.HasPrefix(option, "--timeout-start=") {
			if timeout := strings.TrimPrefix(option, "--timeout-start="); timeout != "" {
				fmt.Sscanf(timeout, "%d", &config.TimeoutStart)
			}
		} else if strings.HasPrefix(option, "--timeout-stop=") {
			if timeout := strings.TrimPrefix(option, "--timeout-stop="); timeout != "" {
				fmt.Sscanf(timeout, "%d", &config.TimeoutStop)
			}
		} else if strings.HasPrefix(option, "--after=") {
			service := strings.TrimPrefix(option, "--after=")
			config.After = append(config.After, service)
		} else if strings.HasPrefix(option, "--wants=") {
			service := strings.TrimPrefix(option, "--wants=")
			config.Wants = append(config.Wants, service)
		} else if strings.HasPrefix(option, "--requires=") {
			service := strings.TrimPrefix(option, "--requires=")
			config.Requires = append(config.Requires, service)
		} else {
			return nil, fmt.Errorf("unknown option: %s", option)
		}
	}

	// 验证工作目录
	if config.WorkDir != "" {
		if !filepath.IsAbs(config.WorkDir) {
			abs, err := filepath.Abs(config.WorkDir)
			if err != nil {
				return nil, fmt.Errorf("failed to get absolute path for workdir: %v", err)
			}
			config.WorkDir = abs
		}
		if _, err := os.Stat(config.WorkDir); os.IsNotExist(err) {
			return nil, fmt.Errorf("working directory does not exist: %s", config.WorkDir)
		}
	}

	return config, nil
}

// 创建 systemd 服务
func createSystemdService(config *ServiceConfig) error {
	serviceName := fmt.Sprintf("autostart-%s", config.Name)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)

	// 构建 Unit 段
	unitSection := "[Unit]\n"
	unitSection += fmt.Sprintf("Description=%s\n", config.Description)

	if len(config.After) > 0 {
		unitSection += fmt.Sprintf("After=%s\n", strings.Join(config.After, " "))
	}
	if len(config.Wants) > 0 {
		unitSection += fmt.Sprintf("Wants=%s\n", strings.Join(config.Wants, " "))
	}
	if len(config.Requires) > 0 {
		unitSection += fmt.Sprintf("Requires=%s\n", strings.Join(config.Requires, " "))
	}
	unitSection += "\n"

	// 构建 Service 段
	serviceSection := "[Service]\n"
	serviceSection += "Type=simple\n"
	serviceSection += fmt.Sprintf("User=%s\n", config.User)
	if config.Group != "" {
		serviceSection += fmt.Sprintf("Group=%s\n", config.Group)
	}
	if config.WorkDir != "" {
		serviceSection += fmt.Sprintf("WorkingDirectory=%s\n", config.WorkDir)
	}
	serviceSection += fmt.Sprintf("ExecStart=%s\n", config.ExecStart)
	serviceSection += fmt.Sprintf("Restart=%s\n", config.Restart)
	serviceSection += fmt.Sprintf("RestartSec=%d\n", config.RestartSec)
	serviceSection += fmt.Sprintf("KillMode=%s\n", config.KillMode)
	serviceSection += fmt.Sprintf("KillSignal=%s\n", config.KillSignal)
	serviceSection += fmt.Sprintf("TimeoutStartSec=%d\n", config.TimeoutStart)
	serviceSection += fmt.Sprintf("TimeoutStopSec=%d\n", config.TimeoutStop)
	serviceSection += "StandardOutput=journal\n"
	serviceSection += "StandardError=journal\n"
	serviceSection += fmt.Sprintf("SyslogIdentifier=%s\n", serviceName)

	// 添加环境变量
	for key, value := range config.Env {
		serviceSection += fmt.Sprintf("Environment=%s=%s\n", key, value)
	}
	serviceSection += "\n"

	// 构建 Install 段
	installSection := "[Install]\n"
	installSection += "WantedBy=multi-user.target\n"

	// 组合完整的服务文件内容
	serviceContent := unitSection + serviceSection + installSection

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

// 保存服务配置
func saveServiceConfig(config *ServiceConfig) error {
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

// 显示服务日志
func showServiceLogs() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: logs <service-name> [lines]")
		return
	}

	serviceName := fmt.Sprintf("autostart-%s", os.Args[2])
	lines := "50" // 默认显示50行

	if len(os.Args) >= 4 {
		lines = os.Args[3]
	}

	cmd := exec.Command("journalctl", "-u", serviceName, "-n", lines, "--no-pager")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// 编辑服务
func editService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: edit <service-name>")
		return
	}

	serviceName := os.Args[2]
	configFile := fmt.Sprintf("/etc/autostart-manager/%s.json", serviceName)

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Error: Service '%s' configuration not found\n", serviceName)
		return
	}

	// 读取现有配置
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error: Failed to read config file: %v\n", err)
		return
	}

	var config ServiceConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error: Failed to parse config file: %v\n", err)
		return
	}

	// 显示当前配置
	fmt.Printf("Current configuration for service '%s':\n", serviceName)
	fmt.Printf("Description: %s\n", config.Description)
	fmt.Printf("ExecStart: %s\n", config.ExecStart)
	fmt.Printf("WorkDir: %s\n", config.WorkDir)
	fmt.Printf("User: %s\n", config.User)
	fmt.Printf("Restart: %s\n", config.Restart)
	fmt.Printf("RestartSec: %d\n", config.RestartSec)

	if len(config.Env) > 0 {
		fmt.Println("Environment Variables:")
		for k, v := range config.Env {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	fmt.Println("\nTo modify the configuration, edit the JSON file directly:")
	fmt.Printf("  sudo nano %s\n", configFile)
	fmt.Println("\nAfter editing, recreate the service:")
	fmt.Printf("  sudo %s remove %s\n", ToolName, serviceName)
	fmt.Printf("  sudo %s add %s \"<new-exec-start>\" [new-options]\n", ToolName, serviceName)
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
