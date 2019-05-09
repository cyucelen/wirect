package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-ffmt/ffmt"

	"github.com/cyucelen/wirect/model"
	"github.com/cyucelen/wirect/test"
	testutil "github.com/cyucelen/wirect/test/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateCrowdAPI(t *testing.T) {
	expectedInterval := 2 * time.Minute
	crowdAPI := CreateCrowdAPI(&test.InMemoryDB{}, SetCrowdCalculationInterval(expectedInterval))

	assert.NotNil(t, crowdAPI)
	assert.Equal(t, expectedInterval, crowdAPI.Interval)
}

func TestGetCurrentCrowd(t *testing.T) {
	mockClock := clock.NewMock()
	mockClock.Add(1 * time.Hour)

	now := mockClock.Now()
	snifferMAC := "11:22:00:33:44:55"
	mockPacketDB := createDBContainsPacketsOfTwoUniquePerson(now, snifferMAC)

	crowdAPI := CreateCrowdAPI(mockPacketDB, SetCrowdClock(mockClock), SetCrowdCalculationInterval(5*time.Minute))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	c, rec := createTestContext(req)
	c.SetPath("/sniffers/:snifferMAC/crowd")
	c.SetParamNames("snifferMAC")
	c.SetParamValues(url.QueryEscape(snifferMAC))
	crowdAPI.GetCrowd(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualCrowd []model.Crowd
	json.NewDecoder(rec.Body).Decode(&actualCrowd)
	expectedCrowd := []model.Crowd{{Count: 2, Time: now}}
	assert.Equal(t, expectedCrowd, actualCrowd)
}

func TestGetCrowdBetweenDates(t *testing.T) {
	mockClock := clock.NewMock()
	mockClock.Add(1 * time.Hour)
	now := mockClock.Now()
	snifferMAC := "11:22:00:33:44:55"
	mockPacketDB := createDBContainsPacketsOfTwoUniquePerson(now, snifferMAC)

	crowdAPI := CreateCrowdAPI(mockPacketDB, SetCrowdClock(mockClock), SetCrowdCalculationInterval(5*time.Minute))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	from := now.Add(-20 * time.Second)
	ffmt.Puts(from)
	until := now.Add(-6 * time.Second)
	forEvery := time.Second * 10
	testutil.AddGetCrowdBetweenDatesQueries(req, from, until, forEvery)

	c, rec := createTestContext(req)
	c.SetPath("/sniffers/:snifferMAC/crowd")
	c.SetParamNames("snifferMAC")
	c.SetParamValues(url.QueryEscape(snifferMAC))
	crowdAPI.GetCrowd(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualCrowd []model.Crowd
	json.NewDecoder(rec.Body).Decode(&actualCrowd)
	expectedCrowd := []model.Crowd{
		{Count: 0, Time: from}, {Count: 2, Time: from.Add(forEvery)}, {Count: 2, Time: until},
	}
	assert.Equal(t, expectedCrowd, actualCrowd)
}

func createDBContainsPacketsOfTwoUniquePerson(now time.Time, snifferMAC string) *test.InMemoryDB {
	db := &test.InMemoryDB{}

	packets := []model.Packet{
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(-15 * time.Second).Unix(), // 02:59:45
			RSSI:       23.4,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "00:11:CC:CC:44:55",
			Timestamp:  now.Add(-10 * time.Second).Unix(), // 02:59:50
			RSSI:       44,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(-7 * time.Second).Unix(), // 02:59:53
			RSSI:       333,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(-5 * time.Second).Unix(), // 02:59:55
			RSSI:       1.2232,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Unix(), // 03:00:00
			RSSI:       1.2,
			SnifferMAC: snifferMAC,
		},
	}
	for _, packet := range packets {
		db.CreatePacket(&packet)
	}

	return db
}
