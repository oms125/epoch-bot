package bot

import (
	"fmt"
	disc "github.com/bwmarrin/discordgo"
)

func (b *Bot) EquipCommand() (*disc.ApplicationCommand, Handler) {
		cmd := &disc.ApplicationCommand{
		Name:        "equip",
		Description: "Equip a piece of gear from your armory",
		Options: []*disc.ApplicationCommandOption {
			{
				Type: disc.ApplicationCommandOptionInteger,
				Name: "armory_slot",
				Description: "The number of the armory slot that the item you want to equip is located in",
				Required: true,
				MinValue: &MIN_VAL,
			},
			{
				Type: disc.ApplicationCommandOptionString,
				Name: "slot",
				Description: "The slot to equip the item to (not necessary for items that can only equip to a single slot)",
				Required: false,
				MinValue: &MIN_VAL,
			},
		},
	}
	return cmd,
		func(s *disc.Session, i *disc.InteractionCreate) {
			p, _ := b.Game.GetPlayer(i.Member)
			var err error
			armSlot := int(i.ApplicationCommandData().Options[0].IntValue())
			if len(i.ApplicationCommandData().Options) > 1 {
				equipSlot := i.ApplicationCommandData().Options[1].StringValue()
				err = p.Equip(armSlot, equipSlot)
			} else {
				err = p.Equip(armSlot, "")
			}

			if err != nil {
				err = Embed(s, i, &disc.MessageEmbed{
					Author: GetAuthor(i.Member),
					Color: 0xc70404,
					Title: "Equip",
					Description: err.Error(),
				})
			} else {
				err = Embed(s, i, &disc.MessageEmbed{
					Author: GetAuthor(i.Member),
					Color: 0x04c20d,
					Title: "Equip",
					Description: fmt.Sprintf("Successfully equipped **%s**", p.Arm[armSlot].Name),
				})
			}

			if err != nil { Error(cmd.Name, err) }
		}
}