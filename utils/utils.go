package utils

import (
	"strings"
)

func GetWhitespaces(x int) string {
	return strings.Repeat(" ", x)
}

func MapStringString2String(m map[string]string, kvSep, entrySep string) string {
	sarr := []string{}
	for k, v := range m {
		sarr = append(sarr, k+kvSep+v)
	}
	return strings.Join(sarr, entrySep)
}
