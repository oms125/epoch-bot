package game

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
	. "github.com/oms125/epoch-bot/game/player"
)

type Game struct {
	DB *sqlx.DB
	ActivePlayers map[string]*Player
}

func NewGame() *Game {
	log.Println("Initializing database...")
	db, err := sqlx.Connect("sqlite", "epoch.db")
	if err != nil { 
		log.Fatal("Failed to initialize database: ", err)
	} else {
		log.Println("Database initialized")
	}

	setupSQL := `
	PRAGMA journal_mode=WAL;
	PRAGMA busy_timeout=3000;
	PRAGMA synchronous=NORMAL;
	PRAGMA foreign_keys=ON;
	`
	if _, err := db.Exec(setupSQL); err != nil {
		log.Fatalf("Failed to set pragmas: %v", err)
	}

	return &Game {
		DB: db,
		ActivePlayers: make(map[string]*Player),
	}
}

func (g *Game) InitTables() {
	tables := []string {
		//Player Table
		`CREATE TABLE IF NOT EXISTS players (
			id TEXT PRIMARY KEY,
			lvl INTEGER DEFAULT 1,
			inv_size INTEGER DEFAULT 50,
			arm_size INTEGER DEFAULT 50
		);`,
		//Inventory Table
		`CREATE TABLE IF NOT EXISTS inventory (
			user_id TEXT, ` +
			FormatTableFields(InventoryFields) + `,
			PRIMARY KEY (user_id, item_id)
		);`,
		//Armory Table
		`CREATE TABLE IF NOT EXISTS armory (
			column INTEGER PRIMARY KEY,
			user_id TEXT, ` +
			FormatTableFields(ArmoryFields) + 
		`);`,
	}

	for _, table := range tables {
		if _, err := g.DB.Exec(table); err != nil {
			log.Fatalf("Error setting up database table: %s, %v", table, err)
		}
	}
}