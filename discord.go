package main

import (
	"fmt"
	"time"
	"github.com/Indev450/ganymede/kart"
	"github.com/bwmarrin/discordgo"
)

// For now, commands don't need any parameters, so they only take current server connection
type CommandFunc = func(config *DiscordBotConfig, connection *kart.KartConnection) string

type Command struct {
	Func CommandFunc
	Description string
}

type DiscordBotConfig struct {
	Token string

	Commands map[string]Command

	// Ignore player with that name
	SeedPlayer string

	// Force gametype in status be this string instead, for example "hangout"
	StatusGametype string
}

func buildBotStatus(config *DiscordBotConfig, connection *kart.KartConnection) (status string, text string) {
	info, ok := connection.GetServerInfo()

	if !ok {
		return "dnd", "until you'll help me"
	}

	if len(info.Players) == 0 {
		return "online", "an empty map"
	}

	numplayers := len(info.Players)

	if len(config.SeedPlayer) > 0 {
		for _, player := range(info.Players) {
			if player.Name == config.SeedPlayer {
				numplayers--
				break
			}
		}
	}

	gametype := info.Gametype

	if len(config.StatusGametype) > 0 {
		gametype = config.StatusGametype
	}

	return "online", fmt.Sprintf("%d players %s", numplayers, gametype)
}

func botStatusUpdateThread(session *discordgo.Session, config *DiscordBotConfig, connection *kart.KartConnection) {
	usd := discordgo.UpdateStatusData {
		Status: "idle",
	}

	err := session.UpdateStatusComplex(usd)

	if err != nil {
		fmt.Println("Failed to update status:", err)
	}

	for {
		time.Sleep(5*time.Second)

		status, text := buildBotStatus(config, connection)

		usd = discordgo.UpdateStatusData {
			Status: status,
			Activities: []*discordgo.Activity{{
				Name: text,
				Type: discordgo.ActivityTypeWatching,
			}},
		}

		err := session.UpdateStatusComplex(usd)

		if err != nil {
			fmt.Println("Failed to update status:", err)
		}
	}
}

func startDiscordSession(config *DiscordBotConfig, connection *kart.KartConnection) (ok bool) {
	session, _ := discordgo.New(config.Token)

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()
		command, ok := config.Commands[data.Name]

		if !ok {
			fmt.Println("Unknown command:", data.Name)
			return
		}

		reply := command.Func(config, connection)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Content: reply,
			},
		})

		if err != nil {
			fmt.Println("Could not respond to command:", err)
		}
	})

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as %s\n", r.User.String())

		var commands []*discordgo.ApplicationCommand

		for name, command := range config.Commands {
			commands = append(commands, &discordgo.ApplicationCommand {
				Name: name,
				Description: command.Description,
			})
		}

		_, err := s.ApplicationCommandBulkOverwrite(r.Application.ID, "", commands)

		if err != nil {
			fmt.Println("Failed to register application commands:", err)
			// TODO - this makes bot kinda useless, probably need to handle this better?
			// Not that its supposed to happen normally...
		}
	})

	err := session.Open()

	if err != nil {
		fmt.Println("Failed to open session:", err)

		return false
	}

	go botStatusUpdateThread(session, config, connection)

	return true
}
