package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/internal/dao"
	"omniscient/internal/model/do"
	"omniscient/internal/model/entity"
	"strings"
)

type SJpid struct{}

func Jpid() *SJpid {
	return &SJpid{}
}

// GetByPid 根据PID获取项目信息
func (s *SJpid) GetByPid(ctx context.Context, pid int) (jpid *entity.Jpid, err error) {
	err = dao.Jpid.Ctx(ctx).Where("pid", pid).Scan(&jpid)
	return
}

// Update 更新项目信息
func (s *SJpid) Update(ctx context.Context, jpid *entity.Jpid) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{
			"status": jpid.Status,
		}).
		Where("id", jpid.Id).
		Update()
	return err
}

// GetList 获取项目列表
func (s *SJpid) GetList(ctx context.Context) (list []*entity.Jpid, err error) {
	err = dao.Jpid.Ctx(ctx).
		Order("id DESC").
		Scan(&list)
	return
}

// UpdateStatus 更新项目状态
func (s *SJpid) UpdateStatus(ctx context.Context, pid int, status int) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{"status": status}).
		Where("pid", pid).
		Update()
	return err
}

// AutoRegister 自动注册和更新Java进程
func (s *SJpid) AutoRegister(ctx context.Context, processes []*entity.LinuxPid) (total, updated, created int, err error) {
	// 获取已存在的项目信息
	var existingProjects []*entity.Jpid
	if err = dao.Jpid.Ctx(ctx).Scan(&existingProjects); err != nil {
		return 0, 0, 0, err
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

	total = len(processes)
	// 处理每个进程
	for _, process := range processes {
		var existingProject *entity.Jpid
		var matchedPort string

		// 检查端口是否已存在
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
			// 更新已存在记录
			if err := s.updateExistingProject(ctx, existingProject, process); err != nil {
				g.Log().Warningf(ctx, "更新项目失败 [Port:%s, PID:%d]: %v",
					matchedPort, process.Pid, err)
				continue
			}
			updated++
		} else {
			// 创建新记录
			if err := s.createNewProject(ctx, process); err != nil {
				g.Log().Warningf(ctx, "插入新项目失败 [Ports:%s, PID:%d]: %v",
					process.Ports, process.Pid, err)
				continue
			}
			created++
		}
	}

	return total, updated, created, nil
}

// updateExistingProject 更新已存在的项目
func (s *SJpid) updateExistingProject(ctx context.Context, existing *entity.Jpid, process *entity.LinuxPid) error {
	_, err := dao.Jpid.Ctx(ctx).Data(g.Map{
		"pid":     process.Pid,
		"catalog": process.Catalog,
		"run":     process.Run,
		"status":  1,
	}).Where("id", existing.Id).Update()
	return err
}

// createNewProject 创建新项目
func (s *SJpid) createNewProject(ctx context.Context, process *entity.LinuxPid) error {
	_, err := dao.Jpid.Ctx(ctx).Data(do.Jpid{
		Name:    process.Name,
		Ports:   process.Ports,
		Pid:     process.Pid,
		Catalog: process.Catalog,
		Run:     process.Run,
		Status:  1,
	}).Insert()
	return err
}

// UpdateProject 更新项目信息
func (s *SJpid) UpdateProject(ctx context.Context, pid int, run, description string) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{
			"run":         run,
			"description": description,
		}).
		Where("pid", pid).
		Update()
	return err
}
