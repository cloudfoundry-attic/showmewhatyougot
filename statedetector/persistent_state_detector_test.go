package statedetector_test

import (
	"errors"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
	"github.com/masters-of-cats/showmewhatyougot/statedetector/statedetectorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PersistentStateDetector", func() {
	var (
		counter                 int
		persistentStateDetector statedetector.StateDetector
		currentStateDetector    *statedetectorfakes.FakeStateDetector
	)

	BeforeEach(func() {
		counter = 3
		currentStateDetector = new(statedetectorfakes.FakeStateDetector)
		persistentStateDetector = statedetector.NewPersistentStateDetector(counter, currentStateDetector)
	})

	Describe("when a process is in the state for at least 'count' times", func() {
		BeforeEach(func() {
			currentStateDetector.PidsReturns([]int{100}, nil)
		})

		It("reports the process as in persistent state after 'count' times", func() {
			pids, err := persistentStateDetector.Pids()
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			pids, err = persistentStateDetector.Pids()
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			pids, err = persistentStateDetector.Pids()
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(ConsistOf(100))
		})
	})

	Describe("when a process is in the state for more than 'count' times but they are not consecutive", func() {
		It("doesn't report the process as persistent", func() {
			currentStateDetector.PidsReturns([]int{100}, nil)
			pids, err := persistentStateDetector.Pids()
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			currentStateDetector.PidsReturns([]int{100}, nil)
			pids, err = persistentStateDetector.Pids()
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			currentStateDetector.PidsReturns([]int{1000}, nil)
			pids, err = persistentStateDetector.Pids()
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			currentStateDetector.PidsReturns([]int{100}, nil)
			pids, err = persistentStateDetector.Pids()
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())
		})
	})

	Describe("when the current process state fails", func() {
		BeforeEach(func() {
			currentStateDetector.PidsReturns(nil, errors.New("failed"))
		})

		It("returns an error", func() {
			_, err := persistentStateDetector.Pids()
			Expect(err).To(MatchError(ContainSubstring("failed")))
		})
	})

})
