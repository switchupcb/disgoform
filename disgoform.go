package disgoform

import "github.com/switchupcb/disgo"

var (
	// GlobalApplicationCommands represents the global application commands of the bot (max: 111).
	//
	// https://discord.com/developers/docs/interactions/application-commands#registering-a-command
	GlobalApplicationCommands []disgo.CreateGlobalApplicationCommand

	// GuildApplicationCommands represents the guild application commands of the bot.
	//
	// https://discord.com/developers/docs/interactions/application-commands#making-a-guild-command
	GuildApplicationCommands []disgo.CreateGuildApplicationCommand
)
