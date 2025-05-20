package jpid

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	v1 "omniscient/api/jpid/v1"
	"omniscient/internal/service"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func (c *ControllerV1) StopProject(ctx context.Context, req *v1.StopProjectReq) (res *v1.StopProjectRes, err error) {
	// Validate input
	if req == nil || req.Pid <= 0 {
		return nil, gerror.New("无效的项目ID")
	}

	// Get project details
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		return nil, gerror.Wrapf(err, "获取项目信息失败: pid=%d", req.Pid)
	}
	if jpid == nil {
		return nil, gerror.Newf("项目不存在: pid=%d", req.Pid)
	}

	// Create command based on project type
	var cmd *exec.Cmd
	if jpid.Way == 1 {
		// Docker方式停止 - 直接使用docker命令，不需要依赖项目目录
		cmd = exec.Command("docker", "stop", jpid.Name)
		g.Log().Info(ctx, "准备停止Docker项目",
			"pid", jpid.Pid,
			"name", jpid.Name,
			"catalog", jpid.Catalog,
		)
	} else {
		// PID方式停止 - 确保PID存在并且有效
		if jpid.Pid <= 0 {
			return nil, gerror.Newf("无效的进程PID: %d", jpid.Pid)
		}
		cmd = exec.Command("kill", "-9", strconv.Itoa(jpid.Pid))

		// 仅对非Docker方式检查目录
		if jpid.Catalog == "" {
			return nil, gerror.New("项目目录为空")
		}
		if _, err := os.Stat(jpid.Catalog); os.IsNotExist(err) {
			return nil, gerror.Wrapf(err, "项目目录不存在: %s", jpid.Catalog)
		}

		// 仅对非Docker方式设置工作目录
		cmd.Dir = jpid.Catalog
	}

	// 设置环境变量
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PROJECT_NAME=%s", jpid.Name),
		fmt.Sprintf("PROJECT_PID=%d", jpid.Pid),
	)

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	if err = cmd.Start(); err != nil {
		return nil, gerror.Wrapf(err, "启动停止命令失败: %s", stderr.String())
	}

	// 创建一个channel用于等待命令执行完成
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// 设置超时时间
	var cmdErr error
	select {
	case <-time.After(30 * time.Second):
		if killErr := cmd.Process.Kill(); killErr != nil {
			g.Log().Warning(ctx, "终止超时进程失败", "error", killErr)
		}
		cmdErr = gerror.New("停止命令执行超时")
	case waitErr := <-done:
		cmdErr = waitErr
	}

	// 检查命令执行结果
	if cmdErr != nil {
		g.Log().Error(ctx, "停止命令执行失败",
			"pid", jpid.Pid,
			"name", jpid.Name,
			"error", cmdErr,
			"stderr", stderr.String(),
		)
		return nil, gerror.Wrapf(cmdErr, "停止命令执行失败: %s", stderr.String())
	}

	// 更新项目状态 - 移到成功执行命令后
	if err = service.Jpid().UpdateStatus(ctx, req.Pid, 0); err != nil {
		return nil, gerror.Wrapf(err, "更新项目状态失败: pid=%d", req.Pid)
	}

	// 记录执行日志
	g.Log().Info(ctx, "项目停止成功",
		"pid", jpid.Pid,
		"name", jpid.Name,
		"output", stdout.String(),
	)

	return &v1.StopProjectRes{Message: "停止成功"}, nil
}
