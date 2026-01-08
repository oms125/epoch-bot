package game

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	_ "modernc.org/sqlite"
)

type Player struct {
	ID string
	Lvl int
	InvSize int
	ArmSize int

	Inv map[int]*Material
	Arm []*Tool
}

//Game Logic
func (g *Game) GetPlayer(ID string) (*Player, error) {
	p, ok := g.ActivePlayers[ID]
	if !ok {
		return g.loadPlayer(ID)
	}
	return p, nil
}

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

func (p *Player) AddTool(tool *Tool) error {
	if len(p.Arm) >= p.ArmSize {
		return &ArmFullError{}
	}
	p.Arm = append(p.Arm, tool)
	return nil
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

func (p *Player) RemoveTool(idx int) {
	p.Arm = slices.Delete(p.Arm, idx, idx+1)
}

func (p *Player) GetInventory() string {
	var invString strings.Builder; 
	//invString.WriteString("*Inventory*\n")
	size := len(p.Inv)
	fmt.Fprintf(&invString, "*Inventory:* (%d/%d)\n", size, p.InvSize)
	for _, mat := range p.Inv {
		fmt.Fprintf(&invString, "%s: %d\n", mat.Name, mat.Quantity)
	}
	return invString.String()
}

func (p *Player) GetArmory() string {
	var armString strings.Builder; 
	//armString.WriteString("*Armory* ")
	size := len(p.Arm)
	fmt.Fprintf(&armString, "*Armory:* (%d/%d)\n", size, p.ArmSize)
	for _, tool := range p.Arm {
		fmt.Fprintf(&armString, "%s: (%d/%d)\n", tool.Name, tool.Durability, tool.MaxDurability)
	}
	return armString.String()
}

//Database Logic
func (g *Game) loadPlayer(ID string) (*Player, error) {
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
			if err != nil { return nil, err }
			return g.loadPlayer(ID)
		}
		return nil, err 
	}
	//Load Player Inventory
	err = g.loadInventory(p)
	if err != nil { return nil, err }
	err = g.loadArmory(p)
	if err != nil { return nil, err }

	g.ActivePlayers[ID] = p
	return p, nil
}

func (g *Game) addPlayer(ID string) error {
	query := `INSERT INTO players (id) VALUES (?)`
	_, err := g.DB.Exec(query, ID)
	if err != nil {
		return fmt.Errorf("Could not add player: %s, %v", ID, err)
	}
	return nil
}

func (g *Game) SavePlayer(ID string) error {
	p, ok := g.ActivePlayers[ID]
	if !ok { return nil }

	query := `
	UPDATE players SET
		lvl = ?
	WHERE id = ?`

	_, err := g.DB.Exec(query, p.Lvl, ID)
	if err != nil { return err }

	err = g.saveInventory(p)
	if err != nil { return err }
	err = g.saveArmory(p)
	if err != nil { return err }

	return nil
}