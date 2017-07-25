package statedetector_test

import (
	"os"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CurrentStateDetector", func() {
	var currentStateDetector statedetector.StateDetector

	Describe("RunPS", func() {
		Context("when a process with the given state is detected", func() {
			BeforeEach(func() {
				currentStateDetector = statedetector.NewCurrentStateDetector("S")
			})

			It("returns an array of PIDS", func() {
				pids, _, err := currentStateDetector.RunPS()
				Expect(err).NotTo(HaveOccurred())
				Expect(pids).NotTo(BeEmpty())
			})

			It("returns an array of proccesses", func() {
				_, proccesses, err := currentStateDetector.RunPS()
				Expect(err).NotTo(HaveOccurred())
				Expect(proccesses).NotTo(BeEmpty())
			})
		})

		Context("when there are no processes in the given state", func() {
			BeforeEach(func() {
				currentStateDetector = statedetector.NewCurrentStateDetector("Q")
			})

			It("returns empty arrays", func() {
				pids, proccesses, err := currentStateDetector.RunPS()
				Expect(err).NotTo(HaveOccurred())
				Expect(pids).To(BeEmpty())
				Expect(proccesses).To(BeEmpty())
			})
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				currentStateDetector = statedetector.NewCurrentStateDetector("S")
				Expect(os.Setenv("PATH", "kitten")).To(Succeed())
			})

			It("returns an error", func() {
				pids, proccesses, err := currentStateDetector.RunPS()
				Expect(pids).To(BeEmpty())
				Expect(proccesses).To(BeEmpty())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Running current state detector"))
			})
		})
	})
})
