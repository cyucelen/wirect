package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"gitlab.com/wirect/wirect-server/model"
)

type PacketDatabase interface {
	CreatePacket(packet *model.Packet) error
	GetPacketsBySniffer(snifferMAC string) []model.Packet
	GetPacketsBySnifferSince(snifferMAC string, since time.Time) []model.Packet
}

type PacketAPI struct {
	DB PacketDatabase
}

func (p *PacketAPI) CreatePackets(ctx echo.Context) {
	var snifferPackets []model.SnifferPacket

	if err := ctx.Bind(&snifferPackets); err != nil {
		ctx.JSON(http.StatusNotFound, "")
		return
	}

	validSnifferPackets := filterValidSnifferPackets(snifferPackets)

	if len(validSnifferPackets) == 0 {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	for _, validSnifferPacket := range validSnifferPackets {
		packet := toPacket(&validSnifferPacket)

		if err := p.DB.CreatePacket(packet); err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	ctx.JSON(http.StatusCreated, snifferPackets)
}

func filterValidSnifferPackets(snifferPackets []model.SnifferPacket) []model.SnifferPacket {
	validSnifferPackets := []model.SnifferPacket{}
	for _, snifferPacket := range snifferPackets {
		if isSnifferPacketValid(snifferPacket) {
			validSnifferPackets = append(validSnifferPackets, snifferPacket)
		}
	}
	return validSnifferPackets
}

func isSnifferPacketValid(snifferPacket model.SnifferPacket) bool {
	return snifferPacket.MAC != "" && snifferPacket.SnifferMAC != "" && snifferPacket.Timestamp != time.Time{}
}

func toPacket(snifferPacket *model.SnifferPacket) *model.Packet {
	return &model.Packet{
		MAC:        snifferPacket.MAC,
		Timestamp:  snifferPacket.Timestamp,
		RSSI:       snifferPacket.RSSI,
		SnifferMAC: snifferPacket.SnifferMAC,
	}
}
