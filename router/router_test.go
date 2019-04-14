package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) BeforeTest(string, string) {
	router := Create(&test.InMemoryDB{})
	s.server = httptest.NewServer(router)
}

func (s *IntegrationSuite) AfterTest(string, string) {
	s.server.Close()
}

func (s *IntegrationSuite) TestCreateSniffer() {
	snifferPayload := `{"MAC":"00:11:22:33:44:55","name":"library_sniffer","location":"library"}`
	s.SendCreateSnifferRequest(snifferPayload)

	snifferPayload = `{"MAC":"02:02:02:02:02:02","name":"room_sniffer","location":"room"}`
	s.SendCreateSnifferRequest(snifferPayload)

	actualSniffers := s.SendGetSniffersRequest()
	expectedSniffers := []model.Sniffer{
		{MAC: "00:11:22:33:44:55", Name: "library_sniffer", Location: "library"},
		{MAC: "02:02:02:02:02:02", Name: "room_sniffer", Location: "room"},
	}

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *IntegrationSuite) TestUpdateSniffer() {
	snifferPayload := `{"MAC":"00:11:22:33:44:55","name":"library_sniffer","location":"library"}`
	s.SendCreateSnifferRequest(snifferPayload)

	snifferPayload = `{"MAC":"02:02:02:02:02:02","name":"room_sniffer","location":"room"}`
	s.SendCreateSnifferRequest(snifferPayload)

	newName := "copy_center_sniffer"
	newLocation := "copy_center"
	bodyTemplate := `{"MAC":"02:02:02:02:02:02","name":"%s","location":"%s"}`
	s.SendUpdateSnifferRequest(fmt.Sprintf(bodyTemplate, newName, newLocation))

	actualSniffers := s.SendGetSniffersRequest()
	expectedSniffers := []model.Sniffer{
		{MAC: "00:11:22:33:44:55", Name: "library_sniffer", Location: "library"},
		{MAC: "02:02:02:02:02:02", Name: newName, Location: newLocation},
	}

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *IntegrationSuite) TestGetCrowd() {
	snifferMAC := "01:01:01:01:01:01"
	snifferPayload := `{"MAC":"` + snifferMAC + `","name":"library_sniffer","location":"library"}`
	s.SendCreateSnifferRequest(snifferPayload)

	now := time.Now()

	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now, RSSI: 1.2, SnifferMAC: snifferMAC},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(5 * time.Second), RSSI: 1.2232, SnifferMAC: snifferMAC},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(7 * time.Second), RSSI: 333, SnifferMAC: snifferMAC},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(10 * time.Second), RSSI: 44, SnifferMAC: snifferMAC},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(15 * time.Second), RSSI: 23.4, SnifferMAC: snifferMAC},
	}

	for _, packet := range packets {
		packetJSON, _ := json.Marshal(packet)
		s.SendCreatePacketRequest(string(packetJSON))
	}

	actualCrowd := s.SendGetCrowdRequest(now, snifferMAC)
	expectedCrowd := 2

	assert.Equal(s.T(), expectedCrowd, actualCrowd.Count)
}

func (s *IntegrationSuite) NewRequest(method, resource, body string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", s.server.URL, resource), strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(s.T(), err)

	return req
}

func (s *IntegrationSuite) SendRequest(method, endpoint, payload string) *http.Response {
	req := s.NewRequest(method, endpoint, payload)
	res, err := client.Do(req)
	assert.Nil(s.T(), err)

	return res
}

func (s *IntegrationSuite) SendGetCrowdRequest(since time.Time, snifferMAC string) model.Crowd {
	req := s.NewRequest(http.MethodGet, "crowd", "")
	testutil.AddGetCrowdRequestHeaders(req, since, snifferMAC)
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var crowd model.Crowd
	json.NewDecoder(res.Body).Decode(&crowd)

	return crowd
}

func (s *IntegrationSuite) SendCreatePacketRequest(payload string) {
	res := s.SendRequest(http.MethodPost, "packet", payload)

	var actualResponse model.Packet
	json.NewDecoder(res.Body).Decode(&actualResponse)

	var expectedResponse model.Packet
	json.NewDecoder(strings.NewReader(payload)).Decode(&expectedResponse)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) SendGetSniffersRequest() []model.Sniffer {
	res := s.SendRequest(http.MethodGet, "sniffer", "")
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var sniffers []model.Sniffer
	json.NewDecoder(res.Body).Decode(&sniffers)

	return sniffers
}

func (s *IntegrationSuite) SendCreateSnifferRequest(payload string) {
	res := s.SendRequest(http.MethodPost, "sniffer", payload)

	var actualResponse model.Sniffer
	json.NewDecoder(res.Body).Decode(&actualResponse)
	var expectedResponse model.Sniffer
	json.NewDecoder(strings.NewReader(payload)).Decode(&expectedResponse)
	assert.Equal(s.T(), expectedResponse, actualResponse)

	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) SendUpdateSnifferRequest(payload string) {
	res := s.SendRequest(http.MethodPut, "sniffer", payload)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)
}
