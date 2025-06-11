package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"autostart/internal/config"
	"autostart/internal/utils"
)

const ToolName = "autostart"

func ListAutostartServices() {
	fmt.Printf("Autostart Services managed by %s:\n", ToolName)
	fmt.Println("============================================")
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service", "--no-pager")
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
				statusCmd := exec.Command("systemctl", "is-active", "autostart-"+serviceName)
				statusOutput, _ := statusCmd.Output()
				activeStatus := strings.TrimSpace(string(statusOutput))
				description := getServiceDescription(serviceName)
				fmt.Printf("%-20s %-12s %-12s %-30s\n", serviceName, autostartStatus, activeStatus, description)
				found = true
			}
		}
	}
	if !found {
		fmt.Println("No autostart services found.")
		fmt.Println("")
		fmt.Printf("Create your first service with: sudo %s add <name> <command>\n", ToolName)
	} else {
		fmt.Println("")
		fmt.Println("LEGEND:")
		fmt.Println("  enabled/disabled - Autostart on boot")
		fmt.Println("  active/inactive  - Current running status")
		fmt.Println("")
		fmt.Printf("Use '%s status <name>' for detailed status\n", ToolName)
		fmt.Printf("Use '%s logs <name>' to view service logs\n", ToolName)
	}
}

func getServiceDescription(serviceName string) string {
	configFile := filepath.Join(config.ConfigDir, serviceName+".json")
	if data, err := ioutil.ReadFile(configFile); err == nil {
		var cfg config.ServiceConfig
		if json.Unmarshal(data, &cfg) == nil && cfg.Description != "" {
			if len(cfg.Description) > 28 {
				return cfg.Description[:25] + "..."
			}
			return cfg.Description
		}
	}
	return "-"
}

func EnableService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: enable <service-name>")
		return
	}
	serviceName := os.Args[2]
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	if !serviceExists(serviceName) {
		fmt.Printf("Error: Service '%s' does not exist\n", serviceName)
		fmt.Printf("Use '%s list' to see available services\n", ToolName)
		return
	}
	cmd := exec.Command("systemctl", "enable", fullServiceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Failed to enable service: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return
	}
	fmt.Printf("✓ Service '%s' enabled for autostart on boot\n", serviceName)
	showServiceBriefStatus(serviceName)
}

func DisableService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: disable <service-name>")
		return
	}
	serviceName := os.Args[2]
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	if !serviceExists(serviceName) {
		fmt.Printf("Error: Service '%s' does not exist\n", serviceName)
		fmt.Printf("Use '%s list' to see available services\n", ToolName)
		return
	}
	cmd := exec.Command("systemctl", "disable", fullServiceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Failed to disable service: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return
	}
	fmt.Printf("✓ Service '%s' disabled from autostart on boot\n", serviceName)
	fmt.Printf("Note: Service is still running if it was started. Use '%s stop %s' to stop it.\n", ToolName, serviceName)
	showServiceBriefStatus(serviceName)
}

func AddAutostartService() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: add <service-name> <exec-start> [options]")
		fmt.Println("Example: add myapp \"java -jar /path/to/app.jar\" --workdir=/path/to --user=myuser")
		return
	}
	serviceName := os.Args[2]
	execStart := os.Args[3]
	options := os.Args[4:]
	if serviceExists(serviceName) {
		fmt.Printf("Error: Service '%s' already exists\n", serviceName)
		fmt.Printf("Use '%s remove %s' to remove it first, or '%s edit %s' to modify it\n",
			ToolName, serviceName, ToolName, serviceName)
		return
	}
	configObj, err := utils.ParseAddOptions(serviceName, execStart, options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if err := os.MkdirAll(config.ConfigDir, 0755); err != nil {
		fmt.Printf("Error: Failed to create config directory: %v\n", err)
		return
	}
	err = createSystemdService(configObj)
	if err != nil {
		fmt.Printf("Error: Failed to create systemd service: %v\n", err)
		return
	}
	err = config.SaveServiceConfig(configObj)
	if err != nil {
		fmt.Printf("Warning: Failed to save service config: %v\n", err)
	}
	fmt.Printf("✓ Service '%s' added successfully!\n", serviceName)
	fmt.Printf("  Command: %s\n", configObj.ExecStart)
	fmt.Printf("  User: %s\n", configObj.User)
	fmt.Printf("  Working Directory: %s\n", configObj.WorkDir)
	fmt.Printf("  Restart Policy: %s\n", configObj.Restart)
	fmt.Println("")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  %s enable %s     # Enable autostart on boot\n", ToolName, serviceName)
	fmt.Printf("  %s start %s      # Start the service now\n", ToolName, serviceName)
	fmt.Printf("  %s status %s     # Check service status\n", ToolName, serviceName)
}

