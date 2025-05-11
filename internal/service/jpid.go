package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/internal/dao"
	"omniscient/internal/model/do"
	"omniscient/internal/model/entity"
	"omniscient/internal/util/javaprocess"
	"omniscient/internal/util/system"
	"os/exec"
	"strings"
	"time"
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
	workerName := system.GetWorkerName()
	_, err := dao.Jpid.Ctx(ctx).Data(do.Jpid{
		Name:    process.Name,
		Ports:   process.Ports,
		Pid:     process.Pid,
		Catalog: process.Catalog,
		Run:     process.Run,
		Status:  1,
		Worker:  workerName,
	}).Insert()
	return err
}

// UpdateProject 更新项目基础信息
func (s *SJpid) UpdateInfo(ctx context.Context, pid int, script, description string) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{
			"script":      script,
			"description": description,
		}).
		Where("pid", pid).
		Update()
	return err
}

// UpdatePid 更新项目的 PID
func (s *SJpid) UpdatePid(ctx context.Context, oldPid int, newPid int) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{
			"pid":    newPid,
			"status": 1,
		}).
		Where("pid", oldPid).
		Update()
	return err
}

// FindNewPid 查找新的 PID
func (s *SJpid) FindNewPid(ctx context.Context, project *entity.Jpid) (int, error) {
	// 等待一些时间让进程完全启动
	time.Sleep(2 * time.Second)

	// 获取当前运行的 Java 进程
	processes, err := javaprocess.GetJavaProcesses()
	if err != nil {
		return 0, err
	}

	// 通过项目名称和端口匹配新进程
	for _, process := range processes {
		// 检查进程名称
		if process.Name == project.Name {
			// 检查端口匹配
			projectPorts := strings.Split(project.Ports, ",")
			processPorts := strings.Split(process.Ports, ",")

			if portListsMatch(projectPorts, processPorts) {
				return process.Pid, nil
			}
		}
	}

	return 0, gerror.New("未找到匹配的新进程")
}

// portListsMatch 检查端口列表是否匹配
func portListsMatch(ports1, ports2 []string) bool {
	if len(ports1) == 0 || len(ports2) == 0 {
		return false
	}

	portMap := make(map[string]bool)
	for _, port := range ports1 {
		if port != "" {
			portMap[port] = true
		}
	}

	for _, port := range ports2 {
		if port != "" && portMap[port] {
			return true
		}
	}

	return false
}

// 在 internal/service/jpid.go 中添加 Stop 方法的优化
func (s *SJpid) Stop(ctx context.Context, pid int) error {
	// 先尝试使用 SIGTERM 信号优雅终止
	cmd := exec.Command("kill", "-15", fmt.Sprintf("%d", pid))
	if err := cmd.Run(); err == nil {
		// 等待一段时间检查进程是否终止
		time.Sleep(500 * time.Millisecond)
		if !s.IsProcessRunning(pid) {
			return nil
		}
	}

	// 如果优雅终止失败，使用 SIGKILL 强制终止
	cmd = exec.Command("kill", "-9", fmt.Sprintf("%d", pid))
	if err := cmd.Run(); err != nil {
		return gerror.Wrapf(err, "终止进程失败: %d", pid)
	}

	return nil
}

// 添加检查进程是否运行的辅助方法
func (s *SJpid) IsProcessRunning(pid int) bool {
	cmd := exec.Command("kill", "-0", fmt.Sprintf("%d", pid))
	return cmd.Run() == nil
}
