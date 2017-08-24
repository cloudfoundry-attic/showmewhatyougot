package statedetector_test

import (
	"errors"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
	"github.com/masters-of-cats/showmewhatyougot/statedetector/statedetectorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BinaryProcessStateCounter", func() {
	var (
		processStateCounter statedetector.ProcessStateCounter
		commandRunner       *statedetectorfakes.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = new(statedetectorfakes.FakeCommandRunner)
		processStateCounter = statedetector.NewBinaryProcessStateCounter(commandRunner, "/hello")
	})

	Describe("Run", func() {
		It("executes the binary with the correct arguments", func() {
			Expect(processStateCounter.Run(10)).To(Succeed())

			Expect(commandRunner.RunCallCount()).To(Equal(1))
			cmd := commandRunner.RunArgsForCall(0)
			Expect(cmd.Args).To(Equal([]string{"/hello", "10"}))
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				err := processStateCounter.Run(10)
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})
	})
})
