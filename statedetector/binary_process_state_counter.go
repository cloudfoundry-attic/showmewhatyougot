package statedetector

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type binaryProcessStateCounter struct {
	path string
}

func NewBinaryProcessStateCounter(binPath string) ProcessStateCounter {
	return &binaryProcessStateCounter{
		path: binPath,
	}
}

func (b *binaryProcessStateCounter) Run(count int) error {
	cmd := exec.Command(b.path, strconv.Itoa(count))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Running process state counter: %s.", err.Error())
	}

	return nil
}
