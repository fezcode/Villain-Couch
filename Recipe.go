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

	bake.Task("build", "Builds the binary for multiple platforms", func(ctx *gobake.Context) error {
		ctx.Log("Building %s v%s...", bake.Info.Name, bake.Info.Version)

		targets := []struct {
			os   string
			arch string
		}{
			{"linux", "amd64"},
			{"linux", "arm64"},
			{"windows", "amd64"},
			{"windows", "arm64"},
			{"darwin", "amd64"},
			{"darwin", "arm64"},
		}

		err := ctx.Mkdir("build")
		if err != nil {
			return err
		}

		ldflags := fmt.Sprintf("-X main.Version=%s", bake.Info.Version)

		for _, t := range targets {
			output := "build/" + bake.Info.Name + "-" + t.os + "-" + t.arch
			if t.os == "windows" {
				output += ".exe"
			}

			ctx.Env = []string{
				"CGO_ENABLED=0",
				"GOOS=" + t.os,
				"GOARCH=" + t.arch,
			}

			err := ctx.Run("go", "build", "-ldflags", ldflags, "-o", output, "./agent/src/")
			if err != nil {
				return err
			}
		}
		return nil
	})

	bake.Task("test", "Runs tests", func(ctx *gobake.Context) error {
		return ctx.Run("go", "test", "./...")
	})

	bake.Task("lint", "Runs linter", func(ctx *gobake.Context) error {
		return ctx.Run("golangci-lint", "run")
	})

	return nil
}
