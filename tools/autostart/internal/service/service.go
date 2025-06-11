package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"autostart/internal/config"
	"autostart/internal/utils"
)

const ToolName = "autostart"

// ServiceManager 服务管理器
type ServiceManager struct{}

// NewServiceManager 创建新的服务管理器
func NewServiceManager() *ServiceManager {
	return &ServiceManager{}
}

// ListAutostartServices 列出所有自启服务
func (sm *ServiceManager) ListAutostartServices() error {
	fmt.Printf("Autostart Services managed by %s:\n", ToolName)
	fmt.Println("============================================")

	services, err := sm.getAutostartServices()
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	if len(services) == 0 {
		fmt.Println("No autostart services found.")
		fmt.Println("")
		fmt.Printf("Create your first service with: sudo %s add <name> <command>\n", ToolName)
		return nil
	}

	fmt.Printf("%-20s %-12s %-12s %-30s\n", "SERVICE", "AUTOSTART", "STATUS", "DESCRIPTION")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, svc := range services {
		fmt.Printf("%-20s %-12s %-12s %-30s\n",
			svc.Name, svc.AutostartStatus, svc.ActiveStatus, svc.Description)
	}

	fmt.Println("")
	fmt.Println("LEGEND:")
	fmt.Println("  enabled/disabled - Autostart on boot")
	fmt.Println("  active/inactive  - Current running status")
	fmt.Println("")
	fmt.Printf("Use '%s status <name>' for detailed status\n", ToolName)
	fmt.Printf("Use '%s logs <name>' to view service logs\n", ToolName)

	return nil
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name            string
	AutostartStatus string
	ActiveStatus    string
	Description     string
}

// getAutostartServices 获取所有自启服务信息
func (sm *ServiceManager) getAutostartServices() ([]ServiceInfo, error) {
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var services []ServiceInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "autostart-") || !strings.Contains(line, ".service") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		serviceName := strings.TrimSuffix(parts[0], ".service")
		serviceName = strings.TrimPrefix(serviceName, "autostart-")
		autostartStatus := parts[1]

		activeStatus := sm.getServiceActiveStatus(serviceName)
		description := sm.getServiceDescription(serviceName)

		services = append(services, ServiceInfo{
			Name:            serviceName,
			AutostartStatus: autostartStatus,
			ActiveStatus:    activeStatus,
			Description:     description,
		})
	}

	return services, scanner.Err()
}

// getServiceActiveStatus 获取服务激活状态
func (sm *ServiceManager) getServiceActiveStatus(serviceName string) string {
	cmd := exec.Command("systemctl", "is-active", "autostart-"+serviceName)
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output))
}

// getServiceDescription 获取服务描述
func (sm *ServiceManager) getServiceDescription(serviceName string) string {
	configFile := filepath.Join(config.ConfigDir, serviceName+".json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return "-"
	}

	var cfg config.ServiceConfig
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.Description == "" {
		return "-"
	}
	// 截断长描述
	if len(cfg.Description) > 28 {
		return cfg.Description[:25] + "..."
	}
	return cfg.Description
}

// EnableService 启用服务自启动
func (sm *ServiceManager) EnableService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("service '%s' does not exist", serviceName)
	}

	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	cmd := exec.Command("systemctl", "enable", fullServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable service: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✓ Service '%s' enabled for autostart on boot\n", serviceName)
	sm.showServiceBriefStatus(serviceName)
	return nil
}

// DisableService 禁用服务自启动
func (sm *ServiceManager) DisableService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("service '%s' does not exist", serviceName)
	}

	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	cmd := exec.Command("systemctl", "disable", fullServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to disable service: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✓ Service '%s' disabled from autostart on boot\n", serviceName)
	fmt.Printf("Note: Service is still running if it was started. Use '%s stop %s' to stop it.\n", ToolName, serviceName)
	sm.showServiceBriefStatus(serviceName)
	return nil
}

// AddAutostartService 添加自启服务
func (sm *ServiceManager) AddAutostartService(serviceName, execStart string, options []string) error {
	if sm.serviceExists(serviceName) {
		return fmt.Errorf("service '%s' already exists. Use '%s remove %s' to remove it first, or '%s edit %s' to modify it",
			serviceName, ToolName, serviceName, ToolName, serviceName)
	}

	configObj, err := utils.ParseAddOptions(serviceName, execStart, options)
	if err != nil {
		return fmt.Errorf("failed to parse options: %w", err)
	}

	if err := sm.createSystemdService(configObj); err != nil {
		return fmt.Errorf("failed to create systemd service: %w", err)
	}

	if err := config.SaveServiceConfig(configObj); err != nil {
		fmt.Printf("Warning: Failed to save service config: %v\n", err)
	}

	sm.printServiceAddedInfo(configObj)
	return nil
}

// printServiceAddedInfo 打印服务添加成功信息
func (sm *ServiceManager) printServiceAddedInfo(cfg *config.ServiceConfig) {
	fmt.Printf("✓ Service '%s' added successfully!\n", cfg.Name)
	fmt.Printf("  Command: %s\n", cfg.ExecStart)
	fmt.Printf("  User: %s\n", cfg.User)
	fmt.Printf("  Working Directory: %s\n", cfg.WorkDir)
	fmt.Printf("  Restart Policy: %s\n", cfg.Restart)
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Printf("  %s enable %s     # Enable autostart on boot\n", ToolName, cfg.Name)
	fmt.Printf("  %s start %s      # Start the service now\n", ToolName, cfg.Name)
	fmt.Printf("  %s status %s     # Check service status\n", ToolName, cfg.Name)
}

