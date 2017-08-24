package statedetector

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type BinaryProcessStateCounter struct {
	path          string
	commandRunner CommandRunner
}

func NewBinaryProcessStateCounter(commandRunner CommandRunner, binPath string) ProcessStateCounter {
	return &BinaryProcessStateCounter{
		path:          binPath,
		commandRunner: commandRunner,
	}
}

func (b *BinaryProcessStateCounter) Run(count int) error {
	cmd := exec.Command(b.path, strconv.Itoa(count))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := b.commandRunner.Run(cmd)
	if err != nil {
		return fmt.Errorf("Running process state counter: %s.", err.Error())
	}

	return nil
}
