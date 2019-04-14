package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	testutil "github.com/cyucelen/wirect/test/util"

	"github.com/cyucelen/wirect/model"
	"github.com/cyucelen/wirect/test"
	"github.com/stretchr/testify/assert"
)

func TestGetCrowd(t *testing.T) {
	now := time.Now()
	snifferMAC := "11:22:00:33:44:55"
	mockPacketDB := createDBContainsPacketsOfTwoUniquePerson(now, snifferMAC)

	crowdAPI := &CrowdAPI{mockPacketDB}

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
			Timestamp:  now,
			RSSI:       1.2,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(5 * time.Second),
			RSSI:       1.2232,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(7 * time.Second),
			RSSI:       333,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "00:11:CC:CC:44:55",
			Timestamp:  now.Add(10 * time.Second),
			RSSI:       44,
			SnifferMAC: snifferMAC,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  now.Add(15 * time.Second),
			RSSI:       23.4,
			SnifferMAC: snifferMAC,
		},
	}
	for _, packet := range packets {
		db.CreatePacket(&packet)
	}

	return db
}
