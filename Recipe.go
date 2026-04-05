package main

import (
	"fmt"
	"github.com/fezcode/gobake"
)

func createEngine() *gobake.Engine {
	engine := gobake.New()
	_ = engine.LoadRecipeInfo("recipe.piml")

	engine.Task("build", "Build the agent", func(ctx *gobake.Context) error {
		ldflags := fmt.Sprintf("-X main.Version=%s", engine.Info.Version)
		output := "bin/villain-couch.exe"
		return ctx.Run("go", "build", "-ldflags", ldflags, "-o", output, "./agent/src/")
	})

	engine.Task("test", "Run tests", func(ctx *gobake.Context) error {
		return ctx.Run("go", "test", "./...")
	})

	engine.Task("lint", "Run linter", func(ctx *gobake.Context) error {
		return ctx.Run("golangci-lint", "run")
	})

	return engine
}
