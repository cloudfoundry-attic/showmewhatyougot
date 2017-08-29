package statedetector_test

import (
	"errors"

	"code.cloudfoundry.org/showmewhatyougot/statedetector"
	"code.cloudfoundry.org/showmewhatyougot/statedetector/statedetectorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BinaryEventEmitter", func() {
	var (
		eventEmitter  statedetector.EventEmitter
		commandRunner *statedetectorfakes.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = new(statedetectorfakes.FakeCommandRunner)
		eventEmitter = statedetector.NewBinaryEventEmitter(commandRunner, "/hello")
	})

	Describe("Run", func() {
		It("executes the binary with the correct arguments", func() {
			Expect(eventEmitter.Run()).To(Succeed())

			Expect(commandRunner.RunCallCount()).To(Equal(1))
			cmd := commandRunner.RunArgsForCall(0)
			Expect(cmd.Args).To(Equal([]string{"/hello"}))
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				err := eventEmitter.Run()
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})
	})
})
