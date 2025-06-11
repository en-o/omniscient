package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"autostart/internal/service"
	"autostart/internal/utils"
)

const (
	ToolName = "autostart"
	Version  = "0.0.2"
)

func Execute() {
	// 检查操作系统支持
	fmt.Printf("system：%s \n", runtime.GOOS)
	if runtime.GOOS != "linux" {
		fmt.Printf("Error: This tool currently only supports Linux systems (current: %s)\n", runtime.GOOS)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		utils.PrintHelp()
		return
	}

	command := os.Args[1]

	// 检查是否有足够权限
	if utils.NeedsRoot(command) && os.Geteuid() != 0 {
		fmt.Printf("Error: This operation requires root privileges. Please run with sudo:\n")
		fmt.Printf("  sudo %s %s\n", os.Args[0], strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}

	switch command {
	case "list", "ls":
		service.ListAutostartServices()
	case "add", "create", "install":
		service.AddAutostartService()
	case "remove", "rm", "delete", "uninstall":
		service.RemoveAutostartService()
	case "enable":
		service.EnableService()
	case "disable":
		service.DisableService()
	case "status":
		service.ShowServiceStatus()
	case "start":
		service.StartService()
	case "stop":
		service.StopService()
	case "restart":
		service.RestartService()
	case "logs":
		service.ShowServiceLogs()
	case "edit":
		service.EditService()
	case "version", "-v", "--version":
		fmt.Printf("%s v%s\n", ToolName, Version)
	case "help", "-h", "--help":
		utils.PrintHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		utils.PrintHelp()
		os.Exit(1)
	}
}
