package utils

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func LabelsFromWorkload(workload *unstructured.Unstructured) (*v1.LabelSelector, error) {
	switch workload.GetKind() {
	case "Deployment":
		var deploy appsv1.Deployment
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(workload.Object, &deploy); err != nil {
			return nil, err
		}
		return v1.SetAsLabelSelector(labels.Set(deploy.Spec.Selector.MatchLabels)), nil
	case "ReplicaSet":
		var rs appsv1.ReplicaSet
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(workload.Object, &rs); err != nil {
			return nil, err
		}
		return v1.SetAsLabelSelector(labels.Set(rs.Spec.Selector.MatchLabels)), nil
	case "StatefulSet":
		var sts appsv1.StatefulSet
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(workload.Object, &sts); err != nil {
			return nil, err
		}
		return v1.SetAsLabelSelector(labels.Set(sts.Spec.Selector.MatchLabels)), nil
	case "Pod":
		var pod corev1.Pod
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(workload.Object, &pod); err != nil {
			return nil, err
		}
		return v1.SetAsLabelSelector(labels.Set(pod.GetLabels())), nil
	}

	return nil, fmt.Errorf("workload of type '%s' is not supported", workload.GetKind())
}
