package statedetector

import (
	"fmt"
	"os"
	"os/exec"
)

type BinaryEventEmitter struct {
	path          string
	commandRunner CommandRunner
}

func NewBinaryEventEmitter(commandRunner CommandRunner, binPath string) EventEmitter {
	return &BinaryEventEmitter{
		path:          binPath,
		commandRunner: commandRunner,
	}
}

func (b *BinaryEventEmitter) Run() error {
	cmd := exec.Command(b.path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := b.commandRunner.Run(cmd)
	if err != nil {
		return fmt.Errorf("Running event emitter: %s.", err.Error())
	}

	return nil
}
