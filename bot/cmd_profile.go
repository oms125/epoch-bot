package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) ProfileCommand() (*discordgo.ApplicationCommand, Handler) {
	return &discordgo.ApplicationCommand {
		Name: "profile",
		Description: "View your player profile",
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var msg string
		p, err := b.Game.GetPlayer(i.Member.User.ID)
		if err != nil {
			msg = "Unable to fetch profile data at this time"
		} else {
			msg = fmt.Sprintf("Level: %d", p.Lvl)
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
		if err != nil { log.Printf("Failed interaction for command: profile, %v", err) }
	}
}