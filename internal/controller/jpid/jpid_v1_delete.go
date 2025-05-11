package jpid

import (
	"context"
	"omniscient/internal/service"

	"omniscient/api/jpid/v1"
)

func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = service.Jpid().Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteRes{
		Message: "删除成功",
	}, nil
}
