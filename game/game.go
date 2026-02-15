package game

import (
	"log"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Game struct {
	DB *sql.DB
	ActivePlayers map[string]*Player
}

//Game Logic


//Database Logic
func New() *Game {
	log.Println("Initializing database...")
	db, err := sql.Open("sqlite", "epoch.db")
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
			user_id TEXT,
			item_id INTEGER,
			quantity INTEGER,
			PRIMARY KEY (user_id, item_id)
		);`,
		//Armory Table
		`CREATE TABLE IF NOT EXISTS armory (
			column INTEGER PRIMARY KEY,
			user_id TEXT,
			item_id INTEGER,
			durability INTEGER
		);`,
	}

	for _, table := range tables {
		if _, err := g.DB.Exec(table); err != nil {
			log.Printf("Error setting up database table: %s, %v", table, err)
		}
	}
}