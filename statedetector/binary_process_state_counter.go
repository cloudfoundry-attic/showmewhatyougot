package statedetector

import (
	"bytes"
	"fmt"
	"os/exec"
)

type binaryProcessStateCounter struct {
	path string
}

func NewBinaryProcessStateCounter(binPath string) ProcessStateCounter {
	return &binaryProcessStateCounter{
		path: binPath,
	}
}

func (b *binaryProcessStateCounter) Run() error {
	cmd := exec.Command(b.path)

	stdoutBuffer := bytes.NewBuffer([]byte{})
	stderrBuffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Running process state counter: %s: %s - %s", err.Error(), stdoutBuffer.String(), stderrBuffer.String())
	}

	return nil
}
