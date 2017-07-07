package statedetector

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type binaryProcessStateReporter struct {
	path string
}

func NewBinaryProcessStateReporter(binPath string) ProcessStateReporter {
	return &binaryProcessStateReporter{
		path: binPath,
	}
}

func (b *binaryProcessStateReporter) Run(pidList []int, processesList []string) error {
	pidListArgs := []string{}
	for pid := range pidList {
		pidListArgs = append(pidListArgs, strconv.Itoa(pid))
	}

	args := []string{
		strings.Join(pidListArgs, " "),
		strings.Join(processesList, " "),
	}

	cmd := exec.Command(b.path, args...)

	stdoutBuffer := bytes.NewBuffer([]byte{})
	stderrBuffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Running process state reporter: %s: %s - %s", err.Error(), stdoutBuffer.String(), stderrBuffer.String())
	}

	return nil
}
