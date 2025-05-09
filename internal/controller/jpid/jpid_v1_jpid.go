package jpid

import (
	"context"
	"omniscient/api/jpid/v1"
	"omniscient/internal/service"
)

func (c *ControllerV1) Jpid(ctx context.Context, req *v1.JpidReq) (res *v1.JpidRes, err error) {
	// 初始化 res , 不初始化会出现 invalid memory address or nil pointer dereference
	res = &v1.JpidRes{}

	res.List, err = service.Jpid().GetList(ctx)
	return
}
