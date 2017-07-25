package statedetector_test

import (
	"github.com/masters-of-cats/showmewhatyougot/statedetector"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PersistentStateDetector", func() {
	var (
		counter                 int
		persistentStateDetector statedetector.StateDetector
	)

	BeforeEach(func() {
		counter = 3
		persistentStateDetector = statedetector.NewPersistentStateDetector(counter)
	})

	Describe("when a process is in the state for at least 'count' times", func() {
		It("reports the pids as in persistent state after 'count' times", func() {
			pids, err := persistentStateDetector.Pids([]int{100})
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			pids, err = persistentStateDetector.Pids([]int{100})
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			pids, err = persistentStateDetector.Pids([]int{100})
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(ConsistOf(100))
		})
	})

	Describe("when a pid is in the state for more than 'count' times but they are not consecutive", func() {
		It("doesn't report the pid as persistent", func() {
			pids, err := persistentStateDetector.Pids([]int{100})
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			pids, err = persistentStateDetector.Pids([]int{100})
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			pids, err = persistentStateDetector.Pids([]int{666})
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())

			pids, err = persistentStateDetector.Pids([]int{100})
			Expect(err).NotTo(HaveOccurred())
			Expect(pids).To(BeEmpty())
		})
	})
})
