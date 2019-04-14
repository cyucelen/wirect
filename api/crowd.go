package api

import (
	"net/http"
	"time"

	"github.com/cyucelen/wirect/model"
	"github.com/labstack/echo"
)

type CrowdDatabase interface {
	PacketDatabase
}

type CrowdAPI struct {
	DB CrowdDatabase
}

func (c *CrowdAPI) GetCrowd(ctx echo.Context) error {
	snifferMAC := ctx.QueryParam("sniffer")
	since := ctx.QueryParam("since")

	sinceTime, _ := time.Parse(since, time.RFC3339)

	packets := c.DB.GetPacketsBySnifferSince(snifferMAC, sinceTime)
	crowd := model.Crowd{Count: getUniquePersonCount(packets)}

	ctx.JSON(http.StatusOK, crowd)
	return nil
}

func getUniquePersonCount(packets []model.Packet) int {
	uniquePerson := make(map[string]bool)

	for _, packet := range packets {
		uniquePerson[packet.MAC] = true
	}

	return len(uniquePerson)
}
