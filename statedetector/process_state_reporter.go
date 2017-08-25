package statedetector

import (
	"fmt"
	"os"
	"time"
)

type processStateReporter struct {
	dataPath string
	env      func(string) string
	time     func() time.Time
}

func NewProcessStateReporter(dataPath string, timeFunc func() time.Time, envFunc func(string) string) ProcessStateReporter {
	return &processStateReporter{
		dataPath: dataPath,
		env:      envFunc,
		time:     timeFunc,
	}
}

func (r *processStateReporter) Run(pidList []int, processesList []string) (err error) {
	err = r.createDataDirectory()
	return err
}

func (r *processStateReporter) createDataDirectory() error {
	return os.MkdirAll(fmt.Sprintf(
		"%s/%s-%s-debug-info-%s",
		r.dataPath,
		r.env("CELL_NAME"),
		r.env("CELL_ID"),
		r.time().Format("2006-01-02_15-04-05"),
	), 700)
}
