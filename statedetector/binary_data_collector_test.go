package statedetector_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"code.cloudfoundry.org/showmewhatyougot/statedetector"
	"code.cloudfoundry.org/showmewhatyougot/statedetector/statedetectorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BinaryDataCollector", func() {
	const (
		timestamp = "2017-08-25_17-34-27"
		cellName  = "cell-name"
		cellId    = "cell-id"
	)

	var (
		dataPath         string
		instanceInfoPath string
		dataDir          string
		dataCollector    statedetector.DataCollector
		commandRunner    *statedetectorfakes.FakeCommandRunner
	)

	timeFunc := func() time.Time {
		t, err := time.Parse("2006-01-02_15-04-05", timestamp)
		Expect(err).NotTo(HaveOccurred())
		return t
	}

	BeforeEach(func() {
		var err error
		dataPath, err = ioutil.TempDir("", "process_state_reporter_test")
		instanceInfoPath, err = ioutil.TempDir("", "instance")
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(path.Join(instanceInfoPath, "id"), []byte(cellId), 0755)).To(Succeed())
		Expect(ioutil.WriteFile(path.Join(instanceInfoPath, "name"), []byte(cellName), 0755)).To(Succeed())

		dataDir = fmt.Sprintf("%s/%s-%s-debug-info-%s", dataPath, cellName, cellId, timestamp)
		commandRunner = new(statedetectorfakes.FakeCommandRunner)
	})

	JustBeforeEach(func() {
		var err error
		dataCollector, err = statedetector.NewBinaryDataCollector(commandRunner, "/hello", dataPath, instanceInfoPath, timeFunc)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("NewBinaryDataCollector", func() {
		It("it fails when the instance name cannot be read", func() {
			Expect(os.RemoveAll(path.Join(instanceInfoPath, "name"))).To(Succeed())
			_, err := statedetector.NewBinaryDataCollector(commandRunner, "/hello", dataPath, instanceInfoPath, timeFunc)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("Unable to read instance name file")))
		})

		It("it fails when the instance id cannot be read", func() {
			Expect(os.RemoveAll(path.Join(instanceInfoPath, "id"))).To(Succeed())
			_, err := statedetector.NewBinaryDataCollector(commandRunner, "/hello", dataPath, instanceInfoPath, timeFunc)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("Unable to read instance id file")))
		})
	})

	Describe("Run", func() {
		It("creates and destroys the data directory around the execution of the binary", func() {
			commandHasRun := false
			commandRunner.RunStub = func(_ *exec.Cmd) error {
				Expect(dataDir).To(BeADirectory())
				commandHasRun = true
				return nil
			}

			Expect(dataCollector.Run([]int{100, 101}, []string{"foo", "bar"})).To(Succeed())
			Expect(commandHasRun).To(BeTrue())

			Expect(dataDir).ToNot(BeAnExistingFile())
		})

		Context("when the data directory can not be created", func() {
			BeforeEach(func() {
				dataPath = fmt.Sprintf("%s/non-existent-directory", dataPath)
			})

			It("does not run the command and returns an error", func() {
				Expect(dataCollector.Run([]int{100, 101}, []string{"foo", "bar"})).To(MatchError(ContainSubstring("Creating data directory")))
				Expect(commandRunner.RunCallCount()).To(Equal(0))
			})
		})

		It("executes the binary with the correct arguments", func() {
			Expect(dataCollector.Run([]int{100, 101}, []string{"foo", "bar"})).To(Succeed())

			Expect(commandRunner.RunCallCount()).To(Equal(1))
			cmd := commandRunner.RunArgsForCall(0)
			Expect(cmd.Args).To(Equal([]string{"/hello", "100 101", "foo\nbar", dataDir}))
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				commandRunner.RunReturns(errors.New("failed"))
			})

			It("returns an error", func() {
				err := dataCollector.Run([]int{100, 101}, []string{"foo", "bar"})
				Expect(commandRunner.RunCallCount()).To(Equal(1))
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})
	})
})
