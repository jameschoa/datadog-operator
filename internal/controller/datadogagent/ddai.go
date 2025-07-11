// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package datadogagent

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	apicommon "github.com/DataDog/datadog-operator/api/datadoghq/common"
	datadoghqv1alpha1 "github.com/DataDog/datadog-operator/api/datadoghq/v1alpha1"
	datadoghqv2alpha1 "github.com/DataDog/datadog-operator/api/datadoghq/v2alpha1"
	"github.com/DataDog/datadog-operator/internal/controller/datadogagent/global"
	"github.com/DataDog/datadog-operator/internal/controller/datadogagent/object"
	"github.com/DataDog/datadog-operator/internal/controller/datadogagent/override"
	"github.com/DataDog/datadog-operator/pkg/constants"
	"github.com/DataDog/datadog-operator/pkg/controller/utils/comparison"
)

func (r *Reconciler) generateDDAIFromDDA(dda *datadoghqv2alpha1.DatadogAgent) (*datadoghqv1alpha1.DatadogAgentInternal, error) {
	ddai := &datadoghqv1alpha1.DatadogAgentInternal{}
	// Object meta
	if err := generateObjMetaFromDDA(dda, ddai, r.scheme); err != nil {
		return nil, err
	}
	// Spec
	if err := generateSpecFromDDA(dda, ddai); err != nil {
		return nil, err
	}

	// Set hash
	if _, err := comparison.SetMD5GenerationAnnotation(&ddai.ObjectMeta, ddai.Spec, constants.MD5DDAIDeploymentAnnotationKey); err != nil {
		return nil, err
	}

	return ddai, nil
}

func generateObjMetaFromDDA(dda *datadoghqv2alpha1.DatadogAgent, ddai *datadoghqv1alpha1.DatadogAgentInternal, scheme *runtime.Scheme) error {
	ddai.ObjectMeta = metav1.ObjectMeta{
		Name:        dda.Name,
		Namespace:   dda.Namespace,
		Labels:      getDDAILabels(dda),
		Annotations: dda.Annotations,
	}
	if err := object.SetOwnerReference(dda, ddai, scheme); err != nil {
		return err
	}
	return nil
}

func generateSpecFromDDA(dda *datadoghqv2alpha1.DatadogAgent, ddai *datadoghqv1alpha1.DatadogAgentInternal) error {
	ddai.Spec = *dda.Spec.DeepCopy()
	global.SetGlobalFromDDA(dda, ddai.Spec.Global)
	override.SetOverrideFromDDA(dda, &ddai.Spec)
	return nil
}

// getDDAILabels adds the following labels to the DDAI:
// - all DDA labels
// - agent.datadoghq.com/datadogagent: <dda-name>
func getDDAILabels(dda metav1.Object) map[string]string {
	labels := make(map[string]string)
	for k, v := range dda.GetLabels() {
		labels[k] = v
	}
	labels[apicommon.DatadogAgentNameLabelKey] = dda.GetName()
	return labels
}
