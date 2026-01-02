package slash

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Test() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
		Name: "test",
		Description: "just a test",
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: "Response",
				},
			},
		)
		if err != nil { log.Printf("Message failed to send for command: test")}
	}
}