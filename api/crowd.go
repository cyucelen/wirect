package api

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/cyucelen/wirect/model"
	"github.com/labstack/echo"
)

type Option func(*CrowdAPI)

type CrowdDatabase interface {
	PacketDatabase
}

type CrowdAPI struct {
	DB                CrowdDatabase
	Interval          time.Duration
	intervalInSeconds int64
	clock             clock.Clock
}

const defaultCalculationInterval = 5 * time.Minute

func CreateCrowdAPI(db CrowdDatabase, options ...Option) *CrowdAPI {
	crowdAPI := &CrowdAPI{DB: db, Interval: defaultCalculationInterval, clock: clock.New()}

	for i := range options {
		options[i](crowdAPI)
	}

	crowdAPI.intervalInSeconds = int64(crowdAPI.Interval / time.Second)

	return crowdAPI
}

func (c *CrowdAPI) GetCrowd(ctx echo.Context) error {
	now := c.clock.Now()
	snifferMAC, _ := url.QueryUnescape(ctx.Param("snifferMAC"))

	params := ctx.QueryParams()
	from, fromExists := params["from"]
	until, untilExists := params["until"]
	forEvery, forExists := params["for"]

	if fromExists && untilExists && forExists {
		c.getCrowdBetweenDates(ctx, snifferMAC, from[0], until[0], forEvery[0])
	}

	count := c.DB.GetUniqueMACCountBySnifferBetweenDates(snifferMAC, now.Unix()-c.intervalInSeconds, now.Unix())
	crowd := []model.Crowd{{Count: count, Time: now}}
	ctx.JSON(http.StatusOK, crowd)

	return nil
}

func (c *CrowdAPI) getCrowdBetweenDates(ctx echo.Context, snifferMAC, from, until, forEvery string) error {
	fromUnix, _ := strconv.ParseInt(from, 10, 64)
	untilUnix, _ := strconv.ParseInt(until, 10, 64)
	forEverySeconds, _ := strconv.ParseInt(forEvery, 10, 64)

	crowd := []model.Crowd{}

	for t := fromUnix; t < untilUnix; t += forEverySeconds {
		crowd = append(crowd, c.getCrowd(snifferMAC, t))
	}
	crowd = append(crowd, c.getCrowd(snifferMAC, untilUnix))
	ctx.JSON(http.StatusOK, crowd)

	return nil
}

func (c *CrowdAPI) getCrowd(snifferMAC string, when int64) model.Crowd {
	count := c.DB.GetUniqueMACCountBySnifferBetweenDates(snifferMAC, when-c.intervalInSeconds, when)
	return model.Crowd{
		Count: count,
		Time:  time.Unix(when, 0),
	}
}

func SetCrowdCalculationInterval(interval time.Duration) Option {
	return func(crowdAPI *CrowdAPI) {
		crowdAPI.Interval = interval
	}
}

func SetCrowdClock(clock clock.Clock) Option {
	return func(crowdAPI *CrowdAPI) {
		crowdAPI.clock = clock
	}
}
