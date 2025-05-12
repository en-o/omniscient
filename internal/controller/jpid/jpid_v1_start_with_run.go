package jpid

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"io"
	"omniscient/api/jpid/v1"
	"omniscient/internal/service"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func (c *ControllerV1) StartWithRun(ctx context.Context, req *v1.StartWithRunReq) (res *v1.StartWithRunRes, err error) {
	// 获取响应写入器
	r := g.RequestFromCtx(ctx)
	w := r.Response.Writer

	// 设置 SSE 相关的 header
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取项目信息并进行验证
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		sendSSEMessage(w, "error", err.Error())
		return nil, err
	}
	if jpid == nil {
		sendSSEMessage(w, "error", "项目不存在")
		return nil, gerror.New("项目不存在")
	}
	if jpid.Run == "" {
		sendSSEMessage(w, "error", "run命令为空")
		return nil, gerror.New("run命令为空")
	}

	// 发送启动提示
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 正在启动项目: %s\x1b[0m", jpid.Name))
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;34m==> 工作目录: %s\x1b[0m", jpid.Catalog))
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;34m==> 执行命令: %s\x1b[0m", jpid.Run))
	sendSSEMessage(w, "output", "\x1b[1;33m==> 开始执行...\x1b[0m\n")

	var processPid int

	if req.Background {
		// 后台运行模式
		// 先删除已存在的 nohup.log
		nohupPath := fmt.Sprintf("%s/nohup.log", jpid.Catalog)
		if err = os.Remove(nohupPath); err != nil && !os.IsNotExist(err) {
			sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 警告: 清理旧日志文件失败: %v\x1b[0m", err))
		}

		// 生成临时PID文件路径
		pidFilePath := fmt.Sprintf("%s/.pid_temp", jpid.Catalog)
		if err = os.Remove(pidFilePath); err != nil && !os.IsNotExist(err) {
			// 忽略错误，只是清理
		}

		// 检测是否为Java命令
		isJavaCmd := strings.Contains(strings.ToLower(jpid.Run), "java ")

		// 创建启动脚本文件
		startScriptPath := fmt.Sprintf("%s/.start_script.sh", jpid.Catalog)
		var scriptContent string

		if isJavaCmd {
			// Java应用需要特殊处理来获取真正的PID
			scriptContent = fmt.Sprintf(`#!/bin/bash
cd "%s"
nohup %s > nohup.log 2>&1 &
SHELL_PID=$!
echo $SHELL_PID > "%s"

# 等待Java进程启动，最多等待10秒
MAX_ATTEMPTS=20
ATTEMPT=0
JAVA_PID=""

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
  # 查找所有与命令相关的Java进程
  POTENTIAL_PIDS=$(ps -eo pid,ppid,command | grep -v grep | grep '[j]ava' | awk '$2 == "'$SHELL_PID'" {print $1}')
  
  # 如果找不到作为SHELL_PID子进程的Java进程，则查找所有最近启动的匹配命令特征的Java进程
  if [ -z "$POTENTIAL_PIDS" ]; then
    # 提取命令中的关键特征词来匹配进程
    CMD_FEATURES=$(echo "%s" | tr ' ' '\n' | grep -v '^-' | grep -v '^java$' | head -3 | tr '\n' '|')
    if [ ! -z "$CMD_FEATURES" ]; then
      # 使用特征词查找匹配的Java进程，按启动时间排序，取最新的
      POTENTIAL_PIDS=$(ps -eo pid,etime,command | grep -v grep | grep '[j]ava' | grep -E "$CMD_FEATURES" | sort -k 2 | head -1 | awk '{print $1}')
    fi
  fi
  
  if [ ! -z "$POTENTIAL_PIDS" ]; then
    # 如果找到了多个PID，取第一个（通常是主进程）
    JAVA_PID=$(echo "$POTENTIAL_PIDS" | head -1)
    echo $JAVA_PID > "%s"
    break
  fi
  
  ATTEMPT=$((ATTEMPT+1))
  sleep 0.5
done

if [ -z "$JAVA_PID" ]; then
  echo "无法确定Java进程PID" >> nohup.log
else
  echo "已获取Java进程PID: $JAVA_PID" >> nohup.log
fi
`, jpid.Catalog, jpid.Run, pidFilePath, jpid.Run, pidFilePath)
		} else {
			// 非Java命令，直接使用bash创建的子进程PID
			scriptContent = fmt.Sprintf(`#!/bin/bash
cd "%s"
nohup %s > nohup.log 2>&1 &
PID=$!
echo $PID > "%s"
`, jpid.Catalog, jpid.Run, pidFilePath)
		}

		// 写入启动脚本
		if err = os.WriteFile(startScriptPath, []byte(scriptContent), 0755); err != nil {
			sendSSEMessage(w, "error", "创建启动脚本失败："+err.Error())
			return nil, gerror.Wrap(err, "创建启动脚本失败")
		}

		// 执行启动脚本
		cmd := exec.Command("bash", startScriptPath)
		cmd.Dir = jpid.Catalog
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("PROJECT_NAME=%s", jpid.Name),
			fmt.Sprintf("PROJECT_PID=%d", jpid.Pid),
			"LANG=en_US.UTF-8", // 添加UTF-8支持
		)

		if err = cmd.Run(); err != nil {
			sendSSEMessage(w, "error", "启动失败："+err.Error())
			// 清理启动脚本
			_ = os.Remove(startScriptPath)
			return nil, gerror.Wrap(err, "启动失败")
		}

		// 清理启动脚本
		_ = os.Remove(startScriptPath)

		// 等待PID文件写入完成，对于Java程序给予充足时间
		waitTime := 3 * time.Second
		if isJavaCmd {
			waitTime = 10 * time.Second
		}
		time.Sleep(waitTime)

		// 读取PID文件获取真实PID
		pidBytes, err := os.ReadFile(pidFilePath)
		if err != nil {
			sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 警告: 读取PID文件失败: %v，将无法追踪进程\x1b[0m", err))
		} else {
			pidStr := strings.TrimSpace(string(pidBytes))
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 警告: 解析PID失败: %v\x1b[0m", err))
			} else {
				processPid = pid
				// 验证PID是否有效
				process, err := os.FindProcess(processPid)
				if err != nil || process == nil {
					sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 警告: 获取到的PID %d 无效\x1b[0m", processPid))
					processPid = 0
				} else {
					// 在Unix系统上，验证进程是否真实存在
					err = process.Signal(syscall.Signal(0))
					if err != nil {
						sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 警告: PID %d 对应的进程不存在\x1b[0m", processPid))
						processPid = 0
					} else {
						// 验证PID是否对应预期的进程类型
						if isJavaCmd {
							// 检查PID是否为Java进程
							cmd := exec.Command("ps", "-p", pidStr, "-o", "command=")
							output, _ := cmd.Output()
							if !strings.Contains(strings.ToLower(string(output)), "java") {
								sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 警告: PID %d 不是Java进程\x1b[0m", processPid))
								// 不将processPid置零，使用找到的PID，即使不是Java进程
							}
						}
						sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 获取到进程 PID: %d\x1b[0m", processPid))
					}
				}
			}
		}

		// 如果无法通过PID文件获取有效PID，尝试使用ps命令查找进程
		if processPid == 0 {
			sendSSEMessage(w, "output", "\x1b[1;33m==> 尝试通过进程列表查找PID...\x1b[0m")

			// 提取命令的关键部分作为搜索特征
			cmdParts := strings.Split(jpid.Run, " ")
			var searchFeature string
			if len(cmdParts) > 0 {
				// 使用命令的第一个非空参数作为特征
				for _, part := range cmdParts {
					if part != "" && part != "nohup" && part != "&" {
						searchFeature = part
						break
					}
				}
			}

			if searchFeature != "" {
				// 使用ps命令查找最近启动的匹配进程
				cmd := exec.Command("bash", "-c", fmt.Sprintf("ps -eo pid,etime,command | grep -v grep | grep '%s' | sort -k 2 | head -1 | awk '{print $1}'", searchFeature))
				output, err := cmd.Output()
				if err == nil && len(output) > 0 {
					pidStr := strings.TrimSpace(string(output))
					pid, err := strconv.Atoi(pidStr)
					if err == nil && pid > 0 {
						processPid = pid
						sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 通过进程搜索找到PID: %d\x1b[0m", processPid))
					}
				}
			}
		}

		// 清理临时PID文件
		_ = os.Remove(pidFilePath)

		// 创建一个停止标记通道
		stopChan := make(chan bool)

		// 创建一个 goroutine 用于实时读取日志内容
		go func() {
			defer func() {
				if r := recover(); r != nil {
					g.Log().Error(ctx, "日志监控异常:", r)
				}
			}()

			// 日志文件可能不存在，持续等待直到它存在
			var logFile *os.File
			for attempts := 0; attempts < 30; attempts++ {
				var openErr error
				logFile, openErr = os.Open(nohupPath)
				if openErr == nil {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}

			if logFile == nil {
				sendSSEMessage(w, "output", "\x1b[1;31m==> 无法打开日志文件，放弃日志监控\x1b[0m")
				return
			}

			defer func() {
				if logFile != nil {
					_ = logFile.Close()
				}
			}()

			// 使用适当的buffer大小
			reader := bufio.NewReader(logFile)
			var line string
			var readErr error

			for {
				select {
				case <-stopChan:
					sendSSEMessage(w, "output", "\x1b[1;32m==> 日志监控已停止\x1b[0m")
					return
				default:
					line, readErr = reader.ReadString('\n')
					if readErr != nil && readErr != io.EOF {
						sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;31m==> 读取日志错误: %v\x1b[0m", readErr))
						time.Sleep(500 * time.Millisecond)
						continue
					}

					if line != "" {
						// 去除尾部的换行符
						line = strings.TrimRight(line, "\n")
						if w != nil {
							sendSSEMessage(w, "output", line)
						}
					}

					// 如果到达文件末尾，等待更多内容
					if readErr == io.EOF {
						// 检查进程是否还在运行
						if processPid > 0 {
							// 检查进程是否存在
							process, err := os.FindProcess(processPid)
							if err != nil || process == nil {
								sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;31m==> 进程 %d 不存在，停止监控\x1b[0m", processPid))
								return
							}

							// 在Unix系统上，需要发送信号0来检查进程是否存在
							err = process.Signal(syscall.Signal(0))
							if err != nil {
								sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;31m==> 进程 %d 已终止，停止监控\x1b[0m", processPid))
								return
							}
						}

						time.Sleep(500 * time.Millisecond)
					}
				}
			}
		}()

		// 后台运行模式下，60秒后自动关闭SSE连接
		go func() {
			time.Sleep(60 * time.Second)
			close(stopChan)
		}()
	} else {
		// 直接运行模式
		// ----- 修改这一部分以解决CentOS下PID获取问题 -----

		// 生成临时PID文件路径和临时脚本文件路径
		pidFilePath := fmt.Sprintf("%s/.pid_direct_temp", jpid.Catalog)
		wrapperScriptPath := fmt.Sprintf("%s/.run_wrapper.sh", jpid.Catalog)

		// 清理可能存在的旧文件
		_ = os.Remove(pidFilePath)
		_ = os.Remove(wrapperScriptPath)

		// 检测是否为Java命令
		isJavaCmd := strings.Contains(strings.ToLower(jpid.Run), "java ")

		// 创建一个包装脚本来执行命令并捕获正确的PID
		var scriptContent string

		if isJavaCmd {
			// 对Java程序做特殊处理
			scriptContent = fmt.Sprintf(`#!/bin/bash
# 记录当前进程的PID（脚本本身）
SCRIPT_PID=$$
echo $SCRIPT_PID > "%s"

# 执行实际命令
cd "%s"
%s &

# 给Java进程一点时间启动
sleep 1

# 查找Java进程
CMD_FEATURES=$(echo "%s" | tr ' ' '\n' | grep -v '^-' | grep -v '^java$' | head -3 | tr '\n' '|')
if [ ! -z "$CMD_FEATURES" ]; then
  # 使用特征词查找匹配的Java进程
  JAVA_PID=$(ps -eo pid,ppid,etime,command | grep -v grep | grep '[j]ava' | grep -E "$CMD_FEATURES" | sort -k 3 | head -1 | awk '{print $1}')
  if [ ! -z "$JAVA_PID" ]; then
    echo $JAVA_PID > "%s"
  fi
fi

# 等待前台进程完成
wait
`, pidFilePath, jpid.Catalog, jpid.Run, jpid.Run, pidFilePath)
		} else {
			// 非Java命令的处理
			scriptContent = fmt.Sprintf(`#!/bin/bash
# 记录当前进程的PID（脚本本身）
SCRIPT_PID=$$
echo $SCRIPT_PID > "%s"

# 执行实际命令
cd "%s"
%s &
ACTUAL_PID=$!

# 使用实际PID更新PID文件
echo $ACTUAL_PID > "%s"

# 等待前台进程完成
wait $ACTUAL_PID
EXIT_CODE=$?
exit $EXIT_CODE
`, pidFilePath, jpid.Catalog, jpid.Run, pidFilePath)
		}

		// 写入包装脚本
		if err = os.WriteFile(wrapperScriptPath, []byte(scriptContent), 0755); err != nil {
			sendSSEMessage(w, "error", "创建包装脚本失败："+err.Error())
			return nil, gerror.Wrap(err, "创建包装脚本失败")
		}

		// 执行包装脚本
		cmd := exec.Command("bash", wrapperScriptPath)
		cmd.Dir = jpid.Catalog
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("PROJECT_NAME=%s", jpid.Name),
			fmt.Sprintf("PROJECT_PID=%d", jpid.Pid),
			"LANG=en_US.UTF-8", // 添加UTF-8支持
		)

		// 创建输出管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			sendSSEMessage(w, "error", "创建输出管道失败："+err.Error())
			_ = os.Remove(wrapperScriptPath)
			return nil, gerror.Wrap(err, "创建输出管道失败")
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			sendSSEMessage(w, "error", "创建错误输出管道失败："+err.Error())
			_ = os.Remove(wrapperScriptPath)
			return nil, gerror.Wrap(err, "创建错误输出管道失败")
		}

		// 启动命令
		if err = cmd.Start(); err != nil {
			sendSSEMessage(w, "error", "启动失败："+err.Error())
			_ = os.Remove(wrapperScriptPath)
			return nil, gerror.Wrap(err, "启动失败")
		}

		// 等待一小段时间让PID文件写入
		time.Sleep(500 * time.Millisecond)

		// 读取PID文件获取真实PID
		var realPidFound bool
		for attempts := 0; attempts < 10; attempts++ {
			pidBytes, pidErr := os.ReadFile(pidFilePath)
			if pidErr == nil && len(pidBytes) > 0 {
				pidStr := strings.TrimSpace(string(pidBytes))
				pid, pidErr := strconv.Atoi(pidStr)
				if pidErr == nil && pid > 0 {
					processPid = pid
					realPidFound = true

					// 如果是Java命令，额外验证进程类型
					if isJavaCmd {
						checkCmd := exec.Command("ps", "-p", pidStr, "-o", "command=")
						output, _ := checkCmd.Output()
						if strings.Contains(strings.ToLower(string(output)), "java") {
							sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 获取到Java进程 PID: %d\x1b[0m", processPid))
						} else {
							// 找到的不是Java进程，继续等待
							realPidFound = false
						}
					} else {
						sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 获取到进程 PID: %d\x1b[0m", processPid))
					}

					if realPidFound {
						break
					}
				}
			}

			// 如果未找到有效PID，短暂等待后重试
			if !realPidFound {
				time.Sleep(300 * time.Millisecond)
			}
		}

		// 如果依然无法获取正确的PID，尝试使用进程搜索
		if !realPidFound {
			sendSSEMessage(w, "output", "\x1b[1;33m==> 无法从PID文件获取进程ID，尝试进程搜索...\x1b[0m")

			// 提取命令的关键部分作为搜索特征
			cmdParts := strings.Split(jpid.Run, " ")
			var searchFeature string
			if len(cmdParts) > 0 {
				// 使用命令的第一个非空参数作为特征
				for _, part := range cmdParts {
					if part != "" && part != "nohup" && part != "&" {
						searchFeature = part
						break
					}
				}
			}

			if searchFeature != "" {
				// 使用ps命令查找最近启动的匹配进程
				searchCmd := exec.Command("bash", "-c", fmt.Sprintf("ps -eo pid,etime,command | grep -v grep | grep '%s' | sort -k 2 | head -1 | awk '{print $1}'", searchFeature))
				output, err := searchCmd.Output()
				if err == nil && len(output) > 0 {
					pidStr := strings.TrimSpace(string(output))
					pid, err := strconv.Atoi(pidStr)
					if err == nil && pid > 0 {
						processPid = pid
						sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 通过进程搜索找到PID: %d\x1b[0m", processPid))
						realPidFound = true
					}
				}
			}
		}

		// 如果找到了有效PID，更新项目PID
		if realPidFound && processPid > 0 {
			if err = service.Jpid().UpdatePid(ctx, jpid.Pid, processPid); err != nil {
				sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: 更新 PID 失败\x1b[0m")
				g.Log().Warning(ctx, "更新PID失败", err)
			} else {
				sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 已更新 PID: %d -> %d\x1b[0m", jpid.Pid, processPid))
			}
		} else {
			// 使用脚本PID作为备选方案
			processPid = cmd.Process.Pid
			sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 无法获取实际进程PID，使用脚本PID: %d\x1b[0m", processPid))

			if err = service.Jpid().UpdatePid(ctx, jpid.Pid, processPid); err != nil {
				sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: 更新 PID 失败\x1b[0m")
				g.Log().Warning(ctx, "更新PID失败", err)
			}
		}

		// 处理输出
		done := make(chan error, 1)
		var outputBuffer bytes.Buffer
		go func() {
			scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
			for scanner.Scan() {
				line := scanner.Text()
				outputBuffer.WriteString(line + "\n")
				sendSSEMessage(w, "output", line)
			}
			done <- cmd.Wait()

			// 清理临时文件
			_ = os.Remove(wrapperScriptPath)
			_ = os.Remove(pidFilePath)
		}()

		// 直接运行模式：等待命令执行完成
		err = <-done
		if err != nil {
			sendSSEMessage(w, "error", "\x1b[1;31m==> 执行失败："+err.Error()+"\x1b[0m")
			return nil, gerror.Wrapf(err, "执行失败: %s", outputBuffer.String())
		}

		// 发送完成消息
		sendSSEMessage(w, "output", "\n\x1b[1;32m==> 执行完成!\x1b[0m")
		sendSSEMessage(w, "complete", "执行完成")

		return &v1.StartWithRunRes{
			Message: "启动成功",
			Output:  outputBuffer.String(),
		}, nil
	}

	// 如果是后台运行模式且有获取到PID，则更新项目PID
	if req.Background && processPid > 0 {
		if err = service.Jpid().UpdatePid(ctx, jpid.Pid, processPid); err != nil {
			sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: 更新 PID 失败\x1b[0m")
			g.Log().Warning(ctx, "更新PID失败", err)
		} else {
			sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 已更新 PID: %d -> %d\x1b[0m", jpid.Pid, processPid))
		}
	}

	// 发送完成消息
	sendSSEMessage(w, "output", "\x1b[1;32m==> 后台运行模式已启动\x1b[0m")
	sendSSEMessage(w, "complete", "执行完成")

	return &v1.StartWithRunRes{
		Message: "启动成功",
		Output:  "",
	}, nil
}
