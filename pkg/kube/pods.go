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

func (c *Clients) GetPodParent(podName, ns string) (metav1.Object, error) {
	cs := c.GetClientset()
	pod, err := cs.CoreV1().Pods(ns).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// if there is now owner, the pod is probably standalone
	if len(pod.GetOwnerReferences()) == 0 {
		return pod, nil
	}

	// owner := pod.GetOwnerReferences()[0]

	return nil, nil
}
