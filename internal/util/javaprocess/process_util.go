package javaprocess

import (
	"fmt"
	"io/ioutil"
	"omniscient/internal/model/entity"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
)

// JavaProcessInfo holds information about a Java process
type JavaProcessInfo struct {
	Name     string
	Pid      int
	Command  string
	Ports    []string
	JarDir   string
	IsDocker bool
}

// GetJavaProcesses returns a list of running Java processes
func GetJavaProcesses() ([]*entity.LinuxPid, error) {
	// 使用ps命令查找Java进程
	cmd := exec.Command("bash", "-c", "ps -ef | grep java")
	output, err := cmd.Output()
	if err != nil {
		return nil, gerror.New("执行命令出错: " + err.Error())
	}

	// 将输出转换为字符串并按行分割
	processes := strings.Split(string(output), "\n")
	result := make([]*entity.LinuxPid, 0)

	// 遍历并解析每个进程信息
	for _, process := range processes {
		// 跳过空行和grep进程本身
		if process == "" || strings.Contains(process, "grep java") {
			continue
		}

		if info := parseJavaProcess(process); info != nil {
			// 只保留包含-jar参数的真正Java项目进程
			if !containsJarParam(info.Command) {
				continue
			}

			// 非Docker进程需要检查是否有TCP连接
			if !info.IsDocker {
				if len(info.Ports) == 0 {
					// 跳过没有TCP连接的非Docker进程
					continue
				}
			}

			// Docker值：1表示docker容器，2表示普通JDK进程
			dockerVal := 2
			if info.IsDocker {
				dockerVal = 1
			}

			linuxPid := &entity.LinuxPid{
				Name:    info.Name,
				Pid:     info.Pid,
				Run:     info.Command,
				Ports:   strings.Join(info.Ports, ","),
				Catalog: info.JarDir,
				Docker:  dockerVal,
			}
			result = append(result, linuxPid)
		}
	}

	return result, nil
}

// 检查命令行是否包含-jar参数
func containsJarParam(command string) bool {
	return strings.Contains(command, "-jar")
}

// parseJavaProcess parses a single Java process line and returns process information
func parseJavaProcess(processLine string) *JavaProcessInfo {
	fields := strings.Fields(processLine)
	if len(fields) < 8 {
		return nil
	}

	pid, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil
	}

	command := strings.Join(fields[7:], " ")
	name := extractJavaProjectName(command)

	if name == "unknown" {
		return nil
	}

	// 检查是否为Docker中的进程
	isDocker := checkIfDockerProcess(pid)

	ports := getTCPPorts(pid)
	cmdLinePorts := extractPortFromCommand(command)
	allPorts := mergePorts(ports, cmdLinePorts)
	jarDir := getJarDirectory(pid, command)

	return &JavaProcessInfo{
		Name:     name,
		Pid:      pid,
		Command:  command,
		Ports:    allPorts,
		JarDir:   jarDir,
		IsDocker: isDocker,
	}
}

// checkIfDockerProcess checks if a process is running inside a Docker container
func checkIfDockerProcess(pid int) bool {
	cgroupPath := fmt.Sprintf("/proc/%d/cgroup", pid)

	content, err := ioutil.ReadFile(cgroupPath)
	if err != nil {
		return false
	}

	// 检查cgroup内容是否包含docker标识
	return strings.Contains(string(content), "docker")
}

// 已移除不需要的函数

// GetTCPPorts returns all TCP ports used by a process
func getTCPPorts(pid int) []string {
	// 首先尝试使用ss命令，只获取TCP端口
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ss -tnlp | grep %d", pid))
	output, err := cmd.Output()

	// 如果ss命令失败，尝试使用netstat
	if err != nil {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("netstat -tnlp | grep %d", pid))
		output, err = cmd.Output()
		if err != nil {
			return nil
		}
	}

	return parsePortsFromOutput(string(output))
}

// parsePortsFromOutput parses ports from command output
func parsePortsFromOutput(output string) []string {
	var ports []string
	portMap := make(map[string]bool)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		re := regexp.MustCompile(`[:\*](\d+)\s+(?:.*?)\s+LISTEN`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			port := matches[1]
			if !portMap[port] {
				ports = append(ports, port)
				portMap[port] = true
			}
		}
	}

	return ports
}

// ExtractPortFromCommand extracts port numbers from Java command line arguments
func extractPortFromCommand(command string) []string {
	var ports []string
	portMap := make(map[string]bool)

	portPatterns := []string{
		`--server\.port=(\d+)`,
		`-Dserver\.port=(\d+)`,
		`port\s*=\s*(\d+)`,
		`-Dport=(\d+)`,
		`PORT=(\d+)`,
	}

	for _, pattern := range portPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(command)
		if len(matches) > 1 {
			port := matches[1]
			if !portMap[port] {
				ports = append(ports, port)
				portMap[port] = true
			}
		}
	}

	return ports
}

// MergePorts merges and deduplicates two slices of port numbers
func mergePorts(ports1, ports2 []string) []string {
	portMap := make(map[string]bool)
	var result []string

	for _, port := range ports1 {
		if !portMap[port] {
			portMap[port] = true
			result = append(result, port)
		}
	}

	for _, port := range ports2 {
		if !portMap[port] {
			portMap[port] = true
			result = append(result, port)
		}
	}

	return result
}

// ExtractJavaProjectName extracts the project name from a Java command
func extractJavaProjectName(command string) string {
	if strings.Contains(command, "-jar") {
		parts := strings.Split(command, "-jar")
		if len(parts) > 1 {
			remaining := parts[1]
			for _, part := range strings.Fields(remaining) {
				if strings.HasSuffix(part, ".jar") {
					return part[strings.LastIndex(part, "/")+1:]
				}
			}
		}
	}
	return "unknown"
}

// GetJarDirectory returns the directory containing the JAR file
func getJarDirectory(pid int, command string) string {
	jarPath := extractJarPath(command)
	if jarPath == "" {
		return ""
	}

	if !filepath.IsAbs(jarPath) {
		if workDir := getProcessWorkDir(pid); workDir != "" {
			jarPath = filepath.Join(workDir, jarPath)
		}
	}

	return filepath.Dir(jarPath)
}

// getProcessWorkDir returns the working directory of a process
func getProcessWorkDir(pid int) string {
	// Try pwdx first
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pwdx %d | cut -d: -f2", pid))
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	// Fall back to /proc/pid/cwd
	cmd = exec.Command("bash", "-c", fmt.Sprintf("readlink /proc/%d/cwd", pid))
	output, err = cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	return ""
}

// extractJarPath extracts the JAR file path from a Java command
func extractJarPath(command string) string {
	if strings.Contains(command, "-jar") {
		parts := strings.Split(command, "-jar")
		if len(parts) > 1 {
			remaining := parts[1]
			for _, part := range strings.Fields(remaining) {
				if strings.HasSuffix(part, ".jar") {
					return part
				}
			}
		}
	}
	return ""
}
