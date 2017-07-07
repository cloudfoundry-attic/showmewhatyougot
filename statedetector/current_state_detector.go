package statedetector

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

func NewCurrentStateDetector(state string) *currentStateDetector {
	return &currentStateDetector{
		state: state,
	}
}

type currentStateDetector struct {
	state string
}

func (p *currentStateDetector) Pids() ([]int, error) {
	cmd := exec.Command("ps", "axho", "pid,state")
	stdoutBuffer := bytes.NewBuffer([]byte{})
	stderrBuffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Running current state detector: %s: %s - %s", err.Error(), stdoutBuffer.String(), stderrBuffer.String())
	}

	pids := []int{}
	scanner := bufio.NewScanner(stdoutBuffer)
	for scanner.Scan() {
		var (
			pid   int
			state string
		)

		_, _ = fmt.Sscanf(scanner.Text(), "%d %s", &pid, &state)
		if state == p.state {
			pids = append(pids, pid)
		}
	}

	return pids, nil
}
