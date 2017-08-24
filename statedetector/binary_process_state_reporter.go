package statedetector

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type binaryProcessStateReporter struct {
	path          string
	commandRunner CommandRunner
}

func NewBinaryProcessStateReporter(commandRunner CommandRunner, binPath string) ProcessStateReporter {
	return &binaryProcessStateReporter{
		path:          binPath,
		commandRunner: commandRunner,
	}
}

func (b *binaryProcessStateReporter) Run(pidList []int, processesList []string) error {
	pidListArgs := []string{}
	for _, pid := range pidList {
		pidListArgs = append(pidListArgs, strconv.Itoa(pid))
	}

	args := []string{
		strings.Join(pidListArgs, " "),
		strings.Join(processesList, "\n"),
	}

	cmd := exec.Command(b.path, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := b.commandRunner.Run(cmd)
	if err != nil {
		return fmt.Errorf("Running process state reporter: %s", err.Error())
	}

	return nil
}
