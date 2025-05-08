package jpid

import (
	"context"
	"fmt"
	"omniscient/internal/model/entity"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"omniscient/api/jpid/v1"
)

func (c *ControllerV1) Online(ctx context.Context, req *v1.OnlineReq) (res *v1.OnlineRes, err error) {
	// 使用ps命令查找Java进程
	cmd := exec.Command("bash", "-c", "ps -ef | grep java")
	output, err := cmd.Output()
	if err != nil {
		return nil, gerror.New("执行命令出错: " + err.Error())
	}

	// 将输出转换为字符串并按行分割
	processes := strings.Split(string(output), "\n")

	res = &v1.OnlineRes{
		List: make([]*entity.LinuxPid, 0),
	}

	// 遍历并解析每个进程信息
	for _, process := range processes {
		// 跳过空行
		if process == "" {
			continue
		}

		// 跳过grep进程本身
		if strings.Contains(process, "grep java") {
			continue
		}

		// 分割进程信息
		fields := strings.Fields(process)
		if len(fields) >= 8 {
			pid, err := strconv.Atoi(fields[1])
			if err != nil {
				continue // 跳过无效的PID
			}

			// 获取完整命令
			command := strings.Join(fields[7:], " ")

			// 尝试从命令中提取项目名称
			name := extractJavaProjectName(command)

			// 跳过无法识别项目名称的进程
			if name == "unknown" {
				continue
			}

			// 获取进程占用的TCP端口
			ports := getTCPPorts(pid)

			// 从命令行参数中提取可能的端口号
			cmdLinePorts := extractPortFromCommand(command)

			// 合并两种方式获取的端口号
			allPorts := mergePorts(ports, cmdLinePorts)

			// 获取jar文件所在目录
			jarDir := getJarDirectory(pid, command)

			linuxPid := &entity.LinuxPid{
				Name:    name,
				Pid:     pid,
				Run:     command,
				Ports:   allPorts,
				Catalog: jarDir,
			}

			res.List = append(res.List, linuxPid)
		}
	}

	return res, nil
}

// 获取进程占用的TCP端口
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

	var ports []string
	portMap := make(map[string]bool) // 用于去重

	// 处理命令输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		// 使用正则表达式匹配本地地址的端口号
		// 匹配格式如 ":8080" 或 "0.0.0.0:8080" 或 "*:8080"
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

// 从命令行中提取端口号
func extractPortFromCommand(command string) []string {
	var ports []string
	portMap := make(map[string]bool) // 用于去重

	// 常见的端口参数模式
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

// 合并端口列表并去重
func mergePorts(ports1, ports2 []string) []string {
	portMap := make(map[string]bool)
	var result []string

	// 添加所有端口到map中以去重
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

// extractJavaProjectName 从Java命令中提取项目名称
func extractJavaProjectName(command string) string {
	// 尝试从-jar参数后面提取名称
	if strings.Contains(command, "-jar") {
		parts := strings.Split(command, "-jar")
		if len(parts) > 1 {
			// 获取-jar之后的部分
			remaining := parts[1]
			// 跳过JVM参数
			for _, part := range strings.Fields(remaining) {
				if strings.HasSuffix(part, ".jar") {
					// 从完整路径中提取文件名
					fileName := part[strings.LastIndex(part, "/")+1:]
					return fileName // 直接返回带.jar后缀的文件名
				}
			}
		}
	}
	return "unknown"
}

// getJarDirectory 获取jar文件所在目录
func getJarDirectory(pid int, command string) string {
	// 从命令行提取jar文件路径
	jarPath := extractJarPath(command)
	if jarPath == "" {
		return ""
	}

	// 如果是相对路径，尝试获取进程的工作目录
	if !filepath.IsAbs(jarPath) {
		// 尝试通过pwdx命令获取进程的工作目录
		cmd := exec.Command("bash", "-c", fmt.Sprintf("pwdx %d | cut -d: -f2", pid))
		output, err := cmd.Output()
		if err == nil {
			workDir := strings.TrimSpace(string(output))
			jarPath = filepath.Join(workDir, jarPath)
		} else {
			// 如果pwdx失败，尝试通过/proc/pid/cwd获取
			cmd = exec.Command("bash", "-c", fmt.Sprintf("readlink /proc/%d/cwd", pid))
			output, err = cmd.Output()
			if err == nil {
				workDir := strings.TrimSpace(string(output))
				jarPath = filepath.Join(workDir, jarPath)
			}
		}
	}

	// 获取jar文件的目录
	return filepath.Dir(jarPath)
}

// extractJarPath 从命令行中提取jar文件路径
func extractJarPath(command string) string {
	if strings.Contains(command, "-jar") {
		parts := strings.Split(command, "-jar")
		if len(parts) > 1 {
			// 获取-jar之后的部分
			remaining := parts[1]
			// 跳过JVM参数
			for _, part := range strings.Fields(remaining) {
				if strings.HasSuffix(part, ".jar") {
					return part
				}
			}
		}
	}
	return ""
}
