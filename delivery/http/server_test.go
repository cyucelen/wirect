package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"

	"github.com/cyucelen/wirect/model"
	"github.com/cyucelen/wirect/test"
	testutil "github.com/cyucelen/wirect/test/util"
)

var client = &http.Client{}

type IntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	clock  clock.Clock
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) BeforeTest(string, string) {
	mockClock := clock.NewMock()
	mockClock.Add(12 * time.Hour)
	s.clock = mockClock
	tick = mockClock
	server := Create(&test.InMemoryDB{})
	s.server = httptest.NewServer(server)
}

func (s *IntegrationSuite) AfterTest(string, string) {
	s.server.Close()
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
	bodyTemplate := `{"MAC":"02:02:02:02:02:02","name":"%s","location":"%s"}`
	s.sendUpdateSnifferRequest(fmt.Sprintf(bodyTemplate, newName, newLocation))

	actualSniffers := s.sendGetSniffersRequest()
	expectedSniffers := []model.Sniffer{
		{MAC: "00:11:22:33:44:55", Name: "library_sniffer", Location: "library"},
		{MAC: "02:02:02:02:02:02", Name: newName, Location: newLocation},
	}

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *IntegrationSuite) TestGetCrowd() {
	snifferMAC := "01:01:01:01:01:01"
	snifferPayload := `{"MAC":"` + snifferMAC + `","name":"library_sniffer","location":"library"}`
	s.sendCreateSnifferRequest(snifferPayload)

	now := s.clock.Now()
	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-15 * time.Second).Unix(), RSSI: 23.4, SnifferMAC: snifferMAC},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(-10 * time.Second).Unix(), RSSI: 44, SnifferMAC: snifferMAC},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-7 * time.Second).Unix(), RSSI: 333, SnifferMAC: snifferMAC},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-5 * time.Second).Unix(), RSSI: 1.2232, SnifferMAC: snifferMAC},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Unix(), RSSI: 1.2, SnifferMAC: snifferMAC},
	}

	for _, packet := range packets {
		packetJSON, _ := json.Marshal(packet)
		s.sendCreatePacketRequest(string(packetJSON))
	}

	actualCrowd := s.sendGetCrowdRequest(now, snifferMAC)
	expectedCrowd := 2
	assert.Equal(s.T(), expectedCrowd, actualCrowd.Count)

	now = now.Add(1 * time.Minute)
	packets = []model.Packet{
		{MAC: "CC:FF:CC:FF:CC:FF", Timestamp: now.Add(-35 * time.Second).Unix(), RSSI: 44, SnifferMAC: snifferMAC},
		{MAC: "DD:CC:DD:CC:DD:CC", Timestamp: now.Add(-25 * time.Second).Unix(), RSSI: 23.4, SnifferMAC: snifferMAC},
	}

	packetsJSON, _ := json.Marshal(packets)
	s.sendCreatePacketsRequest(string(packetsJSON))

	actualCrowd = s.sendGetCrowdRequest(now, snifferMAC)
	expectedCrowd = 4

	assert.Equal(s.T(), expectedCrowd, actualCrowd.Count)
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

func (s *IntegrationSuite) sendGetCrowdRequest(since time.Time, snifferMAC string) model.Crowd {
	req := s.newRequest(http.MethodGet, "crowd", "")

	testutil.AddGetCrowdRequestHeaders(req, since, snifferMAC)
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)
	var crowd model.Crowd
	json.NewDecoder(res.Body).Decode(&crowd)

	return crowd
}

func (s *IntegrationSuite) sendCreatePacketRequest(payload string) {
	res := s.sendRequest(http.MethodPost, "packets", payload)

	var actualResponse model.Packet
	json.NewDecoder(res.Body).Decode(&actualResponse)

	var expectedResponse model.Packet
	json.NewDecoder(strings.NewReader(payload)).Decode(&expectedResponse)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) sendCreatePacketsRequest(payload string) {
	res := s.sendRequest(http.MethodPost, "packets-collection", payload)

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

func (s *IntegrationSuite) sendUpdateSnifferRequest(payload string) {
	res := s.sendRequest(http.MethodPut, "sniffers", payload)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)
}
