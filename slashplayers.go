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
