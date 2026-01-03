package bot

import (
	"log"
	"os"
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	DB *sql.DB
	ID string
}

func New(db *sql.DB) *Bot {
	botToken, ok := os.LookupEnv("EPOCH_BOT_TOKEN")
	if !ok { log.Fatal("No Discord Bot Token: EPOCH_BOT_TOKEN") }
	botID, ok := os.LookupEnv("EPOCH_BOT_ID")
	if !ok { log.Fatal("No Discord Bot ID: EPOCH_BOT_ID") }

	session, err := discordgo.New("Bot " + botToken)
	if err != nil { log.Fatal("Failed to start bot: ", err) }

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	session.State.MaxMessageCount = 10

	return &Bot {
		Session: session,
		DB: db,
		ID: botID,
	}
}