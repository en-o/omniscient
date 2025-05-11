// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package jpid

import (
	"context"

	"omniscient/api/jpid/v1"
)

type IJpidV1 interface {
	Jpid(ctx context.Context, req *v1.JpidReq) (res *v1.JpidRes, err error)
	Online(ctx context.Context, req *v1.OnlineReq) (res *v1.OnlineRes, err error)
	AutoRegister(ctx context.Context, req *v1.AutoRegisterReq) (res *v1.AutoRegisterRes, err error)
	StopProject(ctx context.Context, req *v1.StopProjectReq) (res *v1.StopProjectRes, err error)
	StartWithScript(ctx context.Context, req *v1.StartWithScriptReq) (res *v1.StartWithScriptRes, err error)
	StartWithRun(ctx context.Context, req *v1.StartWithRunReq) (res *v1.StartWithRunRes, err error)
	UpdateProject(ctx context.Context, req *v1.UpdateProjectReq) (res *v1.UpdateProjectRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
}
