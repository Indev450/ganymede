package kart

import (
	"time"
	"sync"
	"net"
	"fmt"
	"errors"
	"os"
)

// TODO - there is "KartConnection", and there is net.UDPConn, which is called connection too
// this might be a bit confusing...
type KartConnection struct {
	proto KartProtocol
	info KartServerInfo
	last_update_time time.Time
	mu sync.Mutex
}

// Returns true if time since last server info update is too long
func (connection *KartConnection) IsInfoExpired() bool {
	return time.Since(connection.last_update_time) > 10*time.Second
}

// If server info is valid returns copy of it and true, otherwise returns empty info and false
func (connection *KartConnection) GetServerInfo() (info_copy KartServerInfo, valid bool) {
	connection.mu.Lock()
	defer connection.mu.Unlock()

	// Too long time passed since last update, server might be down/unreachable
	if connection.IsInfoExpired() {
		return KartServerInfo {}, false
	}

	return connection.info.Copy(), true
}

func (connection *KartConnection) updateServerInfo(packet []byte) {
	connection.mu.Lock()
	defer connection.mu.Unlock()

	if connection.proto.UpdateServerInfo(packet, &connection.info) {
		connection.last_update_time = time.Now()
	}
}

func serverInfoThread(address string, conn *net.UDPConn, connection *KartConnection) {
	buffer := make([]byte, 65535)

	var last_ask_info time.Time

	for {
		// Ask server for info every 5 seconds
		if time.Since(last_ask_info) > time.Second*5 {
			// Update timestamp
			last_ask_info = time.Now()

			packet := connection.proto.AskServerInfo()
			AddChecksum(packet)

			_, err := conn.Write(packet)

			// Server may be offline
			if err != nil {
				fmt.Printf("Error sending to %s: %v\n", address, err)
			}
		}

		conn.SetReadDeadline(time.Now().Add(time.Second/10))

		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			// Nothing from server yet but its fine
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			}

			// Server may be offline
			fmt.Printf("Error reading from %s: %v\n", address, err)
			continue
		}

		// No errors and we got actual packet, "lets go check it aeout" (c)
		connection.updateServerInfo(buffer[:n])
	}
}

// Starts a connection thread and returns reference to it
func StartKartConnection(address string, proto KartProtocol) *KartConnection {
	connection := &KartConnection {
		proto: proto,
	}

	addr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		fmt.Printf("Error resolving %s: %v\n", address, err)
		return nil
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Printf("Error contacting %s: %v\n", address, err)
		return nil
	}

	go serverInfoThread(address, conn, connection)

	return connection
}
