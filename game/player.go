package game

import (
	"database/sql"
	"fmt"
	disc "github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
	. "github.com/oms125/epoch-bot/game/player"
)

func ErrorMsg(msg string) string {
	return "ERROR; " + msg + "\n%w"
}

//Player
func (g *Game) GetPlayer(member *disc.Member) (*Player, error) {
	return g.GetPlayerByID(member.User.ID)
}

func (g *Game) GetPlayerByID(ID string) (*Player, error) {
	p, ok := g.ActivePlayers[ID]
	if !ok {
		return g.loadPlayer(ID)
	}
	return p, nil
}

func (g *Game) SavePlayer(member *disc.Member) error {
	return g.savePlayer(member.User.ID)
}

func (g *Game) loadPlayer(ID string) (*Player, error) {
	errMsg := ErrorMsg("Failed to load player %s")
	//Load Player Data
	p := &Player {}
	query := `SELECT id, lvl, inv_size, arm_size FROM players WHERE id = ?`
	err := g.DB.QueryRow(query, ID).Scan(
		&p.ID,
		&p.Lvl,
		&p.InvSize,
		&p.ArmSize,
	)
	if err != nil { 
		if err == sql.ErrNoRows {
			err = g.addPlayer(ID)
			if err != nil { return nil, fmt.Errorf(errMsg, ID, err) }
			return g.loadPlayer(ID)
		}
		return nil, err 
	}
	//Load Player Inventory
	err = g.loadInventory(p)
	if err != nil { return nil, fmt.Errorf(errMsg, ID, err) }
	err = g.loadArmory(p)
	if err != nil { return nil, fmt.Errorf(errMsg, ID, err) }
	p.BuildEquipment()

	g.ActivePlayers[ID] = p
	return p, nil
}

func (g *Game) addPlayer(ID string) error {
	errMsg := ErrorMsg("Failed to add player %s")

	query := `INSERT INTO players (id) VALUES (?)`
	_, err := g.DB.Exec(query, ID)
	if err != nil {
		return fmt.Errorf(errMsg, ID, err)
	}
	return nil
}

func (g *Game) savePlayer(ID string) error {
	errMsg := ErrorMsg("Failed to save player %s")

	p, ok := g.ActivePlayers[ID]
	if !ok { return nil }

	query := `
	UPDATE players SET
		lvl = ?
	WHERE id = ?`

	_, err := g.DB.Exec(query, p.Lvl, ID)
	if err != nil { return fmt.Errorf(errMsg, ID, err) }
	
	err = g.saveInventory(p)
	if err != nil { return fmt.Errorf(errMsg, ID, err) }
	err = g.saveArmory(p)
	if err != nil { return fmt.Errorf(errMsg, ID, err) }

	return nil
}

//Armory
type (
	loadTool struct {
		ID int `db:"item_id"`
		ToolMetadata `db:",inline"`
	}
	saveTool struct {
		UserID string `db:"user_id"`
		ItemID int `db:"item_id"`
		ToolMetadata `db:",inline"`
	}
)

func (g *Game) loadArmory(p *Player) error {
	errMsg := ErrorMsg("Failed to load armory for player %s")

	query := fmt.Sprintf("SELECT %s FROM armory WHERE user_id = ?", FormatFields(ArmoryFields, false))
	
	rows, err := g.DB.Queryx(query, p.ID)
	if err != nil { return fmt.Errorf(errMsg, p.ID, err) }
	defer rows.Close()
	
	arm := make([]*Tool, 0, p.ArmSize)

	for rows.Next() {
		tool := &loadTool{}
		rows.StructScan(&tool)
		data := ITEM_DATA[tool.ID].(*ToolData)
		arm = append(arm, &Tool{
			ToolData: data,
			ToolMetadata: &tool.ToolMetadata,
		})
	}
	
	p.Arm = arm
	return nil
}

func (g *Game) saveArmory(p *Player) error {
	errMsg := ErrorMsg("Failed to save armory for player %s")

	tx, err := g.DB.Beginx()
	if err != nil { return fmt.Errorf(errMsg, p.ID, err) }
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM armory WHERE user_id = ?`, p.ID)
	if err != nil { return fmt.Errorf(errMsg, p.ID, err) }

	query := fmt.Sprintf(`INSERT INTO armory (user_id, %s) VALUES (:user_id, %s)`,
		FormatFields(ArmoryFields, false), FormatFields(ArmoryFields, true))
	for _, item := range p.Arm {
		meta := &saveTool{
			UserID: p.ID,
			ItemID: item.ID,
			ToolMetadata: *item.ToolMetadata,
		}
		_, err = tx.NamedExec(query, meta)
		if err != nil { return fmt.Errorf(errMsg, p.ID, err) }
	}

	return tx.Commit()
}

//Inventory
type (
	loadMaterial struct {
		ID int `db:"item_id"`
		MaterialMetadata `db:",inline"`
	}
	saveMaterial struct {
		UserID string `db:"user_id"`
		ItemID int `db:"item_id"`
		MaterialMetadata `db:",inline"`
	}
)

func (g *Game) loadInventory(p *Player) error {
	errMsg := ErrorMsg("Failed to load inventory for player %s")

	query := fmt.Sprintf("SELECT %s FROM inventory WHERE user_id = ?", FormatFields(InventoryFields, false))
	
	rows, err := g.DB.Queryx(query, p.ID)
	if err != nil { return fmt.Errorf(errMsg, p.ID, err) }
	defer rows.Close()

	inv := make(map[int]*Material)

	for rows.Next() {
		mat := &loadMaterial{}
		rows.StructScan(&mat)
		data := ITEM_DATA[mat.ID].(*MaterialData)
		inv[mat.ID] = &Material{
			MaterialData: data,
			MaterialMetadata: &mat.MaterialMetadata,
		}
		
	}

	p.Inv = inv
	return nil
}

func (g *Game) saveInventory(p *Player) error {
	errMsg := ErrorMsg("Failed to save inventory for player %s")

	tx, err := g.DB.Beginx()
	if err != nil { return fmt.Errorf(errMsg, p.ID, err) }
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM inventory WHERE user_id = ?`, p.ID)
	if err != nil { return fmt.Errorf(errMsg, p.ID, err) }

	query := fmt.Sprintf(`INSERT INTO inventory (user_id, %s) VALUES (:user_id, %s)`,
		FormatFields(InventoryFields, false), FormatFields(InventoryFields, true))
	for id, item := range p.Inv {
		meta := &saveMaterial{
			UserID: p.ID,
			ItemID: id,
			MaterialMetadata: *item.MaterialMetadata,
		}
		_, err = tx.NamedExec(query, meta)
		if err != nil { return fmt.Errorf(errMsg, p.ID, err) }
	}

	return tx.Commit()
}