// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Jpid is the golang structure for table jpid.
type Jpid struct {
	Id          int    `json:"id"          orm:"id"          description:""`                        //
	Name        string `json:"name"        orm:"name"        description:"java项目名"`                 // java项目名
	Ports       string `json:"ports"       orm:"ports"       description:"运行端口,多个逗号隔开"`             // 运行端口,多个逗号隔开
	Pid         int    `json:"pid"         orm:"pid"         description:"pid"`                     // pid
	Catalog     string `json:"catalog"     orm:"catalog"     description:"运行目录"`                    // 运行目录
	Run         string `json:"run"         orm:"run"         description:"原生启动命令"`                  // 原生启动命令
	Script      string `json:"script"      orm:"script"      description:"sh脚本启动命令"`                // sh脚本启动命令
	Worker      string `json:"worker"      orm:"worker"      description:"服务器"`                     // 服务器
	Status      int    `json:"status"      orm:"status"      description:"状态[1:启动，0:停止]"`           // 状态[1:启动，0:停止]
	Description string `json:"description" orm:"description" description:"项目描述"`                    // 项目描述
	Way         int    `json:"way"      orm:"way"      description:"启动方式[1:docker, 2:jdk]"` // 启动方式[1:docker, 2:jdk]
}

// ps -ef | grep java
//
//	nohup LinuxPid.run  >/dev/null 2>&1
type LinuxPid struct {
	Name    string `json:"name"        orm:"name"        description:"java项目名"`                 // java项目名
	Pid     int    `json:"pid"         orm:"pid"         description:"pid"`                     // pid
	Run     string `json:"run"         orm:"run"         description:"原生启动命令"`                  // 原生启动命令
	Ports   string `json:"ports"    orm:"ports"       description:"占用的端口"`                      // 占用的端口
	Catalog string `json:"catalog"     orm:"catalog"     description:"运行目录[jar文件所在目录]"`         // 运行目录[jar文件所在目录]
	Worker  string `json:"worker"     orm:"worker"     description:"服务器"`                       // 服务器
	Way     int    `json:"way"      orm:"way"      description:"启动方式[1:docker, 2:jdk]"` // 启动方式[1:docker, 2:jdk]
}
