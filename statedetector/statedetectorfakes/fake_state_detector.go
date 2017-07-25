// This file was generated by counterfeiter
package statedetectorfakes

import (
	"sync"

	"github.com/masters-of-cats/showmewhatyougot/statedetector"
)

type FakeStateDetector struct {
	PidsStub        func([]int) ([]int, error)
	pidsMutex       sync.RWMutex
	pidsArgsForCall []struct {
		arg1 []int
	}
	pidsReturns struct {
		result1 []int
		result2 error
	}
	pidsReturnsOnCall map[int]struct {
		result1 []int
		result2 error
	}
	RunPSStub        func() ([]int, []string, error)
	runPSMutex       sync.RWMutex
	runPSArgsForCall []struct{}
	runPSReturns     struct {
		result1 []int
		result2 []string
		result3 error
	}
	runPSReturnsOnCall map[int]struct {
		result1 []int
		result2 []string
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStateDetector) Pids(arg1 []int) ([]int, error) {
	var arg1Copy []int
	if arg1 != nil {
		arg1Copy = make([]int, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.pidsMutex.Lock()
	ret, specificReturn := fake.pidsReturnsOnCall[len(fake.pidsArgsForCall)]
	fake.pidsArgsForCall = append(fake.pidsArgsForCall, struct {
		arg1 []int
	}{arg1Copy})
	fake.recordInvocation("Pids", []interface{}{arg1Copy})
	fake.pidsMutex.Unlock()
	if fake.PidsStub != nil {
		return fake.PidsStub(arg1)
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

func (fake *FakeStateDetector) PidsArgsForCall(i int) []int {
	fake.pidsMutex.RLock()
	defer fake.pidsMutex.RUnlock()
	return fake.pidsArgsForCall[i].arg1
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

func (fake *FakeStateDetector) RunPS() ([]int, []string, error) {
	fake.runPSMutex.Lock()
	ret, specificReturn := fake.runPSReturnsOnCall[len(fake.runPSArgsForCall)]
	fake.runPSArgsForCall = append(fake.runPSArgsForCall, struct{}{})
	fake.recordInvocation("RunPS", []interface{}{})
	fake.runPSMutex.Unlock()
	if fake.RunPSStub != nil {
		return fake.RunPSStub()
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.runPSReturns.result1, fake.runPSReturns.result2, fake.runPSReturns.result3
}

func (fake *FakeStateDetector) RunPSCallCount() int {
	fake.runPSMutex.RLock()
	defer fake.runPSMutex.RUnlock()
	return len(fake.runPSArgsForCall)
}

func (fake *FakeStateDetector) RunPSReturns(result1 []int, result2 []string, result3 error) {
	fake.RunPSStub = nil
	fake.runPSReturns = struct {
		result1 []int
		result2 []string
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeStateDetector) RunPSReturnsOnCall(i int, result1 []int, result2 []string, result3 error) {
	fake.RunPSStub = nil
	if fake.runPSReturnsOnCall == nil {
		fake.runPSReturnsOnCall = make(map[int]struct {
			result1 []int
			result2 []string
			result3 error
		})
	}
	fake.runPSReturnsOnCall[i] = struct {
		result1 []int
		result2 []string
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeStateDetector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pidsMutex.RLock()
	defer fake.pidsMutex.RUnlock()
	fake.runPSMutex.RLock()
	defer fake.runPSMutex.RUnlock()
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
