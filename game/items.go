package game

var items map[int]Item

type Item struct {
	ID int
	Name string
	MaxQuantity int
	MaxDurability int
}

type ItemMetadata struct {
	Quantity int
	Durability int
}

type Inventory struct {
	Items map[int]ItemMetadata
}

//Game Logic

//Init Items
const (
	ROCK = iota
	STICK
	SWORD
)

func init() {
	items = map[int]Item {
		ROCK: makeMaterial(ROCK, "Rock"),
		STICK: makeMaterial(STICK, "Stick"),
		SWORD: makeTool(SWORD, "Sword", 100),
	}
}

func makeMaterial(id int, name string) Item {
	return Item {
		ID: id,
		Name: name,
		MaxQuantity: 99,
		MaxDurability: -1,
	}
}

func makeTool(id int, name string, maxDurability int) Item {
	return Item {
		ID: id,
		Name: name,
		MaxQuantity: 1,
		MaxDurability: maxDurability,
	}
}


//Database Logic
func (g *Game) loadInventory(p *Player) error {
	query := `SELECT item_id, quantity, durability FROM inventory WHERE user_id = ?`
	rows, err := g.DB.Query(query, p.ID)
	if err != nil { return err }
	defer rows.Close()

	items := make(map[int]ItemMetadata)

	for rows.Next() {
		var (
			id int
			quantity int
			durability int
		)
		rows.Scan(&id, &quantity, &durability)
		items[id] = ItemMetadata{
			Quantity: quantity,
			Durability: durability,
		}
	}

	p.Inv = Inventory {
		Items: items,
	}
	return nil
}

func (g *Game) saveInventory(p *Player) error {
	tx, err := g.DB.Begin()
	if err != nil { return err }
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM inventory WHERE user_id = ?`, p.ID)
	if err != nil { return err }

	stmt, err := tx.Prepare(`INSERT INTO inventory (item_id, user_id, quantity, durability) VALUES (?, ?, ?, ?)`)
	if err != nil { return err }

	for item, meta := range p.Inv.Items {
		_, err := stmt.Exec(item, p.ID, meta.Quantity, meta.Durability)
		if err != nil { return err }
	}

	return tx.Commit()
}
