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

	// 根据运行模式构建命令
	var cmdStr string
	if req.Background {
		// 后台运行模式
		// 先删除已存在的 nohup.log
		nohupPath := fmt.Sprintf("%s/nohup.log", jpid.Catalog)
		if err = os.Remove(nohupPath); err != nil && !os.IsNotExist(err) {
			sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 警告: 清理旧日志文件失败: %v\x1b[0m", err))
		}
		cmdStr = fmt.Sprintf("cd %s && nohup %s > nohup.log 2>&1 &", jpid.Catalog, jpid.Run)
		// 创建一个 goroutine 用于实时读取日志内容
		go func() {
			// 等待文件创建
			time.Sleep(time.Second)

			// 打开并持续读取日志文件
			logFile, err := os.Open(nohupPath)
			if err != nil {
				sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;31m==> 打开日志文件失败: %v\x1b[0m", err))
				return
			}
			defer func(logFile *os.File) {
				err := logFile.Close()
				if err != nil {
					sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;31m==> 关闭日志文件失败: %v\x1b[0m", err))
				} else {
					sendSSEMessage(w, "output", "\x1b[1;32m==> 日志文件已关闭\x1b[0m")
				}
			}(logFile)

			scanner := bufio.NewScanner(logFile)
			for scanner.Scan() {
				sendSSEMessage(w, "output", scanner.Text())
			}
		}()
	} else {
		// 直接运行模式
		cmdStr = fmt.Sprintf("cd %s && %s", jpid.Catalog, jpid.Run)
	}

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
	processPid := cmd.Process.Pid
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

	// 根据运行模式处理执行完成
	if req.Background {
		// 后台运行模式设置超时
		select {
		case <-time.After(60 * time.Second):
			// 后台运行模式下不处理超时，直接返回成功
			sendSSEMessage(w, "output", "\x1b[1;32m==> 后台运行模式已启动\x1b[0m")
			sendSSEMessage(w, "complete", "执行完成")
			return &v1.StartWithRunRes{
				Message: "启动成功",
				Output:  outputBuffer.String(),
			}, nil
		case err = <-done:
			if err != nil {
				sendSSEMessage(w, "error", "\x1b[1;31m==> 执行失败："+err.Error()+"\x1b[0m")
				return nil, gerror.Wrapf(err, "执行失败: %s", outputBuffer.String())
			}
		}
	} else {
		// 直接运行模式：等待前端主动关闭，不设置超时
		err = <-done
		if err != nil {
			sendSSEMessage(w, "error", "\x1b[1;31m==> 执行失败："+err.Error()+"\x1b[0m")
			return nil, gerror.Wrapf(err, "执行失败: %s", outputBuffer.String())
		}
	}

	// 发送完成消息
	sendSSEMessage(w, "output", "\n\x1b[1;32m==> 执行完成!\x1b[0m")
	sendSSEMessage(w, "complete", "执行完成")

	return &v1.StartWithRunRes{
		Message: "启动成功",
		Output:  outputBuffer.String(),
	}, nil
}
