package test

import (
	"sort"

	"github.com/cyucelen/wirect/model"
)

type InMemoryDB struct {
	Packets  []model.Packet
	Sniffers []model.Sniffer
}

func (i *InMemoryDB) CreatePacket(packet *model.Packet) error {
	i.Packets = append(i.Packets, *packet)
	return nil
}

func (i *InMemoryDB) GetPacketsBySniffer(snifferMAC string) []model.Packet {
	filteredPackets := []model.Packet{}

	for _, packet := range i.Packets {
		if packet.SnifferMAC == snifferMAC {
			filteredPackets = append(filteredPackets, packet)
		}
	}

	return filteredPackets
}

func (i *InMemoryDB) GetPacketsBySnifferSince(snifferMAC string, since int64) []model.Packet {
	filteredPackets := []model.Packet{}

	for _, packet := range i.Packets {
		if packet.SnifferMAC == snifferMAC && packet.Timestamp >= since {
			filteredPackets = append(filteredPackets, packet)
		}
	}

	return sortByPacketsTime(filteredPackets)
}

func (i *InMemoryDB) CreateSniffer(sniffer *model.Sniffer) error {
	i.Sniffers = append(i.Sniffers, *sniffer)
	return nil
}

func (i *InMemoryDB) GetSniffers() []model.Sniffer {
	return i.Sniffers
}

func (i *InMemoryDB) UpdateSniffer(sniffer *model.Sniffer) error {
	for index := range i.Sniffers {
		if i.Sniffers[index].MAC == sniffer.MAC {
			i.Sniffers[index] = *sniffer
		}
	}
	return nil
}

func sortByPacketsTime(s []model.Packet) []model.Packet {
	sc := make([]model.Packet, len(s))
	copy(sc, s)

	sort.Slice(sc, func(i int, j int) bool {
		if s[i].Timestamp < s[j].Timestamp {
			return true
		}
		return false
	})

	return sc
}
