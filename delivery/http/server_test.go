package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/cyucelen/wirect/model"
	"github.com/cyucelen/wirect/test"
	testutil "github.com/cyucelen/wirect/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var client = &http.Client{}

type IntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	clock  clock.Clock
	db     Database
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) BeforeTest(string, string) {
	mockClock := clock.NewMock()
	mockClock.Add(12 * time.Hour)
	s.clock = mockClock
	tick = mockClock
	s.db = &test.InMemoryDB{}
	server := Create(s.db)
	s.server = httptest.NewServer(server)
}

func (s *IntegrationSuite) AfterTest(string, string) {
	s.server.Close()
}

func (s *IntegrationSuite) setCurrentTime(time time.Time) {
	mockClock := clock.NewMock()
	mockClock.Set(time)
	tick = mockClock
	s.clock = mockClock
	server := Create(s.db)
	s.server = httptest.NewServer(server)
}

func (s *IntegrationSuite) TestCreateSniffer() {
	snifferPayload := `{"MAC":"00:11:22:33:44:55","name":"library_sniffer","location":"library"}`
	s.sendCreateSnifferRequest(snifferPayload)

	snifferPayload = `{"MAC":"02:02:02:02:02:02","name":"room_sniffer","location":"room"}`
	s.sendCreateSnifferRequest(snifferPayload)

	actualSniffers := s.sendGetSniffersRequest()
	expectedSniffers := []model.Sniffer{
		{MAC: "00:11:22:33:44:55", Name: "library_sniffer", Location: "library"},
		{MAC: "02:02:02:02:02:02", Name: "room_sniffer", Location: "room"},
	}

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *IntegrationSuite) TestUpdateSniffer() {
	snifferPayload := `{"MAC":"00:11:22:33:44:55","name":"library_sniffer","location":"library"}`
	s.sendCreateSnifferRequest(snifferPayload)

	snifferPayload = `{"MAC":"02:02:02:02:02:02","name":"room_sniffer","location":"room"}`
	s.sendCreateSnifferRequest(snifferPayload)

	newName := "copy_center_sniffer"
	newLocation := "copy_center"
	bodyTemplate := `{"name":"%s","location":"%s"}`
	s.sendUpdateSnifferRequest("02:02:02:02:02:02", fmt.Sprintf(bodyTemplate, newName, newLocation))

	actualSniffers := s.sendGetSniffersRequest()
	expectedSniffers := []model.Sniffer{
		{MAC: "00:11:22:33:44:55", Name: "library_sniffer", Location: "library"},
		{MAC: "02:02:02:02:02:02", Name: newName, Location: newLocation},
	}

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *IntegrationSuite) TestGetCurrentCrowd() {
	snifferMAC := "01:01:01:01:01:01"
	snifferPayload := `{"MAC":"` + snifferMAC + `","name":"library_sniffer","location":"library"}`
	s.sendCreateSnifferRequest(snifferPayload)

	now := s.clock.Now()
	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-15 * time.Second).Unix(), RSSI: 23.4},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(-10 * time.Second).Unix(), RSSI: 44},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-7 * time.Second).Unix(), RSSI: 333},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-5 * time.Second).Unix(), RSSI: 1.2232},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Unix(), RSSI: 1.2},
	}

	for _, packet := range packets {
		packetJSON, _ := json.Marshal(packet)
		s.sendCreatePacketRequest(snifferMAC, string(packetJSON))
	}

	actualCrowd := s.sendGetCurrentCrowdRequest(snifferMAC)
	expectedCrowd := []model.Crowd{{Count: 2, Time: now}}
	assert.Equal(s.T(), expectedCrowd, actualCrowd)

	now = now.Add(1 * time.Minute)
	s.setCurrentTime(now)
	packets = []model.Packet{
		{MAC: "CC:FF:CC:FF:CC:FF", Timestamp: now.Add(-35 * time.Second).Unix(), RSSI: 44},
		{MAC: "DD:CC:DD:CC:DD:CC", Timestamp: now.Add(-25 * time.Second).Unix(), RSSI: 23.4},
	}

	packetsJSON, _ := json.Marshal(packets)
	s.sendCreatePacketsRequest(snifferMAC, string(packetsJSON))

	actualCrowd = s.sendGetCurrentCrowdRequest(snifferMAC)
	expectedCrowd = []model.Crowd{{Count: 4, Time: now}}

	assert.Equal(s.T(), expectedCrowd, actualCrowd)
}

func (s *IntegrationSuite) TestGetCrowdBetweenDates() {
	snifferMAC := "01:01:01:01:01:01"
	snifferPayload := `{"MAC":"` + snifferMAC + `","name":"library_sniffer","location":"library"}`
	s.sendCreateSnifferRequest(snifferPayload)

	now := s.clock.Now()
	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-15 * time.Second).Unix(), RSSI: 23.4},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(-10 * time.Second).Unix(), RSSI: 44},
		{MAC: "CC:BB:22:11:44:55", Timestamp: now.Add(-7 * time.Second).Unix(), RSSI: 333},
		{MAC: "DD:BB:22:11:44:55", Timestamp: now.Add(-5 * time.Second).Unix(), RSSI: 1.2232},
		{MAC: "EE:BB:22:11:44:55", Timestamp: now.Unix(), RSSI: 1.2},
	}

	for _, packet := range packets {
		packetJSON, _ := json.Marshal(packet)
		s.sendCreatePacketRequest(snifferMAC, string(packetJSON))
	}

	from := now.Add(-20 * time.Second)
	until := now.Add(-8 * time.Second)
	forEvery := 10 * time.Second
	actualCrowd := s.sendGetCrowdBetweenDatesRequest(from, until, forEvery, snifferMAC)
	expectedCrowd := []model.Crowd{
		{Count: 0, Time: from}, {Count: 2, Time: from.Add(forEvery)}, {Count: 2, Time: until},
	}
	assert.Equal(s.T(), expectedCrowd, actualCrowd)
}

