package api

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/labstack/echo"

	"github.com/cyucelen/wirect/model"
)

type PacketDatabase interface {
	CreatePacket(packet *model.Packet) error
	GetPacketsBySniffer(snifferMAC string) []model.Packet
	GetPacketsBySnifferSince(snifferMAC string, since int64) []model.Packet
	GetPacketsBySnifferBetweenDates(snifferMAC string, from, until int64) []model.Packet
	GetUniqueMACCountBySnifferBetweenDates(snifferMAC string, from, until int64) int
}

type PacketAPI struct {
	DB PacketDatabase
}

func (p *PacketAPI) CreatePacket(ctx echo.Context) error {
	var snifferPacket model.SnifferPacket
	if err := ctx.Bind(&snifferPacket); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return err
	}

	if !isSnifferPacketValid(snifferPacket) {
		ctx.JSON(http.StatusBadRequest, nil)
		return errors.New("")
	}

	snifferMAC, err := getSnifferMAC(ctx)
	if err != nil {
		return err
	}
	packet := toPacket(&snifferPacket, snifferMAC)

	if err := p.DB.CreatePacket(packet); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return err
	}

	ctx.JSON(http.StatusCreated, snifferPacket)
	return nil
}

func (p *PacketAPI) CreatePackets(ctx echo.Context) error {
	var snifferPackets []model.SnifferPacket

	if err := ctx.Bind(&snifferPackets); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return err
	}

	validSnifferPackets := filterValidSnifferPackets(snifferPackets)

	if len(validSnifferPackets) == 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return errors.New("")
	}

	snifferMAC, err := getSnifferMAC(ctx)
	if err != nil {
		return err
	}

	for _, validSnifferPacket := range validSnifferPackets {
		packet := toPacket(&validSnifferPacket, snifferMAC)

		if err := p.DB.CreatePacket(packet); err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return err
		}
	}

	ctx.JSON(http.StatusCreated, snifferPackets)
	return nil
}

func getSnifferMAC(ctx echo.Context) (string, error) {
	snifferMAC, err := url.QueryUnescape(ctx.Param("snifferMAC"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return "", err
	}

	if snifferMAC == "" {
		ctx.JSON(http.StatusNotFound, nil)
		return "", errors.New("")
	}

	return snifferMAC, nil
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
	return snifferPacket.MAC != "" && snifferPacket.Timestamp != 0
}

func toPacket(snifferPacket *model.SnifferPacket, snifferMAC string) *model.Packet {
	return &model.Packet{
		MAC:        snifferPacket.MAC,
		Timestamp:  snifferPacket.Timestamp,
		RSSI:       snifferPacket.RSSI,
		SnifferMAC: snifferMAC,
	}
}
