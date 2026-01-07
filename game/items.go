package game

import "fmt"

var (
	//Constant Item Data Structures
	ITEMS map[int]*Item
	TOOLS map[int]*Tool
)

//Constant Item Data
type (
	Item struct {
		ID int
		Name string
	}
	Tool struct {
		Item
		MaxDurability int
	}
)

//Metadata
type (
	ToolMetadata struct {
		Durability int
	}
)

//Errors
type (
	InvFullError struct {}
	ItemMaxError struct { 
		Item *Item
		Overflow int
	}
)
func (err *InvFullError) Error() string {
	return "Inventory Full"
}
func (err *ItemMaxError) Error() string {
	return fmt.Sprintf("Max Stack Size Reached for %s: %d", err.Item.Name, err.Overflow)
}

//Game Logic


//Init Items
const (
	ROCK = iota
	STICK
	SWORD
)

func init() {
	ITEMS = map[int]*Item {
		ROCK: {ID: ROCK, Name: "Rock"},
		STICK: {ID: STICK, Name: "Stick"},
	}
	TOOLS = map[int]*Tool {
		SWORD: {Item: Item{ID: SWORD, Name: "Sword"}, MaxDurability: 100},
	}
}


//Database Logic
func (g *Game) loadInventory(p *Player) error {
	query := `SELECT item_id, quantity FROM inventory WHERE user_id = ?`
	rows, err := g.DB.Query(query, p.ID)
	if err != nil { return err }
	defer rows.Close()

	items := make(map[*Item]int)

	for rows.Next() {
		var (
			id int
			quantity int
		)
		rows.Scan(&id, &quantity)
		items[ITEMS[id]] = quantity
	}

	p.Inv = items
	return nil
}

func (g *Game) saveInventory(p *Player) error {
	tx, err := g.DB.Begin()
	if err != nil { return err }
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM inventory WHERE user_id = ?`, p.ID)
	if err != nil { return err }

	stmt, err := tx.Prepare(`INSERT INTO inventory (item_id, user_id, quantity) VALUES (?, ?, ?)`)
	if err != nil { return err }

	for item, quantity := range p.Inv {
		_, err := stmt.Exec(p.ID, item, quantity)
		if err != nil { return err }
	}

	return tx.Commit()
}
