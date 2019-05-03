package model

// SnifferPacket holds information about a Packet which was sent by a Sniffer
type SnifferPacket struct {
	MAC       string  `json:"MAC"`
	Timestamp int64   `json:"timestamp"`
	RSSI      float64 `json:"RSSI"`
}

// Packet represents database schema of a collected data
type Packet struct {
	ID         uint `gorm:"AUTO_INCREMENT"`
	MAC        string
	Timestamp  int64
	RSSI       float64
	Sniffer    Sniffer `gorm:"foreignkey:SnifferMAC"`
	SnifferMAC string
}
