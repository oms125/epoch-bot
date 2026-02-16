package player

import (
	"fmt"
	"strings"
)

type Player struct {
	ID string
	Lvl int
	InvSize int
	ArmSize int

	Inv Inventory
	Arm Armory
	Equipment
}

//Game Logic


//Helpers
func FormatFields(fields []Field, colon bool) string {
	var str strings.Builder
	for i, field := range fields {
		if (colon) {
			fmt.Fprint(&str, ":")
		}
		if (i < len(fields)-1) {
			fmt.Fprintf(&str, "%s, ", field.Name)
		} else {
			fmt.Fprint(&str, field.Name)
		}
	}
	return str.String()
}

func FormatTableFields(fields []Field) string {
	var str strings.Builder
	for i, field := range fields {
		if (i < len(fields)-1) {
			fmt.Fprintf(&str, "%s %s, ", field.Name, field.Type)
		} else {
			fmt.Fprintf(&str, "%s %s", field.Name, field.Type)
		}
	}
	return str.String()
}