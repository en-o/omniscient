package common

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"omniscient/internal/service"
)

// PrintWelcomeInfo 打印欢迎信息
func PrintWelcomeInfo(ctx g.Ctx) {
	port := g.Cfg().MustGet(ctx, "server.address").String()
	if port == "" {
		port = ":7777"
	}
	if port[0] == ':' {
		port = port[1:]
	}

	// 显示数据库类型信息
	dbType := g.Cfg().MustGet(ctx, "database.type", "mysql").String()

	g.Log().Info(ctx, "===================================")
	g.Log().Info(ctx, "欢迎使用项目管理系统")
	g.Log().Infof(ctx, "数据库类型: %s", dbType)
	g.Log().Info(ctx, "后台管理地址: http://localhost:"+port+"/html/pm.html")
	g.Log().Info(ctx, "===================================")
}

// PrintShellHelp 打印帮助信息
func PrintShellHelp() {
	fmt.Println("Service Management Commands (requires root privileges):")
	fmt.Println("Usage: sudo omniscient sh <command>")
	fmt.Println("  status     - Show service status")
	fmt.Println("  enable     - Enable service auto-start")
	fmt.Println("  disable    - Disable service auto-start")
	fmt.Println("  start      - Start service")
	fmt.Println("  stop       - Stop service")
	fmt.Println("  restart    - Restart service")
	fmt.Println("  reload     - Reload service")
	fmt.Println("  install    - Install systemd service")
	fmt.Println("  uninstall  - Uninstall systemd service")
	fmt.Println("  config [file] - Set default config file or show current config")
}

// InitDatabase 初始化数据库
func InitDatabase(ctx g.Ctx) error {
	g.Log().Info(ctx, "开始初始化数据库...")

	dbManager := service.NewDatabaseManager()

	// 初始化数据库连接
	if err := dbManager.Initialize(ctx); err != nil {
		g.Log().Errorf(ctx, "数据库连接初始化失败: %v", err)
		return err
	}

	// 创建必要的数据表（智能检查，避免重复创建）
	if err := dbManager.CreateTables(ctx); err != nil {
		g.Log().Errorf(ctx, "数据表创建失败: %v", err)
		return err
	}

	g.Log().Info(ctx, "数据库初始化完成")
	return nil
}

// ShowDatabaseInfo 显示数据库信息
func ShowDatabaseInfo(ctx g.Ctx) error {
	dbManager := service.NewDatabaseManager()

	// 初始化数据库连接
	if err := dbManager.Initialize(ctx); err != nil {
		return err
	}

	// 获取数据库信息
	info, err := dbManager.GetDatabaseInfo(ctx)
	if err != nil {
		return err
	}

	g.Log().Info(ctx, "===== 数据库信息 =====")
	g.Log().Infof(ctx, "数据库类型: %v", info["type"])
	g.Log().Infof(ctx, "连接字符串: %v", info["link"])
	if version, ok := info["version"]; ok {
		g.Log().Infof(ctx, "数据库版本: %v", version)
	}
	if database, ok := info["database"]; ok {
		g.Log().Infof(ctx, "数据库名: %v", database)
	}
	if dbFile, ok := info["database_file"]; ok {
		g.Log().Infof(ctx, "数据库文件: %v", dbFile)
	}
	if tables, ok := info["tables"]; ok {
		g.Log().Infof(ctx, "现有数据表: %v", tables)
	}
	g.Log().Info(ctx, "=====================")

	return nil
}
