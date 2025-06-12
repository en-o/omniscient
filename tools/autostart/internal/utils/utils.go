package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"autostart/internal/config"
)

// NeedsRoot 检查命令是否需要root权限
func NeedsRoot(command string) bool {
	rootCommands := map[string]bool{
		"add": true, "create": true, "install": true,
		"remove": true, "rm": true, "delete": true, "uninstall": true,
		"enable": true, "disable": true,
		"start": true, "stop": true, "restart": true,
		"edit":             true,
		"exists":           true,
		"check":            true,
		"install-global":   true,
		"uninstall-global": true,
	}
	return rootCommands[command]
}

// PrintHelp 打印帮助信息
func PrintHelp() {
	const version = "0.0.2"
	fmt.Printf("autostart v%s - Universal Service Autostart Management Tool\n\n", version)
	fmt.Printf("true systemctl service name: autostart-<name>\n\n")

	fmt.Println("USAGE:")
	fmt.Println("  autostart <command> [arguments]\n")

	fmt.Println("COMMANDS:")
	fmt.Println("  Service Management:")
	fmt.Println("    list, ls                                    - List all autostart services")
	fmt.Println("    add, create, install <name> <exec> [opts]  - Add service to autostart")
	fmt.Println("    remove, rm, delete, uninstall <name>       - Remove service from autostart")
	fmt.Println("    edit <name>                                 - Edit service configuration")
	fmt.Println("")

	fmt.Println("  Service Control:")
	fmt.Println("    enable <name>                               - Enable service autostart")
	fmt.Println("    disable <name>                              - Disable service autostart")
	fmt.Println("    start <name>                                - Start service")
	fmt.Println("    stop <name>                                 - Stop service")
	fmt.Println("    restart <name>                              - Restart service")
	fmt.Println("    status <name>                               - Show service status")
	fmt.Println("    logs <name> [lines]                         - Show service logs")
	fmt.Println("")
	fmt.Println("  Service Query:")
	fmt.Println("    exists, check <name>                      - Check if service exists")
	fmt.Println("")
	fmt.Println("  Tool Management:")
	fmt.Println("    install-global                            - Install to global environment")
	fmt.Println("    uninstall-global                          - Uninstall from global environment")
	fmt.Println("    version, -v, --version                      - Show version")
	fmt.Println("    help, -h, --help                            - Show this help")
	fmt.Println("")

	printAddOptions()
	printExamples()
}

// printAddOptions 打印添加选项帮助
func printAddOptions() {
	fmt.Println("ADD OPTIONS:")
	options := []string{
		"--workdir=<path>          - Working directory (default: auto-detect)",
		"--user=<username>         - Run as user (default: root)",
		"--group=<groupname>       - Run as group (default: user's primary group)",
		"--description=<text>      - Service description",
		"--env=<KEY=VALUE>         - Environment variable (repeatable)",
		"--restart=<policy>        - Restart policy: always|on-failure|no (default: always)",
		"--restart-sec=<seconds>   - Restart delay in seconds (default: 5)",
		"--kill-mode=<mode>        - Kill mode: control-group|process|mixed (default: control-group)",
		"--kill-signal=<signal>    - Kill signal (default: SIGTERM)",
		"--timeout-start=<seconds> - Start timeout (default: 90)",
		"--timeout-stop=<seconds>  - Stop timeout (default: 90)",
		"--after=<service>         - Start after service (repeatable)",
		"--wants=<service>         - Wants service (repeatable)",
		"--requires=<service>      - Requires service (repeatable)",
	}

	for _, opt := range options {
		fmt.Printf("  %s\n", opt)
	}
	fmt.Println("")
}

