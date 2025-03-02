# Manage a Discord Bot's Application Commands

[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge&logo=appveyor&logo=appveyor)](https://pkg.go.dev/github.com/switchupcb/disgoform)
[![License](https://img.shields.io/github/license/switchupcb/disgoform.svg?style=for-the-badge)](https://github.com/switchupcb/disgoform/blob/main/LICENSE)

Use `disgoform` to stop wasting development time sending application command updates to the Discord API.

## What is Disgoform?
Disgoform is a tool used to manage the application commands of your Discord Bot.

You run `disgoform` as a program which imports a Go module to synchronize a Discord Bot's declared application commands with the Discord API.

```
go get github.com/switchupcb/disgoform@v0.10.0
```

## Table of Contents

| Topic                                                      | Categories                                                                                                                                        |
| :--------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| [How do you use Disgoform?](#how-do-you-use-disgoform)     | [Define Client](#1-define-your-client), [Declare commands](#2-define-your-application-commands), [Sync](#3-synchronize-your-application-commands) |
| [What else can Disgoform do?](#what-else-can-disgoform-do) | [Reverse Sync](#what-else-can-disgoform-do)                                                                                                       |

## How do you use Disgoform?

View the [main example](_example\main.go) for the example `.go` file.

_You can download the Go programming language [here](https://go.dev/learn/)._

### 1. Define your client.

Disgoform uses the [Discord HTTP REST API](https://github.com/switchupcb/disgo/blob/v10/_contribution/concepts/REQUESTS.md) to update application commands from a Go program.

```go
bot := &disgo.Client{
    ApplicationID:  "APPID",
    Authentication: disgo.BotToken("TOKEN"), // or BearerToken("TOKEN")
    Config:         disgo.DefaultConfig(),
}
```

### 2. Define your application commands.

Read the [Discord API Documentation](https://discord.com/developers/docs/interactions/application-commands#application-commands) for more information about Application Commands.

_TIP: Use a `Go: Fill struct` macro (e.g., `Ctrl .` on VSCode) to declare each command faster._

```go
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

    // Command...
    {
        NameLocalizations:        &map[string]string{},
        Description:              nil,
        DescriptionLocalizations: &map[string]string{},
        DefaultMemberPermissions: nil,
        Type:                     nil,
        NSFW:                     nil,
        Name:                     "",
        Options:                  []*disgo.ApplicationCommandOption{},
        IntegrationTypes:         []disgo.Flag{},
        Contexts:                 []disgo.Flag{},
    },
}
```

Define Guild Application Commands using the same format.

```go
disgoform.GuildApplicationCommands = []disgo.CreateGuildApplicationCommand{
    // Command...
    {
        NameLocalizations:        &map[string]string{},
	Description:              nil,
	DescriptionLocalizations: &map[string]string{},
	DefaultMemberPermissions: nil,
	Type:                     nil,
	NSFW:                     nil,
	GuildID:                  "",
	Name:                     "",
	Options:                  []*disgo.ApplicationCommandOption{},
    },
}
```

_NOTE: The commands in this example are sourced from [Disgo examples](https://github.com/switchupcb/disgo/tree/v10/_examples/command)._

### 3. Synchronize your application commands.
Synchronize the Discord Bot's defined application commands with the Discord API.

```go
// Use disgoform.Sync to synchronize Global and Guild application commands.
//
// Use disgoform.SyncGlobalApplicationCommands to only synchronize global application commands.
// Use disgoform.SyncGuildApplicationCommands to only synchronize guild application commands.
if err := disgoform.Sync(bot); err != nil {
    log.Printf("can't synchronize application commands with Discord: %v", err)
}
```

Use `go build -o disgoform` to build the executable binary, then run `disgoform` from the command line.

```
> disgoform
Synchronizing Global Application Commands...
Synchronized Global Application Commands.
Synchronizing Guild Application Commands...
Synchronized Guild Application Commands.
```

## What else can Disgoform do?
You can also generate a `disgoform` `config.go` file using `disgoform.SyncConfig`.

**Here is an example.**

```go
bot := &disgo.Client{
    ApplicationID:  "APPID",
    Authentication: disgo.BotToken("TOKEN"), // or BearerToken("TOKEN")
    Config:         disgo.DefaultConfig(),
}

guildIDs := []string{"...", "...", "..."}

// Use disgoform.SyncConfig to output a synchronization file from the Discord Bot's current state.
output, err := disgoform.SyncConfig(bot, guildIDs)
if err != nil {
    log.Printf("can't output Discord application command configuration file: %v", err)

    return
}

fmt.Println(output)
```

Use `go build -o config` to build the executable binary, then run it from the command line and redirect the output to a file.

```
config > config.go
```

_You can implement this unimplemented feature by [contributing](/_contribution/CONTRIBUTING.md)._
