/*
	/players command implementation
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
	"strings"
	"github.com/Indev450/ganymede/kart"
)

func joinNames(builder *strings.Builder, players []*kart.KartPlayer, verb string) {
	if len(players) == 0 {
		fmt.Fprintf(builder, "nobody is %s", verb)
	} else if len(players) == 1 {
		fmt.Fprintf(builder, "%s is %s", players[0].Name, verb)
	} else {
		for i, player := range players[:len(players)-1] {
			fmt.Fprintf(builder, "%s", player.Name)

			if i == len(players) - 2 {
				fmt.Fprintf(builder, " and ")
			} else {
				fmt.Fprintf(builder, ", ")
			}
		}

		fmt.Fprintf(builder, "%s are %s", players[len(players)-1].Name, verb)
	}
}

func doSlashPlayers(info *kart.KartServerInfo, verb string) string {
	var builder strings.Builder

	if len(info.Players) == 0 {
		fmt.Fprintf(&builder, "nobody is %s, map is %s.", verb, info.Maptitle)

		return builder.String()
	}

	var playing, watching []*kart.KartPlayer

	for _, player := range info.Players {
		if player.Spectator {
			watching = append(watching, &player)
		} else {
			playing = append(playing, &player)
		}
	}

	if len(playing) == 0 {
		joinNames(&builder, watching, "watching")
	} else {
		joinNames(&builder, playing, verb)

		if len(watching) > 0 {
			fmt.Fprintf(&builder, ", ")
			joinNames(&builder, watching, "watching")
		}
	}

	fmt.Fprintf(&builder, " at %s.", info.Maptitle)

	return builder.String()
}
