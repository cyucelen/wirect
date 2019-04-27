package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/benbjohnson/clock"
	"github.com/cyucelen/wirect/model"
)

func TestGetTime(t *testing.T) {
	mockClock := clock.NewMock()
	mockClock.Add(time.Hour) //  Jan 1, 1970, 01:00:00 -> 3600000ms

	timeAPI := &TimeAPI{Clock: mockClock}

	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	c, rec := createTestContext(req)

	err := timeAPI.GetTime(c)
	assert.Nil(t, err)

	expectedTime := mockClock.Now().Unix()

	var actualTime model.Time
	json.NewDecoder(rec.Body).Decode(&actualTime)
	assert.Equal(t, expectedTime, actualTime.Now)
	assert.Equal(t, http.StatusOK, rec.Code)
}
