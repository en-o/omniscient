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

type OnlineReq struct {
	g.Meta `path:"/findOnline" tags:"Java" method:"get" summary:"查询在线的java项目列表"`
}
type OnlineRes struct {
	List []*entity.LinuxPid `json:"list" dc:"在线java项目"`
}

type AutoRegisterReq struct {
	g.Meta `path:"/auto/register" tags:"Java" method:"get" summary:"自动注册在线的java项目列表"`
}
type AutoRegisterRes struct{}
