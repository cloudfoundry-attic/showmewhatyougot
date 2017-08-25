package statedetector_test

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProcessStateReporter", func() {

	const (
		timestamp = "2017-08-25_17-34-27"
		cellName  = "cell-name"
		cellId    = "cell-id"
	)

	var (
		processStateReporter statedetector.ProcessStateReporter
		pidList              = []int{100, 101}
		processesList        = []string{"foo", "bar"}
		dataPath             string
	)

	fakeEnv := map[string]string{
		"CELL_NAME": cellName,
		"CELL_ID":   cellId,
	}

	timeFunc := func() time.Time {
		t, err := time.Parse("2006-01-02_15-04-05", timestamp)
		Expect(err).NotTo(HaveOccurred())
		return t
	}

	envFunc := func(key string) string {
		return fakeEnv[key]
	}

	BeforeEach(func() {
		var err error
		dataPath, err = ioutil.TempDir("", "process_state_reporter_test")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		processStateReporter = statedetector.NewProcessStateReporter(dataPath, timeFunc, envFunc)
	})

	Describe("Run", func() {
		It("creates the debug data directory", func() {
			Expect(processStateReporter.Run(pidList, processesList)).To(Succeed())
			Expect(fmt.Sprintf("%s/%s-%s-debug-info-%s", dataPath, cellName, cellId, timestamp)).To(BeADirectory())
		})

		Context("when the data directory can not be created", func() {
			BeforeEach(func() {
				dataPath = "/unwritable/directory"
			})

			It("fails with an error", func() {
				Expect(processStateReporter.Run(pidList, processesList)).NotTo(Succeed())
			})
		})
	})
})
