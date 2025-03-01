package main

import (
	"log"
	"os"

	"github.com/switchupcb/disgo"
	"github.com/switchupcb/disgoform"
)

// Environment Variables.
var (
	// token represents the Discord Bot's token.
	token = os.Getenv("TOKEN")

	// appid represents the Discord Bot's ApplicationID.
	//
	// Use Developer Mode to find it OR call GetCurrentUser (request) in your program and set it programmatically.
	appid = os.Getenv("APPID")
)

func main() {
	bot := &disgo.Client{
		ApplicationID:  appid,
		Authentication: disgo.BotToken(token),
		Config:         disgo.DefaultConfig(),
	}

	disgoform.GlobalApplicationCommands = []disgo.CreateGlobalApplicationCommand{
		// Command 1
		{
			Name:        "main",
			Description: disgo.Pointer("A basic command."),
			Type:        disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
		},

		// Command 2
		{
			Name:        "followup",
			Description: disgo.Pointer("Showcase multiple types of interaction responses."),
		},

		// Command 3
		{
			Name:        "autocomplete",
			Description: disgo.Pointer("Learn about autocompletion."),
			Options: []*disgo.ApplicationCommandOption{
				{
					Name:        "freewill",
					Description: "Do you have it?",
					Type:        disgo.FlagApplicationCommandOptionTypeSTRING,
					Required:    disgo.Pointer(true),
					Choices: []*disgo.ApplicationCommandOptionChoice{
						{
							Name:  "Yes",
							Value: "y",
						},
						{
							Name:  "No",
							Value: "n",
						},
					},
				},
				{
					Name:         "confirm",
					Description:  "Confirm your answer.",
					Type:         disgo.FlagApplicationCommandOptionTypeSTRING,
					Required:     disgo.Pointer(true),
					Autocomplete: disgo.Pointer(true),
				},
			},
		},
	}

	if err := disgoform.Sync(bot); err != nil {
		log.Printf("can't synchronize application commands with Discord: %v", err)
	}
}
