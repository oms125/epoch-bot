package bot

import (
	//"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) ArmoryCommand() (*discordgo.ApplicationCommand, Handler) {
	return &discordgo.ApplicationCommand {
		Name: "armory",
		Description: "View your armory",
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var msg string
		p, err := b.Game.GetPlayer(i.Member.User.ID)
		if err != nil {
			msg = "Unable to fetch player data at this time"
		} else {
			msg = p.GetArmory()
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
		if err != nil { log.Printf("Failed interaction for command: armory, %v", err) }
	}
}