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
	Run(int) error
}

type ProcessStateReporter interface {
	Run(pidList []int, processesList []string) error
}

type XfsTracer interface {
	Run() error
	Start() error
	Stop() error
}

type StateDetector interface {
	Pids([]int) ([]int, error)
	RunPS() ([]int, []string, error)
}

func NewShowMeWhatYouGot(
	processStateCounter ProcessStateCounter,
	processStateReporter ProcessStateReporter,
	xfsTracer XfsTracer,
	persistentStateDetector StateDetector,
	currentStateDetector StateDetector,
) *ShowMeWhatYouGot {
	return &ShowMeWhatYouGot{
		processStateCounter:     processStateCounter,
		processStateReporter:    processStateReporter,
		xfsTracer:               xfsTracer,
		persistentStateDetector: persistentStateDetector,
		currentStateDetector:    currentStateDetector,
	}
}

type ShowMeWhatYouGot struct {
	processStateCounter     ProcessStateCounter
	processStateReporter    ProcessStateReporter
	xfsTracer               XfsTracer
	persistentStateDetector StateDetector
	currentStateDetector    StateDetector

	alreadyReported bool
}

func (s *ShowMeWhatYouGot) Run() error {
	currentPids, currentProcesses, err := s.currentStateDetector.RunPS()

	if err := s.processStateCounter.Run(len(currentPids)); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to publish state counter (%s)\n", err.Error())
	}

	persistentPids, err := s.persistentStateDetector.Pids(currentPids)
	if err != nil {
		return err
	}

	if len(persistentPids) == 0 {
		return nil
	}

	if !s.alreadyReported {
		s.alreadyReported = true

		if err := s.processStateReporter.Run(persistentPids, currentProcesses); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed report states (%s)\n", err.Error())
		}
	}

	if err := s.xfsTracer.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed run xfs tracer (%s)\n", err.Error())
	}

	return nil
}
