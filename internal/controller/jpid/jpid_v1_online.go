package jpid

import (
	"context"
	"omniscient/internal/model/entity"
	"os/exec"
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

			linuxPid := &entity.LinuxPid{
				Name: name,
				Pid:  pid,
				Run:  command,
			}

			res.List = append(res.List, linuxPid)
		}
	}

	return res, nil
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
					// 移除.jar扩展名和版本号
					baseName := strings.TrimSuffix(fileName, ".jar")
					// 移除版本号部分（如 -1.0-SNAPSHOT）
					if idx := strings.LastIndex(baseName, "-"); idx != -1 {
						// 检查是否是版本号
						versionPart := baseName[idx:]
						if strings.Contains(versionPart, ".") || strings.Contains(versionPart, "SNAPSHOT") {
							baseName = baseName[:idx]
						}
					}
					return baseName
				}
			}
		}
	}
	return "unknown"
}
