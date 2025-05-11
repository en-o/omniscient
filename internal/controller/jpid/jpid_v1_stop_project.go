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
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		return nil, err
	}
	if jpid == nil {
		return nil, gerror.New("项目不存在")
	}

	// 构建命令执行环境
	cmd := exec.Command("kill", "-9", strconv.Itoa(jpid.Pid))

	// 设置工作目录
	cmd.Dir = jpid.Catalog

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
		return nil, gerror.Wrap(err, "启动脚本失败")
	}
	// 更新项目状态
	if err = service.Jpid().UpdateStatus(ctx, req.Pid, 0); err != nil {
		return nil, gerror.Wrap(err, "更新状态失败")
	}

	// 创建一个channel用于等待命令执行完成
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// 设置超时时间
	select {
	case <-time.After(30 * time.Second):
		if err = cmd.Process.Kill(); err != nil {
			return nil, gerror.Wrap(err, "终止超时进程失败")
		}
		return nil, gerror.New("脚本执行超时")
	case err = <-done:
		if err != nil {
			return nil, gerror.Wrapf(err, "脚本执行失败: %s", stderr.String())
		}
	}

	// 记录执行日志
	g.Log().Info(ctx, "脚本执行成功",
		"pid", jpid.Pid,
		"name", jpid.Name,
		"output", stdout.String(),
	)
	return &v1.StopProjectRes{Message: "停止成功"}, nil
}
