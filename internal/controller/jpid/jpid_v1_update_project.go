package jpid

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"omniscient/api/jpid/v1"
	"omniscient/internal/service"
)

// 更新基础信息
func (c *ControllerV1) UpdateProject(ctx context.Context, req *v1.UpdateProjectReq) (res *v1.UpdateProjectRes, err error) {
	// 验证项目是否存在
	jpid, err := service.Jpid().GetByPid(ctx, req.Pid)
	if err != nil {
		return nil, err
	}
	if jpid == nil {
		return nil, gerror.New("项目不存在")
	}

	// 更新项目信息
	err = service.Jpid().UpdateInfo(ctx, req.Pid, req.Script, req.Description)
	if err != nil {
		return nil, gerror.Wrap(err, "更新项目失败")
	}

	return &v1.UpdateProjectRes{Message: "更新成功"}, nil
}
