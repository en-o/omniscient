package jpid

import (
	"bufio"
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
	// 1. 获取项目信息
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		return nil, gerror.Wrapf(err, "获取项目信息失败: pid=%d", req.Pid)
	}
	if jpid == nil {
		return nil, gerror.Newf("项目不存在: pid=%d", req.Pid)
	}

	// 2. 验证是否为Docker项目
	if jpid.Way != 1 {
		return nil, gerror.New("非Docker项目，无法使用Docker启动")
	}

	// 3. 根据reset参数决定使用的命令
	var cmd *exec.Cmd
	if req.Reset {
		// 重启项目
		cmd = exec.Command("docker", "compose", "restart", jpid.Name)
	} else {
		// 启动项目
		cmd = exec.Command("docker", "compose", "start", jpid.Name)
	}

	// 4. 设置命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, gerror.Wrap(err, "无法创建标准输出管道")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, gerror.Wrap(err, "无法创建错误输出管道")
	}

	// 5. 启动命令
	if err = cmd.Start(); err != nil {
		return nil, gerror.Wrap(err, "启动Docker命令失败")
	}

	// 6. 准备SSE响应
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("无法获取请求上下文")
	}

	// 7. 设置SSE头
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("X-Accel-Buffering", "no") // Nginx特殊设置，防止缓冲

	// 8. 发送初始消息
	r.Response.WriteExit("event: message\ndata: {\"message\": \"开始执行Docker命令\", \"output\": \"\"}\n\n")
	r.Response.Flush()

	// 9. 创建合并的输出读取器
	outputReader := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(outputReader)

	// 10. 创建done通道，用于命令完成通知
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// 11. 实时发送输出
	go func() {
		// 实时发送输出
		for scanner.Scan() {
			line := scanner.Text()
			r.Response.Write(fmt.Sprintf("event: message\ndata: {\"message\": \"执行中\", \"output\": %q}\n\n", line))
			r.Response.Flush()
		}

		// 等待命令完成
		var cmdErr error
		select {
		case <-time.After(5 * time.Minute): // 5分钟超时
			if err := cmd.Process.Kill(); err != nil {
				g.Log().Warning(ctx, "终止超时进程失败", "error", err)
			}
			cmdErr = gerror.New("Docker命令执行超时")
		case err := <-done:
			cmdErr = err
		}

		// 发送最终结果
		if cmdErr != nil {
			// 命令执行失败
			r.Response.Write(fmt.Sprintf("event: message\ndata: {\"message\": \"执行失败: %s\", \"output\": \"\"}\n\n", cmdErr.Error()))
			g.Log().Error(ctx, "Docker命令执行失败",
				"pid", jpid.Pid,
				"name", jpid.Name,
				"error", cmdErr,
			)
		} else {
			// 命令执行成功，更新项目状态
			if updateErr := service.Jpid().UpdateStatus(ctx, req.Pid, 1); updateErr != nil {
				g.Log().Error(ctx, "更新项目状态失败",
					"pid", jpid.Pid,
					"name", jpid.Name,
					"error", updateErr,
				)
				r.Response.Write("event: message\ndata: {\"message\": \"Docker命令执行成功，但更新状态失败\", \"output\": \"\"}\n\n")
			} else {
				r.Response.Write("event: message\ndata: {\"message\": \"执行成功\", \"output\": \"\"}\n\n")
				g.Log().Info(ctx, "Docker命令执行成功",
					"pid", jpid.Pid,
					"name", jpid.Name,
					"reset", req.Reset,
				)
			}
		}

		// 发送关闭事件
		r.Response.Write("event: close\ndata: {}\n\n")
		r.Response.Flush()
	}()

	// 12. 告诉GoFrame不要自动结束响应
	return nil, r.Response.Buffer()
}
