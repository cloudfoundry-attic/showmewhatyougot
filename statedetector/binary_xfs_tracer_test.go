package statedetector_test

import (
	"errors"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
	"github.com/masters-of-cats/showmewhatyougot/statedetector/statedetectorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BinaryXFSTracer", func() {
	var (
		xfsTracer     statedetector.XfsTracer
		commandRunner *statedetectorfakes.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = new(statedetectorfakes.FakeCommandRunner)
		xfsTracer = statedetector.NewBinaryXfsTracer(commandRunner, "/hello")
	})

	Describe("Run", func() {
		It("executes the binary", func() {
			Expect(xfsTracer.Run()).To(Succeed())

			Expect(commandRunner.RunCallCount()).To(Equal(1))
			cmd := commandRunner.RunArgsForCall(0)
			Expect(cmd.Args).To(Equal([]string{"/hello", "extract"}))
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				err := xfsTracer.Run()
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})
	})

	Describe("Start", func() {
		It("executes the binary", func() {
			Expect(xfsTracer.Start()).To(Succeed())

			Expect(commandRunner.RunCallCount()).To(Equal(1))
			cmd := commandRunner.RunArgsForCall(0)
			Expect(cmd.Args).To(Equal([]string{"/hello", "start"}))
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				err := xfsTracer.Start()
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})
	})

	Describe("Stop", func() {
		It("executes the binary", func() {
			Expect(xfsTracer.Stop()).To(Succeed())

			Expect(commandRunner.RunCallCount()).To(Equal(1))
			cmd := commandRunner.RunArgsForCall(0)
			Expect(cmd.Args).To(Equal([]string{"/hello", "stop"}))
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				err := xfsTracer.Stop()
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})
	})
})
