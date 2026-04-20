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
