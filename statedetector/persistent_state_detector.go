package statedetector

func NewPersistentStateDetector(persistentStateCountThreshold int, currentStateDetector StateDetector) *persistentStateDetector {
	return &persistentStateDetector{
		persistentStateCountThreshold: persistentStateCountThreshold,
		currentStateDetector:          currentStateDetector,
		persistingProcesses:           map[int]int{},
	}
}

type persistentStateDetector struct {
	currentStateDetector          StateDetector
	persistentStateCountThreshold int
	persistingProcesses           map[int]int
}

func (p *persistentStateDetector) Pids() ([]int, error) {
	persistentPids := []int{}
	newPersistingProcesses := map[int]int{}

	pids, err := p.currentStateDetector.Pids()
	if err != nil {
		return nil, err
	}

	for _, pid := range pids {
		newPersistingProcesses[pid] = p.persistingProcesses[pid] + 1

		if newPersistingProcesses[pid] >= p.persistentStateCountThreshold {
			persistentPids = append(persistentPids, pid)
		}
	}

	p.persistingProcesses = newPersistingProcesses
	return persistentPids, nil
}

func (p *persistentStateDetector) Processes() ([]string, error) {
	return []string{}, nil
}
