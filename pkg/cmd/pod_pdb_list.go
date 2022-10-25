package cmd

import (
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	SortByNameKey  = "name"
	SortByCountKey = "count"
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

func (lst *PodPDBList) sortBy(key string) error {
	switch key {
	case SortByNameKey:
		sort.Slice(lst.Items, func(i, j int) bool {
			return strings.Compare(lst.Items[i].Pod.GetName(), lst.Items[j].Pod.GetName()) >= 0
		})
	case SortByCountKey:
		sort.Slice(lst.Items, func(i, j int) bool {
			return len(lst.Items[i].getPDBList()) > len(lst.Items[j].getPDBList())
		})
	default:
		return fmt.Errorf("cannot sort by '%s'", key)
	}
	return nil
}

func (lst *PodPDBList) toMetaTable() runtime.Object {
	rows := make([]metav1.TableRow, 0, len(lst.Items))
	for _, entry := range lst.Items {
		rows = append(rows, metav1.TableRow{
			Cells: []interface{}{
				entry.OwnerName(),
				entry.Pod.GetName(),
				entry.getPDBList(),
				len(entry.Pdbs),
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
			{Name: "PDBs", Type: "string", Description: "hello"},
			{Name: "CNT", Type: "integer", Description: "Amount of pdbs"},
		},
		Rows: rows,
	}
}
