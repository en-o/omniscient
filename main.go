package main

import (
	_ "omniscient/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"omniscient/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
