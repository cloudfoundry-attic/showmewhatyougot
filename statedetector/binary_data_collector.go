package statedetector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type binaryDataCollector struct {
	badappsBinPath       string
	dataCollectorBinPath string
	commandRunner        CommandRunner
	dataPath             string
	instanceInfoPath     string
	time                 func() time.Time
	instanceID           string
	instanceName         string
	dataDirectoryPath    string
}

func NewBinaryDataCollector(
	commandRunner CommandRunner,
	badappsBinPath string,
	dataCollectorBinPath string,
	dataPath string,
	instanceInfoPath string,
	timeFunc func() time.Time,
) (DataCollector, error) {
	b := &binaryDataCollector{
		badappsBinPath:       badappsBinPath,
		dataCollectorBinPath: dataCollectorBinPath,
		commandRunner:        commandRunner,
		dataPath:             dataPath,
		instanceInfoPath:     instanceInfoPath,
		time:                 timeFunc,
	}

	return b, b.getInstanceInformation()
}

func (b *binaryDataCollector) Run(processesList []string) (string, error) {
	if err := b.createDataDirectory(); err != nil {
		return "", fmt.Errorf("Creating data directory: %s", err.Error())
	}
	defer b.deleteDataDirectory()

	var allErrors bytes.Buffer
	if err := b.writeProcessesFile(processesList); err != nil {
		allErrors.WriteString(err.Error() + "\n")
	}

	applicationsJson, err := b.runCommand(b.badappsBinPath)
	if err != nil {
		allErrors.WriteString(err.Error() + "\n")
	}

	if err := b.writeApplicationsFile(applicationsJson); err != nil {
		allErrors.WriteString(err.Error() + "\n")
	}

	collectedDataPath, err := b.runCommand(b.dataCollectorBinPath, b.dataDirectoryPath)
	if err != nil {
		allErrors.WriteString(err.Error() + "\n")
	}

	if allErrors.Len() != 0 {
		err = errors.New(allErrors.String())
	}

	return collectedDataPath, err
}

func (b *binaryDataCollector) runCommand(binaryPath string, args ...string) (string, error) {
	cmdStdout := new(bytes.Buffer)

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = io.MultiWriter(os.Stdout, cmdStdout)
	cmd.Stderr = os.Stderr

	err := b.commandRunner.Run(cmd)
	if err != nil {
		return "", fmt.Errorf("Running command failed: %s", err.Error())
	}

	return strings.TrimSpace(cmdStdout.String()), nil
}

func (b *binaryDataCollector) createDataDirectory() error {
	return os.MkdirAll(b.generateDataDirectoryPath(), 0700)
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

func (b *binaryDataCollector) writeApplicationsFile(jsonString string) error {
	var jsonBytes bytes.Buffer
	err := json.Indent(&jsonBytes, []byte(jsonString), "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(b.dataDirectoryPath, "applications.json"), jsonBytes.Bytes(), 0600)
}

func (b *binaryDataCollector) writeProcessesFile(processesList []string) error {
	file, err := os.OpenFile(filepath.Join(b.dataDirectoryPath, "dstate_processes"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, process := range processesList {
		_, _ = fmt.Fprintln(file, process)
	}

	return nil
}
