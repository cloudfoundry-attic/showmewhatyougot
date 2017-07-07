package statedetector

import (
	"fmt"
	"os"
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
	cmd := b.command("extract")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Running xfs tracer: %s", err.Error())
	}

	return nil
}

func (b *binaryXfsTracer) Start() error {
	cmd := b.command("start")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Starting xfs tracer: %s", err.Error())
	}

	return nil
}

func (b *binaryXfsTracer) Stop() error {
	cmd := b.command("stop")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Stopping xfs tracer: %s", err.Error())
	}

	return nil
}

func (b *binaryXfsTracer) command(action string) *exec.Cmd {
	cmd := exec.Command(b.path, action)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
