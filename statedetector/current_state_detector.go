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

func (p *currentStateDetector) RunPS() ([]int, []string, error) {
	cmd := exec.Command("ps", "axho", "pid,state,comm")
	stdoutBuffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("Running current state detector: %s: %s", err.Error(), stdoutBuffer.String())
	}

	lines := []string{}
	pids := []int{}
	scanner := bufio.NewScanner(stdoutBuffer)
	for scanner.Scan() {
		var (
			pid   int
			state string
			comm  string
		)

		line := scanner.Text()
		_, _ = fmt.Sscanf(line, "%d %s %s", &pid, &state, &comm)
		if state == p.state {
			pids = append(pids, pid)
			lines = append(lines, line)
		}
	}

	return pids, lines, nil
}

func (p *currentStateDetector) Pids(pids []int) ([]int, error) {
	return []int{}, nil
}
