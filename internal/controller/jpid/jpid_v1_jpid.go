package jpid

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"omniscient/api/jpid/v1"
)

func (c *ControllerV1) Jpid(ctx context.Context, req *v1.JpidReq) (res *v1.JpidRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
