package disgoform

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"

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
		return fmt.Errorf("Sync: %w", err)
	}

	log.Println("Synchronized Global Application Commands.")

	log.Println("Synchronizing Guild Application Commands...")

	if err := SyncGuildApplicationCommands(bot); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	log.Println("Synchronized Guild Application Commands.")

	return nil
}

// SyncGlobalApplicationCommands synchronizes Global application commands.
func SyncGlobalApplicationCommands(bot *disgo.Client) error {
	// parse the defined command list into a map of names to application commands.
	definedCommandMap := make(map[string]disgo.CreateGlobalApplicationCommand, len(GlobalApplicationCommands))

	for _, definedCommand := range GlobalApplicationCommands {
		if definedCommand.Name == "" {
			return errors.New("SyncGlobalApplicationCommands: cannot define application command with empty name")
		}

		if _, ok := definedCommandMap[definedCommand.Name]; ok {
			return fmt.Errorf("SyncGlobalApplicationCommands: more than one command exists with name %q", definedCommand.Name)
		}

		definedCommandMap[definedCommand.Name] = definedCommand
	}

	// get the bot's current Global Application Command State.
	getGlobalApplicatonCommands := &disgo.GetGlobalApplicationCommands{
		WithLocalizations: disgo.Pointer(true),
	}

	currentCommands, err := getGlobalApplicatonCommands.Send(bot)
	if err != nil {
		return fmt.Errorf("SyncGlobalApplicationCommands: %w", err)
	}

	// parse the current command list into a map of names to application commands.
	currentCommandMap := make(map[string]disgo.CreateGlobalApplicationCommand, len(currentCommands))
	currentCommandIDMap := make(map[string]string, len(currentCommands))

	for _, currentCommand := range currentCommands {
		currentCommandIDMap[currentCommand.Name] = currentCommand.ID
		currentCommandMap[currentCommand.Name] = disgo.CreateGlobalApplicationCommand{
			NameLocalizations:        currentCommand.NameLocalizations,
			Description:              &currentCommand.Description,
			DescriptionLocalizations: currentCommand.DescriptionLocalizations,
			DefaultMemberPermissions: nil,
			Type:                     currentCommand.Type,
			NSFW:                     currentCommand.NSFW,
			Name:                     currentCommand.Name,
			Options:                  currentCommand.Options,
			IntegrationTypes:         currentCommand.IntegrationTypes,
			Contexts:                 nil,
		}

		if currentCommand.DefaultMemberPermissions != nil {
			copied := currentCommandMap[currentCommand.Name]
			copied.DefaultMemberPermissions = &currentCommand.DefaultMemberPermissions
			currentCommandMap[currentCommand.Name] = copied
		}

		if currentCommand.Contexts != nil {
			copied := currentCommandMap[currentCommand.Name]
			copied.Contexts = *currentCommand.Contexts
			currentCommandMap[currentCommand.Name] = copied
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
					CommandID:                currentCommandIDMap[definedCommand.Name],
					Options:                  definedCommand.Options,
				}

				if _, err := request.Send(bot); err != nil {
					return fmt.Errorf("SyncGlobalApplicationCommands: cannot update current application command %q: %w", name, err)
				}
			}

			delete(currentCommandMap, definedCommand.Name)
			delete(currentCommandIDMap, definedCommand.Name)
			disgo.Logger.Info().Msgf("SyncGlobalApplicationCommands: global application command updated: %q", definedCommand.Name)

			continue
		}

		// definedCommand name does not exist on Discord, so create it.
		if _, err := definedCommand.Send(bot); err != nil {
			return fmt.Errorf("SyncGlobalApplicationCommands: cannot create defined application command %q: %w", name, err)
		}

		disgo.Logger.Info().Msgf("SyncGlobalApplicationCommands: global application command created: %q", definedCommand.Name)
	}

	// delete existing current application commands that aren't defined.
	for _, currentCommand := range currentCommandMap {
		request := &disgo.DeleteGlobalApplicationCommand{
			CommandID: currentCommandIDMap[currentCommand.Name],
		}

		if err := request.Send(bot); err != nil {
			return fmt.Errorf("SyncGlobalApplicationCommands: cannot delete current application command %q: %w", currentCommand.Name, err)
		}

		disgo.Logger.Info().Msgf("SyncGlobalApplicationCommands: global application command deleted: %q", currentCommand.Name)
	}

	return nil
}

