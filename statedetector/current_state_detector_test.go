package statedetector_test

import (
	"errors"
	"os/exec"

	"code.cloudfoundry.org/showmewhatyougot/statedetector"
	"code.cloudfoundry.org/showmewhatyougot/statedetector/statedetectorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CurrentStateDetector", func() {
	var (
		currentStateDetector statedetector.StateDetector
		commandRunner        *statedetectorfakes.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = new(statedetectorfakes.FakeCommandRunner)
	})

	Describe("DetectedProcesses", func() {
		Context("when a process with the given state is detected", func() {
			BeforeEach(func() {
				currentStateDetector = statedetector.NewCurrentStateDetector(commandRunner, "S")

				commandRunner.RunStub = func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("100 S hello\n"))
					cmd.Stdout.Write([]byte("101 S good-bye\n"))
					return nil
				}
			})

			It("returns an array of PIDS", func() {
				pids, _, err := currentStateDetector.DetectedProcesses()
				Expect(err).NotTo(HaveOccurred())
				Expect(pids).To(Equal([]int{100, 101}))
			})

			It("returns an array of proccesses", func() {
				_, proccesses, err := currentStateDetector.DetectedProcesses()
				Expect(err).NotTo(HaveOccurred())
				Expect(proccesses).To(Equal([]string{"100 S hello", "101 S good-bye"}))
			})
		})

		Context("when there are no processes in the given state", func() {
			BeforeEach(func() {
				currentStateDetector = statedetector.NewCurrentStateDetector(commandRunner, "S")
			})

			It("returns empty arrays", func() {
				pids, proccesses, err := currentStateDetector.DetectedProcesses()
				Expect(err).NotTo(HaveOccurred())
				Expect(pids).To(BeEmpty())
				Expect(proccesses).To(BeEmpty())
			})
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				currentStateDetector = statedetector.NewCurrentStateDetector(commandRunner, "S")
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				pids, proccesses, err := currentStateDetector.DetectedProcesses()
				Expect(pids).To(BeEmpty())
				Expect(proccesses).To(BeEmpty())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Running current state detector"))
			})
		})
	})
})
