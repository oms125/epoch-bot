package bot

import (
	disc "github.com/bwmarrin/discordgo"
	. "github.com/oms125/epoch-bot/game"
)

func (b *Bot) ArmoryCommand() (*disc.ApplicationCommand, Handler) {
	cmd := &disc.ApplicationCommand{
		Name:        "armory",
		Description: "View your armory",
	}
	return cmd,
		func(s *disc.Session, i *disc.InteractionCreate) {
			_, err := b.Game.GetPlayer(i.Member)

			embed, file := armoryMsg(b, i)
			err = Embed(s, i, embed, file...)

			if err != nil {
				Error(cmd.Name, err)
			}
		}
}

func armoryMsg(b *Bot, i *disc.InteractionCreate) (*disc.MessageEmbed, []*disc.File) {
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
		Color: 0x999999,
		Title: "Armory",
		Description: p.Arm.String(),
	},
	[]*disc.File{
		{
			Name: img.Name,
			Reader: img.File(),
		},
	}
}