// printExamples 打印示例
func printExamples() {
	fmt.Println("EXAMPLES:")
	examples := []struct {
		desc string
		cmds []string
	}{
		{
			"Add and start a Java application",
			[]string{
				`sudo autostart add myapp "java -jar /path/to/app.jar" --workdir=/path/to`,
				`sudo autostart enable myapp`,
				`sudo autostart start myapp`,
			},
		},
		{
			"Add a Python application with custom user and environment",
			[]string{
				`sudo autostart add pyapp "python3 /path/to/app.py" --user=www-data --env=PYTHONPATH=/path/to --env=DEBUG=true`,
			},
		},
		{
			"Add a Node.js application with custom restart policy",
			[]string{
				`sudo autostart add nodeapp "node /path/to/app.js" --restart=on-failure --restart-sec=10`,
			},
		},
		{
			"Add a service with dependencies",
			[]string{
				`sudo autostart add webapp "./webapp" --after=network.target --after=mysql.service --requires=mysql.service`,
			},
		},
		{
			"Common workflow",
			[]string{
				`autostart list                # List all services`,
				`autostart status myapp        # Check service status`,
				`autostart logs myapp 100      # View recent logs`,
				`sudo autostart restart myapp  # Restart service`,
			},
		},
	}

	for _, example := range examples {
		fmt.Printf("  # %s\n", example.desc)
		for _, cmd := range example.cmds {
			fmt.Printf("  %s\n", cmd)
		}
		fmt.Println("")
	}
}

// ParseAddOptions 解析添加选项
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

// inferWorkingDirectory 推断工作目录
func inferWorkingDirectory(execStart string) string {
	parts := strings.Fields(execStart)
	if len(parts) == 0 {
		return getCurrentDir()
	}

	for _, part := range parts {
		if strings.Contains(part, "/") && !strings.HasPrefix(part, "-") && filepath.IsAbs(part) {
			if _, err := os.Stat(part); err == nil {
				return filepath.Dir(part)
			}
		}
	}

	return getCurrentDir()
}

// getCurrentDir 获取当前目录
func getCurrentDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "/tmp"
}

// parseOption 解析单个选项
func parseOption(cfg *config.ServiceConfig, option string) error {
	if !strings.HasPrefix(option, "--") {
		return fmt.Errorf("invalid option format: %s (must start with --)", option)
	}

	key, value, err := parseOptionKeyValue(option)
	if err != nil {
		return err
	}

	return applyOptionToConfig(cfg, key, value)
}

// parseOptionKeyValue 解析选项的键值对
func parseOptionKeyValue(option string) (string, string, error) {
	parts := strings.SplitN(option[2:], "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid option format: %s (use --key=value)", option)
	}

	key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if key == "" {
		return "", "", fmt.Errorf("empty option key in: %s", option)
	}
	if value == "" {
		return "", "", fmt.Errorf("empty option value in: %s", option)
	}

	return key, value, nil
}

// optionHandler 定义选项处理函数类型
type optionHandler func(*config.ServiceConfig, string) error

// applyOptionToConfig 将选项应用到配置中
func applyOptionToConfig(cfg *config.ServiceConfig, key, value string) error {
	// 选项处理映射表
	optionHandlers := map[string]optionHandler{
		"workdir":       handleWorkdirOption,
		"user":          handleUserOption,
		"group":         handleGroupOption,
		"description":   handleDescriptionOption,
		"env":           handleEnvOption,
		"restart":       handleRestartOption,
		"restart-sec":   handleRestartSecOption,
		"kill-mode":     handleKillModeOption,
		"kill-signal":   handleKillSignalOption,
		"timeout-start": handleTimeoutStartOption,
		"timeout-stop":  handleTimeoutStopOption,
		"after":         handleAfterOption,
		"wants":         handleWantsOption,
		"requires":      handleRequiresOption,
	}

	handler, exists := optionHandlers[key]
	if !exists {
		return fmt.Errorf("unknown option: --%s", key)
	}

	return handler(cfg, value)
}

// 各种选项处理函数
func handleWorkdirOption(cfg *config.ServiceConfig, value string) error {
	cfg.WorkDir = value
	return nil
}

func handleUserOption(cfg *config.ServiceConfig, value string) error {
	if !isValidUsername(value) {
		return fmt.Errorf("invalid username: %s", value)
	}
	cfg.User = value
	return nil
}

