package database

import (
	"time"

	"github.com/cyucelen/wirect/model"
)

func (g *GormDatabase) CreatePacket(packet *model.Packet) error {
	return g.DB.Create(packet).Error
}

func (g *GormDatabase) GetPacketsBySniffer(snifferMAC string) []model.Packet {
	var packets []model.Packet
	g.DB.Order("timestamp asc").Where("sniffer_mac = ?", snifferMAC).Find(&packets)
	return packets
}

func (g *GormDatabase) GetPacketsBySnifferSince(snifferMAC string, since time.Time) []model.Packet {
	var packets []model.Packet
	g.DB.Order("timestamp asc").Where("sniffer_mac = ? AND timestamp >= ?", snifferMAC, since).Find(&packets)
	return packets
}
