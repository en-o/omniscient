package jpid

import (
	"context"
	"omniscient/internal/dao"
	"omniscient/internal/model/do"

	"omniscient/api/jpid/v1"
)

func (c *ControllerV1) Jpid(ctx context.Context, req *v1.JpidReq) (res *v1.JpidRes, err error) {
	dao.Jpid.Ctx(ctx).Where(do.Jpid{}).Scan(&res.List)
	return
}
