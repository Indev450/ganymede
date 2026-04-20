/*
	SRB2Kart discord bot, takes in kart server address, discord bot token and other options and runs forever
	Copyright (C) 2026 Indev

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"os"
	"github.com/Indev450/ganymede/kart"
)

func ensureEnv(name string) string {
	value, ok := os.LookupEnv(name)

	if !ok {
		fmt.Println("Missing environment variable:", name)
		os.Exit(1)
	}

	return value
}

func defaultEnv(name, def string) string {
	value, ok := os.LookupEnv(name)

	if !ok { return def }

	return value
}

// :3
func copyleft() {
	fmt.Println("Copyright (C) 2026 Indev")
	fmt.Println("License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>")
	fmt.Println("This is free software: you are free to change and redistribute it.")
	fmt.Println("There is NO WARRANTY, to the extent permitted by law.")
}

func main() {
	// Target server for bot
	address := ensureEnv("SRB2KART_ADDRESS")

	// Which protocol to use?
	protocol_name := defaultEnv("SRB2KART_PROTOCOL", "srb2kart-16p")

	// Override for gametype in status ("Watching n players %s")
	status_gametype := os.Getenv("SRB2KART_STATUS_GAMEMODE")

	// Override for gametype /players ("Player 1, PLayer 2 and Player 3 are %s at Green Hills Zone")
	slashplayers_gametype := os.Getenv("SRB2KART_SLASHPLAYERS_GAMEMODE")

	// Ignore this player for /players and number of players in status
	seedplayer := os.Getenv("SRB2KART_SEEDPLAYER")

	// File to open to fetch server gamemodes for /gamemode command
	// Gamemodes are written one per line
	gamemodefile := os.Getenv("SRB2KART_GAMEMODEFILE")

	// Token
	token := "Bot " + ensureEnv("DISCORD_TOKEN")

	// Prefix for commands, for example DISCORD_CMDPREFIX="blan-" will turn /players and /gamemode into
	// /blan-players and /blan-gamemode, useful if you have multiple bots for different games on server
	cmdprefix := os.Getenv("DISCORD_CMDPREFIX")

	proto := kart.GetProtocol(protocol_name)

	if proto == nil {
		fmt.Println("Unknown protocol:", protocol_name)
		return
	}

	copyleft()

	connection := kart.StartKartConnection(address, proto)

	// Couldn't setup server connection, bot will not be able to do anything!
	if connection == nil {
		return
	}

	config := DiscordBotConfig {
		Token: token,
		Commands: make(map[string]Command),

		SeedPlayer: seedplayer,
		StatusGametype: status_gametype,
	}

	// /players
	config.Commands[cmdprefix + "players"] = Command {
		Func: func(config *DiscordBotConfig, connection *kart.KartConnection) string {
			info, ok := connection.GetServerInfo()

			if !ok {
				return "Server is currently unavailable."
			}

			verb := slashplayers_gametype

			if len(verb) == 0 {
				verb = info.GetGametypeVerb()
			}

			return doSlashPlayers(&info, verb)
		},
		Description: "Get player info",
	}

	// /gamemode
	if len(gamemodefile) > 0 {
		config.Commands[cmdprefix + "gamemode"] = Command {
			Func: func(config *DiscordBotConfig, connection *kart.KartConnection) string {
				return doSlashGamemode(gamemodefile)
			},
			Description: "Get current gamemode",
		}
	}

	// Starting discord session failed, bot will not be able to do anything!
	if !startDiscordSession(&config, connection) {
		return
	}

	// Run until canceled manually
	for {}
}
