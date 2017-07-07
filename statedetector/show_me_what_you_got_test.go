package statedetector_test

import (
	"errors"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
	"github.com/masters-of-cats/showmewhatyougot/statedetector/statedetectorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShowMeWhatYouGot", func() {

	var (
		showMeWhatYouGot        *statedetector.ShowMeWhatYouGot
		processStateCounter     *statedetectorfakes.FakeProcessStateCounter
		processStateReporter    *statedetectorfakes.FakeProcessStateReporter
		persistentStateDetector *statedetectorfakes.FakeStateDetector
		xfsTracer               *statedetectorfakes.FakeXfsTracer
	)

	BeforeEach(func() {
		processStateCounter = new(statedetectorfakes.FakeProcessStateCounter)
		processStateReporter = new(statedetectorfakes.FakeProcessStateReporter)
		persistentStateDetector = new(statedetectorfakes.FakeStateDetector)
		xfsTracer = new(statedetectorfakes.FakeXfsTracer)

		showMeWhatYouGot = statedetector.NewShowMeWhatYouGot(processStateCounter, processStateReporter, xfsTracer, persistentStateDetector)
	})

	Context("when no processes are detected", func() {
		BeforeEach(func() {
			persistentStateDetector.PidsReturns([]int{}, nil)
		})

		It("does run count", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(processStateCounter.RunCallCount()).To(Equal(1))
		})

		It("does not run xfs trace", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(xfsTracer.RunCallCount()).To(BeZero())
		})

		It("does not run report", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(processStateReporter.RunCallCount()).To(BeZero())
		})
	})

	Context("when processes are detected", func() {
		BeforeEach(func() {
			persistentStateDetector.PidsReturns([]int{
				10, 100, 50, 25,
			}, nil)
		})

		It("does run count", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(processStateCounter.RunCallCount()).To(Equal(1))
		})

		It("does run xfs trace", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(xfsTracer.RunCallCount()).To(Equal(1))
		})

		It("does run report", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(processStateReporter.RunCallCount()).To(Equal(1))
		})

		It("reports the correct pids", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(processStateReporter.RunCallCount()).To(Equal(1))

			pids, _ := processStateReporter.RunArgsForCall(0)
			Expect(pids).To(ConsistOf(10, 100, 50, 25))
		})

		Context("when processes are detected a second time", func() {
			It("still runs count", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(processStateCounter.RunCallCount()).To(Equal(2))
			})

			It("still runs xfs trace", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(xfsTracer.RunCallCount()).To(Equal(2))
			})

			It("doesn't run report anymore", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(processStateReporter.RunCallCount()).To(Equal(1))
			})
		})

		Context("when process state reporter returns an error", func() {
			BeforeEach(func() {
				processStateReporter.RunReturns(errors.New("failed"))
			})

			It("doesn't fail", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
			})
		})

		Context("when xfs tracer returns an error", func() {
			BeforeEach(func() {
				xfsTracer.RunReturns(errors.New("failed"))
			})

			It("doesn't fail", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
			})
		})
	})

	Context("when process state counter returns an error", func() {
		BeforeEach(func() {
			processStateCounter.RunReturns(errors.New("failed"))
		})

		It("doesn't fail", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
		})
	})

	Context("when persistent state detector returns an error", func() {
		BeforeEach(func() {
			persistentStateDetector.PidsReturns(nil, errors.New("failed to detect"))
		})

		It("fails", func() {
			err := showMeWhatYouGot.Run()
			Expect(err).To(MatchError(ContainSubstring("failed to detect")))
		})
	})

})
