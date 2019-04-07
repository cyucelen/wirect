package model

import (
	"time"
)

// SnifferPacket holds information about a Packet which was sent by a Sniffer
type SnifferPacket struct {
	MAC        string    `json:"MAC"`
	Timestamp  time.Time `json:"timestamp"`
	RSSI       float64   `json:"RSSI"`
	SnifferMAC string    `json:"snifferMAC"`
}

// Packet represents database schema of a collected data
type Packet struct {
	ID         uint `gorm:"AUTO_INCREMENT"`
	MAC        string
	Timestamp  time.Time
	RSSI       float64
	Sniffer    Sniffer `gorm:"foreignkey:SnifferMAC"`
	SnifferMAC string
}
