package operations

import (
	"os"
	"path/filepath"
	"sort"
	"villain-couch/agent/src/app_info"
	"villain-couch/agent/src/cli"
	"villain-couch/agent/src/options"
	"villain-couch/agent/src/storage"
	"villain-couch/common/logger"
)

type OperationRunner struct {
	ops []Runnable
	err error
}

func New() *OperationRunner {
	return &OperationRunner{}
}

func (opr *OperationRunner) Add(runnable Runnable) {
	opr.ops = append(opr.ops, runnable)
}

func (opr *OperationRunner) Run() *OperationRunner {
	for _, op := range opr.ops {
		if err := op.Run(); err != nil {
			logger.Log.Error(op.DefaultError(), "error", err.Error())
			opr.err = err
			return opr
		}
	}
	return opr
}

func (opr *OperationRunner) Build(cliFlags *cli.CLIFlags, db *storage.DB, opts *options.Options) *OperationRunner {
	opBasics := Operation{Database: db, Options: opts}
	if !cliFlags.AddWorkspace.Empty() {
		dirName := filepath.Base(cliFlags.AddWorkspace.String())
		r := AddWorkspace{Operation: opBasics, DirPath: cliFlags.AddWorkspace.String(), DirName: dirName}
		opr.Add(r)
	}
	if cliFlags.FindNext {
		r := NextEpisode{Operation: opBasics}
		opr.Add(r)
	}
	if cliFlags.Version {
		r := PrintVersion{Version: app_info.VersionInfo}
		opr.Add(r)
	}
	return opr
}

func (opr *OperationRunner) Sort() *OperationRunner {
	sort.SliceStable(opr.ops, func(i, j int) bool {
		return opr.ops[i].Priority() > opr.ops[j].Priority()
	})
	return opr
}

func (opr *OperationRunner) Finalize() {
	if opr.err != nil {
		logger.Log.Error("Operation Runner failed")
		logger.Log.Error(opr.err.Error(), "error", opr.err, "action", "exit")
		os.Exit(1)
	}

	for _, op := range opr.ops {
		op.Finalize()
	}
}
