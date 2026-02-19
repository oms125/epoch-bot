package bot

import (
	"log"
	"os"

	disc "github.com/bwmarrin/discordgo"
	. "github.com/oms125/epoch-bot/game/database"
)

var (
	GuildID string
)

type Bot struct {
	Session *disc.Session
	Game *Game
	ID string
}

func NewBot(g *Game) *Bot {
	botToken, ok := os.LookupEnv("EPOCH_BOT_TOKEN")
	if !ok { log.Fatal("Failed to initialize bot: EPOCH_BOT_TOKEN") }
	botID, ok := os.LookupEnv("EPOCH_BOT_ID")
	if !ok { log.Fatal("Failed to initialize bot: EPOCH_BOT_ID") }
	guildID, ok := os.LookupEnv("EPOCH_BOT_GUILD_ID")
	if !ok { log.Fatal("Failed to initialize bot: EPOCH_BOT_GUILD_ID")}
	GuildID = guildID

	session, err := disc.New("Bot " + botToken)
	if err != nil { log.Fatal("Failed to initialize bot: ", err) }

	session.Identify.Intents = disc.MakeIntent(disc.IntentsAll)
	session.State.MaxMessageCount = 10

	return &Bot {
		Session: session,
		Game: g,
		ID: botID,
	}
}