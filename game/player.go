package game

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Player struct {
	ID string
	Lvl int
	InvSize int

	Inv Inventory
}

//Game Logic
func (g *Game) GetPlayer(ID string) (*Player, error) {
	p, ok := g.ActivePlayers[ID]
	if !ok {
		return g.loadPlayer(ID)
	}
	return p, nil
}

//Database Logic
func (g *Game) loadPlayer(ID string) (*Player, error) {
	//Load Player Data
	p := &Player {}
	query := `SELECT id, lvl, inv_size FROM players WHERE id = ?`
	err := g.DB.QueryRow(query, ID).Scan(
		&p.ID,
		&p.Lvl,
		&p.InvSize,
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
	if err != nil {
		return nil, err
	}
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
	return err
}