func createSystemdService(cfg *config.ServiceConfig) error {
	serviceName := fmt.Sprintf("autostart-%s", cfg.Name)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	serviceContent := utils.BuildServiceContent(cfg, serviceName)
	err := ioutil.WriteFile(servicePath, []byte(serviceContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create service file: %v", err)
	}
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}
	return nil
}

func RemoveAutostartService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: remove <service-name>")
		return
	}
	serviceName := os.Args[2]
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", fullServiceName)
	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		fmt.Printf("Error: Service '%s' does not exist\n", serviceName)
		return
	}
	fmt.Printf("Stopping service '%s'...\n", serviceName)
	stopCmd := exec.Command("systemctl", "stop", fullServiceName)
	stopCmd.Run()
	fmt.Printf("Disabling service '%s'...\n", serviceName)
	disableCmd := exec.Command("systemctl", "disable", fullServiceName)
	disableCmd.Run()
	err := os.Remove(servicePath)
	if err != nil {
		fmt.Printf("Error: Failed to remove service file: %v\n", err)
		return
	}
	configFile := fmt.Sprintf("%s/%s.json", config.ConfigDir, serviceName)
	os.Remove(configFile)
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Warning: Failed to reload systemd: %v\n", err)
	}
	fmt.Printf("✓ Service '%s' removed successfully!\n", serviceName)
}

func ShowServiceStatus() {
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

func ShowServiceLogs() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: logs <service-name> [lines]")
		return
	}
	serviceName := fmt.Sprintf("autostart-%s", os.Args[2])
	lines := "50"
	if len(os.Args) >= 4 {
		lines = os.Args[3]
	}
	cmd := exec.Command("journalctl", "-u", serviceName, "-n", lines, "--no-pager")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func EditService() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: edit <service-name>")
		return
	}
	serviceName := os.Args[2]
	configFile := fmt.Sprintf("%s/%s.json", config.ConfigDir, serviceName)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Error: Service '%s' configuration not found\n", serviceName)
		return
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error: Failed to read config file: %v\n", err)
		return
	}
	var cfg config.ServiceConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("Error: Failed to parse config file: %v\n", err)
		return
	}
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
	fmt.Println("\nTo modify the configuration, edit the JSON file directly:")
	fmt.Printf("  sudo nano %s\n", configFile)
	fmt.Println("\nAfter editing, recreate the service:")
	fmt.Printf("  sudo %s remove %s\n", ToolName, serviceName)
	fmt.Printf("  sudo %s add %s \"<new-exec-start>\" [new-options]\n", ToolName, serviceName)
}

func StartService() {
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
	fmt.Println("\nCurrent status:")
	statusCmd := exec.Command("systemctl", "is-active", serviceName)
	output, _ := statusCmd.Output()
	status := strings.TrimSpace(string(output))
	fmt.Printf("Status: %s\n", status)
}

func StopService() {
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

func RestartService() {
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
	fmt.Println("\nCurrent status:")
	statusCmd := exec.Command("systemctl", "is-active", serviceName)
	output, _ := statusCmd.Output()
	status := strings.TrimSpace(string(output))
	fmt.Printf("Status: %s\n", status)
}

func serviceExists(serviceName string) bool {
	servicePath := fmt.Sprintf("/etc/systemd/system/autostart-%s.service", serviceName)
	_, err := os.Stat(servicePath)
	return err == nil
}

func showServiceBriefStatus(serviceName string) {
	fullServiceName := fmt.Sprintf("autostart-%s", serviceName)
	enabledCmd := exec.Command("systemctl", "is-enabled", fullServiceName)
	enabledOutput, _ := enabledCmd.Output()
	enabledStatus := strings.TrimSpace(string(enabledOutput))
	activeCmd := exec.Command("systemctl", "is-active", fullServiceName)
	activeOutput, _ := activeCmd.Output()
	activeStatus := strings.TrimSpace(string(activeOutput))
	fmt.Printf("Current status: autostart=%s, running=%s\n", enabledStatus, activeStatus)
}
