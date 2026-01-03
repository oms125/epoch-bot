package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

var (
	Commands []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{}
	CommandHandlers map[string]Handler = make(map[string]Handler)
)

func (b *Bot) InitCommands() {
	//Profile Command
	Commands = append(Commands, &discordgo.ApplicationCommand {
		Name: "profile",
		Description: "View your player profile",
	})
	CommandHandlers["profile"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var msg string
		p, err := b.Game.GetPlayer(i.Member.User.ID)
		if err != nil {
			msg = fmt.Sprintf("Unable to fetch profile data at this time: %v", err)
		} else {
			msg = fmt.Sprintf("Level: %d", p.Lvl)
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
		if err != nil { log.Printf("Message failed to send for command: profile")}
	}
	
	//Save Command
	Commands = append(Commands, &discordgo.ApplicationCommand {
		Name: "save",
		Description: "Save your current player state",
	})
	CommandHandlers["save"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var msg string
		err := b.Game.SavePlayer(i.Member.User.ID)
		if err != nil {
			msg = "Failed to save player data at this time"
		} else {
			msg = "Player data saved!"
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
		if err != nil { log.Printf("Message failed to send for command: save")}
	}

	//Register Commands and Handlers
	_, err := b.Session.ApplicationCommandBulkOverwrite(b.ID, "", Commands)
	if err != nil { log.Fatal(err) }

	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			if command, ok := CommandHandlers[data.Name]; ok {
				command(s, i)
			}
		}
	})
}