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

		// 构建启动命令，包含获取真实PID的逻辑
		isJavaCmd := strings.Contains(strings.ToLower(jpid.Run), "java ")

		// 构建一个能获取实际PID的命令
		var cmdStr string
		if isJavaCmd {
			// 对于Java命令，我们需要额外的步骤来获取真正的Java进程PID
			// 先启动程序，然后查找关联的Java进程
			cmdStr = fmt.Sprintf("cd %s && nohup %s > nohup.log 2>&1 & "+
				"echo $! > %s && "+
				"sleep 1 && "+
				"SHELL_PID=$(cat %s) && "+
				"JAVA_PID=$(ps -o pid,ppid -ax | grep -v grep | grep $SHELL_PID | grep '[j]ava' | awk '{print $1}') && "+
				"if [ -n \"$JAVA_PID\" ]; then echo $JAVA_PID > %s; fi",
				jpid.Catalog, jpid.Run, pidFilePath, pidFilePath, pidFilePath)
		} else {
			// 非Java命令，直接保存bash创建的子进程PID
			cmdStr = fmt.Sprintf("cd %s && nohup %s > nohup.log 2>&1 & echo $! > %s",
				jpid.Catalog, jpid.Run, pidFilePath)
		}

		// 执行命令
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Dir = jpid.Catalog
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("PROJECT_NAME=%s", jpid.Name),
			fmt.Sprintf("PROJECT_PID=%d", jpid.Pid),
			"LANG=en_US.UTF-8", // 添加UTF-8支持
		)

		if err = cmd.Run(); err != nil {
			sendSSEMessage(w, "error", "启动失败："+err.Error())
			return nil, gerror.Wrap(err, "启动失败")
		}

		// 等待PID文件写入完成
		time.Sleep(2 * time.Second)

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
				sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 获取到进程 PID: %d\x1b[0m", processPid))
			}
		}

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
		// 直接运行模式，不变
		cmdStr := fmt.Sprintf("cd %s && %s", jpid.Catalog, jpid.Run)

		// 使用 bash 执行命令以支持更复杂的 shell 特性
		cmd := exec.Command("bash", "-c", cmdStr)
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
			return nil, gerror.Wrap(err, "创建输出管道失败")
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			sendSSEMessage(w, "error", "创建错误输出管道失败："+err.Error())
			return nil, gerror.Wrap(err, "创建错误输出管道失败")
		}

		// 启动命令
		if err = cmd.Start(); err != nil {
			sendSSEMessage(w, "error", "启动失败："+err.Error())
			return nil, gerror.Wrap(err, "启动失败")
		}

		// 记录进程 PID
		processPid = cmd.Process.Pid
		sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 进程 PID: %d\x1b[0m", processPid))

		// 更新项目 PID
		if err = service.Jpid().UpdatePid(ctx, jpid.Pid, processPid); err != nil {
			sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: 更新 PID 失败\x1b[0m")
			g.Log().Warning(ctx, "更新PID失败", err)
		} else {
			sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 已更新 PID: %d -> %d\x1b[0m", jpid.Pid, processPid))
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
		}()

		// 直接运行模式：等待前端主动关闭，不设置超时
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
