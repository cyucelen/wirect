package database

import (
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

func (g *GormDatabase) GetPacketsBySnifferSince(snifferMAC string, since int64) []model.Packet {
	var packets []model.Packet
	g.DB.Order("timestamp asc").Where("sniffer_mac = ? AND timestamp >= ?", snifferMAC, since).Find(&packets)
	return packets
}
