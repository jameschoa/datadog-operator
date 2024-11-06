// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package enabledefault

import (
	"testing"

	"github.com/DataDog/datadog-operator/api/datadoghq/v2alpha1"
	v2alpha1test "github.com/DataDog/datadog-operator/api/datadoghq/v2alpha1/test"
	"github.com/stretchr/testify/assert"
)

func Test_getDaemonSetNameFromDatadogAgent(t *testing.T) {
	tests := []struct {
		name              string
		ddaName           string
		overrideAgentName string
		expectedName      string
	}{
		{
			name:              "No override",
			ddaName:           "foo",
			overrideAgentName: "",
			expectedName:      "foo-agent",
		},
		{
			name:              "With override",
			ddaName:           "bar",
			overrideAgentName: "node",
			expectedName:      "node",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dda := v2alpha1test.NewDatadogAgentBuilder().
				WithName(tt.ddaName).
				WithComponentOverride(v2alpha1.NodeAgentComponentName, v2alpha1.DatadogAgentComponentOverride{
					Name: &tt.overrideAgentName,
				}).
				Build()
			dsName := getDaemonSetNameFromDatadogAgent(dda)
			assert.Equal(t, tt.expectedName, dsName)
		})
	}
}

func Test_getDeploymentNameFromDatadogAgent(t *testing.T) {
	tests := []struct {
		name                     string
		ddaName                  string
		overrideClusterAgentName string
		expectedName             string
	}{
		{
			name:                     "No override",
			ddaName:                  "foo",
			overrideClusterAgentName: "",
			expectedName:             "foo-cluster-agent",
		},
		{
			name:                     "With override",
			ddaName:                  "bar",
			overrideClusterAgentName: "dca",
			expectedName:             "dca",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dda := v2alpha1test.NewDatadogAgentBuilder().
				WithName(tt.ddaName).
				WithComponentOverride(v2alpha1.ClusterAgentComponentName, v2alpha1.DatadogAgentComponentOverride{
					Name: &tt.overrideClusterAgentName,
				}).
				Build()
			deployName := getDeploymentNameFromDatadogAgent(dda)
			assert.Equal(t, tt.expectedName, deployName)
		})
	}
}
