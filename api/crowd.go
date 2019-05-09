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

type CrowdParams struct {
	from           int64
	until          int64
	forEverySecond int64
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
	snifferMAC, _ := url.QueryUnescape(ctx.Param("snifferMAC")) // TODO: test error case
	params := c.getCrowdParams(ctx)
	crowd := c.getCrowdBetweenDates(ctx, snifferMAC, params.from, params.until, params.forEverySecond)
	ctx.JSON(http.StatusOK, crowd)
	return nil
}

func (c *CrowdAPI) getCrowdBetweenDates(ctx echo.Context, snifferMAC string, from, until, forEverySeconds int64) []model.Crowd {
	crowd := []model.Crowd{}

	for t := from; t < until; t += forEverySeconds {
		crowd = append(crowd, c.getCrowd(snifferMAC, t))
	}
	crowd = append(crowd, c.getCrowd(snifferMAC, until))

	return crowd
}

func (c *CrowdAPI) getCrowd(snifferMAC string, when int64) model.Crowd {
	count := c.DB.GetUniqueMACCountBySnifferBetweenDates(snifferMAC, when-c.intervalInSeconds, when)
	return model.Crowd{
		Count: count,
		Time:  time.Unix(when, 0),
	}
}

func (c *CrowdAPI) getCrowdParams(ctx echo.Context) CrowdParams {
	params := ctx.QueryParams()
	from, fromExists := params["from"]
	until, untilExists := params["until"]
	forEvery, forExists := params["for"]

	crowdParams := CrowdParams{}

	if fromExists && untilExists && forExists {
		crowdParams.from, _ = strconv.ParseInt(from[0], 10, 64)
		crowdParams.until, _ = strconv.ParseInt(until[0], 10, 64)
		crowdParams.forEverySecond, _ = strconv.ParseInt(forEvery[0], 10, 64)
		return crowdParams
	}

	now := c.clock.Now().Unix()

	crowdParams.from = now
	crowdParams.until = now
	return crowdParams
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
