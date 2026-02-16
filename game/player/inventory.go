package player

import (
	"fmt"
	"strings"
	"errors"
)

func (m *MaterialData) Type() int { return TYPE_MATERIAL }
func (m *MaterialMetadata) Type() int { return TYPE_MATERIAL }
func (m *Material) Type() int { return TYPE_MATERIAL }

type (
	Inventory map[int]*Material

	Material struct {
		Item
		*MaterialData
		*MaterialMetadata
	}
	MaterialData struct {
		ItemData
		ID int
		Name string
	}
	MaterialMetadata struct {
		ItemMetadata
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

//Game Logic
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

func (p *Player) AddMaterial(id int, quantity int) error {
	if quantity <= 0 { return nil }
	item := ITEM_DATA[id].(*MaterialData)
	mat, ok := p.Inv[id]
	if !ok && len(p.Inv) >= p.InvSize {
		return &InvFullError{}
	} else if !ok {
		p.Inv[id] = &Material{
			MaterialData: item,
			MaterialMetadata: &MaterialMetadata{ Quantity: 0 },
		}
		mat = p.Inv[id]
	}
	err := mat.ChangeQuantity(quantity)
	return err
}

func (p *Player) RemoveMaterial(id int, quantity int) error {
	if quantity <= 0 { return nil }
	mat, ok := p.Inv[id]
	if !ok { return nil }
	err := mat.ChangeQuantity(quantity)
	var iee *ItemEmptyError
	if errors.As(err, &iee) {
		delete(p.Inv, id)
	}
	return err
}

//Helpers
func (i Inventory) String() string {
	var invString strings.Builder; 
	for _, mat := range i {
		fmt.Fprintf(&invString, "%s: %d\n", mat.Name, mat.Quantity)
	}
	return invString.String()
}
