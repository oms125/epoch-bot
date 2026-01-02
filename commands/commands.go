package commands

import (
	"github.com/oms125/epoch-bot/commands/slash"
)

func populateSlashCommands() {
	SlashCommands["test"] = slash.Test
}