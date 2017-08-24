package statedetector

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

//go:generate counterfeiter . ProcessStateCounter
//go:generate counterfeiter . ProcessStateReporter
//go:generate counterfeiter . XfsTracer
//go:generate counterfeiter . StateDetector
//go:generate counterfeiter . CommandRunner
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
	DetectedProcesses() ([]int, []string, error)
}

type CommandRunner interface {
	Run(*exec.Cmd) error
}

func NewShowMeWhatYouGot(
	processStateCounter ProcessStateCounter,
	processStateReporter ProcessStateReporter,
	xfsTracer XfsTracer,
	persistentStateDetector StateDetector,
	currentStateDetector StateDetector,
	reporterBackoffDuration time.Duration,
) *ShowMeWhatYouGot {
	return &ShowMeWhatYouGot{
		processStateCounter:     processStateCounter,
		processStateReporter:    processStateReporter,
		xfsTracer:               xfsTracer,
		persistentStateDetector: persistentStateDetector,
		currentStateDetector:    currentStateDetector,
		reporterBackoffDuration: reporterBackoffDuration,
	}
}

type ShowMeWhatYouGot struct {
	processStateCounter     ProcessStateCounter
	processStateReporter    ProcessStateReporter
	xfsTracer               XfsTracer
	persistentStateDetector StateDetector
	currentStateDetector    StateDetector

	reporterBackoffDuration time.Duration
	timeOfLastReport        time.Time
}

func (s *ShowMeWhatYouGot) Run() error {
	currentPids, currentProcesses, err := s.currentStateDetector.DetectedProcesses()

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

	if time.Since(s.timeOfLastReport) > s.reporterBackoffDuration {
		s.timeOfLastReport = time.Now()

		if err := s.processStateReporter.Run(persistentPids, currentProcesses); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed report states (%s)\n", err.Error())
		}
	}

	if err := s.xfsTracer.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed run xfs tracer (%s)\n", err.Error())
	}

	return nil
}
