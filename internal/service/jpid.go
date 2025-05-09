package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/internal/dao"
	"omniscient/internal/model/entity"
)

type sJpid struct{}

func Jpid() *sJpid {
	return &sJpid{}
}

// GetByPid 根据PID获取项目信息
func (s *sJpid) GetByPid(ctx context.Context, pid int) (jpid *entity.Jpid, err error) {
	err = dao.Jpid.Ctx(ctx).Where("pid", pid).Scan(&jpid)
	return
}

// Update 更新项目信息
func (s *sJpid) Update(ctx context.Context, jpid *entity.Jpid) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{
			"status": jpid.Status,
		}).
		Where("id", jpid.Id).
		Update()
	return err
}

// GetList 获取项目列表
func (s *sJpid) GetList(ctx context.Context) (list []*entity.Jpid, err error) {
	err = dao.Jpid.Ctx(ctx).
		Order("id DESC").
		Scan(&list)
	return
}

// UpdateStatus 更新项目状态
func (s *sJpid) UpdateStatus(ctx context.Context, pid int, status int) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{"status": status}).
		Where("pid", pid).
		Update()
	return err
}