func handleGroupOption(cfg *config.ServiceConfig, value string) error {
	if !isValidGroupname(value) {
		return fmt.Errorf("invalid group name: %s", value)
	}
	cfg.Group = value
	return nil
}

func handleDescriptionOption(cfg *config.ServiceConfig, value string) error {
	if len(value) > 256 {
		return fmt.Errorf("description too long (max 256 characters): %d", len(value))
	}
	cfg.Description = value
	return nil
}

func handleEnvOption(cfg *config.ServiceConfig, value string) error {
	return parseAndSetEnvVar(cfg, value)
}

func handleRestartOption(cfg *config.ServiceConfig, value string) error {
	return setRestartPolicy(cfg, value)
}

func handleRestartSecOption(cfg *config.ServiceConfig, value string) error {
	return parseAndSetIntOption(&cfg.RestartSec, value, "restart-sec", 1, 3600)
}

func handleKillModeOption(cfg *config.ServiceConfig, value string) error {
	validModes := []string{"control-group", "process", "mixed", "none"}
	if !contains(validModes, value) {
		return fmt.Errorf("invalid kill mode: %s (valid: %s)", value, strings.Join(validModes, ", "))
	}
	cfg.KillMode = value
	return nil
}

func handleKillSignalOption(cfg *config.ServiceConfig, value string) error {
	validSignals := []string{"SIGTERM", "SIGKILL", "SIGINT", "SIGHUP", "SIGQUIT", "SIGUSR1", "SIGUSR2"}
	if !contains(validSignals, value) {
		return fmt.Errorf("invalid kill signal: %s (valid: %s)", value, strings.Join(validSignals, ", "))
	}
	cfg.KillSignal = value
	return nil
}

func handleTimeoutStartOption(cfg *config.ServiceConfig, value string) error {
	return parseAndSetIntOption(&cfg.TimeoutStart, value, "timeout-start", 1, 600)
}

func handleTimeoutStopOption(cfg *config.ServiceConfig, value string) error {
	return parseAndSetIntOption(&cfg.TimeoutStop, value, "timeout-stop", 1, 600)
}

func handleAfterOption(cfg *config.ServiceConfig, value string) error {
	if !isValidServiceName(value) {
		return fmt.Errorf("invalid service name for after: %s", value)
	}
	cfg.After = appendUnique(cfg.After, value)
	return nil
}

func handleWantsOption(cfg *config.ServiceConfig, value string) error {
	if !isValidServiceName(value) {
		return fmt.Errorf("invalid service name for wants: %s", value)
	}
	cfg.Wants = appendUnique(cfg.Wants, value)
	return nil
}

func handleRequiresOption(cfg *config.ServiceConfig, value string) error {
	if !isValidServiceName(value) {
		return fmt.Errorf("invalid service name for requires: %s", value)
	}
	cfg.Requires = appendUnique(cfg.Requires, value)
	return nil
}

// parseAndSetEnvVar 解析并设置环境变量
func parseAndSetEnvVar(cfg *config.ServiceConfig, envStr string) error {
	key, value, err := parseEnvVarString(envStr)
	if err != nil {
		return err
	}

	if !isValidEnvVarName(key) {
		return fmt.Errorf("invalid environment variable name: %s", key)
	}

	cfg.Env[key] = value
	return nil
}

// parseEnvVarString 解析环境变量字符串
func parseEnvVarString(envStr string) (string, string, error) {
	parts := strings.SplitN(envStr, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid environment variable format: %s (use KEY=VALUE)", envStr)
	}

	key, value := strings.TrimSpace(parts[0]), parts[1] // 值可能包含空格，不要trim
	if key == "" {
		return "", "", fmt.Errorf("empty environment variable name in: %s", envStr)
	}

	return key, value, nil
}

