// Code generated by counterfeiter. DO NOT EDIT.
package cubefakes

import (
	"sync"

	"github.com/julz/cube"
)

type FakeWorkspaceManager struct {
	CreateStub        func(string) (cube.Workspace, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 string
	}
	createReturns struct {
		result1 cube.Workspace
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 cube.Workspace
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWorkspaceManager) Create(arg1 string) (cube.Workspace, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Create", []interface{}{arg1})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createReturns.result1, fake.createReturns.result2
}

func (fake *FakeWorkspaceManager) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeWorkspaceManager) CreateArgsForCall(i int) string {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].arg1
}

func (fake *FakeWorkspaceManager) CreateReturns(result1 cube.Workspace, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 cube.Workspace
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkspaceManager) CreateReturnsOnCall(i int, result1 cube.Workspace, result2 error) {
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 cube.Workspace
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 cube.Workspace
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkspaceManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeWorkspaceManager) recordInvocation(key string, args []interface{}) {
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

var _ cube.WorkspaceManager = new(FakeWorkspaceManager)
