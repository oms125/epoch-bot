package game

import (
	"fmt"
	"log"
)

type (
	Armory []*Tool

	loadTool struct {
		ID int `db:"item_id"`
		ToolMetadata
	}
	saveTool struct {
		UserID string `db:"user_id"`
		ItemID int `db:"item_id"`
		ToolMetadata
	}
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
		Durability int `db:"durability"`
		Equipped bool `db:"equipped"`
	}
)

var armory_fields = []string{
	"item_id",
	"durability",
	"equipped",
}

//Game Logic


//Database Logic
func (g *Game) loadArmory(p *Player) error {
	query := fmt.Sprintf("SELECT %s FROM armory WHERE user_id = ?", FormatFields(armory_fields, false))
	
	rows, err := g.DB.Queryx(query, p.ID)
	if err != nil {
		log.Printf("Failed to load armory for player %s:\n%v", p.ID, err)
		return err 
	}
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
	tx, err := g.DB.Beginx()
	if err != nil { return err }
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM armory WHERE user_id = ?`, p.ID)
	if err != nil { return err }

	query := fmt.Sprintf(`INSERT INTO armory (user_id, %s) VALUES (:user_id, %s)`, FormatFields(armory_fields, false), FormatFields(armory_fields, true))
	for _, item := range p.Arm {
		meta := &saveTool{
			UserID: p.ID,
			ItemID: item.ID,
			ToolMetadata: *item.ToolMetadata,
		}
		_, err = tx.NamedExec(query, meta)
		if err != nil { return err }
	}

	return tx.Commit()
}