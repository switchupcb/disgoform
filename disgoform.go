package disgoform

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/switchupcb/disgo"
)

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

var (
	// Equal returns whether two application commands are equal.
	Equal = reflect.DeepEqual
)

// Sync synchronizes Global and Guild application commands.
func Sync(bot *disgo.Client) error {
	log.Println("Synchronizing Global Application Commands...")
	if err := SyncGlobalApplicationCommands(bot); err != nil {
		return fmt.Errorf("SyncGlobalApplicationCommands: %w", err)
	}
	log.Println("Synchronized Global Application Commands.")

	log.Println("Synchronizing Guild Application Commands...")
	if err := SyncGuildApplicationCommands(bot); err != nil {
		return fmt.Errorf("SyncGuildApplicationCommands: %w", err)
	}
	log.Println("Synchronized Guild Application Commands.")

	return nil
}

// SyncGlobalApplicationCommands synchronizes Global application commands.
func SyncGlobalApplicationCommands(bot *disgo.Client) error {
	// get the bot's current Global Application Command State.
	getGlobalApplicatonCommands := &disgo.GetGlobalApplicationCommands{}
	currentCommands, err := getGlobalApplicatonCommands.Send(bot)
	if err != nil {
		return err
	}

	// parse each command list into a map of names to application commands.
	definedCommandMap := make(map[string]disgo.CreateGlobalApplicationCommand, len(GlobalApplicationCommands))
	for _, definedCommand := range GlobalApplicationCommands {
		if definedCommand.Name == "" {
			return errors.New("cannot define application command with empty name")
		}

		if _, ok := definedCommandMap[definedCommand.Name]; ok {
			return fmt.Errorf("more than one command exists with name %q", definedCommand.Name)
		}

		definedCommandMap[definedCommand.Name] = definedCommand
	}

	currentCommandMap := make(map[string]disgo.CreateGlobalApplicationCommand, len(currentCommands))
	currentCommandIDMap := make(map[string]string, len(currentCommands))
	for _, currentCommand := range currentCommands {
		currentCommandIDMap[currentCommand.Name] = currentCommand.ID
		currentCommandMap[currentCommand.Name] = disgo.CreateGlobalApplicationCommand{
			NameLocalizations:        currentCommand.NameLocalizations,
			Description:              disgo.Pointer(currentCommand.Description),
			DescriptionLocalizations: currentCommand.DescriptionLocalizations,
			DefaultMemberPermissions: nil,
			Type:                     currentCommand.Type,
			NSFW:                     currentCommand.NSFW,
			Name:                     currentCommand.Name,
			Options:                  currentCommand.Options,
			IntegrationTypes:         currentCommand.IntegrationTypes,
			Contexts:                 nil,
		}

		c := currentCommandMap[currentCommand.Name]
		if currentCommand.DefaultMemberPermissions != nil {
			c.DefaultMemberPermissions = &currentCommand.DefaultMemberPermissions
		}

		if currentCommand.Contexts != nil {
			c.Contexts = *currentCommand.Contexts
		}
	}

	// sync the bot's Global Application Command State.
	for name, definedCommand := range definedCommandMap {
		// definedCommand name exists on Discord
		if currentCommand, ok := currentCommandMap[name]; ok {

			// but is not equal to Discord's version, so update it.
			if !Equal(definedCommand, currentCommand) {
				request := &disgo.EditGlobalApplicationCommand{
					Name:                     &definedCommand.Name,
					NameLocalizations:        definedCommand.NameLocalizations,
					Description:              definedCommand.Description,
					DescriptionLocalizations: definedCommand.DescriptionLocalizations,
					DefaultMemberPermissions: definedCommand.DefaultMemberPermissions,
					NSFW:                     definedCommand.NSFW,
					CommandID:                currentCommandIDMap[currentCommand.Name],
					Options:                  definedCommand.Options,
				}

				if _, err := request.Send(bot); err != nil {
					return fmt.Errorf("cannot update current application command %q: %w", name, err)
				}
			}

			delete(currentCommandMap, definedCommand.Name)
			delete(currentCommandIDMap, definedCommand.Name)

			continue
		}

		// definedCommand name does not exist on Discord, so create it.
		if _, err := definedCommand.Send(bot); err != nil {
			return fmt.Errorf("cannot create defined application command %q: %w", name, err)
		}
	}

	// delete existing current application commands that aren't defined.
	for _, currentCommand := range currentCommandMap {
		request := &disgo.DeleteGlobalApplicationCommand{
			CommandID: currentCommandIDMap[currentCommand.Name],
		}

		if err := request.Send(bot); err != nil {
			return fmt.Errorf("cannot delete current command %q: %w", currentCommand.Name, err)
		}
	}

	return nil
}

// SyncGlobalApplicationCommands synchronizes Guild application commands.
func SyncGuildApplicationCommands(bot *disgo.Client) error {
	return nil
}
