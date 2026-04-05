//go:build gobake
package bake_recipe

import (
	"fmt"
	"github.com/fezcode/gobake"
)

func Run(bake *gobake.Engine) error {
	if err := bake.LoadRecipeInfo("recipe.piml"); err != nil {
		return err
	}

	bake.Task("build", "Builds the binary for the agent", func(ctx *gobake.Context) error {
		ldflags := fmt.Sprintf("-X main.Version=%s", bake.Info.Version)
		output := "bin/" + bake.Info.Name + ".exe"

		// Ensure bin/ exists
		_ = ctx.Mkdir("bin")

		return ctx.Run("go", "build", "-ldflags", ldflags, "-o", output, "./agent/src/")
	})

	bake.Task("test", "Runs tests", func(ctx *gobake.Context) error {
		return ctx.Run("go", "test", "./...")
	})

	bake.Task("lint", "Runs linter", func(ctx *gobake.Context) error {
		return ctx.Run("golangci-lint", "run")
	})

	return nil
}
