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
func (s *SJpid) GetList(ctx context.Context, worker string) (list []*entity.Jpid, err error) {
	// 创建查询对象
	query := dao.Jpid.Ctx(ctx)

	// 如果worker为空，使用当前服务器的worker
	if worker == "" {
		worker = system.GetWorkerName()
	}

	// 添加worker条件查询
	err = query.Where("worker", worker).
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

// UpdateStatusById 更新项目状态
func (s *SJpid) UpdateStatusById(ctx context.Context, id int, status int) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{"status": status}).
		Where("id", id).
		Update()
	return err
}

// AutoRegister 自动注册和更新Java进程
func (s *SJpid) AutoRegister(ctx context.Context, processes []*entity.LinuxPid) (total, updated, created int, err error) {
	// 获取当前服务器标识
	currentWorker := system.GetWorkerName()

	// 1. 添加超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 只获取当前服务器的已存在项目信息
	var existingProjects []*entity.Jpid
	if err = dao.Jpid.Ctx(ctx).Where("worker", currentWorker).Scan(&existingProjects); err != nil {
		return 0, 0, 0, err
	}

	// 构建当前运行进程的PID映射
	runningPids := make(map[int]bool)
	for _, process := range processes {
		runningPids[process.Pid] = true
	}

	// 更新已停止的进程状态
	for _, project := range existingProjects {
		if project.Status == 1 && !runningPids[project.Pid] {
			if err := s.UpdateStatus(ctx, project.Pid, 0); err != nil {
				g.Log().Warningf(ctx, "更新已停止项目状态失败 [Worker:%s, PID:%d]: %v",
					currentWorker, project.Pid, err)
			}
		}
	}

	// 构建复合键(worker+port)到项目的映射
	portToProject := make(map[string]*entity.Jpid)
	for _, project := range existingProjects {
		if project.Ports != "" {
			for _, port := range strings.Split(project.Ports, ",") {
				if port != "" {
					// 使用worker+port作为key
					key := fmt.Sprintf("%s:%s", project.Worker, port)
					portToProject[key] = project
				}
			}
		}
	}

	total = len(processes)
	// 处理每个进程
	for _, process := range processes {
		var existingProject *entity.Jpid
		var matchedPort string

		// 检查端口是否已存在于当前服务器
		ports := strings.Split(process.Ports, ",")
		for _, port := range ports {
			if port != "" {
				// 使用当前worker+port构建key
				key := fmt.Sprintf("%s:%s", currentWorker, port)
				if proj, exists := portToProject[key]; exists {
					existingProject = proj
					matchedPort = port
					break
				}
			}
		}

		if existingProject != nil {
			// 更新已存在记录
			if err := s.updateExistingProject(ctx, existingProject, process); err != nil {
				g.Log().Warningf(ctx, "更新项目失败 [Worker:%s, Port:%s, PID:%d]: %v",
					currentWorker, matchedPort, process.Pid, err)
				continue
			}
			updated++
		} else {
			// 创建新记录
			if err := s.createNewProject(ctx, process); err != nil {
				g.Log().Warningf(ctx, "插入新项目失败 [Worker:%s, Ports:%s, PID:%d]: %v",
					currentWorker, process.Ports, process.Pid, err)
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
		"name":    process.Name,
		"catalog": process.Catalog,
		"run":     process.Run,
		"status":  1,
		"worker":  system.GetWorkerName(), // 确保worker字段也更新
		"way":     process.Way,
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
		Way:     process.Way,
	}).Insert()
	return err
}

// UpdateProject 更新项目基础信息
func (s *SJpid) UpdateInfo(ctx context.Context, pid int, script, catalog, description string) error {
	_, err := dao.Jpid.Ctx(ctx).
		Data(g.Map{
			"script":      script,
			"catalog":     catalog,
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

// Delete 删除项目
func (s *SJpid) Delete(ctx context.Context, id int) error {
	// 先获取项目信息
	var jpid *entity.Jpid
	err := dao.Jpid.Ctx(ctx).Where("id", id).Scan(&jpid)
	if err != nil {
		return err
	}
	if jpid == nil {
		return gerror.New("项目不存在")
	}

	// 检查项目状态，只允许删除已停止的项目
	if jpid.Status == 1 {
		return gerror.New("项目正在运行中，请先停止项目后再删除")
	}

	// 执行删除操作
	_, err = dao.Jpid.Ctx(ctx).Where("id", id).Delete()
	return err
}

// UpdateAutostart 更新自启状态并处理自启服务
func (s *SJpid) UpdateAutostart(ctx context.Context, id int, autostart int) error {
	// 获取项目信息
	var jpid *entity.Jpid
	err := dao.Jpid.Ctx(ctx).Where("id", id).Scan(&jpid)
	if err != nil {
		return err
	}
	if jpid == nil {
		return gerror.New("项目不存在")
	}

	// 判断是否为 docker 方式运行
	if jpid.Way == 1 {
		return gerror.New("docker 方式运行的项目不支持设置自启动")
	}

	// 检查autostart命令是否存在
	if !isAutostartInstalled() {
		return gerror.New("请先安装autostart并设置环境变量")
	}

	var autoName = jpid.Name + "_" + jpid.Ports
	g.Log().Info(ctx, "更新自启状态", "pid", jpid.Pid, "autoName", autoName, "autostart", autostart)

	// 检查服务是否已存在
	serviceExists := checkAutostartServiceExists(ctx, autoName)

	if autostart == 1 {
		// 启用自启
		if !serviceExists {
			// 服务不存在，添加新服务
			var execStr string
			// 使用单引号包裹整个命令
			if jpid.Script == "" {
				execStr = "'" + jpid.Run + "'"
			} else {
				execStr = "'" + jpid.Script + " -b false" + "'"
			}

			g.Log().Info(ctx, "添加自启服务", "execStr", execStr)

			// 获取description,如果为空设置默认值
			description := jpid.Description
			if description == "" {
				description = "Service for " + jpid.Name
			}

			// 注册自启（使用改进的命令执行函数）
			err := execCommand(ctx, "autostart", "add", autoName, execStr,
				"--workdir="+jpid.Catalog, "--description="+description)
			if err != nil {
				return gerror.Wrap(err, "注册自启服务失败")
			}
		}

		// 启用自启（无论服务是否已存在都需要确保启用）
		err := execCommand(ctx, "autostart", "enable", autoName)
		if err != nil {
			return gerror.Wrap(err, "启用自启服务失败")
		}

		// 验证自启动服务
		err = s.verifyAutostartService(ctx, autoName)
		if err != nil {
			g.Log().Warning(ctx, "验证自启动服务失败", "error", err)
		}

	} else {
		// 移除自启
		if serviceExists {
			// 使用非交互模式移除服务
			err := s.removeAutostartServiceNonInteractive(ctx, autoName)
			if err != nil {
				return gerror.Wrap(err, "移除自启服务失败")
			}
		} else {
			g.Log().Info(ctx, "自启服务不存在，跳过移除", "autoName", autoName)
		}
	}

	// 更新数据库
	_, err = dao.Jpid.Ctx(ctx).Data(g.Map{
		"autostart": autostart,
	}).Where("id", id).Update()
	return err
}

// removeAutostartServiceNonInteractive 非交互式移除自启服务
func (s *SJpid) removeAutostartServiceNonInteractive(ctx context.Context, autoName string) error {
	// 方法1: 使用echo 'y'通过管道
	cmdCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "bash", "-c", "echo 'y' | autostart rm "+autoName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		g.Log().Warning(ctx, "方法1移除自启服务失败，尝试方法2",
			"error", err, "output", string(output))

		// 方法2: 直接禁用服务
		err3 := execSudoCommand(ctx, "autostart", "disable", autoName)
		if err3 != nil {
			return gerror.Wrapf(err, "所有移除方法都失败了，原始错误: %v", err)
		}
	}

	return nil
}

// verifyAutostartService 验证自启动服务状态
func (s *SJpid) verifyAutostartService(ctx context.Context, autoName string) error {
	cmd := exec.Command("autostart", "ls")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return gerror.Wrap(err, "获取自启动服务列表失败")
	}

	g.Log().Info(ctx, "当前自启动服务列表", "output", string(output))

	// 检查服务是否在列表中
	if !strings.Contains(string(output), autoName) {
		return gerror.New("服务未出现在自启动列表中")
	}

	return nil
}

// checkAutostartServiceExists 检查autostart服务是否存在
func checkAutostartServiceExists(ctx context.Context, autoName string) bool {
	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "autostart", "exists", autoName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		g.Log().Debug(ctx, "检查自启服务状态",
			"autoName", autoName, "error", err, "output", string(output))
		return false
	}

	// 检查输出中是否包含确认信息
	return strings.Contains(string(output), "Service '"+autoName+"' exists") ||
		strings.Contains(string(output), autoName)
}

// 检查autostart命令是否安装
func isAutostartInstalled() bool {
	cmd := exec.Command("which", "autostart")
	err := cmd.Run()
	if err != nil {
		// 也尝试检查常见安装路径
		paths := []string{
			"/usr/local/bin/autostart",
			"/usr/bin/autostart",
			"/opt/autostart/bin/autostart",
		}

		for _, path := range paths {
			if cmd := exec.Command("test", "-x", path); cmd.Run() == nil {
				return true
			}
		}
		return false
	}
	return true
}

// checkSudoNoPassword 检查是否可以无密码执行sudo命令
func checkSudoNoPassword(command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", "-n", command, "--help")
	if err := cmd.Run(); err != nil {
		return gerror.New("无法无密码执行sudo命令，请配置sudo免密:\n" +
			"1. 执行: sudo visudo\n" +
			"2. 添加: " + getCurrentUser() + " ALL=(ALL) NOPASSWD: /*/autostart\n" +
			"3. 保存并退出")
	}
	return nil
}

// getCurrentUser 获取当前用户名
func getCurrentUser() string {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return "username" // 默认值
	}
	return strings.TrimSpace(string(output))
}

// execSudoCommand 执行命令，带超时和错误处理
func execCommand(ctx context.Context, args ...string) error {
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return gerror.Wrapf(err, "执行命令失败: %s, 输出: %s",
			strings.Join(args, " "), string(output))
	}
	return nil
}

// execSudoCommand 执行sudo命令，带超时和错误处理
func execSudoCommand(ctx context.Context, args ...string) error {
	// 设置命令超时
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 构建完整的sudo命令
	fullArgs := append([]string{"sudo", "-n"}, args...)
	cmd := exec.CommandContext(cmdCtx, fullArgs[0], fullArgs[1:]...)

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return gerror.Wrapf(err, "执行sudo命令失败: %s, 输出: %s", strings.Join(args, " "), string(output))
	}

	return nil
}
