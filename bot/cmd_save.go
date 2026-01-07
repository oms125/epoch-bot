package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) SaveCommand() (*discordgo.ApplicationCommand, Handler) {
	return &discordgo.ApplicationCommand {
		Name: "save",
		Description: "Save your current player state",
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var msg string
		err := b.Game.SavePlayer(i.Member.User.ID)
		if err != nil {
			msg = "Failed to save player data at this time"
		} else {
			msg = "Player data saved!"
		}
		err = s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: msg,
				},
			},
		)
		if err != nil { log.Printf("Failed interaction for command: save, %v", err) }
	}
}