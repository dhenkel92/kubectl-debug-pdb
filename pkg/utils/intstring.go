package utils

import (
	"strconv"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func StrToIntOrString(str string) *intstr.IntOrString {
	if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		intStr := intstr.FromInt(int(i))
		return &intStr
	}

	intStr := intstr.FromString(str)
	return &intStr
}
