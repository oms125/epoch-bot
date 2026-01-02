package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/oms125/epoch-bot/bot"
	"github.com/oms125/epoch-bot/commands/slash"
)

var (
	SlashCommands map[string]func() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) = make(map[string]func() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)))
	SlashCommandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	ComponentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func Init() {
	slash.Init(&ComponentHandlers)

	populateSlashCommands()

	addSlashCommands()
}

func addSlashCommands() {
	for _, slashCommand := range SlashCommands {
		command, handler := slashCommand()
		_, err := bot.Session.ApplicationCommandCreate(bot.ID, "", command)
		if err != nil { 
			log.Printf("Error registering command: %s, %s", command.Name, err)
		} else {
			log.Printf("Registered command: %s", command.Name)
		}

		SlashCommandHandlers[command.Name] = handler
	}

	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			if command, ok := SlashCommandHandlers[data.Name]; ok {
				command(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if command, ok := ComponentHandlers[i.MessageComponentData().CustomID]; ok {
				command(s, i)
			}
		case discordgo.InteractionModalSubmit:
			if command, ok := ComponentHandlers[i.ModalSubmitData().CustomID]; ok {
				command(s, i)
			}
		}
	})
}