package game

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Player struct {
	ID string
	Lvl int
}

//Game Logic for Players
func (g *Game) LevelUpPlayer(ID string) error {
	p, err := g.GetPlayer(ID)
	if err != nil { return err }
	p.Lvl += 1
	return nil
}

//Database Logic for Players
func (g *Game) GetPlayer(ID string) (*Player, error) {
	p, ok := g.ActivePlayers[ID]
	if ok { return p, nil }

	p = &Player {}
	query := `SELECT id, lvl FROM players WHERE id = ?`
	err := g.DB.QueryRow(query, ID).Scan(
		&p.ID,
		&p.Lvl,
	)
	if err != nil { 
		if err == sql.ErrNoRows {
			err = g.AddPlayer(ID)
			if err != nil { return nil, err }
			return g.GetPlayer(ID)
		}
		return nil, err 
	}
	g.ActivePlayers[ID] = p
	return p, nil
}

func (g *Game) AddPlayer(ID string) error {
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
	return err
}