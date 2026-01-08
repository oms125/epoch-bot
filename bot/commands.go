package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	PERM_ALL int64 = discordgo.PermissionUseApplicationCommands
	PERM_ADMIN int64 = discordgo.PermissionAdministrator
)

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

var (
	Commands []*discordgo.ApplicationCommand
	CommandHandlers map[string]Handler
)

func (b *Bot) InitCommands() {
	Commands = []*discordgo.ApplicationCommand{}
	CommandHandlers = make(map[string]Handler)
	//Profile Command
	cmd, hand := b.ProfileCommand()
	Commands = append(Commands, cmd)
	CommandHandlers["profile"] = hand
	//Save Command
	cmd, hand = b.SaveCommand()
	Commands = append(Commands, cmd)
	CommandHandlers["save"] = hand
	//Get Command
	cmd, hand = b.GetCommand()
	Commands = append(Commands, cmd)
	CommandHandlers["get"] = hand
	//Inventory Command
	cmd, hand = b.InventoryCommand()
	Commands = append(Commands, cmd)
	CommandHandlers["inventory"] = hand
	//Armory Command
	cmd, hand = b.ArmoryCommand()
	Commands = append(Commands, cmd)
	CommandHandlers["armory"] = hand

	//Register
	b.registerCommands()
}

func (b *Bot) registerCommands() {
	//Register Commands and Handlers
	_, err := b.Session.ApplicationCommandBulkOverwrite(b.ID, GuildID, Commands)
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