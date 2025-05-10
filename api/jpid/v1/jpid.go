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
	g.Meta `path:"/jpid/findOnline" tags:"Java" method:"get" summary:"查询在线的java项目列表"`
}
type OnlineRes struct {
	List []*entity.LinuxPid `json:"list" dc:"在线java项目"`
}

type AutoRegisterReq struct {
	g.Meta `path:"/jpid/auto/register" tags:"Java" method:"get" summary:"自动注册在线的java项目列表"`
}

type AutoRegisterRes struct {
	Message string `json:"message" dc:"操作结果"`
	Total   int    `json:"total"   dc:"处理进程数"`
	Updated int    `json:"updated" dc:"更新记录数"`
	Created int    `json:"created" dc:"新增记录数"`
}

type StopProjectReq struct {
	g.Meta `path:"/jpid/stop/:pid" tags:"Java" method:"post" summary:"根据pid停止运行"`
	Pid    int `v:"required|min:1" json:"pid" dc:"进程ID"`
}
type StopProjectRes struct {
	Message string `json:"message" dc:"操作结果"`
}

type StartWithRunReq struct {
	g.Meta `path:"/jpid/start/run/:pid" tags:"Java" method:"post" summary:"原生命令启动"`
	Pid    int `v:"required|min:1" json:"pid" dc:"进程ID"`
}
type StartWithRunRes struct {
	Message string `json:"message" dc:"操作结果"`
	Output  string `json:"output"  dc:"脚本输出"`
}

type StartWithScriptReq struct {
	g.Meta `path:"/jpid/start/script/:pid" tags:"Java" method:"post" summary:"脚本启动"`
	Pid    int `v:"required|min:1" json:"pid" dc:"进程ID"`
}

type StartWithScriptRes struct {
	Message string `json:"message" dc:"操作结果"`
	Output  string `json:"output"  dc:"脚本输出"`
}

type UpdateProjectReq struct {
	g.Meta      `path:"/jpid/update/:pid" tags:"Java" method:"post" summary:"更新项目信息"`
	Pid         int    `v:"required|min:1"      json:"pid"         dc:"进程ID"`
	Run         string `v:"required"            json:"run"         dc:"启动命令"`
	Description string `v:"required"            json:"description" dc:"项目描述"`
}

type UpdateProjectRes struct {
	Message string `json:"message" dc:"操作结果"`
}
