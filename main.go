package main

import (
	_ "omniscient/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"

	"omniscient/internal/cmd"
)

func main() {
	// 创建命令管理器
	command := gcmd.Command{
		Name:  "omniscient",
		Usage: "omniscient [command]",
		Brief: "Omniscient service management tool",
		Description: `
Omniscient is a web service with comprehensive management capabilities.

Available Commands:
  run  - Run the HTTP server (default)
  sh   - Service management shell commands

Examples:
  omniscient              # Run the server (default)
  omniscient run          # Run the server explicitly  
  omniscient sh status    # Show service status
  omniscient sh install   # Install systemd service
  omniscient sh enable    # Enable auto-start
  omniscient sh start     # Start service
`,
	}

	// 添加子命令
	command.AddCommand(&cmd.Run, &cmd.Shell)

	// 如果没有提供参数，默认运行服务器
	ctx := gctx.GetInitCtx()
	args := gctx.GetArgs()

	if len(args) <= 1 {
		// 没有参数时默认运行 run 命令
		cmd.Main.Run(ctx)
	} else {
		// 有参数时使用命令路由
		command.Run(ctx)
	}
}
