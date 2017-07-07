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
	cmd, stdoutBuffer, stderrBuffer := b.command("extract")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Running xfs tracer: %s: %s - %s", err.Error(), stdoutBuffer.String(), stderrBuffer.String())
	}

	return nil
}

func (b *binaryXfsTracer) Start() error {
	cmd, stdoutBuffer, stderrBuffer := b.command("start")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Starting xfs tracer: %s: %s - %s", err.Error(), stdoutBuffer.String(), stderrBuffer.String())
	}

	return nil
}

func (b *binaryXfsTracer) Stop() error {
	cmd, stdoutBuffer, stderrBuffer := b.command("stop")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Stopping xfs tracer: %s: %s - %s", err.Error(), stdoutBuffer.String(), stderrBuffer.String())
	}

	return nil
}

func (b *binaryXfsTracer) command(action string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	cmd := exec.Command(b.path, action)

	stdoutBuffer := bytes.NewBuffer([]byte{})
	stderrBuffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	return cmd, stdoutBuffer, stderrBuffer
}
