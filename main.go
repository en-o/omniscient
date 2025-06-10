package main

import (
	_ "omniscient/internal/packed"
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

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
		// 运行服务器
		cmd.Run.Func(ctx, nil)
	case "sh":
		// 服务管理命令
		cmd.Shell.Func(ctx, nil)
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
  sh       - Service management shell commands

Examples:
  omniscient              # Run the server (default)
  omniscient run          # Run the server explicitly  
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
		command.AddCommand(&cmd.Run, &cmd.Shell)
		command.Run(ctx)
	}
}
