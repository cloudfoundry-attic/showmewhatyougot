package statedetector

import (
	"bytes"
	"fmt"
	"os/exec"
)

type binaryXfsTracer struct {
	path string
}

func NewBinaryXfsTracer(binPath string) XfsTracer {
	return &binaryXfsTracer{
		path: binPath,
	}
}

func (b *binaryXfsTracer) Run() error {
	cmd := exec.Command(b.path, "extract")

	stdoutBuffer := bytes.NewBuffer([]byte{})
	stderrBuffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Running xfs tracer: %s: %s - %s", err.Error(), stdoutBuffer.String(), stderrBuffer.String())
	}

	return nil
}
