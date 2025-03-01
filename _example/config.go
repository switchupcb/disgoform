package main

import (
	"log"

	"github.com/switchupcb/disgo"
	"github.com/switchupcb/disgoform"
)

func main() {
	bot := &disgo.Client{
		ApplicationID:  "TODO",
		Authentication: disgo.BotToken("TOKEN"), // or BearerToken("TOKEN")
		Config:         disgo.DefaultConfig(),
	}

	disgoform.GlobalApplicationCommands = []disgo.CreateGlobalApplicationCommand{
		// Command 1
		disgo.CreateGlobalApplicationCommand{
			Name:        "main",
			Description: disgo.Pointer("A basic command."),
			Type:        disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
		},
		// Command 2
		disgo.CreateGlobalApplicationCommand{
			Name:        "followup",
			Description: disgo.Pointer("Showcase multiple types of interaction responses."),
		},
		// Command 3
		disgo.CreateGlobalApplicationCommand{
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
