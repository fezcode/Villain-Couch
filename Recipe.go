package main

import (
	"fmt"
	"github.com/thefezcodex/gobake/bake"
)

func init() {
	bake.Recipe(func(ctx *bake.Context) {
		ctx.Task("build", func(ctx *bake.Context) {
			ldflags := fmt.Sprintf("-X main.Version=%s", bake.Info.Version)
			output := "bin/villain-couch.exe"
			err := ctx.Run("go", "build", "-ldflags", ldflags, "-o", output, "./agent/src/")
			if err != nil {
				ctx.Log.Error("Build failed: %v", err)
			}
		})

		ctx.Task("test", func(ctx *bake.Context) {
			ctx.Run("go", "test", "./...")
		})

		ctx.Task("lint", func(ctx *bake.Context) {
			ctx.Run("golangci-lint", "run")
		})
	})
}
