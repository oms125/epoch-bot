package bot

import (
	disc "github.com/bwmarrin/discordgo"
	. "github.com/oms125/epoch-bot/game"
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
				Value: p.Equipment.SlotValue("head"),
			},
			{
				Name: "Body",
				Value: p.Equipment.SlotValue("body"),
			},
			{
				Name: "Legs",
				Value: p.Equipment.SlotValue("legs"),
				Inline: true,
			},
			{ Inline: true },
			{
				Name: "Primary",
				Value: p.Equipment.SlotValue("primary"),
				Inline: true,
			},
			{
				Name: "Feet",
				Value: p.Equipment.SlotValue("feet"),
				Inline: true,
			},
			{ Inline: true },
			{
				Name: "Secondary",
				Value: p.Equipment.SlotValue("secondary"),
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