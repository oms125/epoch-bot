package bot

import (
	"log"
	disc "github.com/bwmarrin/discordgo"
)

func (b *Bot) SaveCommand() (*disc.ApplicationCommand, Handler) {
	cmd := &disc.ApplicationCommand{
		Name: "save",
		Description: "Save your current player state",
	}
	return cmd,
	func(s *disc.Session, i *disc.InteractionCreate) {
		var msg string
		err := b.Game.SavePlayer(i.Member)
		if err != nil {
			msg = "Failed to save player data at this time"
			log.Printf("Failed to save data for player %s:\n%v", i.Member.User.ID, err)
		} else {
			msg = "Player data saved!"
		}
		err = s.InteractionRespond(
			i.Interaction,
			&disc.InteractionResponse {
				Type: disc.InteractionResponseChannelMessageWithSource,
				Data: &disc.InteractionResponseData {
					Content: msg,
				},
			},
		)
		if err != nil { log.Printf("Failed interaction for command: save, %v", err) }
	}
}