package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"time"

	"code.cloudfoundry.org/showmewhatyougot/runner"
	"code.cloudfoundry.org/showmewhatyougot/statedetector"
)

func main() {
	var (
		state                        string
		pollingInterval              time.Duration
		reporterBackoffDuration      time.Duration
		commandsTimeout              time.Duration
		alertIntervalThreshold       int
		tracingEnabled               bool
		stateCountReporterBinaryPath string
		dataCollectorBinaryPath      string
		appCollectorBinaryPath       string
		xfsTraceBinaryPath           string
		eventEmitterBinaryPath       string
		pidFilePath                  string
		dataPath                     string
		instanceInfoPath             string
	)

	flag.StringVar(&state, "state", "D", "Type of state to detect")
	flag.DurationVar(&pollingInterval, "polling-interval", 10*time.Second, "Interval between process state checks")
	flag.DurationVar(&reporterBackoffDuration, "reporter-backoff-duration", 10*time.Minute, "Reporting is restricted to one report per backoff duration")
	flag.IntVar(&alertIntervalThreshold, "alert-interval-threshold", 15, "Number of checks before a process is considered in a persistent state")
	flag.BoolVar(&tracingEnabled, "tracing-enabled", false, "Enable XFS Kernel tracing")
	flag.StringVar(&stateCountReporterBinaryPath, "state-count-reporter-path", "", "State process count reporter binary path")
	flag.StringVar(&dataCollectorBinaryPath, "data-collector-path", "", "Data collector binary path")
	flag.StringVar(&appCollectorBinaryPath, "app-collector-path", "", "Application collector binary path")
	flag.StringVar(&xfsTraceBinaryPath, "xfs-trace-path", "", "XFS Trace binary path")
	flag.StringVar(&eventEmitterBinaryPath, "event-emitter-path", "", "Event emitter binary path")
	flag.StringVar(&pidFilePath, "pid-file-path", "", "Path to write out this process's pid file")
	flag.StringVar(&dataPath, "data-path", "", "Path to write out collected data to")
	flag.StringVar(&instanceInfoPath, "instance-info-path", "", "Path to BOSH instance information")
	flag.DurationVar(&commandsTimeout, "commands-timeout", 15*time.Second, "Maximum external command duration")

	flag.Parse()

	if pidFilePath != "" {
		if err := ioutil.WriteFile(pidFilePath, []byte(strconv.Itoa(os.Getpid())), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write pid to '%s': %s\n", pidFilePath, err.Error())
			os.Exit(1)
		}
	}

	commandRunner := runner.New(commandsTimeout)
	stateCountReporter := statedetector.NewBinaryProcessStateCounter(commandRunner, stateCountReporterBinaryPath)
	dataCollector, err := statedetector.NewBinaryDataCollector(
		commandRunner,
		appCollectorBinaryPath,
		dataCollectorBinaryPath,
		dataPath,
		instanceInfoPath,
		time.Now,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to construct data collector: %s\n", err.Error())
		os.Exit(1)
	}

	xfsTracer := statedetector.NewDummyXfsTracer()
	if tracingEnabled {
		xfsTracer = statedetector.NewBinaryXfsTracer(commandRunner, xfsTraceBinaryPath)
	}

	currentStateDetector := statedetector.NewCurrentStateDetector(commandRunner, state)
	persistentStateDetector := statedetector.NewPersistentStateDetector(alertIntervalThreshold)

	eventEmitter := statedetector.NewBinaryEventEmitter(commandRunner, eventEmitterBinaryPath)

	showMeWhatYouGot := statedetector.NewShowMeWhatYouGot(
		stateCountReporter,
		dataCollector,
		xfsTracer,
		persistentStateDetector,
		currentStateDetector,
		eventEmitter,
		reporterBackoffDuration,
		os.Stderr,
	)

	err = xfsTracer.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start XFS Tracer: %s\n", err.Error())
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func(xfsTracer statedetector.XfsTracer) {
		<-sig
		_ = xfsTracer.Stop()
		os.Exit(0)
	}(xfsTracer)

	for {
		err := showMeWhatYouGot.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
		time.Sleep(pollingInterval)
	}
}