func (s *IntegrationSuite) TestGetTotalSniffedMACDaily() {
	snifferMAC := "01:01:01:01:01:01"
	snifferPayload := `{"MAC":"` + snifferMAC + `","name":"library_sniffer","location":"library"}`
	s.sendCreateSnifferRequest(snifferPayload)

	now := s.clock.Now()
	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-15 * time.Second).Unix(), RSSI: 23.4},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(-10 * time.Second).Unix(), RSSI: 44},
		{MAC: "DD:BB:22:11:44:55", Timestamp: now.Add(-7 * time.Second).Unix(), RSSI: 333},
		{MAC: "DD:BB:22:11:44:55", Timestamp: now.Add(-5 * time.Second).Unix(), RSSI: 1.2232},
		{MAC: "EE:BB:22:11:44:55", Timestamp: now.Unix(), RSSI: 1.2},
		{MAC: "FF:FB:F2:F1:F4:F5", Timestamp: now.Add(25 * time.Hour).Unix(), RSSI: 1.2},
	}

	for _, packet := range packets {
		packetJSON, _ := json.Marshal(packet)
		s.sendCreatePacketRequest(snifferMAC, string(packetJSON))
	}

	actualTotalSniffed := s.sendGetTotalSniffedMACDailyRequest(snifferMAC)
	expectedTotalSniffed := model.TotalSniffed{Count: 4}

	assert.Equal(s.T(), expectedTotalSniffed, actualTotalSniffed)
}

func (s *IntegrationSuite) TestGetTime() {
	expectedTime := s.clock.Now().Unix()
	actualTime := s.sendGetTimeRequest()
	assert.Equal(s.T(), expectedTime, actualTime.Now)
}

func (s *IntegrationSuite) newRequest(method, resource, body string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", s.server.URL, resource), strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(s.T(), err)

	return req
}

func (s *IntegrationSuite) sendRequest(method, endpoint, payload string) *http.Response {
	req := s.newRequest(method, endpoint, payload)
	res, err := client.Do(req)
	assert.Nil(s.T(), err)

	return res
}

func (s *IntegrationSuite) sendGetTimeRequest() model.Time {
	res := s.sendRequest(http.MethodGet, "time", "")

	var t model.Time
	json.NewDecoder(res.Body).Decode(&t)

	return t
}

func (s *IntegrationSuite) sendGetCurrentCrowdRequest(snifferMAC string) []model.Crowd {
	resource := fmt.Sprintf("sniffers/%s/stats/crowd", url.QueryEscape(snifferMAC))
	req := s.newRequest(http.MethodGet, resource, "")

	// testutil.AddGetCurrentCrowdQueries(req, since)
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)
	var crowd []model.Crowd
	json.NewDecoder(res.Body).Decode(&crowd)

	return crowd
}

func (s *IntegrationSuite) sendGetCrowdBetweenDatesRequest(from, until time.Time, forEvery time.Duration, snifferMAC string) []model.Crowd {
	resource := fmt.Sprintf("sniffers/%s/stats/crowd", url.QueryEscape(snifferMAC))
	req := s.newRequest(http.MethodGet, resource, "")
	testutil.AddGetCrowdBetweenDatesQueries(req, from, until, forEvery)
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var crowd []model.Crowd
	json.NewDecoder(res.Body).Decode(&crowd)

	return crowd
}

func (s *IntegrationSuite) sendGetTotalSniffedMACDailyRequest(snifferMAC string) model.TotalSniffed {
	resource := fmt.Sprintf("sniffers/%s/stats/total-sniffed/daily", url.QueryEscape(snifferMAC))
	req := s.newRequest(http.MethodGet, resource, "")
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var totalSniffed model.TotalSniffed
	json.NewDecoder(res.Body).Decode(&totalSniffed)

	return totalSniffed
}

func (s *IntegrationSuite) sendCreatePacketRequest(snifferMAC, payload string) {
	resource := fmt.Sprintf("sniffers/%s/packets", snifferMAC)
	res := s.sendRequest(http.MethodPost, resource, payload)

	var actualResponse model.Packet
	json.NewDecoder(res.Body).Decode(&actualResponse)

	var expectedResponse model.Packet
	json.NewDecoder(strings.NewReader(payload)).Decode(&expectedResponse)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) sendCreatePacketsRequest(snifferMAC, payload string) {
	resource := fmt.Sprintf("sniffers/%s/packets-collection", snifferMAC)
	res := s.sendRequest(http.MethodPost, resource, payload)

	var actualResponse []model.Packet
	json.NewDecoder(res.Body).Decode(&actualResponse)

	var expectedResponse []model.Packet
	json.NewDecoder(strings.NewReader(payload)).Decode(&expectedResponse)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) sendGetSniffersRequest() []model.Sniffer {
	res := s.sendRequest(http.MethodGet, "sniffers", "")
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var sniffers []model.Sniffer
	json.NewDecoder(res.Body).Decode(&sniffers)

	return sniffers
}

func (s *IntegrationSuite) sendCreateSnifferRequest(payload string) {
	res := s.sendRequest(http.MethodPost, "sniffers", payload)

	var actualResponse model.Sniffer
	json.NewDecoder(res.Body).Decode(&actualResponse)

	var expectedResponse model.Sniffer
	json.NewDecoder(strings.NewReader(payload)).Decode(&expectedResponse)
	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) sendUpdateSnifferRequest(snifferMAC, payload string) {
	resource := fmt.Sprintf("sniffers/%s", url.QueryEscape(snifferMAC))
	res := s.sendRequest(http.MethodPut, resource, payload)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)
}
