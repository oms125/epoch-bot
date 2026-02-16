package bot

import (
	"log"

	disc "github.com/bwmarrin/discordgo"
)

type Handler func(s *disc.Session, i *disc.InteractionCreate)

var (
	PERM_ALL int64 = disc.PermissionUseApplicationCommands
	PERM_ADMIN int64 = disc.PermissionAdministrator

	MIN_VAL float64 = 0

	Commands []*disc.ApplicationCommand
	CommandHandlers map[string]Handler
)

func (b *Bot) InitCommands() {
	Commands = []*disc.ApplicationCommand{}
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
	//Equipment Command
	cmd, hand = b.EquipmentCommand()
	Commands = append(Commands, cmd)
	CommandHandlers["equipment"] = hand
	//Equip Command
	cmd, hand = b.EquipCommand()
	Commands = append(Commands, cmd)
	CommandHandlers["equip"] = hand

	//Register
	b.registerCommands()
}

func (b *Bot) registerCommands() {
	//Register Commands and Handlers
	_, err := b.Session.ApplicationCommandBulkOverwrite(b.ID, GuildID, Commands)
	if err != nil { log.Fatal(err) }

	b.Session.AddHandler(func(s *disc.Session, i *disc.InteractionCreate) {
		switch i.Type {
		case disc.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			if command, ok := CommandHandlers[data.Name]; ok {
				command(s, i)
			}
		}
	})
}

//General command helper functions
func GetAuthor(m *disc.Member) *disc.MessageEmbedAuthor {
	return &disc.MessageEmbedAuthor{
		Name: m.User.Username,
		IconURL: m.AvatarURL(""),
	}
}

func Embed(s *disc.Session, i *disc.InteractionCreate, msg *disc.MessageEmbed, files ...*disc.File) error {
	return s.InteractionRespond(
		i.Interaction,
		&disc.InteractionResponse{
			Type: disc.InteractionResponseChannelMessageWithSource,
			Data: &disc.InteractionResponseData{
				Embeds: []*disc.MessageEmbed{
					msg,
				},
				Files: files,
			},
		},
	)
}

func Error(cmdName string, err error) {
	log.Printf("Failed interaction for command %s:\n%v", cmdName, err)
}