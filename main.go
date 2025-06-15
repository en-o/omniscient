package main

import (
	"github.com/gogf/gf/v2/frame/g"
	_ "omniscient/internal/packed"
	"omniscient/internal/util/common"
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"

	"omniscient/internal/cmd"
)

func main() {
	// 获取命令行参数
	ctx := gctx.GetInitCtx()
	args := os.Args

	// 直接处理命令路由，不使用 GoFrame 的复杂命令系统
	if len(args) <= 1 {
		// 没有参数时默认运行 run 命令
		cmd.Main.Run(ctx)
		return
	}

	// 手动处理命令路由
	switch args[1] {
	case "run":
		err := cmd.Run.Func(ctx, nil)
		if err != nil {
			g.Log().Error(ctx, "运行服务器失败=========================")
			return
		}
	case "sh":
		// 服务管理命令
		err := cmd.Shell.Func(ctx, nil)
		if err != nil {
			return
		}
	case "dbinfo":
		// 显示数据库信息
		if err := common.ShowDatabaseInfo(ctx); err != nil {
			g.Log().Error(ctx, "获取数据库信息失败:", err)
			return
		}
	default:
		// 使用 GoFrame 命令系统作为备用
		command := gcmd.Command{
			Name:  "omniscient",
			Usage: "omniscient [command] [arguments]",
			Brief: "Omniscient service management tool",
			Description: `
Omniscient is a web service with comprehensive management capabilities.

Available Commands:
  run      - Run the HTTP server (default)
  sh       - Usage: sudo omniscient sh <command> (Service management shell commands)
  dbinfo   - Show database information

Examples:
  omniscient              # Run the server (default)
  omniscient run          # Run the server explicitly  
  omniscient dbinfo       # Show database configuration and status
  omniscient sh status    # Show service status
  omniscient sh install   # Install systemd service
  omniscient sh uninstall # uninstall systemd service
  omniscient sh enable    # Enable auto-start
  omniscient sh disable   # Disable auto-start
  omniscient sh start     # Start service
  omniscient sh stop      # stop service
  omniscient sh reload    # reload service
  omniscient sh restart   # restart service
  omniscient sh config [file] # Set default config file or show current config
`,
		}

		// 添加子命令
		err := command.AddCommand(&cmd.Run, &cmd.Shell)
		if err != nil {
			g.Log().Error(ctx, "子命令运行失败=========================")
			return
		}
		command.Run(ctx)
	}
}
