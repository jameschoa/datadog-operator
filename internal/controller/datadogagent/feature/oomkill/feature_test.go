// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package oomkill

import (
	"testing"

	apicommon "github.com/DataDog/datadog-operator/api/datadoghq/common"
	"github.com/DataDog/datadog-operator/api/datadoghq/v2alpha1"
	apiutils "github.com/DataDog/datadog-operator/api/utils"
	"github.com/DataDog/datadog-operator/internal/controller/datadogagent/component/agent"
	"github.com/DataDog/datadog-operator/internal/controller/datadogagent/feature"
	"github.com/DataDog/datadog-operator/internal/controller/datadogagent/feature/fake"
	"github.com/DataDog/datadog-operator/internal/controller/datadogagent/feature/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func Test_oomKillFeature_Configure(t *testing.T) {
	ddaOOMKillDisabled := v2alpha1.DatadogAgent{
		Spec: v2alpha1.DatadogAgentSpec{
			Features: &v2alpha1.DatadogFeatures{
				OOMKill: &v2alpha1.OOMKillFeatureConfig{
					Enabled: apiutils.NewBoolPointer(false),
				},
			},
		},
	}
	ddaOOMKillEnabled := ddaOOMKillDisabled.DeepCopy()
	{
		ddaOOMKillEnabled.Spec.Features.OOMKill.Enabled = apiutils.NewBoolPointer(true)
	}

	oomKillAgentNodeWantFunc := func(t testing.TB, mgrInterface feature.PodTemplateManagers) {
		mgr := mgrInterface.(*fake.PodTemplateManagers)

		// check security context capabilities
		sysProbeCapabilities := mgr.SecurityContextMgr.CapabilitiesByC[apicommon.SystemProbeContainerName]
		assert.True(
			t,
			apiutils.IsEqualStruct(sysProbeCapabilities, agent.DefaultCapabilitiesForSystemProbe()),
			"System Probe security context capabilities \ndiff = %s",
			cmp.Diff(sysProbeCapabilities, agent.DefaultCapabilitiesForSystemProbe()),
		)

		// check volume mounts
		wantCoreAgentVolMounts := []corev1.VolumeMount{
			{
				Name:      v2alpha1.SystemProbeSocketVolumeName,
				MountPath: v2alpha1.SystemProbeSocketVolumePath,
				ReadOnly:  true,
			},
		}

		wantSystemProbeVolMounts := []corev1.VolumeMount{
			{
				Name:      v2alpha1.ModulesVolumeName,
				MountPath: v2alpha1.ModulesVolumePath,
				ReadOnly:  true,
			},
			{
				Name:      v2alpha1.SrcVolumeName,
				MountPath: v2alpha1.SrcVolumePath,
				ReadOnly:  true,
			},
			{
				Name:      v2alpha1.DebugfsVolumeName,
				MountPath: v2alpha1.DebugfsPath,
				ReadOnly:  false,
			},
			{
				Name:      v2alpha1.SystemProbeSocketVolumeName,
				MountPath: v2alpha1.SystemProbeSocketVolumePath,
				ReadOnly:  false,
			},
		}

		coreAgentVolumeMounts := mgr.VolumeMountMgr.VolumeMountsByC[apicommon.CoreAgentContainerName]
		assert.True(t, apiutils.IsEqualStruct(coreAgentVolumeMounts, wantCoreAgentVolMounts), "Core agent volume mounts \ndiff = %s", cmp.Diff(coreAgentVolumeMounts, wantCoreAgentVolMounts))

		systemProbeVolumeMounts := mgr.VolumeMountMgr.VolumeMountsByC[apicommon.SystemProbeContainerName]
		assert.True(t, apiutils.IsEqualStruct(systemProbeVolumeMounts, wantSystemProbeVolMounts), "System Probe volume mounts \ndiff = %s", cmp.Diff(systemProbeVolumeMounts, wantSystemProbeVolMounts))

		// check volumes
		wantVolumes := []corev1.Volume{
			{
				Name: v2alpha1.ModulesVolumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: v2alpha1.ModulesVolumePath,
					},
				},
			},
			{
				Name: v2alpha1.SrcVolumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: v2alpha1.SrcVolumePath,
					},
				},
			},
			{
				Name: v2alpha1.DebugfsVolumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: v2alpha1.DebugfsPath,
					},
				},
			},
			{
				Name: v2alpha1.SystemProbeSocketVolumeName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}

		volumes := mgr.VolumeMgr.Volumes
		assert.True(t, apiutils.IsEqualStruct(volumes, wantVolumes), "Volumes \ndiff = %s", cmp.Diff(volumes, wantVolumes))

		// check env vars
		wantEnvVars := []*corev1.EnvVar{
			{
				Name:  DDEnableOOMKillEnvVar,
				Value: "true",
			},
			{
				Name:  v2alpha1.DDSystemProbeEnabled,
				Value: "true",
			},
			{
				Name:  v2alpha1.DDSystemProbeSocket,
				Value: v2alpha1.DefaultSystemProbeSocketPath,
			},
		}
		agentEnvVars := mgr.EnvVarMgr.EnvVarsByC[apicommon.CoreAgentContainerName]
		assert.True(t, apiutils.IsEqualStruct(agentEnvVars, wantEnvVars), "Agent envvars \ndiff = %s", cmp.Diff(agentEnvVars, wantEnvVars))

		systemProbeEnvVars := mgr.EnvVarMgr.EnvVarsByC[apicommon.SystemProbeContainerName]
		assert.True(t, apiutils.IsEqualStruct(systemProbeEnvVars, wantEnvVars), "System Probe envvars \ndiff = %s", cmp.Diff(systemProbeEnvVars, wantEnvVars))
	}

	tests := test.FeatureTestSuite{
		{
			Name:          "oom kill not enabled",
			DDA:           ddaOOMKillDisabled.DeepCopy(),
			WantConfigure: false,
		},
		{
			Name:          "oom kill enabled",
			DDA:           ddaOOMKillEnabled,
			WantConfigure: true,
			Agent:         test.NewDefaultComponentTest().WithWantFunc(oomKillAgentNodeWantFunc),
		},
	}

	tests.Run(t, buildOOMKillFeature)
}
