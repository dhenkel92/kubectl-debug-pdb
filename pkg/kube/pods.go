package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Clients) GetNamespacedPods(ns, podName string) (map[string][]corev1.Pod, error) {
	lstOpts := metav1.ListOptions{}
	if podName != "" {
		lstOpts.FieldSelector = fmt.Sprintf("metadata.name=%s", podName)
	}

	podLst, err := c.clientset.CoreV1().Pods(ns).List(context.Background(), lstOpts)
	if err != nil {
		return map[string][]corev1.Pod{}, err
	}

	podRes := make(map[string][]corev1.Pod)
	for _, pod := range podLst.Items {
		ns := pod.GetNamespace()
		if _, ok := podRes[ns]; !ok {
			podRes[ns] = make([]corev1.Pod, 0, 1)
		}
		podRes[ns] = append(podRes[ns], pod)
	}

	return podRes, nil
}

func (c *Clients) GetPodParent(podName string) (metav1.Object, error) {
	return nil, nil
}
