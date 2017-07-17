// This file was generated by counterfeiter
package statedetectorfakes

import (
	"sync"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
)

type FakeStateDetector struct {
	PidsStub        func() ([]int, error)
	pidsMutex       sync.RWMutex
	pidsArgsForCall []struct{}
	pidsReturns     struct {
		result1 []int
		result2 error
	}
	pidsReturnsOnCall map[int]struct {
		result1 []int
		result2 error
	}
	ProcessesStub        func() ([]string, error)
	processesMutex       sync.RWMutex
	processesArgsForCall []struct{}
	processesReturns     struct {
		result1 []string
		result2 error
	}
	processesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStateDetector) Pids() ([]int, error) {
	fake.pidsMutex.Lock()
	ret, specificReturn := fake.pidsReturnsOnCall[len(fake.pidsArgsForCall)]
	fake.pidsArgsForCall = append(fake.pidsArgsForCall, struct{}{})
	fake.recordInvocation("Pids", []interface{}{})
	fake.pidsMutex.Unlock()
	if fake.PidsStub != nil {
		return fake.PidsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.pidsReturns.result1, fake.pidsReturns.result2
}

func (fake *FakeStateDetector) PidsCallCount() int {
	fake.pidsMutex.RLock()
	defer fake.pidsMutex.RUnlock()
	return len(fake.pidsArgsForCall)
}

func (fake *FakeStateDetector) PidsReturns(result1 []int, result2 error) {
	fake.PidsStub = nil
	fake.pidsReturns = struct {
		result1 []int
		result2 error
	}{result1, result2}
}

func (fake *FakeStateDetector) PidsReturnsOnCall(i int, result1 []int, result2 error) {
	fake.PidsStub = nil
	if fake.pidsReturnsOnCall == nil {
		fake.pidsReturnsOnCall = make(map[int]struct {
			result1 []int
			result2 error
		})
	}
	fake.pidsReturnsOnCall[i] = struct {
		result1 []int
		result2 error
	}{result1, result2}
}

func (fake *FakeStateDetector) Processes() ([]string, error) {
	fake.processesMutex.Lock()
	ret, specificReturn := fake.processesReturnsOnCall[len(fake.processesArgsForCall)]
	fake.processesArgsForCall = append(fake.processesArgsForCall, struct{}{})
	fake.recordInvocation("Processes", []interface{}{})
	fake.processesMutex.Unlock()
	if fake.ProcessesStub != nil {
		return fake.ProcessesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.processesReturns.result1, fake.processesReturns.result2
}

func (fake *FakeStateDetector) ProcessesCallCount() int {
	fake.processesMutex.RLock()
	defer fake.processesMutex.RUnlock()
	return len(fake.processesArgsForCall)
}

func (fake *FakeStateDetector) ProcessesReturns(result1 []string, result2 error) {
	fake.ProcessesStub = nil
	fake.processesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeStateDetector) ProcessesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.ProcessesStub = nil
	if fake.processesReturnsOnCall == nil {
		fake.processesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.processesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeStateDetector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pidsMutex.RLock()
	defer fake.pidsMutex.RUnlock()
	fake.processesMutex.RLock()
	defer fake.processesMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeStateDetector) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ statedetector.StateDetector = new(FakeStateDetector)
