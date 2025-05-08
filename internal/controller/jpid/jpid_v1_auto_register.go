package jpid

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"omniscient/api/jpid/v1"
)

func (c *ControllerV1) AutoRegister(ctx context.Context, req *v1.AutoRegisterReq) (res *v1.AutoRegisterRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
