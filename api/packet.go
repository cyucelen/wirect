package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-ffmt/ffmt"
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

func (p *PacketAPI) CreatePacket(ctx echo.Context) error {
	var snifferPacket model.SnifferPacket
	if err := ctx.Bind(&snifferPacket); err != nil {
		ffmt.Puts(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return err
	}

	if !isSnifferPacketValid(snifferPacket) {
		ctx.JSON(http.StatusBadRequest, nil)
		return errors.New("")
	}

	packet := toPacket(&snifferPacket)

	if err := p.DB.CreatePacket(packet); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return err
	}

	ctx.JSON(http.StatusCreated, snifferPacket)
	return nil
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
