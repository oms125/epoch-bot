package bot

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) InventoryCommand() (*discordgo.ApplicationCommand, Handler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "inventory",
		Description: "View your inventory",
	}
	return cmd,
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, err := b.Game.GetPlayer(i.Member)

			embed, file := b.invMessage(i)
			err = Embed(s, i, embed, file...)

			if err != nil {
				Error(cmd.Name, err)
			}
		}
}

func (b *Bot) invMessage(i *discordgo.InteractionCreate) (*discordgo.MessageEmbed, []*discordgo.File) {
	p, _ := b.Game.GetPlayer(i.Member)
	img := Image{
		Name: "inventory_icon.png",
		Type: ICONS,
	}

	return &discordgo.MessageEmbed{
		Author: GetAuthor(i.Member),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: img.Attach(),
		},
		Color: 0x999999,
		Title: "Inventory",
		Description: p.Inv.ToString(),
	},
	[]*discordgo.File{
		{
			Name: img.Name,
			Reader: img.File(),
		},
	}
}