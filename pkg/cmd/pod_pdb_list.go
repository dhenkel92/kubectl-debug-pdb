package cmd

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type PodPDBEntry struct {
	Pod  corev1.Pod
	Pdbs []policyv1.PodDisruptionBudget
}

func (e *PodPDBEntry) getPDBList() string {
	names := make([]string, 0, len(e.Pdbs))
	for _, pdb := range e.Pdbs {
		names = append(names, pdb.GetName())
	}
	return strings.Join(names, ", ")
}

func (e *PodPDBEntry) OwnerName() string {
	for _, owner := range e.Pod.GetOwnerReferences() {
		return owner.Kind
	}
	return ""
}

type PodPDBList struct {
	Items []PodPDBEntry
}

func (lst *PodPDBList) toMetaTable() runtime.Object {
	rows := make([]metav1.TableRow, 0, len(lst.Items))
	for _, entry := range lst.Items {
		rows = append(rows, metav1.TableRow{
			Cells: []interface{}{
				entry.OwnerName(),
				entry.Pod.GetName(),
				len(entry.Pdbs),
				entry.getPDBList(),
			},
			Object: runtime.RawExtension{
				Object: entry.Pod.DeepCopy(),
			},
		})
	}

	return &metav1.Table{
		TypeMeta: metav1.TypeMeta{Kind: "pod", APIVersion: "v1"},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "owner", Type: "string", Format: "name", Description: "owner"},
			{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
			{Name: "CNT", Type: "integer", Description: "Amount of pdbs"},
			{Name: "PDBs", Type: "string", Description: "hello"},
		},
		Rows: rows,
	}
}
