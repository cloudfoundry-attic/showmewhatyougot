package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
)

func main() {
	var (
		state                          string
		pollingInterval                time.Duration
		alertIntervalThreshold         int
		tracingEnabled                 bool
		processStateCounterBinaryPath  string
		processStateReporterBinaryPath string
		xfsTraceBinaryPath             string
		pidFilePath                    string
	)

	flag.StringVar(&state, "state", "D", "Type of state to detect")
	flag.DurationVar(&pollingInterval, "polling-interval", 10*time.Second, "Interval between process state checks")
	flag.IntVar(&alertIntervalThreshold, "alert-interval-threshold", 15, "Number of checks before a process is considered in a persistent state")
	flag.BoolVar(&tracingEnabled, "tracing-enabled", false, "Enable XFS Kernel tracing")
	flag.StringVar(&processStateCounterBinaryPath, "process-state-counter", "", "State process counter binary path")
	flag.StringVar(&processStateReporterBinaryPath, "process-state-reporter", "", "State process reporter binary path")
	flag.StringVar(&xfsTraceBinaryPath, "xfs-trace-path", "", "XFS Trace binary path")
	flag.StringVar(&pidFilePath, "pid-file-path", "", "Path to write out this process's pid file")

	flag.Parse()

	if pidFilePath != "" {
		if err := ioutil.WriteFile(pidFilePath, []byte(strconv.Itoa(os.Getpid())), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write pid to '%s': %s\n", pidFilePath, err.Error())
			os.Exit(1)
		}
	}

	processStateCounter := statedetector.NewBinaryProcessStateCounter(processStateCounterBinaryPath)
	processStateReporter := statedetector.NewBinaryProcessStateReporter(processStateReporterBinaryPath)

	xfsTracer := statedetector.NewDummyXfsTracer()
	if tracingEnabled {
		xfsTracer = statedetector.NewBinaryXfsTracer(xfsTraceBinaryPath)
	}

	currentStateDetector := statedetector.NewCurrentStateDetector(state)
	persistentStateDetector := statedetector.NewPersistentStateDetector(alertIntervalThreshold, currentStateDetector)

	showMeWhatYouGot := statedetector.NewShowMeWhatYouGot(processStateCounter, processStateReporter, xfsTracer, persistentStateDetector)

	err := xfsTracer.Start()
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
