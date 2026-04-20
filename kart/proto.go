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
