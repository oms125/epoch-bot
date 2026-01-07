package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/oms125/epoch-bot/game"
)

type Bot struct {
	Session *discordgo.Session
	Game *game.Game
	ID string
}

func New(g *game.Game) *Bot {
	botToken, ok := os.LookupEnv("EPOCH_BOT_TOKEN")
	if !ok { log.Fatal("Failed to initialize bot: EPOCH_BOT_TOKEN") }
	botID, ok := os.LookupEnv("EPOCH_BOT_ID")
	if !ok { log.Fatal("Failed to initialize bot: EPOCH_BOT_ID") }

	session, err := discordgo.New("Bot " + botToken)
	if err != nil { log.Fatal("Failed to initialize bot: ", err) }

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	session.State.MaxMessageCount = 10

	return &Bot {
		Session: session,
		Game: g,
		ID: botID,
	}
}