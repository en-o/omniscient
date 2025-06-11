package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"autostart/internal/config"
)

func NeedsRoot(command string) bool {
	rootCommands := []string{"add", "create", "install", "remove", "rm", "delete", "uninstall", "enable", "disable", "uninstall", "start", "stop", "restart", "edit"}
	for _, cmd := range rootCommands {
		if command == cmd {
			return true
		}
	}
	return false
}

func PrintHelp() {
	fmt.Printf("%s v%s - Universal Service Autostart Management Tool\n\n", "autostart", "0.0.2")
	fmt.Println("USAGE:")
	fmt.Printf("  %s <command> [arguments]\n\n", "autostart")
	fmt.Println("COMMANDS:")
	fmt.Println("  Service Management:")
	fmt.Println("    list, ls                              - List all autostart services")
	fmt.Println("    add, create, install <name> <exec> [options]   - Add service to autostart")
	fmt.Println("    remove, rm, delete, uninstall <name>             - Remove service from autostart")
	fmt.Println("    edit <name>                           - Edit service configuration")
	fmt.Println("")
	fmt.Println("  Service Control:")
	fmt.Println("    enable <name>                         - Enable service autostart")
	fmt.Println("    disable <name>                        - Disable service autostart")
	fmt.Println("    start <name>                          - Start service")
	fmt.Println("    stop <name>                           - Stop service")
	fmt.Println("    restart <name>                        - Restart service")
	fmt.Println("    status <name>                         - Show service status")
	fmt.Println("    logs <name> [lines]                   - Show service logs")
	fmt.Println("")
	fmt.Println("  Tool Management:")
	fmt.Println("    uninstall                             - Uninstall this tool and all managed services")
	fmt.Println("    version, -v, --version                - Show version")
	fmt.Println("    help, -h, --help                      - Show this help")
	fmt.Println("")
	fmt.Println("ADD OPTIONS:")
	fmt.Println("  --workdir=<path>                        - Working directory (default: auto-detect)")
	fmt.Println("  --user=<username>                       - Run as user (default: root)")
	fmt.Println("  --group=<groupname>                     - Run as group (default: user's primary group)")
	fmt.Println("  --description=<text>                    - Service description")
	fmt.Println("  --env=<KEY=VALUE>                       - Environment variable (repeatable)")
	fmt.Println("  --restart=<policy>                      - Restart policy: always|on-failure|no (default: always)")
	fmt.Println("  --restart-sec=<seconds>                 - Restart delay in seconds (default: 5)")
	fmt.Println("  --kill-mode=<mode>                      - Kill mode: control-group|process|mixed (default: control-group)")
	fmt.Println("  --kill-signal=<signal>                  - Kill signal (default: SIGTERM)")
	fmt.Println("  --timeout-start=<seconds>               - Start timeout (default: 90)")
	fmt.Println("  --timeout-stop=<seconds>                - Stop timeout (default: 90)")
	fmt.Println("  --after=<service>                       - Start after service (repeatable)")
	fmt.Println("  --wants=<service>                       - Wants service (repeatable)")
	fmt.Println("  --requires=<service>                    - Requires service (repeatable)")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Add and start a Java application")
	fmt.Printf("  sudo %s add myapp \"java -jar /path/to/app.jar\" --workdir=/path/to\n", "autostart")
	fmt.Printf("  sudo %s enable myapp\n", "autostart")
	fmt.Printf("  sudo %s start myapp\n", "autostart")
	fmt.Println("")
	fmt.Println("  # Add a Python application with custom user and environment")
	fmt.Printf("  sudo %s add pyapp \"python3 /path/to/app.py\" --user=www-data --env=PYTHONPATH=/path/to --env=DEBUG=true\n", "autostart")
	fmt.Println("")
	fmt.Println("  # Add a Node.js application with custom restart policy")
	fmt.Printf("  sudo %s add nodeapp \"node /path/to/app.js\" --restart=on-failure --restart-sec=10\n", "autostart")
	fmt.Println("")
	fmt.Println("  # Add a service with dependencies")
	fmt.Printf("  sudo %s add webapp \"./webapp\" --after=network.target --after=mysql.service --requires=mysql.service\n", "autostart")
	fmt.Println("")
	fmt.Println("  # Common workflow")
	fmt.Printf("  %s list                    # List all services\n", "autostart")
	fmt.Printf("  %s status myapp            # Check service status\n", "autostart")
	fmt.Printf("  %s logs myapp 100          # View recent logs\n", "autostart")
	fmt.Printf("  sudo %s restart myapp      # Restart service\n", "autostart")
}

