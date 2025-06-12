package jpid

import (
	"context"
	"omniscient/api/jpid/v1"
	"omniscient/internal/service"
)

func (c *ControllerV1) Jpid(ctx context.Context, req *v1.JpidReq) (res *v1.JpidRes, err error) {
	// 初始化 res , 不初始化会出现 invalid memory address or nil pointer dereference
	res = &v1.JpidRes{}

	res.List, err = service.Jpid().GetList(ctx, req.Worker)
	return
}

// UpdateAutostart 更新自启状态
func (c *ControllerV1) UpdateAutostart(ctx context.Context, req *v1.UpdateAutostartReq) (res *v1.UpdateAutostartRes, err error) {
	err = service.Jpid().UpdateAutostart(ctx, req.Id, req.Autostart)
	return
}
