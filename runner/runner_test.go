package runner_test

import (
	"os/exec"
	"time"

	"github.com/masters-of-cats/showmewhatyougot/runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Runner", func() {
	var (
		commandRunner *runner.CommandRunner
		timeout       time.Duration
	)

	BeforeEach(func() {
		timeout = 2 * time.Second
		commandRunner = runner.New(timeout)
	})

	Describe("Run", func() {
		It("runs the command with success", func() {
			cmd := exec.Command("echo", "hello")
			stdoutBuffer := gbytes.NewBuffer()
			cmd.Stdout = stdoutBuffer
			Expect(commandRunner.Run(cmd)).To(Succeed())
			Eventually(stdoutBuffer).Should(gbytes.Say("hello"))
		})

		Context("when the time out is reached", func() {
			It("returns a time out error", func() {
				cmd := exec.Command("sleep", "5")
				startedAt := time.Now()
				err := commandRunner.Run(cmd)
				Expect(time.Since(startedAt)).To(BeNumerically("~", 2*time.Second, 500*time.Millisecond))
				Expect(err).To(MatchError(ContainSubstring("command took more than 2.000000 seconds to finish")))
			})
		})
	})

})
