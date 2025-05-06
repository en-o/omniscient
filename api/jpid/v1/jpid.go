package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/internal/model/entity"
)

type JpidReq struct {
	g.Meta `path:"/jpid" tags:"Java" method:"get" summary:"java 项目列表"`
}
type JpidRes struct {
	List []*entity.Jpid `json:"list" dc:"java 项目列表"`
}
