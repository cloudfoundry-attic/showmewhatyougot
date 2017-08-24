package statedetector

func NewPersistentStateDetector(persistentStateCountThreshold int) *persistentStateDetector {
	return &persistentStateDetector{
		persistentStateCountThreshold: persistentStateCountThreshold,
		persistingProcesses:           map[int]int{},
	}
}

type persistentStateDetector struct {
	persistentStateCountThreshold int
	persistingProcesses           map[int]int
}

func (p *persistentStateDetector) Pids(currentPids []int) ([]int, error) {
	persistentPids := []int{}
	newPersistingProcesses := map[int]int{}

	for _, pid := range currentPids {
		newPersistingProcesses[pid] = p.persistingProcesses[pid] + 1

		if newPersistingProcesses[pid] >= p.persistentStateCountThreshold {
			persistentPids = append(persistentPids, pid)
		}
	}

	p.persistingProcesses = newPersistingProcesses
	return persistentPids, nil
}

func (p *persistentStateDetector) DetectedProcesses() ([]int, []string, error) {
	return []int{}, []string{}, nil
}
