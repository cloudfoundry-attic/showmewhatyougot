package statedetector

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type binaryDataCollector struct {
	path              string
	commandRunner     CommandRunner
	dataPath          string
	instanceInfoPath  string
	time              func() time.Time
	instanceID        string
	instanceName      string
	dataDirectoryPath string
}

func NewBinaryDataCollector(
	commandRunner CommandRunner,
	binPath string,
	dataPath string,
	instanceInfoPath string,
	timeFunc func() time.Time,
) (DataCollector, error) {
	b := &binaryDataCollector{
		path:             binPath,
		commandRunner:    commandRunner,
		dataPath:         dataPath,
		instanceInfoPath: instanceInfoPath,
		time:             timeFunc,
	}

	return b, b.getInstanceInformation()
}

func (b *binaryDataCollector) Run(pidList []int, processesList []string) error {
	if err := b.createDataDirectory(); err != nil {
		return fmt.Errorf("Creating data directory: %s", err.Error())
	}
	defer b.deleteDataDirectory()

	return b.runCommand(pidList, processesList)
}

func (b *binaryDataCollector) runCommand(pidList []int, processesList []string) error {
	pidListArgs := []string{}
	for _, pid := range pidList {
		pidListArgs = append(pidListArgs, strconv.Itoa(pid))
	}

	args := []string{
		strings.Join(pidListArgs, " "),
		strings.Join(processesList, "\n"),
		b.dataDirectoryPath,
	}

	cmd := exec.Command(b.path, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := b.commandRunner.Run(cmd)
	if err != nil {
		return fmt.Errorf("Running process state reporter: %s", err.Error())
	}

	return nil
}

func (b *binaryDataCollector) createDataDirectory() error {
	return os.Mkdir(b.generateDataDirectoryPath(), 700)
}

func (b *binaryDataCollector) deleteDataDirectory() error {
	return os.RemoveAll(b.dataDirectoryPath)
}

func (b *binaryDataCollector) generateDataDirectoryPath() string {
	b.dataDirectoryPath = fmt.Sprintf(
		"%s/%s-%s-debug-info-%s",
		b.dataPath,
		b.instanceName,
		b.instanceID,
		b.time().Format("2006-01-02_15-04-05"),
	)
	return b.dataDirectoryPath
}

func (b *binaryDataCollector) getInstanceInformation() error {
	idBytes, err := ioutil.ReadFile(path.Join(b.instanceInfoPath, "id"))
	if err != nil {
		return fmt.Errorf("Unable to read instance id file: %s", err)
	}

	nameBytes, err := ioutil.ReadFile(path.Join(b.instanceInfoPath, "name"))
	if err != nil {
		return fmt.Errorf("Unable to read instance name file: %s", err)
	}

	b.instanceID = string(idBytes)
	b.instanceName = string(nameBytes)

	return nil
}
