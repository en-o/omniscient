package jpid

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"omniscient/api/jpid/v1"
	"omniscient/internal/dao"
	"omniscient/internal/model/do"
	"omniscient/internal/model/entity"
	"omniscient/internal/util/javaprocess"
	"strings"
)

func (c *ControllerV1) AutoRegister(ctx context.Context, req *v1.AutoRegisterReq) (res *v1.AutoRegisterRes, err error) {
	// 初始化返回结构
	res = &v1.AutoRegisterRes{
		Message: "注册成功",
	}

	// 获取在线Java进程信息
	processes, err := javaprocess.GetJavaProcesses()
	if err != nil {
		return nil, err
	}

	// 获取已存在的项目端口信息
	var existingProjects []*entity.Jpid
	err = dao.Jpid.Ctx(ctx).Scan(&existingProjects)
	if err != nil {
		return nil, err
	}

	existingPorts := make(map[string]bool)
	for _, project := range existingProjects {
		if project.Ports != "" {
			ports := strings.Split(project.Ports, ",")
			for _, port := range ports {
				if port != "" {
					existingPorts[port] = true
				}
			}
		}
	}

	// 注册新的进程信息
	for _, process := range processes {
		// 检查端口是否已存在
		ports := strings.Split(process.Ports, ",")
		hasExistingPort := false
		for _, port := range ports {
			if port != "" && existingPorts[port] {
				hasExistingPort = true
				break
			}
		}

		// 如果端口已存在，跳过该进程
		if hasExistingPort {
			continue
		}

		// 插入新记录
		_, err = dao.Jpid.Ctx(ctx).Data(do.Jpid{
			Name:    process.Name,
			Ports:   process.Ports,
			Pid:     process.Pid,
			Catalog: process.Catalog,
			Run:     process.Run,
			Status:  1, // 设置状态为启动
		}).Insert()

		if err != nil {
			return nil, gerror.Wrap(err, "插入数据失败")
		}
	}

	return res, nil
}
