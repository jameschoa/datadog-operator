package fake

import (
	"testing"

	v1 "k8s.io/api/core/v1"

	"github.com/DataDog/datadog-operator/api/datadoghq/common"
	merger "github.com/DataDog/datadog-operator/internal/controller/datadogagent/merger"
)

// EnvVarManager is an autogenerated mock type for the EnvVarManager type
type EnvFromVarManager struct {
	EnvFromVarsByC map[common.AgentContainerName][]*v1.EnvFromSource

	t testing.TB
}

// AddEnvVar provides a mock function with given fields: newEnvFromVar
func (_m *EnvFromVarManager) AddEnvFromVar(newEnvFromVar *v1.EnvFromSource) {
	_m.t.Logf("AddEnvFromVar %s", newEnvFromVar.SecretRef.Name)
	_m.EnvFromVarsByC[AllContainers] = append(_m.EnvFromVarsByC[AllContainers], newEnvFromVar)
}

// AddEnvFromVarToContainer provides a mock function with given fields: containerName, newEnvFromVar
func (_m *EnvFromVarManager) AddEnvFromVarToContainer(containerName common.AgentContainerName, newEnvFromVar *v1.EnvFromSource) {
	isInitContainer := false
	for _, initContainerName := range initContainerNames {
		if containerName == initContainerName {
			isInitContainer = true
			break
		}
	}
	if !isInitContainer {
		_m.t.Logf("AddEnvFromVar %s", newEnvFromVar.SecretRef.Name)
		_m.EnvFromVarsByC[containerName] = append(_m.EnvFromVarsByC[containerName], newEnvFromVar)
	}
}

// AddEnvFromVarToContainers provides a mock function with given fields: containerNames, newEnvFromVar
func (_m *EnvFromVarManager) AddEnvFromVarToContainers(containerNames []common.AgentContainerName, newEnvFromVar *v1.EnvFromSource) {
	for _, containerName := range containerNames {
		_m.AddEnvFromVarToContainer(containerName, newEnvFromVar)
	}
}

// AddEnvFromVarToInitContainer provides a mock function with given fields: containerName, newEnvFromVar
func (_m *EnvFromVarManager) AddEnvFromVarToInitContainer(containerName common.AgentContainerName, newEnvFromVar *v1.EnvFromSource) {
	for _, initContainerName := range initContainerNames {
		if containerName == initContainerName {
			_m.t.Logf("AddEnvVar to container %s name:%s", containerName, newEnvFromVar.SecretRef.Name)
			_m.EnvFromVarsByC[containerName] = append(_m.EnvFromVarsByC[containerName], newEnvFromVar)
		}
	}
}

// AddEnvFromVarToContainerWithMergeFunc provides a mock function with given fields: containerName, newEnvFromVar, mergeFunc
func (_m *EnvFromVarManager) AddEnvFromVarToContainerWithMergeFunc(containerName common.AgentContainerName, newEnvFromVar *v1.EnvFromSource, mergeFunc merger.EnvFromSourceFromMergeFunction) error {
	found := false
	idFound := 0
	for id, envVar := range _m.EnvFromVarsByC[containerName] {
		if envVar.SecretRef.Name == newEnvFromVar.SecretRef.Name {
			found = true
			idFound = id
		}
	}

	if found {
		var err error
		newEnvFromVar, err = mergeFunc(_m.EnvFromVarsByC[containerName][idFound], newEnvFromVar)
		_m.EnvFromVarsByC[containerName][idFound] = newEnvFromVar
		return err
	}

	_m.EnvFromVarsByC[containerName] = append(_m.EnvFromVarsByC[containerName], newEnvFromVar)
	return nil
}

// AddEnvVarWithMergeFunc provides a mock function with given fields: newEnvFromVar, mergeFunc
func (_m *EnvFromVarManager) AddEnvFromVarWithMergeFunc(newEnvFromVar *v1.EnvFromSource, mergeFunc merger.EnvFromSourceFromMergeFunction) error {
	return _m.AddEnvFromVarToContainerWithMergeFunc(AllContainers, newEnvFromVar, mergeFunc)
}

// NewFakeEnvVarFromManager creates a new instance of EnvFromVarManager. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFakeEnvFromVarManager(t testing.TB) *EnvFromVarManager {
	return &EnvFromVarManager{
		EnvFromVarsByC: make(map[common.AgentContainerName][]*v1.EnvFromSource),
		t:              t,
	}
}
