package statedetector

import (
	"fmt"
	"os"
)

//go:generate counterfeiter . ProcessStateCounter
//go:generate counterfeiter . ProcessStateReporter
//go:generate counterfeiter . XfsTracer
//go:generate counterfeiter . StateDetector
type ProcessStateCounter interface {
	Run() error
}

type ProcessStateReporter interface {
	Run(pidList []int, processesList []string) error
}

type XfsTracer interface {
	Run() error
}

type StateDetector interface {
	Pids() ([]int, error)
}

func NewShowMeWhatYouGot(
	processStateCounter ProcessStateCounter,
	processStateReporter ProcessStateReporter,
	xfsTracer XfsTracer,
	persistentStateDetector StateDetector,
) *ShowMeWhatYouGot {
	return &ShowMeWhatYouGot{
		processStateCounter:     processStateCounter,
		processStateReporter:    processStateReporter,
		xfsTracer:               xfsTracer,
		persistentStateDetector: persistentStateDetector,
	}
}

type ShowMeWhatYouGot struct {
	processStateCounter     ProcessStateCounter
	processStateReporter    ProcessStateReporter
	xfsTracer               XfsTracer
	persistentStateDetector StateDetector

	alreadyReported bool
}

func (s *ShowMeWhatYouGot) Run() error {
	if err := s.processStateCounter.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to publish state counter (%s)\n", err.Error())
	}

	pids, err := s.persistentStateDetector.Pids()
	if err != nil {
		return err
	}

	if len(pids) != 0 {

		if !s.alreadyReported {
			s.alreadyReported = true
			if err := s.processStateReporter.Run(pids, nil); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed report states (%s)\n", err.Error())
			}
		}
		if err := s.xfsTracer.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed run xfs tracer (%s)\n", err.Error())
		}
	}

	return nil
}