// createSystemdService 创建systemd服务
func (sm *ServiceManager) createSystemdService(cfg *config.ServiceConfig) error {
	serviceName := fmt.Sprintf("autostart-%s", cfg.Name)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	// 构建服务文件内容
	serviceContent := utils.BuildServiceContent(cfg, serviceName)
	// 写入服务文件
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	// 重载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	return nil
}

// RemoveAutostartService 移除自启服务
func (sm *ServiceManager) RemoveAutostartService(serviceName string) error {
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", fullServiceName)

	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		return fmt.Errorf("service '%s' does not exist", serviceName)
	}

	if !sm.confirmRemoval(serviceName) {
		fmt.Println("Removal cancelled.")
		return nil
	}

	// 停止并禁用服务
	sm.stopAndDisableService(fullServiceName)

	// 删除服务文件
	if err := os.Remove(servicePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// 删除配置文件
	configFile := filepath.Join(config.ConfigDir, serviceName+".json")
	os.Remove(configFile) // 忽略错误

	// 重载systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Failed to reload systemd: %v\n", err)
	}

	fmt.Printf("✓ Service '%s' removed successfully!\n", serviceName)
	return nil
}

// confirmRemoval 确认移除操作
func (sm *ServiceManager) confirmRemoval(serviceName string) bool {
	fmt.Printf("This will remove service '%s'. Are you sure? (y/N): ", serviceName)

	var response string
	fmt.Scanln(&response)

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// stopAndDisableService 停止并禁用服务
func (sm *ServiceManager) stopAndDisableService(fullServiceName string) {
	fmt.Printf("Stopping service...\n")
	stopCmd := exec.Command("systemctl", "stop", fullServiceName)
	stopCmd.Run() // 忽略错误

	fmt.Printf("Disabling service...\n")
	disableCmd := exec.Command("systemctl", "disable", fullServiceName)
	disableCmd.Run() // 忽略错误
}

// ShowServiceStatus 显示服务状态
func (sm *ServiceManager) ShowServiceStatus(serviceName string) error {
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	cmd := exec.Command("systemctl", "status", fullServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ShowServiceLogs 显示服务日志
func (sm *ServiceManager) ShowServiceLogs(serviceName string, lines string) error {
	if lines == "" {
		lines = "50"
	}
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	cmd := exec.Command("journalctl", "-u", fullServiceName, "-n", lines, "--no-pager")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// EditService 编辑服务配置
func (sm *ServiceManager) EditService(serviceName string) error {
	configFile := filepath.Join(config.ConfigDir, serviceName+".json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("service '%s' configuration not found", serviceName)
	}

	cfg, err := config.LoadServiceConfig(serviceName)
	if err != nil {
		return fmt.Errorf("failed to load service config: %w", err)
	}

	sm.printCurrentConfig(serviceName, cfg)
	return nil
}

// printCurrentConfig 打印当前配置
func (sm *ServiceManager) printCurrentConfig(serviceName string, cfg *config.ServiceConfig) {
	fmt.Printf("Current configuration for service '%s':\n", serviceName)
	fmt.Printf("Description: %s\n", cfg.Description)
	fmt.Printf("ExecStart: %s\n", cfg.ExecStart)
	fmt.Printf("WorkDir: %s\n", cfg.WorkDir)
	fmt.Printf("User: %s\n", cfg.User)
	fmt.Printf("Restart: %s\n", cfg.Restart)
	fmt.Printf("RestartSec: %d\n", cfg.RestartSec)

	if len(cfg.Env) > 0 {
		fmt.Println("Environment Variables:")
		for k, v := range cfg.Env {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	configFile := filepath.Join(config.ConfigDir, serviceName+".json")
	fmt.Println("\nTo modify the configuration, edit the JSON file directly:")
	fmt.Printf("  sudo nano %s\n", configFile)
	fmt.Println("\nAfter editing, recreate the service:")
	fmt.Printf("  sudo %s remove %s\n", ToolName, serviceName)
	fmt.Printf("  sudo %s add %s \"<new-exec-start>\" [new-options]\n", ToolName, serviceName)
}

// StartService 启动服务
func (sm *ServiceManager) StartService(serviceName string) error {
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	cmd := exec.Command("systemctl", "start", fullServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	fmt.Printf("✓ Service '%s' started successfully!\n", serviceName)
	fmt.Println("\nCurrent status:")
	status := sm.getServiceActiveStatus(serviceName)
	fmt.Printf("Status: %s\n", status)
	return nil
}

// StopService 停止服务
func (sm *ServiceManager) StopService(serviceName string) error {
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	cmd := exec.Command("systemctl", "stop", fullServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	fmt.Printf("✓ Service '%s' stopped successfully!\n", serviceName)
	return nil
}

// RestartService 重启服务
func (sm *ServiceManager) RestartService(serviceName string) error {
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	cmd := exec.Command("systemctl", "restart", fullServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	fmt.Printf("✓ Service '%s' restarted successfully!\n", serviceName)
	fmt.Println("\nCurrent status:")
	status := sm.getServiceActiveStatus(serviceName)
	fmt.Printf("Status: %s\n", status)
	return nil
}

// serviceExists 检查服务是否存在
func (sm *ServiceManager) serviceExists(serviceName string) bool {
	servicePath := fmt.Sprintf("/etc/systemd/system/autostart-%s.service", serviceName)
	_, err := os.Stat(servicePath)
	return err == nil
}

// showServiceBriefStatus 显示服务简要状态
func (sm *ServiceManager) showServiceBriefStatus(serviceName string) {
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)

	enabledCmd := exec.Command("systemctl", "is-enabled", fullServiceName)
	enabledOutput, _ := enabledCmd.Output()
	enabledStatus := strings.TrimSpace(string(enabledOutput))

	activeStatus := sm.getServiceActiveStatus(serviceName)
	fmt.Printf("Current status: autostart=%s, running=%s\n", enabledStatus, activeStatus)
}
