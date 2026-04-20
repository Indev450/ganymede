/*
	Server info struct definition
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
