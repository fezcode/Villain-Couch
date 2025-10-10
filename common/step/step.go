package step

import (
	"fmt"
	"vlc-tracker-agent/common/logger"
	"vlc-tracker-agent/common/runtime"
)

type Step struct {
	F func(param ...string) error
	P *string
}

//// Step is a function type representing a single step in a process.
//// It returns an error to indicate success or failure.
//type Step func() error

func RunSteps(steps []Step) error {
	location, _ := runtime.GetCallerGrandparent()
	for i, step := range steps {
		p := ""
		if step.P != nil {
			p = *step.P
		}
		if err := step.F(p); err != nil {
			logger.Log.Error(err.Error(), "location", location, "step", i)
			return fmt.Errorf("could not run step %d: %w", i, err)
		}
	}
	return nil
}
