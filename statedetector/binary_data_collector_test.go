package statedetector_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
		dataCollector, err = statedetector.NewBinaryDataCollector(commandRunner, "/badapps", "/data-collector", dataPath, instanceInfoPath, timeFunc)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("NewBinaryDataCollector", func() {
		It("it fails when the instance name cannot be read", func() {
			Expect(os.RemoveAll(path.Join(instanceInfoPath, "name"))).To(Succeed())
			_, err := statedetector.NewBinaryDataCollector(commandRunner, "/badapps", "/data-collector", dataPath, instanceInfoPath, timeFunc)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("Unable to read instance name file")))
		})

		It("it fails when the instance id cannot be read", func() {
			Expect(os.RemoveAll(path.Join(instanceInfoPath, "id"))).To(Succeed())
			_, err := statedetector.NewBinaryDataCollector(commandRunner, "/badapps", "/data-collector", dataPath, instanceInfoPath, timeFunc)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("Unable to read instance id file")))
		})
	})

	Describe("Run", func() {
		Describe("data directory", func() {
			It("creates and destroys the data directory around the execution of the binary", func() {
				commandHasRun := false
				commandRunner.RunStub = func(cmd *exec.Cmd) error {
					if cmd.Args[0] == "/badapps" {
						cmd.Stdout.Write([]byte("{}"))
						return nil
					}

					Expect(dataDir).To(BeADirectory())
					commandHasRun = true
					return nil
				}

				_, err := dataCollector.Run([]string{"foo", "bar"})
				Expect(err).NotTo(HaveOccurred())
				Expect(commandHasRun).To(BeTrue())

				Expect(dataDir).ToNot(BeAnExistingFile())
			})

			Context("when the data directory can not be created", func() {
				BeforeEach(func() {
					Expect(os.RemoveAll(dataPath)).To(Succeed())
					_, err := os.Create(dataPath)
					Expect(err).NotTo(HaveOccurred())
				})

				It("does not run any commands and returns an error", func() {
					_, err := dataCollector.Run([]string{"foo", "bar"})
					Expect(err).To(MatchError(ContainSubstring("Creating data directory")))
					Expect(commandRunner.RunCallCount()).To(Equal(0))
				})
			})
		})

		Describe("processes list", func() {
			It("writes the processes list to a file in the data path", func() {
				commandRunner.RunStub = func(cmd *exec.Cmd) error {
					if cmd.Args[0] == "/badapps" {
						cmd.Stdout.Write([]byte("{}"))
						return nil
					}

					stuff, err := ioutil.ReadFile(path.Join(dataDir, "dstate_processes"))
					Expect(err).ToNot(HaveOccurred())
					Expect(string(stuff)).To(Equal("foo\nbar\n"))
					return nil
				}

				_, err := dataCollector.Run([]string{"foo", "bar"})
				Expect(err).NotTo(HaveOccurred())

				Expect(commandRunner.RunCallCount()).To(Equal(2), "incorrect number of commands run")
			})

			Context("when the processes list file can't be written", func() {
				BeforeEach(func() {
					Expect(os.Mkdir(dataDir, 0700)).To(Succeed())
					Expect(ioutil.WriteFile(filepath.Join(dataDir, "dstate_processes"), []byte{}, 0600)).To(Succeed())
				})

				It("continues to run all commands, then returns an error", func() {
					_, err := dataCollector.Run([]string{"foo", "bar"})
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("dstate_processes: file exists")))
					Expect(commandRunner.RunCallCount()).To(Equal(2), "incorrect number of commands run")
				})
			})
		})

		Describe("application collection binary", func() {
			It("executes the application collection binary with no additional arguments", func() {
				commandRunner.RunStub = func(cmd *exec.Cmd) error {
					if cmd.Args[0] == "/badapps" {
						cmd.Stdout.Write([]byte("{}"))
					}
					return nil
				}

				_, err := dataCollector.Run([]string{"foo", "bar"})
				Expect(err).NotTo(HaveOccurred())

				Expect(commandRunner.RunCallCount()).To(Equal(2))
				cmd := commandRunner.RunArgsForCall(0)
				Expect(cmd.Args).To(Equal([]string{"/badapps"}))
			})

			Context("when the application collection binary fails", func() {
				BeforeEach(func() {
					commandRunner.RunReturnsOnCall(0, errors.New("application collection failed"))
				})

				It("continues running further commands but returns an error", func() {
					_, err := dataCollector.Run([]string{"foo", "bar"})
					Expect(commandRunner.RunCallCount()).To(Equal(2), "incorrect number of commands invoked")
					Expect(err).To(MatchError(ContainSubstring("application collection failed")))
				})
			})
		})

		Describe("application json", func() {
			It("writes list of running applications to a file in the data path", func() {
				applicationsJsonChecked := false

				commandRunner.RunStub = func(cmd *exec.Cmd) error {
					switch cmd.Args[0] {
					case "/badapps":
						cmd.Stdout.Write([]byte(`{"badapps":"output"}`))

					case "/data-collector":
						jsonBytes, err := ioutil.ReadFile(path.Join(dataDir, "applications.json"))
						Expect(err).ToNot(HaveOccurred())
						Expect(string(jsonBytes)).To(Equal(`{
  "badapps": "output"
}`))
						applicationsJsonChecked = true
					}

					return nil
				}

				_, err := dataCollector.Run([]string{"foo", "bar"})
				Expect(err).NotTo(HaveOccurred())
				Expect(commandRunner.RunCallCount()).To(Equal(2), "incorrect number of commands run")
				Expect(applicationsJsonChecked).To(BeTrue())
			})
		})

		Describe("data collection script", func() {
			It("executes the data collection script with the correct arguments", func() {
				commandRunner.RunStub = func(cmd *exec.Cmd) error {
					if cmd.Args[0] == "/badapps" {
						cmd.Stdout.Write([]byte("{}"))
					}
					return nil
				}
				_, err := dataCollector.Run([]string{"foo", "bar"})
				Expect(err).NotTo(HaveOccurred())

				Expect(commandRunner.RunCallCount()).To(Equal(2))
				cmd := commandRunner.RunArgsForCall(1)
				Expect(cmd.Args).To(Equal([]string{"/data-collector", dataDir}))
			})

			It("returns the path to the collected data returned by the data collection script", func() {
				commandRunner.RunStub = func(cmd *exec.Cmd) error {
					if cmd.Args[0] == "/badapps" {
						cmd.Stdout.Write([]byte("{}"))
						return nil
					}

					cmd.Stdout.Write([]byte("/path/to/data.tar.gz"))
					return nil
				}

				path, err := dataCollector.Run([]string{"foo", "bar"})
				Expect(err).NotTo(HaveOccurred())
				Expect(path).To(Equal("/path/to/data.tar.gz"))
			})

			Context("when the data collection script fails", func() {
				BeforeEach(func() {
					commandRunner.RunReturnsOnCall(1, errors.New("data collection failed"))
				})

				It("returns an error", func() {
					_, err := dataCollector.Run([]string{"foo", "bar"})
					Expect(commandRunner.RunCallCount()).To(Equal(2), "incorrect number of commands invoked")
					Expect(err).To(MatchError(ContainSubstring("data collection failed")))
				})
			})
		})

		Context("when multiple things fail", func() {
			BeforeEach(func() {
				Expect(os.Mkdir(dataDir, 0700)).To(Succeed())
				Expect(ioutil.WriteFile(filepath.Join(dataDir, "dstate_processes"), []byte{}, 0600)).To(Succeed())
				commandRunner.RunReturnsOnCall(0, errors.New("application collection failed"))
				commandRunner.RunReturnsOnCall(1, errors.New("data collection failed"))
			})

			It("returns an error containing information about all the failures", func() {
				_, err := dataCollector.Run([]string{"foo", "bar"})
				Expect(commandRunner.RunCallCount()).To(Equal(2), "incorrect number of commands invoked")
				Expect(err).To(MatchError(ContainSubstring("dstate_processes: file exists")))
				Expect(err).To(MatchError(ContainSubstring("application collection failed")))
				Expect(err).To(MatchError(ContainSubstring("data collection failed")))
			})
		})
	})
})
