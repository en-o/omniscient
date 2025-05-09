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
		return nil, gerror.Wrap(err, "获取Java进程失败")
	}

	// 获取已存在的项目信息
	var existingProjects []*entity.Jpid
	if err = dao.Jpid.Ctx(ctx).Scan(&existingProjects); err != nil {
		return nil, gerror.Wrap(err, "查询已存在项目失败")
	}

	// 构建端口到项目的映射
	portToProject := make(map[string]*entity.Jpid)
	for _, project := range existingProjects {
		if project.Ports != "" {
			for _, port := range strings.Split(project.Ports, ",") {
				if port != "" {
					portToProject[port] = project
				}
			}
		}
	}

	// 处理每个检测到的进程
	for _, process := range processes {
		var existingProject *entity.Jpid
		var matchedPort string

		// 检查进程的端口是否已存在
		ports := strings.Split(process.Ports, ",")
		for _, port := range ports {
			if port != "" {
				if proj, exists := portToProject[port]; exists {
					existingProject = proj
					matchedPort = port
					break
				}
			}
		}

		if existingProject != nil {
			// 更新已存在的记录
			_, err = dao.Jpid.Ctx(ctx).Data(g.Map{
				"pid":        process.Pid,
				"catalog":    process.Catalog,
				"run":        process.Run,
				"status":     1, // 设置状态为启动
				"updated_at": gtime.Now(),
			}).Where("id", existingProject.Id).Update()

			if err != nil {
				g.Log().Warningf(ctx, "更新项目失败 [Port:%s, PID:%d]: %v",
					matchedPort, process.Pid, err)
				continue
			}

			g.Log().Infof(ctx, "更新项目成功 [Port:%s, OldPID:%d, NewPID:%d]",
				matchedPort, existingProject.Pid, process.Pid)
		} else {
			// 插入新记录
			_, err = dao.Jpid.Ctx(ctx).Data(do.Jpid{
				Name:      process.Name,
				Ports:     process.Ports,
				Pid:       process.Pid,
				Catalog:   process.Catalog,
				Run:       process.Run,
				Status:    1,
				CreatedAt: gtime.Now(),
				UpdatedAt: gtime.Now(),
			}).Insert()

			if err != nil {
				g.Log().Warningf(ctx, "插入新项目失败 [Ports:%s, PID:%d]: %v",
					process.Ports, process.Pid, err)
				continue
			}

			g.Log().Infof(ctx, "插入新项目成功 [Ports:%s, PID:%d]",
				process.Ports, process.Pid)
		}
	}

	// 更新结果统计
	res.Message = fmt.Sprintf("注册完成: 处理 %d 个进程", len(processes))
	return res, nil
}
