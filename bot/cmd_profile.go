package bot

import (
	"log"
	disc "github.com/bwmarrin/discordgo"
)

func (b *Bot) ProfileCommand() (*disc.ApplicationCommand, Handler) {
	return &disc.ApplicationCommand {
		Name: "profile",
		Description: "View your player profile",
	},
	func(s *disc.Session, i *disc.InteractionCreate) {
		var msg string
		p, err := b.Game.GetPlayer(i.Member)
		if err != nil {
			msg = "Unable to fetch profile data at this time"
		} else {
			msg = p.Inv.String()
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
		if err != nil { log.Printf("Failed interaction for command: profile, %v", err) }
	}
}