func ParseAddOptions(name, execStart string, options []string) (*config.ServiceConfig, error) {
	cfg := &config.ServiceConfig{
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
	cfg.WorkDir = inferWorkingDirectory(execStart)
	for _, option := range options {
		if err := parseOption(cfg, option); err != nil {
			return nil, err
		}
	}
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func inferWorkingDirectory(execStart string) string {
	parts := strings.Fields(execStart)
	if len(parts) > 0 {
		for _, part := range parts {
			if strings.Contains(part, "/") && !strings.HasPrefix(part, "-") {
				if filepath.IsAbs(part) {
					if _, err := os.Stat(part); err == nil {
						return filepath.Dir(part)
					}
				}
			}
		}
	}
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "/tmp"
}

func parseOption(cfg *config.ServiceConfig, option string) error {
	if strings.HasPrefix(option, "--workdir=") {
		cfg.WorkDir = strings.TrimPrefix(option, "--workdir=")
	} else if strings.HasPrefix(option, "--user=") {
		cfg.User = strings.TrimPrefix(option, "--user=")
	} else if strings.HasPrefix(option, "--group=") {
		cfg.Group = strings.TrimPrefix(option, "--group=")
	} else if strings.HasPrefix(option, "--description=") {
		cfg.Description = strings.TrimPrefix(option, "--description=")
	} else if strings.HasPrefix(option, "--env=") {
		envStr := strings.TrimPrefix(option, "--env=")
		if parts := strings.SplitN(envStr, "=", 2); len(parts) == 2 {
			cfg.Env[parts[0]] = parts[1]
		} else {
			return fmt.Errorf("invalid environment variable format: %s (use KEY=VALUE)", envStr)
		}
	} else if strings.HasPrefix(option, "--restart=") {
		restart := strings.TrimPrefix(option, "--restart=")
		if restart != "always" && restart != "on-failure" && restart != "no" {
			return fmt.Errorf("invalid restart policy: %s (use always, on-failure, or no)", restart)
		}
		cfg.Restart = restart
	} else if strings.HasPrefix(option, "--restart-sec=") {
		if sec := strings.TrimPrefix(option, "--restart-sec="); sec != "" {
			if _, err := fmt.Sscanf(sec, "%d", &cfg.RestartSec); err != nil {
				return fmt.Errorf("invalid restart-sec value: %s", sec)
			}
		}
	} else if strings.HasPrefix(option, "--kill-mode=") {
		cfg.KillMode = strings.TrimPrefix(option, "--kill-mode=")
	} else if strings.HasPrefix(option, "--kill-signal=") {
		cfg.KillSignal = strings.TrimPrefix(option, "--kill-signal=")
	} else if strings.HasPrefix(option, "--timeout-start=") {
		if timeout := strings.TrimPrefix(option, "--timeout-start="); timeout != "" {
			if _, err := fmt.Sscanf(timeout, "%d", &cfg.TimeoutStart); err != nil {
				return fmt.Errorf("invalid timeout-start value: %s", timeout)
			}
		}
	} else if strings.HasPrefix(option, "--timeout-stop=") {
		if timeout := strings.TrimPrefix(option, "--timeout-stop="); timeout != "" {
			if _, err := fmt.Sscanf(timeout, "%d", &cfg.TimeoutStop); err != nil {
				return fmt.Errorf("invalid timeout-stop value: %s", timeout)
			}
		}
	} else if strings.HasPrefix(option, "--after=") {
		service := strings.TrimPrefix(option, "--after=")
		cfg.After = append(cfg.After, service)
	} else if strings.HasPrefix(option, "--wants=") {
		service := strings.TrimPrefix(option, "--wants=")
		cfg.Wants = append(cfg.Wants, service)
	} else if strings.HasPrefix(option, "--requires=") {
		service := strings.TrimPrefix(option, "--requires=")
		cfg.Requires = append(cfg.Requires, service)
	} else {
		return fmt.Errorf("unknown option: %s", option)
	}
	return nil
}

func validateConfig(cfg *config.ServiceConfig) error {
	if cfg.WorkDir != "" {
		if !filepath.IsAbs(cfg.WorkDir) {
			abs, err := filepath.Abs(cfg.WorkDir)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for workdir: %v", err)
			}
			cfg.WorkDir = abs
		}
		if _, err := os.Stat(cfg.WorkDir); os.IsNotExist(err) {
			return fmt.Errorf("working directory does not exist: %s", cfg.WorkDir)
		}
	}
	if cfg.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if strings.Contains(cfg.Name, " ") {
		return fmt.Errorf("service name cannot contain spaces")
	}
	if cfg.ExecStart == "" {
		return fmt.Errorf("exec start command cannot be empty")
	}
	return nil
}

func BuildServiceContent(cfg *config.ServiceConfig, serviceName string) string {
	var content strings.Builder
	content.WriteString("[Unit]\n")
	content.WriteString(fmt.Sprintf("Description=%s\n", cfg.Description))
	if len(cfg.After) > 0 {
		content.WriteString(fmt.Sprintf("After=%s\n", strings.Join(cfg.After, " ")))
	}
	if len(cfg.Wants) > 0 {
		content.WriteString(fmt.Sprintf("Wants=%s\n", strings.Join(cfg.Wants, " ")))
	}
	if len(cfg.Requires) > 0 {
		content.WriteString(fmt.Sprintf("Requires=%s\n", strings.Join(cfg.Requires, " ")))
	}
	content.WriteString("\n")
	content.WriteString("[Service]\n")
	content.WriteString("Type=simple\n")
	content.WriteString(fmt.Sprintf("User=%s\n", cfg.User))
	if cfg.Group != "" {
		content.WriteString(fmt.Sprintf("Group=%s\n", cfg.Group))
	}
	if cfg.WorkDir != "" {
		content.WriteString(fmt.Sprintf("WorkingDirectory=%s\n", cfg.WorkDir))
	}
	content.WriteString(fmt.Sprintf("ExecStart=%s\n", cfg.ExecStart))
	content.WriteString(fmt.Sprintf("Restart=%s\n", cfg.Restart))
	content.WriteString(fmt.Sprintf("RestartSec=%d\n", cfg.RestartSec))
	content.WriteString(fmt.Sprintf("KillMode=%s\n", cfg.KillMode))
	content.WriteString(fmt.Sprintf("KillSignal=%s\n", cfg.KillSignal))
	content.WriteString(fmt.Sprintf("TimeoutStartSec=%d\n", cfg.TimeoutStart))
	content.WriteString(fmt.Sprintf("TimeoutStopSec=%d\n", cfg.TimeoutStop))
	content.WriteString("StandardOutput=journal\n")
	content.WriteString("StandardError=journal\n")
	content.WriteString(fmt.Sprintf("SyslogIdentifier=%s\n", serviceName))
	for key, value := range cfg.Env {
		content.WriteString(fmt.Sprintf("Environment=%s=%s\n", key, value))
	}
	content.WriteString("\n")
	content.WriteString("[Install]\n")
	content.WriteString("WantedBy=multi-user.target\n")
	return content.String()
}
