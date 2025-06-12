package main

import (
	"github.com/gogf/gf/v2/frame/g"
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
		printWelcomeInfo(ctx)
		// 没有参数时默认运行 run 命令
		cmd.Main.Run(ctx)
		return
	}

	// 手动处理命令路由
	switch args[1] {
	case "run":
		printWelcomeInfo(ctx)
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
		err := command.AddCommand(&cmd.Run, &cmd.Shell)
		if err != nil {
			g.Log().Error(ctx, "子命令运行失败=========================")
			return
		}
		command.Run(ctx)
	}
}

// printWelcomeInfo 打印欢迎信息
func printWelcomeInfo(ctx g.Ctx) {
	port := g.Cfg().MustGet(ctx, "server.address").String()
	if port == "" {
		port = ":7777"
	}
	if port[0] == ':' {
		port = port[1:]
	}
	g.Log().Info(ctx, "===================================")
	g.Log().Info(ctx, "欢迎使用项目管理系统")
	g.Log().Info(ctx, "后台管理地址: http://localhost:"+port+"/html/pm.html")
	g.Log().Info(ctx, "===================================")
}
