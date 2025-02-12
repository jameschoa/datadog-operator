// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package fake

import (
	"testing"

	v1 "k8s.io/api/core/v1"

	common "github.com/DataDog/datadog-operator/api/datadoghq/common"
	merger "github.com/DataDog/datadog-operator/internal/controller/datadogagent/merger"
)

// PortManager is an autogenerated mock type for the PortManager type
type PortManager struct {
	PortsByC map[common.AgentContainerName][]*v1.ContainerPort

	t testing.TB
}

// AddPortToContainer provides a mock function with given fields: containerName, newPort
func (_m *PortManager) AddPortToContainer(containerName common.AgentContainerName, newPort *v1.ContainerPort) {
	_m.t.Logf("AddPortToContainer %s: %#v", newPort.Name, newPort.ContainerPort)
	_m.PortsByC[containerName] = append(_m.PortsByC[containerName], newPort)
}

// AddPortToContainerWithMergeFunc provides a mock function with given fields: containerName, newPort, mergeFunc
func (_m *PortManager) AddPortToContainerWithMergeFunc(containerName common.AgentContainerName, newPort *v1.ContainerPort, mergeFunc merger.PortMergeFunction) error {
	found := false
	idFound := 0
	for id, Port := range _m.PortsByC[containerName] {
		if Port.Name == newPort.Name {
			found = true
			idFound = id
		}
	}

	if found {
		var err error
		newPort, err = mergeFunc(_m.PortsByC[containerName][idFound], newPort)
		_m.PortsByC[containerName][idFound] = newPort
		return err
	}

	_m.PortsByC[containerName] = append(_m.PortsByC[containerName], newPort)
	return nil
}

// NewFakePortManager creates a new instance of PortManager. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFakePortManager(t testing.TB) *PortManager {
	return &PortManager{
		PortsByC: make(map[common.AgentContainerName][]*v1.ContainerPort),
		t:        t,
	}
}
