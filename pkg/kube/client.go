package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

type Clients struct {
	clientset kubernetes.Interface
}

func New(conf *genericclioptions.ConfigFlags) (*Clients, error) {
	restConfig, err := conf.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return &Clients{clientset: clientset}, nil
}

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

func (c *Clients) GetNamespacedPDBs(ns string) (map[string][]policyv1.PodDisruptionBudget, error) {
	pdbLst, err := c.clientset.PolicyV1().PodDisruptionBudgets(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return map[string][]policyv1.PodDisruptionBudget{}, err
	}

	pdbRes := make(map[string][]policyv1.PodDisruptionBudget)
	for _, pdb := range pdbLst.Items {
		ns := pdb.GetNamespace()
		if _, ok := pdbRes[ns]; !ok {
			pdbRes[ns] = make([]policyv1.PodDisruptionBudget, 0, 1)
		}
		pdbRes[ns] = append(pdbRes[ns], pdb)
	}

	return pdbRes, nil
}
