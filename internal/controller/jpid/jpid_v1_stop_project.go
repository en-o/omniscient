package jpid

import (
	"context"
	"omniscient/internal/service"
)

func (c *ControllerV1) StopProject(ctx context.Context, req *v1.StopProjectReq) (res *v1.StopProjectRes, err error) {
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		return nil, err
	}
	if jpid == nil {
		return nil, gerror.New("项目不存在")
	}

	// 执行停止命令
	cmd := exec.Command("kill", "-9", strconv.Itoa(jpid.Pid))
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 更新状态
	if err := service.Jpid().UpdateStatus(ctx, req.Pid, 0); err != nil {
		return nil, err
	}

	return &v1.StopProjectRes{Message: "停止成功"}, nil
}
