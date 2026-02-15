package bot

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/oms125/epoch-bot/game"
)

func (b *Bot) GetCommand() (*discordgo.ApplicationCommand, Handler) {
	return &discordgo.ApplicationCommand {
		Name: "get",
		Description: "Get item by ID",
		DefaultMemberPermissions: &PERM_ADMIN,
		Options: []*discordgo.ApplicationCommandOption {
			{
				Type: discordgo.ApplicationCommandOptionInteger,
				Name: "item",
				Description: "The ID of the item you want",
				Required: true,
			},
			{
				Type: discordgo.ApplicationCommandOptionInteger,
				Name: "quantity",
				Description: "How many items you want (only for non-tools)",
				Required: false,
			},
		},
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var msg string
		p, err := b.Game.GetPlayer(i.Member)
		if err != nil {
			msg = "Unable to fetch player data at this time"
		} else {
			id := int(i.ApplicationCommandData().Options[0].IntValue())
			var quantity = 1
			if (len(i.ApplicationCommandData().Options) > 1) {
				quantity = int(i.ApplicationCommandData().Options[1].IntValue())
			}
			err := p.AddItem(id, quantity)
			var ime *game.ItemMaxError
			var ife *game.InvFullError
			if errors.As(err, &ime) {
				msg = ime.Error()
			} else if errors.As(err, &ife) {
				msg = ife.Error()
			} else if err != nil {
				msg = "Unable to add item: " + err.Error()
			} else {
				msg = "Item added!"
			}
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
		if err != nil { log.Printf("Failed interaction for command: get, %v", err) }
	}
}