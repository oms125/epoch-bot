package items

import (
	"fmt"
	"strings"
	. "github.com/oms125/epoch-bot/game"
)

var ITEM_DATA map[int]ItemData

const (
	ROCK = iota
	STICK
	SWORD
)

func init() {
	ITEM_DATA = map[int]ItemData{
		//Materials
		ROCK:  &MaterialData{ID: ROCK, Name: "Rock"},
		STICK: &MaterialData{ID: STICK, Name: "Stick"},

		//Gears
		SWORD: &GearData{
			ID: SWORD, Name: "Sword", MaxDurability: 100, EquipSlot: "tool", Emoji: Emoji{
				Name: "sword", ID: "1474187356482568460",
			}, BaseStats: BaseStats{
				Attack: 5,
			},
		},
	}
}

//Items
type (
	Item interface {
		Type() string
	}
	ItemData interface {
		Type() string
	}

	Field struct {
		Name string
		Type string
	}

	ItemMaxError struct {
		Name     string
		Overflow int
	}
	ItemIDError struct {
		ID int
	}
	InvFullError struct{}
	ArmFullError struct{}
	ItemEmptyError struct{}
)



func (err *ItemIDError) Error() string {
	return fmt.Sprintf("Invalid item ID: %d", err.ID)
}
func (err *InvFullError) Error() string {
	return "Inventory Full"
}
func (err *ArmFullError) Error() string {
	return "Armory Full"
}
func (err *ItemMaxError) Error() string {
	return fmt.Sprintf("Max Stack Size Reached for %s: %d", err.Name, err.Overflow)
}
func (err *ItemEmptyError) Error() string {
	return "Item quantity reduced to 0"
}

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
