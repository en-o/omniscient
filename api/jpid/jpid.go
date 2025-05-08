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
}