// setRestartPolicy 设置重启策略
func setRestartPolicy(cfg *config.ServiceConfig, restart string) error {
	validPolicies := []string{"always", "on-failure", "on-abnormal", "on-watchdog", "on-abort", "no"}
	if !contains(validPolicies, restart) {
		return fmt.Errorf("invalid restart policy: %s (valid: %s)", restart, strings.Join(validPolicies, ", "))
	}
	cfg.Restart = restart
	return nil
}

// parseAndSetIntOption 解析并设置整数选项（带范围检查）
func parseAndSetIntOption(target *int, value, optName string, min, max int) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid %s value: %s (must be an integer)", optName, value)
	}

	if val < min || val > max {
		return fmt.Errorf("%s value out of range: %d (valid range: %d-%d)", optName, val, min, max)
	}

	*target = val
	return nil
}

// 验证配置
func validateConfig(cfg *config.ServiceConfig) error {
	validators := []func(*config.ServiceConfig) error{
		validateServiceName,
		validateExecStart,
		validateWorkDir,
		validateUser,
		validateTimeouts,
		validateDependencies,
	}

	for _, validator := range validators {
		if err := validator(cfg); err != nil {
			return err
		}
	}

	return nil
}

// validateServiceName 验证服务名称
func validateServiceName(cfg *config.ServiceConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	if !isValidServiceName(cfg.Name) {
		return fmt.Errorf("invalid service name: %s (only letters, numbers, hyphens, and underscores allowed)", cfg.Name)
	}

	if len(cfg.Name) > 64 {
		return fmt.Errorf("service name too long: %s (max 64 characters)", cfg.Name)
	}

	return nil
}

// validateExecStart 验证执行命令
func validateExecStart(cfg *config.ServiceConfig) error {
	if cfg.ExecStart == "" {
		return fmt.Errorf("exec start command cannot be empty")
	}

	if len(cfg.ExecStart) > 1024 {
		return fmt.Errorf("exec start command too long (max 1024 characters)")
	}

	return nil
}

// validateWorkDir 验证工作目录
func validateWorkDir(cfg *config.ServiceConfig) error {
	if cfg.WorkDir == "" {
		return nil // 工作目录可以为空
	}

	// 转换为绝对路径
	if !filepath.IsAbs(cfg.WorkDir) {
		abs, err := filepath.Abs(cfg.WorkDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for workdir: %v", err)
		}
		cfg.WorkDir = abs
	}

	// 检查目录是否存在
	if stat, err := os.Stat(cfg.WorkDir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("working directory does not exist: %s", cfg.WorkDir)
		}
		return fmt.Errorf("failed to access working directory: %v", err)
	} else if !stat.IsDir() {
		return fmt.Errorf("working directory is not a directory: %s", cfg.WorkDir)
	}

	return nil
}

// validateUser 验证用户配置
func validateUser(cfg *config.ServiceConfig) error {
	if cfg.User == "" {
		cfg.User = "root" // 设置默认用户
		return nil
	}

	if !isValidUsername(cfg.User) {
		return fmt.Errorf("invalid username: %s", cfg.User)
	}

	return nil
}

// validateTimeouts 验证超时配置
func validateTimeouts(cfg *config.ServiceConfig) error {
	if cfg.TimeoutStart <= 0 {
		return fmt.Errorf("timeout-start must be positive: %d", cfg.TimeoutStart)
	}

	if cfg.TimeoutStop <= 0 {
		return fmt.Errorf("timeout-stop must be positive: %d", cfg.TimeoutStop)
	}

	if cfg.RestartSec < 0 {
		return fmt.Errorf("restart-sec cannot be negative: %d", cfg.RestartSec)
	}

	return nil
}

// validateDependencies 验证依赖配置
func validateDependencies(cfg *config.ServiceConfig) error {
	allDeps := append(append(cfg.After, cfg.Wants...), cfg.Requires...)

	for _, dep := range allDeps {
		if !isValidServiceName(dep) {
			return fmt.Errorf("invalid dependency service name: %s", dep)
		}
	}

	return nil
}

// 构建服务文件内容
func BuildServiceContent(cfg *config.ServiceConfig, serviceName string) string {
	builder := &serviceContentBuilder{cfg: cfg, serviceName: serviceName}
	return builder.build()
}

