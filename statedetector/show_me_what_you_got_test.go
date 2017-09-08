package statedetector_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/showmewhatyougot/statedetector"
	"code.cloudfoundry.org/showmewhatyougot/statedetector/statedetectorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("ShowMeWhatYouGot", func() {

	var (
		showMeWhatYouGot        *statedetector.ShowMeWhatYouGot
		processStateCounter     *statedetectorfakes.FakeProcessStateCounter
		dataCollector           *statedetectorfakes.FakeDataCollector
		persistentStateDetector *statedetectorfakes.FakeStateDetector
		currentStateDetector    *statedetectorfakes.FakeStateDetector
		xfsTracer               *statedetectorfakes.FakeXfsTracer
		eventEmitter            *statedetectorfakes.FakeEventEmitter
		reporterBackoffDuration time.Duration
		errorBuffer             *gbytes.Buffer
	)

	BeforeEach(func() {
		processStateCounter = new(statedetectorfakes.FakeProcessStateCounter)
		dataCollector = new(statedetectorfakes.FakeDataCollector)
		persistentStateDetector = new(statedetectorfakes.FakeStateDetector)
		currentStateDetector = new(statedetectorfakes.FakeStateDetector)
		xfsTracer = new(statedetectorfakes.FakeXfsTracer)
		eventEmitter = new(statedetectorfakes.FakeEventEmitter)
		reporterBackoffDuration = 1 * time.Second
		errorBuffer = gbytes.NewBuffer()

		showMeWhatYouGot = statedetector.NewShowMeWhatYouGot(processStateCounter, dataCollector, xfsTracer, persistentStateDetector, currentStateDetector, eventEmitter, reporterBackoffDuration, errorBuffer)
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
			Expect(dataCollector.RunCallCount()).To(BeZero())
		})

		It("does not emit an event", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(eventEmitter.RunCallCount()).To(BeZero())
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

		Context("whan the process state counter returns an error", func() {
			BeforeEach(func() {
				processStateCounter.RunReturns(errors.New("failed"))
			})

			It("logs an error but doesn't fail", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(processStateCounter.RunCallCount()).To(Equal(1))
				Expect(errorBuffer).To(gbytes.Say("Failed to publish state counter"))
			})
		})

		It("does run xfs trace", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(xfsTracer.RunCallCount()).To(Equal(1))
		})

		Context("when xfs tracer returns an error", func() {
			BeforeEach(func() {
				xfsTracer.RunReturns(errors.New("failed"))
			})

			It("doesn't fail", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(xfsTracer.RunCallCount()).To(Equal(1))
				Expect(errorBuffer).To(gbytes.Say("Failed to run xfs tracer"))
			})
		})

		It("does run report", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(dataCollector.RunCallCount()).To(Equal(1))
		})

		It("reports the correct pids", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(dataCollector.RunCallCount()).To(Equal(1))

			pids, _ := dataCollector.RunArgsForCall(0)
			Expect(pids).To(ConsistOf(10, 100, 50, 25))
		})

		It("reports the correct processes", func() {
			currentStateDetector.DetectedProcessesReturns(nil, []string{"highway to", "the danger", "zone"}, nil)

			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(dataCollector.RunCallCount()).To(Equal(1))

			_, processes := dataCollector.RunArgsForCall(0)
			Expect(processes).To(ConsistOf("highway to", "the danger", "zone"))
		})

		Context("when the data collector returns an error", func() {
			BeforeEach(func() {
				dataCollector.RunReturns("", errors.New("failed"))
			})

			It("doesn't fail", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(dataCollector.RunCallCount()).To(Equal(1))
				Expect(errorBuffer).To(gbytes.Say("Failed to collect debug data"))
			})
		})

		It("emits an event", func() {
			dataCollector.RunReturns("/path/to/data.tar.gz", nil)

			Expect(showMeWhatYouGot.Run()).To(Succeed())

			Expect(eventEmitter.RunCallCount()).To(Equal(1))
			Expect(eventEmitter.RunArgsForCall(0)).To(Equal("/path/to/data.tar.gz"))
		})

		Context("when event emitter returns an error", func() {
			BeforeEach(func() {
				eventEmitter.RunReturns(errors.New("failed"))
			})

			It("logs an error but doesn't fail", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(eventEmitter.RunCallCount()).To(Equal(1))
				Expect(errorBuffer).To(gbytes.Say("Failed to emit an event"))
			})
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

			It("doesn't run report until the reporter is reset", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(dataCollector.RunCallCount()).To(Equal(1))
				time.Sleep(time.Second * 2)
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(dataCollector.RunCallCount()).To(Equal(2))
			})

			It("doesn't emit an event until the reporter is reset", func() {
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(eventEmitter.RunCallCount()).To(Equal(1))
				time.Sleep(time.Second * 2)
				Expect(showMeWhatYouGot.Run()).To(Succeed())
				Expect(eventEmitter.RunCallCount()).To(Equal(2))
			})
		})
	})

	Context("when persistent state detector returns an error", func() {
		BeforeEach(func() {
			persistentStateDetector.PidsReturns(nil, errors.New("failed to detect"))
		})

		It("fails", func() {
			err := showMeWhatYouGot.Run()
			Expect(persistentStateDetector.PidsCallCount()).To(Equal(1))
			Expect(err).To(MatchError(ContainSubstring("failed to detect")))
		})
	})

	Context("when current state detector returns an error", func() {
		BeforeEach(func() {
			currentStateDetector.DetectedProcessesReturns(nil, nil, errors.New("failed to detect"))
		})

		It("doesn't fail", func() {
			Expect(showMeWhatYouGot.Run()).To(Succeed())
			Expect(currentStateDetector.DetectedProcessesCallCount()).To(Equal(1))
		})
	})

})
