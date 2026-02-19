package bot

import (
	"errors"
	"log"

	disc "github.com/bwmarrin/discordgo"
	"github.com/oms125/epoch-bot/game/items"
)

func (b *Bot) GetCommand() (*disc.ApplicationCommand, Handler) {
	return &disc.ApplicationCommand{
			Name:                     "get",
			Description:              "Get item by ID",
			DefaultMemberPermissions: &PERM_ADMIN,
			Options: []*disc.ApplicationCommandOption{
				{
					Type:        disc.ApplicationCommandOptionInteger,
					Name:        "item",
					Description: "The ID of the item you want",
					Required:    true,
					MinValue:    &MIN_VAL,
				},
				{
					Type:        disc.ApplicationCommandOptionInteger,
					Name:        "quantity",
					Description: "How many items you want (only for non-gears)",
					Required:    false,
				},
			},
		},
		func(s *disc.Session, i *disc.InteractionCreate) {
			var msg string
			p, err := b.Game.GetPlayer(i.Member)
			if err != nil {
				msg = "Unable to fetch player data at this time"
			} else {
				id := int(i.ApplicationCommandData().Options[0].IntValue())
				var quantity = 1
				if len(i.ApplicationCommandData().Options) > 1 {
					quantity = int(i.ApplicationCommandData().Options[1].IntValue())
				}
				err := p.AddItem(id, quantity)
				var ime *items.ItemMaxError
				var ife *items.InvFullError
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
				&disc.InteractionResponse{
					Type: disc.InteractionResponseChannelMessageWithSource,
					Data: &disc.InteractionResponseData{
						Content: msg,
					},
				},
			)
			if err != nil {
				log.Printf("Failed interaction for command: get, %v", err)
			}
		}
}
