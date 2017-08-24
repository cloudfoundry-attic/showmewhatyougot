package statedetector_test

import (
	"errors"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
	"github.com/masters-of-cats/showmewhatyougot/statedetector/statedetectorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BinaryProcessStateReporter", func() {
	var (
		processStateReporter statedetector.ProcessStateReporter
		commandRunner        *statedetectorfakes.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = new(statedetectorfakes.FakeCommandRunner)
		processStateReporter = statedetector.NewBinaryProcessStateReporter(commandRunner, "/hello")
	})

	Describe("Run", func() {
		It("executes the binary with the correct arguments", func() {
			Expect(processStateReporter.Run([]int{100, 101}, []string{"foo", "bar"})).To(Succeed())

			Expect(commandRunner.RunCallCount()).To(Equal(1))
			cmd := commandRunner.RunArgsForCall(0)
			Expect(cmd.Args).To(Equal([]string{"/hello", "100 101", "foo\nbar"}))
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				err := processStateReporter.Run([]int{100, 101}, []string{"foo", "bar"})
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})
	})
})
