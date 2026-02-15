package game

import (
	"fmt"
	"strings"
)

func FormatFields(fields []string, colon bool) string {
	var str strings.Builder
	for i, field := range fields {
		if (colon) {
			fmt.Fprint(&str, ":")
		}
		if (i < len(fields)-1) {
			fmt.Fprintf(&str, "%s, ", field)
		} else {
			fmt.Fprint(&str, field)
		}
	}
	return str.String()
}