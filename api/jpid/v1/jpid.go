package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/internal/model/entity"
)

type JpidReq struct {
	g.Meta `path:"/jpid" tags:"Java" method:"get" summary:"java 项目列表"`
	Worker string `dc:"worker名称，为空时查询当前worker的项目" v:"" in:"query"`
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

type StartWithScriptReq struct {
	g.Meta `path:"/jpid/start/script/:pid" method:"get" tags:"Jpid" summary:"脚本启动"`
	Pid    int `v:"required|min:1" json:"pid" dc:"进程ID"`
}

type StartWithScriptRes struct {
	g.Meta  `mime:"text/event-stream"`
	Message string `json:"message" dc:"返回信息"`
	Output  string `json:"output" dc:"执行输出"`
}

type StartWithRunReq struct {
	g.Meta     `path:"/jpid/start/run/:pid" method:"get" tags:"Jpid" summary:"原生命令启动"`
	Pid        int  `v:"required|min:1" json:"pid" dc:"进程ID"`
	Background bool `json:"background" dc:"是否后台运行"`
}

type StartWithRunRes struct {
	g.Meta  `mime:"text/event-stream"`
	Message string `json:"message" dc:"返回信息"`
	Output  string `json:"output" dc:"执行输出"`
}

type UpdateProjectReq struct {
	g.Meta      `path:"/jpid/update/:pid" tags:"Java" method:"post" summary:"更新项目信息"`
	Pid         int    `v:"required|min:1"      json:"pid"         dc:"进程ID"`
	Script      string `v:"required"            json:"script"         dc:"脚本命令"`
	Catalog     string `v:"required"            json:"catalog"         dc:"运行目录[临时设置后面会自动更新]"`
	Description string `v:"required"            json:"description" dc:"项目描述"`
}

type UpdateProjectRes struct {
	Message string `json:"message" dc:"操作结果"`
}

type DeleteReq struct {
	g.Meta `path:"/jpid/delete/:id" tags:"Java" method:"delete" summary:"删除项目"`
	Id     int `v:"required|min:1" in:"path" json:"id" dc:"项目ID"`
}

type DeleteRes struct {
	Message string `json:"message" dc:"操作结果"`
}

type StartWithDockerReq struct {
	g.Meta `path:"/jpid/start/docker/:pid" method:"get" tags:"Jpid" summary:"docker启动"`
	Pid    int  `v:"required|min:1" json:"pid" dc:"进程ID"`
	Reset  bool `json:"reset" dc:"是否重启"`
}

type StartWithDockerRes struct {
	g.Meta  `mime:"text/event-stream"`
	Message string `json:"message" dc:"返回信息"`
	Output  string `json:"output" dc:"执行输出"`
}
