package jpid

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"io"
	"net/http"
	"omniscient/api/jpid/v1"
	"omniscient/internal/service"
	"os"
	"os/exec"
	"time"
)

func (c *ControllerV1) StartWithScript(ctx context.Context, req *v1.StartWithScriptReq) (res *v1.StartWithScriptRes, err error) {

	// 获取响应写入器
	r := g.RequestFromCtx(ctx)
	w := r.Response.Writer

	// 设置 SSE 相关的 header
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 刷新 header
	//w.(http.Flusher).Flush()

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
	if jpid.Script == "" {
		sendSSEMessage(w, "error", "脚本命令为空")
		return nil, gerror.New("脚本命令为空")
	}

	// 发送启动提示
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 正在启动项目: %s\x1b[0m", jpid.Name))
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;34m==> 工作目录: %s\x1b[0m", jpid.Catalog))
	sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;34m==> 执行命令: %s\x1b[0m", jpid.Script))
	sendSSEMessage(w, "output", "\x1b[1;33m==> 开始执行...\x1b[0m\n")

	// 构建命令执行环境
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", jpid.Catalog, jpid.Script))

	// 设置工作目录
	cmd.Dir = jpid.Catalog

	// 设置环境变量
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PROJECT_NAME=%s", jpid.Name),
		fmt.Sprintf("PROJECT_PID=%d", jpid.Pid),
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
		sendSSEMessage(w, "error", "启动脚本失败："+err.Error())
		return nil, gerror.Wrap(err, "启动脚本失败")
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
	select {
	case <-time.After(30 * time.Second):
		if err = cmd.Process.Kill(); err != nil {
			sendSSEMessage(w, "error", "\x1b[1;31m==> 终止超时进程失败："+err.Error()+"\x1b[0m")
			return nil, gerror.Wrap(err, "终止超时进程失败")
		}
		sendSSEMessage(w, "error", "\x1b[1;31m==> 脚本执行超时\x1b[0m")
		return nil, gerror.New("脚本执行超时")
	case err = <-done:
		if err != nil {
			sendSSEMessage(w, "error", "\x1b[1;31m==> 脚本执行失败："+err.Error()+"\x1b[0m")
			return nil, gerror.Wrapf(err, "脚本执行失败: %s", outputBuffer.String())
		}
	}

	// 等待进程启动并获取新的 PID
	sendSSEMessage(w, "output", "\x1b[1;33m==> 正在获取新进程 PID...\x1b[0m")
	newPid, err := service.Jpid().FindNewPid(ctx, jpid)
	if err != nil {
		sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: 无法获取新的 PID\x1b[0m")
		g.Log().Warning(ctx, "无法获取新的PID", err)
	} else if newPid != jpid.Pid {
		if err = service.Jpid().UpdatePid(ctx, jpid.Pid, newPid); err != nil {
			sendSSEMessage(w, "output", "\x1b[1;31m==> 警告: 更新 PID 失败\x1b[0m")
			g.Log().Warning(ctx, "更新PID失败", err)
		} else {
			sendSSEMessage(w, "output", fmt.Sprintf("\x1b[1;32m==> 已更新 PID: %d -> %d\x1b[0m", jpid.Pid, newPid))
		}
	}

	// 发送完成消息
	sendSSEMessage(w, "output", "\n\x1b[1;32m==> 执行完成!\x1b[0m")
	sendSSEMessage(w, "complete", "执行完成")
	return &v1.StartWithScriptRes{
		Message: "启动成功",
		Output:  outputBuffer.String(),
	}, nil
}

// 辅助函数：发送 SSE 消息
func sendSSEMessage(w http.ResponseWriter, event, data string) {
	_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
	if err != nil {
		return
	}
	w.(http.Flusher).Flush()
}
