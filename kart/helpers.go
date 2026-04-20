/*
	Various helpers used across protocols
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
	"fmt"
	"encoding/binary"
)

// Converts maptitle bytes and iszone into go string
func GetMapTitle(maptitle []byte, iszone bool) string {
	var zone string

	if iszone {
		zone = " Zone"
	}

	return fmt.Sprintf("%s%s", ParseNullTerminatedString(maptitle), zone)
}

// Calculates checksum that kart would expect for packet
func GetPacketChecksum(packet []byte) uint32 {
	var ret uint32 = 0x1234567

	for i, val := range packet {
		ret += uint32(val)*uint32(i+1)
	}

	return ret
}

// Checks if checksum from packet matches actual checksum
// NOTE - assumes checksum is always uint32 and is always at the beginning of packet
func VerifyChecksum(packet []byte) bool {
	var checksum uint32

	binary.Decode(packet, binary.LittleEndian, &checksum)

	return checksum == GetPacketChecksum(packet[4:])
}

// Calculates and inserts checksum into packet, use before sending it to kart
// NOTE - assumes checksum is always uint32 and is always at the beginning of packet
func AddChecksum(packet []byte) {
	checksum := GetPacketChecksum(packet[4:])
	binary.Encode(packet, binary.LittleEndian, checksum)
}

// Raw packets return null-terminated strings, but they may also have non-cleared stuff
// after null terminator (for example, "Player 1\0something"), we want to grab only
// bytes that come before null terminator
func ParseNullTerminatedString(cstring []byte) string {
	for i, v := range cstring {
		if v == 0 {
			return string(cstring[:i])
		}
	}

	return string(cstring)
}
