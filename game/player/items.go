package player

import (
	"fmt"
)

//ItemData Structures
var (
	ITEM_DATA map[int]ItemData
)

//Constant Item Type References
const (
	TYPE_MATERIAL = iota
	TYPE_TOOL
)

//Init Items
const (
	ROCK = iota
	STICK
	SWORD
)

func init() {
	ITEM_DATA = map[int]ItemData {
		//Materials
		ROCK: &MaterialData{ID: ROCK, Name: "Rock"},
		STICK: &MaterialData{ID: STICK, Name: "Stick"},

		//Tools
		SWORD: &ToolData{ID: SWORD, Name: "Sword", MaxDurability: 100, Slot: TOOL},
	}
}

//Game Logic
type (
	Item interface {
		Type() int
	}
	ItemData interface {
		Type() int
	}
	ItemMetadata interface {
		Type() int
	}

	Field struct {
		Name string
		Type string
	}
	
	ItemIDError struct {
		ID int
	}
	InvFullError struct {}
	ArmFullError struct {}
	ItemMaxError struct { 
		Name string
		Overflow int
	}
	ItemEmptyError struct {}
)

func (p *Player) AddItem(id int, quantity ...int) error {
	data := ITEM_DATA[id]
	switch data.Type() {
	case TYPE_MATERIAL:
		num := 1
		if len(quantity) > 0 {
			num = quantity[0]
		}
		return p.AddMaterial(id, num)
	case TYPE_TOOL:
		return p.AddTool(NewDefaultTool(id))
	}
	return fmt.Errorf("Invalid item ID")
}

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

