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

	// 发送启动提示
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 正在启动Docker项目: %s\x1b[0m", jpid.Name))

	// 3. 根据reset参数决定使用的命令
	var cmdStr string
	if req.Reset {
		cmdStr = "restart"
		sendSSEMessage(w, "output", "\x1b[1;34m==> 执行操作: 重启容器\x1b[0m")
	} else {
		cmdStr = "start"
		sendSSEMessage(w, "output", "\x1b[1;34m==> 执行操作: 启动容器\x1b[0m")
	}

	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;34m==> 执行命令: docker compose %s %s\x1b[0m", cmdStr, jpid.Name))
	sendSSEMessage(w, "output", "\x1b[1;33m==> 开始执行...\x1b[0m\n")

	// 创建命令
	cmd := exec.Command("docker", cmdStr, jpid.Name)

	// 4. 设置命令输出管道
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

	// 5. 启动命令
	if err = cmd.Start(); err != nil {
		sendSSEMessage(w, "error", "\x1b[1;31m==> 启动Docker命令失败: "+err.Error()+"\x1b[0m")
		return nil, gerror.Wrap(err, "启动Docker命令失败")
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

	// 设置超时和等待执行完成
	var cmdErr error
	select {
	case <-time.After(5 * time.Minute): // 5分钟超时
		if err := cmd.Process.Kill(); err != nil {
			sendSSEMessage(w, "error", "\x1b[1;31m==> 终止超时进程失败: "+err.Error()+"\x1b[0m")
			g.Log().Warning(ctx, "终止超时进程失败", "error", err)
		}
		sendSSEMessage(w, "error", "\x1b[1;31m==> Docker命令执行超时\x1b[0m")
		cmdErr = gerror.New("Docker命令执行超时")
	case err := <-done:
		cmdErr = err
	}

	// 处理命令执行结果
	if cmdErr != nil {
		// 命令执行失败
		sendSSEMessage(w, "error", "\x1b[1;31m==> 执行失败: "+cmdErr.Error()+"\x1b[0m")
		g.Log().Error(ctx, "Docker命令执行失败",
			"pid", jpid.Pid,
			"name", jpid.Name,
			"error", cmdErr,
		)
		return nil, gerror.Wrapf(cmdErr, "Docker执行失败: %s", outputBuffer.String())
	} else {
		// 命令执行成功，更新项目状态
		if updateErr := service.Jpid().UpdateStatus(ctx, req.Pid, 1); updateErr != nil {
			g.Log().Error(ctx, "更新项目状态失败",
				"pid", jpid.Pid,
				"name", jpid.Name,
				"error", updateErr,
			)
			sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: Docker命令执行成功，但更新状态失败\x1b[0m")
		} else {
			sendSSEMessage(w, "output", "\n\x1b[1;32m==> 执行成功!\x1b[0m")
			g.Log().Info(ctx, "Docker命令执行成功",
				"pid", jpid.Pid,
				"name", jpid.Name,
				"reset", req.Reset,
			)
		}
	}

	// 发送完成消息
	sendSSEMessage(w, "complete", "执行完成")

	return &v1.StartWithDockerRes{
		Message: "Docker操作执行完成",
		Output:  outputBuffer.String(),
	}, nil
}
