package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/benbjohnson/clock"

	testutil "github.com/cyucelen/wirect/test/util"

	"github.com/cyucelen/wirect/model"
	"github.com/cyucelen/wirect/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateCrowdAPI(t *testing.T) {
	expectedInterval := 2 * time.Minute
	crowdAPI := CreateCrowdAPI(&test.InMemoryDB{}, SetCrowdCalculationInterval(expectedInterval))

	assert.NotNil(t, crowdAPI)
	assert.Equal(t, expectedInterval, crowdAPI.Interval)
}

func TestGetCrowd(t *testing.T) {
	mockClock := clock.NewMock()
	mockClock.Add(1 * time.Hour)

	now := mockClock.Now()
	snifferMAC := "11:22:00:33:44:55"
	mockPacketDB := createDBContainsPacketsOfTwoUniquePerson(now, snifferMAC)

	crowdAPI := &CrowdAPI{DB: mockPacketDB, Interval: 5 * time.Minute}

	req := httptest.NewRequest(http.MethodGet, "/crowd", nil)
	testutil.AddGetCrowdRequestHeaders(req, now, snifferMAC)

	c, rec := createTestContext(req)
	crowdAPI.GetCrowd(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var crowd model.Crowd
	json.NewDecoder(rec.Body).Decode(&crowd)

	assert.Equal(t, crowd.Count, 2)
}

func createDBContainsPacketsOfTwoUniquePerson(now time.Time, snifferMAC string) *test.InMemoryDB {
	db := &test.InMemoryDB{}

	packets := []model.Packet{
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(-15 * time.Second).Unix(),
			RSSI:       23.4,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "00:11:CC:CC:44:55",
			Timestamp:  now.Add(-10 * time.Second).Unix(),
			RSSI:       44,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(-7 * time.Second).Unix(),
			RSSI:       333,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(-5 * time.Second).Unix(),
			RSSI:       1.2232,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Unix(),
			RSSI:       1.2,
			SnifferMAC: snifferMAC,
		},
	}
	for _, packet := range packets {
		db.CreatePacket(&packet)
	}

	return db
}
