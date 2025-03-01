package tests

import (
	"os"
	"testing"

	"github.com/switchupcb/disgo"
	"github.com/switchupcb/disgoform"
)

// TestSync tests SyncGlobalApplicationCommands() functionality.
func TestSync(t *testing.T) {
	bot := &disgo.Client{
		ApplicationID:  os.Getenv("APPID"),
		Authentication: disgo.BotToken(os.Getenv("TOKEN")),
		Config:         disgo.DefaultConfig(),
	}

	// global defined command reset
	if err := disgoform.Sync(bot); err != nil {
		t.Fatalf("reset: %v", err)
	}

	// global defined command empty name
	disgoform.GlobalApplicationCommands = append(disgoform.GlobalApplicationCommands, disgo.CreateGlobalApplicationCommand{})
	if err := disgoform.Sync(bot); err == nil {
		t.Fatalf("expected error while syncing application command with empty name")
	}

	// global defined command two names same
	disgoform.GlobalApplicationCommands = []disgo.CreateGlobalApplicationCommand{
		{
			Name: "test",
		},
		{
			Name: "test",
		},
	}

	if err := disgoform.Sync(bot); err == nil {
		t.Fatalf("expected error while syncing application command with duplicate name")
	}

	// add global defined command from no state
	disgoform.GlobalApplicationCommands = []disgo.CreateGlobalApplicationCommand{
		{
			Name:        "main",
			Description: disgo.Pointer("A basic command."),
			Type:        disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
		},
	}

	if err := disgoform.Sync(bot); err != nil {
		t.Fatalf("add command: %v", err)
	}

	getGlobalApplicatonCommands := &disgo.GetGlobalApplicationCommands{}
	currentCommands, err := getGlobalApplicatonCommands.Send(bot)
	if err != nil {
		t.Fatalf("add command: confirmation: %v", err)
	}

	if len(currentCommands) != 1 {
		t.Fatal("add command: confirmation: amount of global application commands is not 1", err)
	}

	// global defined command update to 1 command, add 1 command
	disgoform.GlobalApplicationCommands = []disgo.CreateGlobalApplicationCommand{
		{
			Name:        "main",
			Description: disgo.Pointer("A basic command test."),
			Type:        disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
		},
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
		t.Fatalf("add command and update command: %v", err)
	}

	currentCommands, err = getGlobalApplicatonCommands.Send(bot)
	if err != nil {
		t.Fatalf("add command and update command: confirmation: %v", err)
	}

	if len(currentCommands) != 2 {
		t.Fatal("add command and update command: confirmation: amount of global application commands is not 2", err)
	}

	// global defined command delete all
	disgoform.GlobalApplicationCommands = []disgo.CreateGlobalApplicationCommand{}
	if err := disgoform.Sync(bot); err != nil {
		t.Fatalf("reset: %v", err)
	}

	currentCommands, err = getGlobalApplicatonCommands.Send(bot)
	if err != nil {
		t.Fatalf("delete all commands: confirmation: %v", err)
	}

	if len(currentCommands) != 0 {
		t.Fatal("delete all commands amount of global application commands is not 0", err)
	}
}
