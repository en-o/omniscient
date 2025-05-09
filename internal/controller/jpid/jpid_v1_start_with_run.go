package jpid

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/api/jpid/v1"
	"omniscient/internal/model/entity"
)

func (c *ControllerV1) StartWithRun(ctx context.Context, req *v1.StartWithRunReq) (res *v1.StartWithRunRes, err error) {
	r := g.RequestFromCtx(ctx)

	var data entity.Jpid
	err = r.ParseForm(&data)
	if err != nil {
		r.Response.WriteJsonExit(map[string]string{"message": "参数错误", "error": err.Error()})
	}

	//cmd := exec.Command("sh", "-c", "cd "+data.Catalog+" && "+data.Run)
	//if err := cmd.Run(); err != nil {
	//	r.Response.WriteJsonExit(map[string]string{"message": "启动失败", "error": err.Error()})
	//}
	//r.Response.WriteJsonExit(map[string]string{"message": "启动成功"})
	r.Response.WriteJson(data)
	return
}
