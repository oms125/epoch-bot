package bot

import (
	disc "github.com/bwmarrin/discordgo"
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
				Value: p.Equipment["head"].Name,
			},
			{
				Name: "Body",
				Value: p.Equipment["body"].Name,
			},
			{
				Name: "Legs",
				Value: p.Equipment["legs"].Name,
				Inline: true,
			},
			{ Inline: true },
			{
				Name: "Primary",
				Value: p.Equipment["primary"].Name,
				Inline: true,
			},
			{
				Name: "Feet",
				Value: p.Equipment["feet"].Name,
				Inline: true,
			},
			{ Inline: true },
			{
				Name: "Secondary",
				Value: p.Equipment["secondary"].Name,
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