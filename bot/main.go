package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
 	Session *discordgo.Session
	ID string
)

func init() {
	botToken, ok := os.LookupEnv("EPOCH_BOT_TOKEN")
	if !ok {
		log.Fatal("No Discord Bot Token Set: EPOCH_BOT_TOKEN")
	}
	botID, ok := os.LookupEnv("EPOCH_BOT_ID")
	if !ok {
		log.Fatal("No Discord Bot ID Set: EPOCH_BOT_ID")
	}

	ID = botID

	session, err := discordgo.New("Bot "+botToken)
	if err != nil { log.Fatal(err) }
	Session = session

	Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	Session.State.MaxMessageCount = 10
}