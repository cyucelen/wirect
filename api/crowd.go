package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cyucelen/wirect/model"
	"github.com/labstack/echo"
)

type Option func(*CrowdAPI)

type CrowdDatabase interface {
	PacketDatabase
}

type CrowdAPI struct {
	DB       CrowdDatabase
	Interval time.Duration
}

const defaultCalculationInterval = 5 * time.Minute

func CreateCrowdAPI(db CrowdDatabase, options ...Option) *CrowdAPI {
	crowdAPI := &CrowdAPI{DB: db, Interval: defaultCalculationInterval}

	for i := range options {
		options[i](crowdAPI)
	}

	return crowdAPI
}

func (c *CrowdAPI) GetCrowd(ctx echo.Context) error {
	snifferMAC := ctx.QueryParam("sniffer")
	since, _ := strconv.ParseInt(ctx.QueryParam("since"), 10, 64)

	intervalInSeconds := int64(c.Interval / time.Second)

	packets := c.DB.GetPacketsBySnifferSince(snifferMAC, since-intervalInSeconds)

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

func SetCrowdCalculationInterval(interval time.Duration) Option {
	return func(crowdAPI *CrowdAPI) {
		crowdAPI.Interval = interval
	}
}
