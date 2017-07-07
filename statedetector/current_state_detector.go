package statedetector

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Running current state detector: %s: %s", err.Error(), stdoutBuffer.String())
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
