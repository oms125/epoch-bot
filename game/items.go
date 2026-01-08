package game

import "fmt"

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
		ROCK: &MaterialData{ID: ROCK, Name: "Rock"},
		STICK: &MaterialData{ID: STICK, Name: "Stick"},
		SWORD: &ToolData{ID: SWORD, Name: "Sword", MaxDurability: 100},
	}
}

//Game Logic
type (
	//Item Interfaces
	Item interface {
		Type() int
	}
	ItemData interface {
		Type() int
	}
	ItemMetadata interface {
		Type() int
	}
	//Material Structures
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
		Quantity int
	}
	//Tool Structures
	Tool struct {
		Item
		*ToolData
		*ToolMetadata
	}
	ToolData struct {
		ItemData
		ID int
		Name string
		MaxDurability int
	}
	ToolMetadata struct {
		ItemMetadata
		Durability int
	}
	//Errors
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

//Material Functions
func (m *MaterialData) Type() int { return TYPE_MATERIAL }
func (m *MaterialMetadata) Type() int { return TYPE_MATERIAL }
func (m *Material) Type() int { return TYPE_MATERIAL }
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

//Tool Functions
func (t *ToolData) Type() int { return TYPE_TOOL }
func (t *ToolMetadata) Type() int { return TYPE_TOOL }
func (t *Tool) Type() int { return TYPE_TOOL }

func NewDefaultTool(id int) *Tool {
	data := ITEM_DATA[id].(*ToolData)
	return &Tool{
		ToolData: data,
		ToolMetadata: &ToolMetadata{
			Durability: data.MaxDurability,
		},
	}
}

//Error Functions
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

//Database Logic
func (g *Game) loadInventory(p *Player) error {
	query := `SELECT item_id, quantity FROM inventory WHERE user_id = ?`
	rows, err := g.DB.Query(query, p.ID)
	if err != nil { return err }
	defer rows.Close()

	inv := make(map[int]*Material)

	for rows.Next() {
		var (
			id int
			quantity int
		)
		rows.Scan(&id, &quantity)
		data := ITEM_DATA[id].(*MaterialData)
		inv[id] = &Material{
			MaterialData: data,
			MaterialMetadata: &MaterialMetadata{ Quantity: quantity },
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

	for id, item := range p.Inv {
		_, err = tx.Exec(`INSERT INTO inventory (user_id, item_id, quantity) VALUES (?, ?, ?)`, p.ID, id, item.Quantity)
		if err != nil { return err }
	}

	return tx.Commit()
}

func (g *Game) loadArmory(p *Player) error {
	query := `SELECT item_id, durability FROM armory WHERE user_id = ?`
	rows, err := g.DB.Query(query, p.ID)
	if err != nil { return err }
	defer rows.Close()

	arm := make([]*Tool, 0, p.ArmSize)

	for rows.Next() {
		var (
			id int
			durability int
		)
		rows.Scan(&id, &durability)
		data := ITEM_DATA[id].(*ToolData)
		arm = append(arm, &Tool{
			ToolData: data,
			ToolMetadata: &ToolMetadata{ Durability: durability },
		})
		
	}

	p.Arm = arm
	return nil
}

func (g *Game) saveArmory(p *Player) error {
	tx, err := g.DB.Begin()
	if err != nil { return err }
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM armory WHERE user_id = ?`, p.ID)
	if err != nil { return err }

	for _, item := range p.Arm {
		_, err = tx.Exec(`INSERT INTO armory (user_id, item_id, durability) VALUES (?, ?, ?)`, p.ID, item.ID, item.Durability)
		if err != nil { return err }
	}

	return tx.Commit()
}