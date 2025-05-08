package jpid

import (
	"context"
	"omniscient/api/jpid/v1"
	"omniscient/internal/model/entity"
	"omniscient/internal/util/javaprocess"
)

func (c *ControllerV1) Online(ctx context.Context, req *v1.OnlineReq) (res *v1.OnlineRes, err error) {
	// 初始化返回结构
	res = &v1.OnlineRes{
		List: make([]*entity.LinuxPid, 0),
	}

	// 获取Java进程信息
	processes, err := javaprocess.GetJavaProcesses()
	if err != nil {
		return nil, err
	}

	// 转换进程信息到API响应格式
	for _, p := range processes {
		linuxPid := &entity.LinuxPid{
			Name:    p.Name,
			Pid:     p.Pid,
			Run:     p.Run,
			Ports:   p.Ports,
			Catalog: p.Catalog,
		}
		res.List = append(res.List, linuxPid)
	}

	return res, nil
}
