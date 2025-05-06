// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Jpid is the golang structure of table jpid for DAO operations like Where/Data.
type Jpid struct {
	g.Meta      `orm:"table:jpid, do:true"`
	Id          interface{} //
	Name        interface{} // java项目名
	Port        interface{} // 运行端口
	Pid         interface{} // pid
	Catalog     interface{} // 运行目录
	Run         interface{} // 运行脚本（sh命令
	Description interface{} // 项目描述
}
