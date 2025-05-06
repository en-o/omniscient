package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type JpidReq struct {
	g.Meta `path:"/jpid" tags:"Java" method:"get" summary:"java 项目"`
}
type JpidRes struct{}
