package bot

import (
	disc "github.com/bwmarrin/discordgo"
	. "github.com/oms125/epoch-bot/game/player"
)

func (b *Bot) EquipmentCommand() (*disc.ApplicationCommand, Handler) {
	cmd := &disc.ApplicationCommand{
		Name:        "equipment",
		Description: "View your currently equipped gear",
	}
	return cmd,
		func(s *disc.Session, i *disc.InteractionCreate) {
			msg, file := equipmentMsg(b, i)
			err := Embed(s, i, msg, file...)

			if err != nil { Error(cmd.Name, err) }
		}
}

func equipmentMsg(b *Bot, i *disc.InteractionCreate) (*disc.MessageEmbed, []*disc.File) {
	p, _ := b.Game.GetPlayer(i.Member)
	img := Image{
		Name: "armory_icon.png",
		Type: ICONS,
	}

	return &disc.MessageEmbed{
		Author: GetAuthor(i.Member),
		Thumbnail: &disc.MessageEmbedThumbnail{
			URL: img.Attach(),
		},
		Title: "Equipment",
		Color: 0x999999,
		Fields: []*disc.MessageEmbedField{
			{
				Name: "Head",
				Value: slotValue(p, HEAD),
			},
			{
				Name: "Body",
				Value: slotValue(p, BODY),
			},
			{
				Name: "Legs",
				Value: slotValue(p, LEGS),
				Inline: true,
			},
			{ Inline: true },
			{
				Name: "Primary",
				Value: slotValue(p, PRIMARY),
				Inline: true,
			},
			{
				Name: "Feet",
				Value: slotValue(p, FEET),
				Inline: true,
			},
			{ Inline: true },
			{
				Name: "Secondary",
				Value: slotValue(p, SECONDARY),
				Inline: true,
			},
		},
	},
	[]*disc.File{
		{
			Name: img.Name,
			Reader: img.File(),
		},
	}
}

func slotValue(p *Player, slot int) string {
	equipment := p.Equipment[slot]
	if equipment == nil {
		return "Nothing equipped"
	}
	return equipment.Name
}