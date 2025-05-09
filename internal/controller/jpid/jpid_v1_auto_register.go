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

	// 获取在线Java进程信息
	processes, err := javaprocess.GetJavaProcesses()
	if err != nil {
		return nil, gerror.Wrap(err, "获取Java进程失败")
	}

	// 调用service层处理注册逻辑
	total, updated, created, err := service.Jpid().AutoRegister(ctx, processes)
	if err != nil {
		return nil, err
	}

	return &v1.AutoRegisterRes{
		Message: fmt.Sprintf("注册完成: 处理 %d 个进程, 更新 %d 个, 新增 %d 个",
			total, updated, created),
		Total:   total,
		Updated: updated,
		Created: created,
	}, nil
}
