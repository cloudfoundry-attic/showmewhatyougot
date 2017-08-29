package runner // import "code.cloudfoundry.org/showmewhatyougot/runner"

import (
	"fmt"
	"os/exec"
	"time"
)

func New(timeout time.Duration) *CommandRunner {
	return &CommandRunner{timeout}
}

type CommandRunner struct {
	timeOut time.Duration
}

func (c *CommandRunner) Run(cmd *exec.Cmd) error {
	if c.timeOut == 0 {
		return cmd.Run()
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		errChan <- cmd.Wait()
		close(errChan)
	}()

	select {
	case runErr := <-errChan:
		return runErr

	case <-time.After(c.timeOut):
		return fmt.Errorf("command took more than %f seconds to finish", c.timeOut.Seconds())

	default:
		return nil
	}
}