// SyncGuildApplicationCommands synchronizes Guild application commands.
//
// WARNING: This function connects and disconnects from the Discord Gateway.
func SyncGuildApplicationCommands(bot *disgo.Client) error {
	// lock represents a lock used to confirm the synchronization is run once.
	var lock sync.Mutex

	// run tracks whether a guild application command synchronization operation is running.
	run := false

	// parse the defined guild command list into a map of GuildIDs to a map of names to guild application commands.
	definedCommandGuildIDMap := make(map[string]map[string]disgo.CreateGuildApplicationCommand)

	for _, definedCommand := range GuildApplicationCommands {
		if definedCommand.GuildID == "" {
			return fmt.Errorf("SyncGuildApplicationCommands: cannot define guild application command with name %q using empty guild id", definedCommand.Name)
		}

		if _, ok := definedCommandGuildIDMap[definedCommand.GuildID]; !ok {
			definedCommandGuildIDMap[definedCommand.GuildID] = make(map[string]disgo.CreateGuildApplicationCommand)
		}

		if definedCommand.Name == "" {
			return fmt.Errorf("SyncGuildApplicationCommands: cannot define guild application command for guild %q using empty name", definedCommand.GuildID)
		}

		if _, ok := definedCommandGuildIDMap[definedCommand.GuildID][definedCommand.Name]; ok {
			return fmt.Errorf("SyncGuildApplicationCommands: more than one command exists with name %q for guild %q", definedCommand.Name, definedCommand.GuildID)
		}

		definedCommandGuildIDMap[definedCommand.GuildID][definedCommand.Name] = definedCommand
	}

	// Connect to the Discord Gateway to receive a ready event which contains all of the guilds the bot is in.
	// https://discord.com/developers/docs/events/gateway-events#ready
	if bot.Handlers == nil {
		bot.Handlers = new(disgo.Handlers)
	}

	if bot.Sessions == nil {
		bot.Sessions = disgo.NewSessionManager()
	}

	// s represents a Session used to connect to the Discord Gateway.
	s := disgo.NewSession()

	// err represents an error used to return any errors experienced during synchronization.
	var err error

	if e := bot.Handle(disgo.FlagGatewayEventNameReady, func(r *disgo.Ready) {
		lock.Lock()
		if run {
			lock.Unlock()

			return
		}

		defer func() {
			if disconnectErr := s.Disconnect(); disconnectErr != nil {
				disgo.Logger.Error().Err(disconnectErr).Msg("SyncGuildApplicationCommands: disconnection")
			}

			run = true
			lock.Unlock()
		}()

		for _, guild := range r.Guilds {
			if guild == nil {
				err = errors.New("SyncGuildApplicationCommands: impossible")

				return
			}

			// get the bot's current Guild Application Command State.
			getGuildApplicatonCommands := &disgo.GetGuildApplicationCommands{
				WithLocalizations: disgo.Pointer(true),
				GuildID:           guild.ID,
			}

			currentCommands, e2 := getGuildApplicatonCommands.Send(bot)
			if e2 != nil {
				err = e2

				return
			}

			// parse the current guild command list into a map of names to application commands.
			currentCommandMap := make(map[string]disgo.CreateGuildApplicationCommand, len(currentCommands))
			currentCommandIDMap := make(map[string]string, len(currentCommands))
			for _, currentCommand := range currentCommands {
				currentCommandIDMap[currentCommand.Name] = currentCommand.ID
				currentCommandMap[currentCommand.Name] = disgo.CreateGuildApplicationCommand{
					NameLocalizations:        currentCommand.NameLocalizations,
					Description:              &currentCommand.Description,
					DescriptionLocalizations: currentCommand.DescriptionLocalizations,
					DefaultMemberPermissions: &currentCommand.DefaultMemberPermissions,
					Type:                     currentCommand.Type,
					NSFW:                     currentCommand.NSFW,
					GuildID:                  *currentCommand.GuildID,
					Name:                     currentCommand.Name,
					Options:                  currentCommand.Options,
				}
			}

			// sync the bot's Guild Application Command State.
			for name, definedCommand := range definedCommandGuildIDMap[guild.ID] {
				// definedCommand name exists on Discord
				if currentCommand, ok := currentCommandMap[name]; ok {
					// but is not equal to Discord's version, so update it.
					if !Equal(definedCommand, currentCommand) {
						request := &disgo.EditGuildApplicationCommand{
							Name:                     &definedCommand.Name,
							NameLocalizations:        definedCommand.NameLocalizations,
							Description:              definedCommand.Description,
							DescriptionLocalizations: definedCommand.DescriptionLocalizations,
							DefaultMemberPermissions: definedCommand.DefaultMemberPermissions,
							NSFW:                     definedCommand.NSFW,
							GuildID:                  definedCommand.GuildID,
							CommandID:                currentCommandIDMap[definedCommand.Name],
							Options:                  definedCommand.Options,
						}

						if _, e3 := request.Send(bot); e3 != nil {
							err = fmt.Errorf("SyncGuildApplicationCommands: cannot update current guild %q application command %q: %w", definedCommand.GuildID, name, e3)

							return
						}
					}

					delete(currentCommandMap, definedCommand.Name)
					delete(currentCommandIDMap, definedCommand.Name)
					disgo.Logger.Info().Msgf("SyncGuildApplicationCommands: SyncGuildApplicationCommands: guild %q application command updated: %q", definedCommand.GuildID, definedCommand.Name)

					continue
				}

				// definedCommand name does not exist on Discord, so create it.
				if _, e3 := definedCommand.Send(bot); e3 != nil {
					err = fmt.Errorf("SyncGuildApplicationCommands: cannot create defined guild %q application command %q: %w", definedCommand.GuildID, name, e3)

					return
				}

				disgo.Logger.Info().Msgf("SyncGuildApplicationCommands: SyncGuildApplicationCommands: guild %q application command created: %q", definedCommand.GuildID, definedCommand.Name)
			}

			// delete existing current guild application commands that aren't defined.
			for _, currentCommand := range currentCommandMap {
				request := &disgo.DeleteGuildApplicationCommand{
					GuildID:   currentCommand.GuildID,
					CommandID: currentCommandIDMap[currentCommand.Name],
				}

				if e3 := request.Send(bot); e3 != nil {
					err = fmt.Errorf("SyncGuildApplicationCommands: cannot delete current guild %q application command %q: %w", currentCommand.GuildID, currentCommand.Name, e3)

					return
				}

				disgo.Logger.Info().Msgf("SyncGuildApplicationCommands: guild %q application command deleted: %q", currentCommand.GuildID, currentCommand.Name)
			}
		} // for each guild
	}); e != nil {
		return fmt.Errorf("SyncGuildApplicationCommands: %w", e)
	}

	if err := s.Connect(bot); err != nil {
		return fmt.Errorf("SyncGuildApplicationCommands: %w", err)
	}

	_, _ = s.Wait()

	if err != nil {
		return fmt.Errorf("SyncGuildApplicationCommands: %w", err)
	}

	return nil
}
