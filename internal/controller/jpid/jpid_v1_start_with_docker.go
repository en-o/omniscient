package jpid

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	v1 "omniscient/api/jpid/v1"
	"omniscient/internal/service"
)

func (c *ControllerV1) StartWithDocker(ctx context.Context, req *v1.StartWithDockerReq) (res *v1.StartWithDockerRes, err error) {
	// 获取响应写入器
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("无法获取请求上下文")
	}
	w := r.Response.Writer

	// 设置 SSE 相关的 header
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no") // Nginx特殊设置，防止缓冲

	// 1. 获取项目信息
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 获取项目信息失败: "+err.Error()+"\x1b[0m")
		return nil, gerror.Wrapf(err, "获取项目信息失败: pid=%d", req.Pid)
	}
	if jpid == nil {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 项目不存在\x1b[0m")
		return nil, gerror.Newf("项目不存在: pid=%d", req.Pid)
	}

	// 2. 验证是否为Docker项目
	if jpid.Way != 1 {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 非Docker项目，无法使用Docker启动\x1b[0m")
		return nil, gerror.New("非Docker项目，无法使用Docker启动")
	}
	// 3. 验证项目状态
	if jpid.Status == 1 && !req.Reset {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 容器已在运行中，如需重启请使用重启功能\x1b[0m")
		return nil, gerror.New("容器已在运行中")
	}

	// 发送启动提示
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 正在启动Docker项目: %s\x1b[0m", jpid.Name))

	// 4. 根据reset参数决定使用的命令
	var cmdStr string
	if req.Reset {
		cmdStr = "restart"
		sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 正在重启容器: %s\x1b[0m", jpid.Name))
	} else {
		cmdStr = "start"
		sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;33m==> 正在启动新容器: %s\x1b[0m", jpid.Name))
	}

	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;34m==> 执行命令: docker %s %s\x1b[0m", cmdStr, jpid.Name))
	sendSSEMessage(w, "output", "\x1b[1;33m==> 开始执行...\x1b[0m\n")

	// 创建命令,允许超时控制
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", cmdStr, jpid.Name)

	// 5. 设置命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 无法创建标准输出管道: "+err.Error()+"\x1b[0m")
		return nil, gerror.Wrap(err, "无法创建标准输出管道")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 无法创建错误输出管道: "+err.Error()+"\x1b[0m")
		return nil, gerror.Wrap(err, "无法创建错误输出管道")
	}

	// 启动命令
	if err = cmd.Start(); err != nil {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 启动Docker命令失败: "+err.Error()+"\x1b[0m")
		return nil, gerror.Wrap(err, "启动Docker命令失败")
	}

	//  处理输出
	var outputBuffer bytes.Buffer
	outputChan := make(chan struct{})

	// 使用单独的程序来捕获输出
	go func() {
		defer close(outputChan)

		// 将标准输出和标准输入组合到一个阅读器中
		reader := io.MultiReader(stdout, stderr)
		scanner := bufio.NewScanner(reader)

		// Set a larger buffer for scanner if needed
		const maxScannerSize = 1024 * 1024 // 1MB
		scannerBuf := make([]byte, maxScannerSize)
		scanner.Buffer(scannerBuf, maxScannerSize)

		for scanner.Scan() {
			line := scanner.Text()
			outputBuffer.WriteString(line + "\n")

			// Send each line as an SSE message and flush immediately
			sendSSEMessage(w, "output", line)
		}

		if err := scanner.Err(); err != nil {
			outputBuffer.WriteString("Scanner error: " + err.Error() + "\n")
			sendSSEMessage(w, "error", "\x1b[1;31m==> 读取输出错误: "+err.Error()+"\x1b[0m")
		}
	}()

	// Wait for command to complete and handle errors
	var cmdErr error

	// Create a channel for the command result
	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	// Wait for either output completion, command completion, or timeout
	select {
	case <-ctx.Done(): // Context timeout
		cmdErr = ctx.Err()
		sendSSEMessage(w, "error", "\x1b[1;31m==> Docker命令执行超时\x1b[0m")

		// Force kill the process if it's still running
		if cmd.Process != nil {
			if killErr := cmd.Process.Kill(); killErr != nil {
				g.Log().Warning(ctx, "终止超时进程失败", "error", killErr)
			}
		}
	case err := <-cmdDone: // Command completed
		cmdErr = err
	}

	// Wait for output goroutine to finish, with a timeout
	select {
	case <-outputChan:
		// Output processing completed normally
	case <-time.After(3 * time.Second):
		// Output processing timed out, but continue anyway
		g.Log().Warning(ctx, "输出处理未完成，已超时")
	}

	// Handle command execution result
	if cmdErr != nil {
		// Command execution failed
		sendSSEMessage(w, "error", "\x1b[1;31m==> 执行失败: "+cmdErr.Error()+"\x1b[0m")
		g.Log().Error(ctx, "Docker命令执行失败",
			"pid", jpid.Pid,
			"name", jpid.Name,
			"error", cmdErr,
		)
		sendSSEMessage(w, "complete", "执行失败")
		return nil, gerror.Wrapf(cmdErr, "Docker执行失败: %s", outputBuffer.String())
	} else {
		// 命令执行成功，更新项目状态
		if updateErr := service.Jpid().UpdateStatusById(ctx, jpid.Id, 1); updateErr != nil {
			g.Log().Error(ctx, "更新项目状态失败",
				"pid", jpid.Pid,
				"name", jpid.Name,
				"error", updateErr,
			)
			sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: Docker命令执行成功，但更新状态失败\x1b[0m")
		} else {
			g.Log().Info(ctx, "Docker命令执行成功",
				"pid", jpid.Pid,
				"name", jpid.Name,
				"reset", req.Reset,
			)
		}
		// 发送完成消息
		var message string
		if req.Reset {
			message = "重启"
		} else {
			message = "启动"
		}
		logMessage := fmt.Sprintf("Docker命令执行成功，pid: %d, name: %s, reset: %v", jpid.Pid, jpid.Name, req.Reset)
		sendSSEMessage(w, "output", "\n\x1b[1;32m==> "+message+"成功!\x1b[0m\n"+logMessage)

		// 发送完成消息
		sendSSEMessage(w, "complete", "执行完成")
	}

	return &v1.StartWithDockerRes{
		Message: "Docker操作执行完成",
		Output:  outputBuffer.String(),
	}, nil
}
