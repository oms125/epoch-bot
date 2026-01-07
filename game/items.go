package game

import "fmt"

//Constant Item Data Structures
var (
	ITEM_DATA map[int]ItemData
)

//Constant Item Type References
const (
	MATERIAL_TYPE = iota
	TOOL_TYPE
)

//Data
type (
	ItemData interface {
		Type() int
	}
	MaterialData struct {
		ItemData
		ID int
		Name string
	}
	ToolData struct {
		ItemData
		ID int
		Name string
		MaxDurability int
	}
)
func (m MaterialData) Type() int { return MATERIAL_TYPE }
func (t ToolData) Type() int { return TOOL_TYPE }

//Metadata
type (
	ItemMetadata interface {
		Type() int
	}
	MaterialMetadata struct {
		ItemMetadata
		Quantity int
	}
	ToolMetadata struct {
		ItemMetadata
		Durability int
	}
)
func (m MaterialMetadata) Type() int { return MATERIAL_TYPE }
func (t ToolMetadata) Type() int { return TOOL_TYPE }

//Item
type (
	Item interface {
		Type() int
		GetID() int
		GetName() string
		GetData() ItemData
		GetMetadata() ItemMetadata
	}
	Material struct {
		Item
		*MaterialData
		*MaterialMetadata
	}
	Tool struct {
		Item
		*ToolData
		*ToolMetadata
	}
)
func (m Material) Type() int { return MATERIAL_TYPE }
func (m Material) GetID() int { return m.MaterialData.ID }
func (m Material) GetName() string { return m.MaterialData.Name }
func (m Material) GetQuantity() int { return m.MaterialMetadata.Quantity }
func (m Material) GetData() ItemData { return m.MaterialData }
func (m Material) GetMetadata() ItemMetadata { return m.MaterialMetadata }
func (t Tool) Type() int { return TOOL_TYPE }
func (t Tool) GetID() int { return t.ToolData.ID }
func (t Tool) GetName() string { return t.ToolData.Name }
func (t Tool) GetDurability() int { return t.ToolMetadata.Durability }
func (t Tool) GetMaxDurability() int { return t.ToolData.MaxDurability }
func (t Tool) GetData() ItemData { return t.ToolData }
func (t Tool) GetMetadata() ItemMetadata { return t.ToolMetadata }

//Errors
type (
	InvFullError struct {}
	ItemMaxError struct { 
		Name string
		Overflow int
	}
)
func (err *InvFullError) Error() string {
	return "Inventory Full"
}
func (err *ItemMaxError) Error() string {
	return fmt.Sprintf("Max Stack Size Reached for %s: %d", err.Name, err.Overflow)
}

//Init Items
const (
	ROCK = iota
	STICK
	SWORD
)

func init() {
	ITEM_DATA = map[int]ItemData {
		ROCK: &MaterialData{ID: ROCK, Name: "Rock"},
		STICK: &MaterialData{ID: STICK, Name: "Stick"},
		SWORD: &ToolData{ID: SWORD, Name: "Sword", MaxDurability: 100},
	}
}

//Database Logic
func (g *Game) loadInventory(p *Player) error {
	query := `SELECT item_id, quantity, durability FROM inventory WHERE user_id = ?`
	rows, err := g.DB.Query(query, p.ID)
	if err != nil { return err }
	defer rows.Close()

	inv := make([]Item, 0, p.InvSize)

	for rows.Next() {
		var (
			id int
			quantity int
			durability int
		)
		rows.Scan(&id, &quantity, &durability)
		data := ITEM_DATA[id]
		switch data.Type() {
		case MATERIAL_TYPE:
			inv = append(inv, &Material{
				MaterialData: ITEM_DATA[id].(*MaterialData),
				MaterialMetadata: &MaterialMetadata{ Quantity: quantity },
			})
		case TOOL_TYPE:
			inv = append(inv, &Tool{
				ToolData: ITEM_DATA[id].(*ToolData),
				ToolMetadata: &ToolMetadata{
					Durability: durability,
				},
			})
		}
	}

	p.Inv = inv
	return nil
}

func (g *Game) saveInventory(p *Player) error {
	tx, err := g.DB.Begin()
	if err != nil { return err }
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM inventory WHERE user_id = ?`, p.ID)
	if err != nil { return err }

	for _, item := range p.Inv {
		switch item.Type() {
		case MATERIAL_TYPE:
			_, err = tx.Exec(`INSERT INTO inventory (user_id, item_id, quantity) VALUES (?, ?, ?)`, p.ID, item.GetID(), item.(*Material).GetQuantity())
			if err != nil { return err }
		case TOOL_TYPE:
			_, err = tx.Exec(`INSERT INTO inventory (user_id, item_id, durability) VALUES (?, ?, ?)`, p.ID, item.GetID(), item.(*Tool).GetDurability())
			if err != nil { return err }
		}
	}

	return tx.Commit()
}