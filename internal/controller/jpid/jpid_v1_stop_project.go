package jpid

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"os/exec"

	"omniscient/api/jpid/v1"
)

func (c *ControllerV1) StopProject(ctx context.Context, req *v1.StopProjectReq) (res *v1.StopProjectRes, err error) {
	r := g.RequestFromCtx(ctx)
	pid := r.GetForm("pid").String()
	cmd := exec.Command("kill", "-9", pid)
	if err := cmd.Run(); err != nil {
		r.Response.WriteJsonExit(map[string]string{"message": "停止失败", "error": err.Error()})
	}
	r.Response.WriteJsonExit(map[string]string{"message": "停止成功"})
	return
}