// serviceContentBuilder 服务内容构建器
type serviceContentBuilder struct {
	cfg         *config.ServiceConfig
	serviceName string
	content     strings.Builder
}

// build 构建服务文件内容
func (b *serviceContentBuilder) build() string {
	b.buildUnitSection()
	b.buildServiceSection()
	b.buildInstallSection()
	return b.content.String()
}

// buildUnitSection 构建 Unit 段
func (b *serviceContentBuilder) buildUnitSection() {
	b.content.WriteString("[Unit]\n")
	b.content.WriteString(fmt.Sprintf("Description=%s\n", b.cfg.Description))

	b.writeServiceList("After", b.cfg.After)
	b.writeServiceList("Wants", b.cfg.Wants)
	b.writeServiceList("Requires", b.cfg.Requires)

	b.content.WriteString("\n")
}

// buildServiceSection 构建 Service 段
func (b *serviceContentBuilder) buildServiceSection() {
	b.content.WriteString("[Service]\n")
	b.content.WriteString("Type=simple\n")
	b.content.WriteString(fmt.Sprintf("User=%s\n", b.cfg.User))

	if b.cfg.Group != "" {
		b.content.WriteString(fmt.Sprintf("Group=%s\n", b.cfg.Group))
	}

	if b.cfg.WorkDir != "" {
		b.content.WriteString(fmt.Sprintf("WorkingDirectory=%s\n", b.cfg.WorkDir))
	}

	b.content.WriteString(fmt.Sprintf("ExecStart=%s\n", b.normalizeExecStart(b.cfg.ExecStart, b.cfg.WorkDir)))
	b.content.WriteString(fmt.Sprintf("Restart=%s\n", b.cfg.Restart))
	b.content.WriteString(fmt.Sprintf("RestartSec=%d\n", b.cfg.RestartSec))
	b.content.WriteString(fmt.Sprintf("KillMode=%s\n", b.cfg.KillMode))
	b.content.WriteString(fmt.Sprintf("KillSignal=%s\n", b.cfg.KillSignal))
	b.content.WriteString(fmt.Sprintf("TimeoutStartSec=%d\n", b.cfg.TimeoutStart))
	b.content.WriteString(fmt.Sprintf("TimeoutStopSec=%d\n", b.cfg.TimeoutStop))
	b.content.WriteString("StandardOutput=journal\n")
	b.content.WriteString("StandardError=journal\n")
	b.content.WriteString(fmt.Sprintf("SyslogIdentifier=%s\n", b.cfg.Name))

	b.writeEnvironmentVariables()
	b.content.WriteString("\n")
}

// buildInstallSection 构建 Install 段
func (b *serviceContentBuilder) buildInstallSection() {
	b.content.WriteString("[Install]\n")
	b.content.WriteString("WantedBy=multi-user.target\n")
}

// writeServiceList 写入服务列表
func (b *serviceContentBuilder) writeServiceList(key string, services []string) {
	if len(services) > 0 {
		b.content.WriteString(fmt.Sprintf("%s=%s\n", key, strings.Join(services, " ")))
	}
}

// writeEnvironmentVariables 写入环境变量
func (b *serviceContentBuilder) writeEnvironmentVariables() {
	for key, value := range b.cfg.Env {
		// 对包含特殊字符的值进行引用
		if needsQuoting(value) {
			value = fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\\\""))
		}
		b.content.WriteString(fmt.Sprintf("Environment=%s=%s\n", key, value))
	}
}

// 工具函数

// contains 检查字符串切片是否包含指定值
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// appendUnique 向切片添加唯一元素
func appendUnique(slice []string, item string) []string {
	if !contains(slice, item) {
		return append(slice, item)
	}
	return slice
}

// isValidUsername 验证用户名格式
func isValidUsername(username string) bool {
	if len(username) == 0 || len(username) > 32 {
		return false
	}

	// 用户名只能包含字母、数字、下划线和连字符，且不能以数字开头
	for i, r := range username {
		if i == 0 && (r >= '0' && r <= '9') {
			return false
		}
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-') {
			return false
		}
	}

	return true
}

