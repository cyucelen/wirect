package server

import (
	"github.com/benbjohnson/clock"
	"github.com/cyucelen/wirect/api"
	"github.com/labstack/echo"
)

type Database interface {
	api.PacketDatabase
	api.SnifferDatabase
}

var tick = clock.New()

const snifferMACParam = "/:snifferMAC"
const sniffersEndpoint = "/sniffers"
const packetsEndpoint = sniffersEndpoint + snifferMACParam + "/packets"
const packetsCollectionEndpoint = sniffersEndpoint + snifferMACParam + "/packets-collection"
const updateSnifferEndpoint = sniffersEndpoint + snifferMACParam
const crowdEndpoint = sniffersEndpoint + snifferMACParam + "/stats/crowd"
const timeEndpoint = "/time"

func Create(db Database) *echo.Echo {
	e := echo.New()
	createPacketEndpoints(e, db)
	createSnifferEndpoints(e, db)
	createCrowdEndpoints(e, db)
	createTimeEndpoint(e)

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
	e.PUT(updateSnifferEndpoint, snifferAPI.UpdateSniffer)
}

func createCrowdEndpoints(e *echo.Echo, db Database) {
	crowdAPI := api.CreateCrowdAPI(db, api.SetCrowdClock(tick))
	e.GET(crowdEndpoint, crowdAPI.GetCrowd)
}

func createTimeEndpoint(e *echo.Echo) {
	timeAPI := api.TimeAPI{Clock: tick}
	e.GET(timeEndpoint, timeAPI.GetTime)
}
