package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"

	"gitlab.com/wirect/wirect-server/model"
	"gitlab.com/wirect/wirect-server/test"
)

var client = &http.Client{}

type IntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	db     *test.InMemoryDB
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) BeforeTest(string, string) {
	s.db = &test.InMemoryDB{}
	router := Create(s.db)
	s.server = httptest.NewServer(router)
}

func (s *IntegrationSuite) AfterTest(string, string) {
	s.server.Close()
}

func (s *IntegrationSuite) TestCreatePacket() {
	body := `
	{
		"MAC":"00:11:22:33:44:55",
		"timestamp":"2019-04-13T12:10:23Z",
		"RSSI":122.4,
		"snifferMAC":"FF:FF:FF:FF:AA:AA"
	}`

	req := s.NewRequest(http.MethodPost, "packet", body)

	res, err := client.Do(req)
	assert.Nil(s.T(), err)

	var actualResponse model.SnifferPacket
	json.NewDecoder(res.Body).Decode(&actualResponse)

	var expectedResponse model.SnifferPacket
	json.NewDecoder(strings.NewReader(body)).Decode(&expectedResponse)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) TestCreateSniffer() {
	body := `
	{
		"MAC":"00:11:22:33:44:55",
		"name":"library_sniffer",
		"location":"library"
	}
	`
	req := s.NewRequest(http.MethodPost, "sniffer", body)

	res, err := client.Do(req)
	assert.Nil(s.T(), err)

	var actualResponse model.Sniffer
	json.NewDecoder(res.Body).Decode(&actualResponse)

	var expectedResponse model.Sniffer
	json.NewDecoder(strings.NewReader(body)).Decode(&expectedResponse)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
}

func (s *IntegrationSuite) TestGetSniffers() {
	expectedSniffers := []model.Sniffer{
		{
			MAC:      "01:01:01:01:01:01",
			Name:     "lib_sniffer",
			Location: "library",
		},
		{
			MAC:      "02:02:02:02:02:02",
			Name:     "room_sniffer",
			Location: "room",
		},
	}
	s.createSniffers(expectedSniffers)

	req := s.NewRequest(http.MethodGet, "sniffer", "")
	res, err := client.Do(req)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	var actualSniffers []model.Sniffer
	json.NewDecoder(res.Body).Decode(&actualSniffers)

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *IntegrationSuite) TestUpdateSniffer() {
	sniffers := []model.Sniffer{
		{
			MAC:      "01:01:01:01:01:01",
			Name:     "lib_sniffer",
			Location: "library",
		},
		{
			MAC:      "02:02:02:02:02:02",
			Name:     "room_sniffer",
			Location: "room",
		},
	}
	s.createSniffers(sniffers)

	newName := "copy_center_sniffer"
	newLocation := "copy_center"

	bodyTemplate := `
		{
			"MAC":"02:02:02:02:02:02",
			"name":"%s",
			"location":"%s"
		}
	`

	req := s.NewRequest(http.MethodPut, "sniffer", fmt.Sprintf(bodyTemplate, newName, newLocation))
	res, err := client.Do(req)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	expectedSniffers := make([]model.Sniffer, len(sniffers))
	copy(expectedSniffers, sniffers)
	expectedSniffers[1].Name = newName
	expectedSniffers[1].Location = newLocation

	actualSniffers := s.db.Sniffers

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *IntegrationSuite) createSniffers(sniffers []model.Sniffer) {
	for _, sniffer := range sniffers {
		s.db.CreateSniffer(&sniffer)
	}
}

func (s *IntegrationSuite) NewRequest(method, resource, body string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", s.server.URL, resource), strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(s.T(), err)

	return req
}
