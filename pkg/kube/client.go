package kube

import (
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	GetClientset() kubernetes.Interface
	NewBuilder() *resource.Builder

	GetPodParent(string, string) (metav1.Object, error)
	GetNamespacedPods(string, string) (map[string][]corev1.Pod, error)
	GetNamespacedPDBs(string) (map[string][]policyv1.PodDisruptionBudget, error)
}

type Clients struct {
	clientset kubernetes.Interface
	conf      genericclioptions.RESTClientGetter
}

func New(conf genericclioptions.RESTClientGetter) (Interface, error) {
	restConfig, err := conf.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return &Clients{clientset: clientset, conf: conf}, nil
}

func (c *Clients) GetClientset() kubernetes.Interface {
	return c.clientset
}

func (c *Clients) NewBuilder() *resource.Builder {
	return resource.NewBuilder(c.conf)
}
