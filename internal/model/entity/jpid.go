// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Jpid is the golang structure for table jpid.
type Jpid struct {
	Id          int    `json:"id"          orm:"id"          description:""`              //
	Name        string `json:"name"        orm:"name"        description:"java项目名"`       // java项目名
	Port        int    `json:"port"        orm:"port"        description:"运行端口"`          // 运行端口
	Pid         int    `json:"pid"         orm:"pid"         description:"pid"`           // pid
	Catalog     string `json:"catalog"     orm:"catalog"     description:"运行目录"`          // 运行目录
	Run         string `json:"run"         orm:"run"         description:"运行脚本（sh命令"`     // 运行脚本（sh命令
	Status      int    `json:"status"      orm:"status"      description:"状态[1:启动，0:停止]"` // 状态[1:启动，0:停止]
	Description string `json:"description" orm:"description" description:"项目描述"`          // 项目描述
}

// ps -ef | grep java
type LinuxPid struct {
	Name  string   `json:"name"        orm:"name"        description:"java项目名"`   // java项目名
	Pid   int      `json:"pid"         orm:"pid"         description:"pid"`       // pid
	Run   string   `json:"run"         orm:"run"         description:"运行脚本（sh命令"` // 运行脚本（sh命令
	Ports []string `json:"ports"    orm:"ports"       description:"占用的端口"`        // 占用的端口
}
