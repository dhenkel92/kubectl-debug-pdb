package kube

import (
	"context"

	"github.com/dhenkel92/kubectl-utils/pkg/utils"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

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

func NewPDB(workload *unstructured.Unstructured, ls *metav1.LabelSelector) policyv1.PodDisruptionBudget {
	var pdb policyv1.PodDisruptionBudget
	falseVar := false

	pdb.ObjectMeta.Name = utils.UniqueName(workload.GetName())
	pdb.ObjectMeta.Namespace = workload.GetNamespace()
	pdb.Spec.Selector = ls
	pdb.GetObjectKind().SetGroupVersionKind(policyv1.SchemeGroupVersion.WithKind("PodDisruptionBudget"))
	pdb.OwnerReferences = []metav1.OwnerReference{
		{
			Kind:               workload.GetKind(),
			APIVersion:         workload.GetAPIVersion(),
			Name:               workload.GetName(),
			Controller:         &falseVar,
			BlockOwnerDeletion: &falseVar,
			UID:                workload.GetUID(),
		},
	}

	return pdb
}