// isValidGroupname 验证组名格式
func isValidGroupname(groupname string) bool {
	return isValidUsername(groupname) // 组名和用户名使用相同的验证规则
}

// isValidServiceName 验证服务名格式
func isValidServiceName(serviceName string) bool {
	if len(serviceName) == 0 || len(serviceName) > 64 {
		return false
	}

	// 服务名只能包含字母、数字、下划线、连字符和点
	for _, r := range serviceName {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.') {
			return false
		}
	}

	return true
}

// isValidEnvVarName 验证环境变量名格式
func isValidEnvVarName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// 环境变量名只能包含字母、数字和下划线，且不能以数字开头
	for i, r := range name {
		if i == 0 && (r >= '0' && r <= '9') {
			return false
		}
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}

	return true
}

// needsQuoting 检查值是否需要引用
func needsQuoting(value string) bool {
	return strings.ContainsAny(value, " \t\n\"'\\$`|&;<>(){}[]?*~")
}

// RedHatBased 判断是否是 RedHat 系列系统
func RedHatBased() bool {
	_, err := os.Stat("/etc/redhat-release")
	return err == nil
}

// 结合工作目录生成绝对路径
func (b *serviceContentBuilder) normalizeExecStart(execStart, workDir string) string {
	parts := strings.Fields(execStart)
	if len(parts) == 0 {
		return execStart
	}

	// 获取第一个参数（可执行文件）
	executable := parts[0]

	// 如果已经是绝对路径，直接返回
	if filepath.IsAbs(executable) {
		return execStart
	}

	// 处理常见的可执行文件
	var absoluteExec string

	switch {
	case executable == "sh" || strings.HasSuffix(executable, ".sh"):
		// 处理 shell 脚本
		if executable == "sh" && len(parts) > 1 {
			// 如果是 "sh script.sh" 格式
			scriptPath := parts[1]
			if !filepath.IsAbs(scriptPath) {
				if workDir != "" {
					scriptPath = filepath.Join(workDir, scriptPath)
				} else {
					// 尝试获取脚本的绝对路径
					if abs, err := filepath.Abs(scriptPath); err == nil {
						scriptPath = abs
					}
				}
			}
			absoluteExec = "/bin/sh"
			parts[1] = scriptPath
		} else {
			// 直接执行 .sh 文件
			if workDir != "" {
				absoluteExec = filepath.Join(workDir, executable)
			} else {
				if abs, err := filepath.Abs(executable); err == nil {
					absoluteExec = abs
				} else {
					absoluteExec = executable
				}
			}
		}

	case executable == "java":
		// Java 命令使用系统路径
		if javaPath, err := exec.LookPath("java"); err == nil {
			absoluteExec = javaPath
		} else {
			absoluteExec = "/usr/bin/java"
		}

	case executable == "python" || executable == "python3":
		// Python 命令使用系统路径
		if pythonPath, err := exec.LookPath(executable); err == nil {
			absoluteExec = pythonPath
		} else {
			absoluteExec = "/usr/bin/" + executable
		}

	case executable == "node":
		// Node.js 命令使用系统路径
		if nodePath, err := exec.LookPath("node"); err == nil {
			absoluteExec = nodePath
		} else {
			absoluteExec = "/usr/bin/node"
		}

	default:
		// 其他情况，尝试查找可执行文件
		if execPath, err := exec.LookPath(executable); err == nil {
			absoluteExec = execPath
		} else {
			// 如果找不到，尝试相对于工作目录
			if workDir != "" {
				absoluteExec = filepath.Join(workDir, executable)
			} else {
				if abs, err := filepath.Abs(executable); err == nil {
					absoluteExec = abs
				} else {
					absoluteExec = executable
				}
			}
		}
	}

	// 重建命令
	parts[0] = absoluteExec
	return strings.Join(parts, " ")
}
