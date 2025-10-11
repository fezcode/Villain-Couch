package operations

import (
	"fmt"
	"os"
)

type PrintVersion struct {
	Version string
}

func (a PrintVersion) Priority() int {
	return OrderHigh
}

func (a PrintVersion) DefaultError() string { return "" }

func (a PrintVersion) Name() string {
	return "Print Version"
}

func (a PrintVersion) Run() error {
	fmt.Printf("Version: %s", a.Version)
	os.Exit(0)
	return nil
}

func (a PrintVersion) Finalize() {}
