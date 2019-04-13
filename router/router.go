package router

import (
	"github.com/labstack/echo"
	"gitlab.com/wirect/wirect-server/api"
)

type Database interface {
	api.PacketDatabase
	api.SnifferDatabase
}

const packetEndpoint = "/packet"
const snifferEndpoint = "/sniffer"

func Create(db Database) *echo.Echo {
	e := echo.New()
	createPacketEndpoints(e, db)
	createSnifferEndpoints(e, db)
	return e
}

func createPacketEndpoints(e *echo.Echo, db Database) {
	packetAPI := api.PacketAPI{DB: db}
	e.POST(packetEndpoint, packetAPI.CreatePacket)
}

func createSnifferEndpoints(e *echo.Echo, db Database) {
	snifferAPI := api.SnifferAPI{DB: db}
	e.GET(snifferEndpoint, snifferAPI.GetSniffers)
	e.POST(snifferEndpoint, snifferAPI.CreateSniffer)
	e.PUT(snifferEndpoint, snifferAPI.UpdateSniffer)
}
