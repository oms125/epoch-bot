package items

import (
	"fmt"
	"errors"
	"strings"
	. "github.com/oms125/epoch-bot/game"
)

const TYPE_MATERIAL = "material"

type (
	Inventory map[int]*Material

	Material struct {
		*MaterialData
		MaterialMetadata
	}
	MaterialData struct {
		ID int `db:"item_id"`
		Name string
		Emoji
	}
	MaterialMetadata struct {
		Quantity int `db:"quantity"`
	}
)

var InventoryFields = []Field{
	{
		Name: "item_id",
		Type: "INTEGER",
	},
	{
		Name: "quantity",
		Type: "INTEGER",
	},
}

//Materials
func (m *Material) Type() string { return TYPE_MATERIAL }
func (m *MaterialData) Type() string { return TYPE_MATERIAL }

func (m *Material) ChangeQuantity(n int) error { 
	total := m.MaterialMetadata.Quantity + n
	if (total > 99) {
		m.MaterialMetadata.Quantity = 99
	} else if (total <= 0) {
		return &ItemEmptyError{}
	} else {
		m.MaterialMetadata.Quantity = total
	}
	return nil
}

//Inventory
func (i Inventory) AddMaterial(id int, quantity int, invSize int) error {
	if quantity <= 0 { return nil }
	item := ITEM_DATA[id].(*MaterialData)
	mat, ok := i[id]
	if !ok && len(i) >= invSize {
		return &InvFullError{}
	} else if !ok {
		i[id] = &Material{
			MaterialData: item,
			MaterialMetadata: MaterialMetadata{ Quantity: 0 },
		}
		mat = i[id]
	}
	err := mat.ChangeQuantity(quantity)
	return err
}

func (i Inventory) RemoveMaterial(id int, quantity int) error {
	if quantity <= 0 { return nil }
	mat, ok := i[id]
	if !ok { return nil }
	err := mat.ChangeQuantity(quantity)
	var iee *ItemEmptyError
	if errors.As(err, &iee) {
		delete(i, id)
	}
	return err
}

func (i Inventory) String() string {
	var invString strings.Builder; 
	for _, mat := range i {
		fmt.Fprintf(&invString, "%s: %d\n", mat.Name, mat.Quantity)
	}
	return invString.String()
}