package kart

import (
	"strings"
	"fmt"
)

// Simplified player info entry
type KartPlayer struct {
	Name string
	Spectator bool
}

// Current information that is known about server
type KartServerInfo struct {
	Players []KartPlayer
	Maptitle string
	Gametype string
}

func (info *KartServerInfo) Copy() KartServerInfo {
	ret := *info
	copy(ret.Players, info.Players)
	return ret
}

// Return verb for strings like
// "Watching n players racing"
// "Watching n players battling"
// "Watching n players playing some gamemode"
// etc
func (info *KartServerInfo) GetGametypeVerb() string {
	gametype := strings.ToLower(info.Gametype)

	if gametype == "race" {
		return "racing"
	}

	if gametype == "battle" {
		return "battling"
	}

	return fmt.Sprintf("playing %s", gametype)
}
