/*
	Interface for different protocols for srb2-based games
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

type KartProtocol interface {
	// Should return true if packet updated anything in server info
	UpdateServerInfo(packet []byte, info *KartServerInfo) bool

	// Should format and return askinfo packet
	AskServerInfo() []byte
}

// Choose protocol based by name
func GetProtocol(name string) KartProtocol {
	switch name {
	case "srb2kart-16p":
		return VanillaProtocol {}
	case "blankart":
		return BlankartProtocol {}
	case "ringracers-16p":
		return RingracersProtocol {}
	default:
		return nil
	}
}
