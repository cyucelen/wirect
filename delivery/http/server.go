package server

import (
	"github.com/cyucelen/wirect/api"
	"github.com/labstack/echo"
)

type Database interface {
	api.PacketDatabase
	api.SnifferDatabase
}

const packetsEndpoint = "/packets"
const packetsCollectionEndpoint = "/packets-collection"
const sniffersEndpoint = "/sniffers"
const crowdEndpoint = "/crowd"

func Create(db Database) *echo.Echo {
	e := echo.New()
	createPacketEndpoints(e, db)
	createSnifferEndpoints(e, db)
	createCrowdEndpoints(e, db)

	return e
}

func createPacketEndpoints(e *echo.Echo, db Database) {
	packetAPI := api.PacketAPI{DB: db}
	e.POST(packetsEndpoint, packetAPI.CreatePacket)
	e.POST(packetsCollectionEndpoint, packetAPI.CreatePackets)
}

func createSnifferEndpoints(e *echo.Echo, db Database) {
	snifferAPI := api.SnifferAPI{DB: db}
	e.GET(sniffersEndpoint, snifferAPI.GetSniffers)
	e.POST(sniffersEndpoint, snifferAPI.CreateSniffer)
	e.PUT(sniffersEndpoint, snifferAPI.UpdateSniffer)
}

func createCrowdEndpoints(e *echo.Echo, db Database) {
	crowdAPI := &api.CrowdAPI{DB: db}
	e.GET(crowdEndpoint, crowdAPI.GetCrowd)
}
