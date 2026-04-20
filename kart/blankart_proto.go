/*
	Blankart protocol implementation
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
	"io"
	"fmt"
	"bytes"
	"encoding/binary"
)

type BlankartProtocol struct {}

func blankartReadServerInfo(packet io.Reader, info *KartServerInfo) bool {
	var serverinfo struct {
		_ uint8 // _255
		_ uint8 // packetversion
		_ [16]byte // application
		_ uint8 // version
		_ uint8 // subversion
		_ [4]byte // commit
		_ uint8 // numberofplayer
		_ uint8 // maxplayer
		_ uint8 // refusereason
		Gametypename [24]byte
		_ uint8 // modifiedgame
		_ uint8 // cheatsenabled
		_ uint8 // kartvars
		_ uint8 // fileneedednum
		_ uint32 // time
		_ uint32 // leveltime
		_ [32]byte // servername
		Maptitle [33]byte
		_ uint64 // maphash
		_ uint8 // actnum
		Iszone uint8
		// Don't care whats left :p
	}

	err := binary.Read(packet, binary.LittleEndian, &serverinfo)

	if err != nil {
		fmt.Println("Error reading serverinfo:", err)
		return false
	}

	info.Maptitle = GetMapTitle(serverinfo.Maptitle[:], serverinfo.Iszone != 0)
	info.Gametype = ParseNullTerminatedString(serverinfo.Gametypename[:])

	return true
}

func blankartReadPlayerInfo(packet io.Reader, info *KartServerInfo) bool {
	var plrinfo [32]struct {
		Node uint8
		Name [21+1]byte
		Team uint8
		Skin uint16
		Data uint8
		Score uint32
		Timeinserver uint16
	}

	err := binary.Read(packet, binary.LittleEndian, &plrinfo)

	if err != nil {
		fmt.Println("Error reading plrinfo:", err)
		return false
	}

	info.Players = info.Players[:0]

	for _, value := range plrinfo {
		if value.Node == 255 {
			continue // Empty slot
		}

		player := KartPlayer {
			Name: ParseNullTerminatedString(value.Name[:]),
			Spectator: value.Team != 0,
		}

		info.Players = append(info.Players, player)
	}

	return true
}

func (p BlankartProtocol) UpdateServerInfo(packet []byte, info *KartServerInfo) bool {
	var doomdata struct {
		Checksum uint32
		_ uint8
		_ uint8
		Packettype uint8
		_ uint8
	}

	reader := bytes.NewReader(packet)

	err := binary.Read(reader, binary.LittleEndian, &doomdata)

	if err != nil {
		fmt.Println("Error reading doomdata:", err)
		return false
	}

	if !VerifyChecksum(packet) {
		fmt.Println("Error reading packet: checksum mismatch")
		return false
	}

	if doomdata.Packettype == 13 { // PT_SERVERINFO
		return blankartReadServerInfo(reader, info)
	} else if doomdata.Packettype == 14 { // PT_PLAYERINFO
		return blankartReadPlayerInfo(reader, info)
	}

	return false
}

func (p BlankartProtocol) AskServerInfo() []byte {
	data := []byte{
		0, 0, 0, 0, // checksum
		0, 0, // padding
		12, // packettype, PT_ASKINFO
		0, // padding
	}

	return data
}
