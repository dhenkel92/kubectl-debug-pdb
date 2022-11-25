package utils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/rand"
)

func UniqueName(withBase string) string {
	return fmt.Sprintf("%s-%s", withBase, rand.String(5))
}
