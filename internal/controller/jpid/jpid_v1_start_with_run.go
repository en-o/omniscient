package jpid

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"omniscient/api/jpid/v1"
	"omniscient/internal/service"
	"os/exec"
)

func (c *ControllerV1) StartWithRun(ctx context.Context, req *v1.StartWithRunReq) (res *v1.StartWithRunRes, err error) {
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		return nil, err
	}
	if jpid == nil {
		return nil, gerror.New("项目不存在")
	}

	// 执行启动命令
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", jpid.Catalog, jpid.Run))
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 更新状态
	if err := service.Jpid().UpdateStatus(ctx, req.Pid, 1); err != nil {
		return nil, err
	}

	return &v1.StartWithRunRes{Message: "启动成功"}, nil
}
