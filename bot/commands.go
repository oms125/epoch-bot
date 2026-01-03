package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	AppCommand *discordgo.ApplicationCommand
	CommandHandler Handler
}

var (
	Commands map[string]Command = make(map[string]Command)
	SlashCommandHandlers map[string]Handler = make(map[string]Handler)
)

func (b *Bot) InitCommands() {
	Commands["test"] = Command {
		AppCommand: &discordgo.ApplicationCommand {
		Name: "test",
		Description: "just a test",
		},
		CommandHandler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
	}

	b.AddCommands()
}

func (b *Bot) AddCommands() {
	for _, command := range Commands {
		_, err := b.Session.ApplicationCommandCreate(b.ID, "", command.AppCommand)
		if err != nil { 
			log.Printf("Error registering command: %s, %s", command.AppCommand.Name, err)
		} else {
			log.Printf("Registered command: %s", command.AppCommand.Name)
		}
		SlashCommandHandlers[command.AppCommand.Name] = command.CommandHandler
	}

	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			if command, ok := SlashCommandHandlers[data.Name]; ok {
				command(s, i)
			}
		}
	})
}