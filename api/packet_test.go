package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/cyucelen/wirect/api/mocks"
	"github.com/cyucelen/wirect/model"
	"github.com/cyucelen/wirect/test"
	"github.com/stretchr/testify/assert"
)

type PacketAPISuite struct {
	suite.Suite
	packetAPI *PacketAPI
	packetDB  *test.InMemoryDB
}

func TestPacketAPISuite(t *testing.T) {
	suite.Run(t, new(PacketAPISuite))
}

func createMockPacketDB() *test.InMemoryDB {
	mockPacketDB := &test.InMemoryDB{}
	return mockPacketDB
}

func createFailingMockPacketDB() *mocks.PacketDatabase {
	mockPacketDB := &mocks.PacketDatabase{}
	mockPacketDB.On("CreatePacket", mock.AnythingOfType("*model.Packet")).Return(errors.New(""))
	return mockPacketDB
}

const defaultTestSnifferMAC = "00:00:00:00:00:00"

func (s *PacketAPISuite) BeforeTest(string, string) {
	s.packetDB = createMockPacketDB()
	s.packetDB.Sniffers = append(s.packetDB.Sniffers, model.Sniffer{MAC: defaultTestSnifferMAC})
	s.packetAPI = &PacketAPI{s.packetDB}
}

func (s *PacketAPISuite) TestCreatePacket() {
	snifferPacket := model.SnifferPacket{
		MAC: "22:44:66:88:AA:CC", Timestamp: time.Now().UTC().Unix(), RSSI: 123,
	}

	rec := sendTestRequestToHandler(defaultTestSnifferMAC, snifferPacket, s.packetAPI.CreatePacket)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	var actualSnifferPacket model.SnifferPacket
	json.NewDecoder(rec.Body).Decode(&actualSnifferPacket)
	assert.Equal(s.T(), snifferPacket, actualSnifferPacket)

	expectedPacket := *toPacket(&snifferPacket, defaultTestSnifferMAC)
	assert.True(s.T(), assert.ObjectsAreEqual(expectedPacket, s.packetDB.Packets[0]))
}

func (s *PacketAPISuite) TestCreatePacketWithEmptyRequiredFields() {
	notValidSnifferPackets := []model.SnifferPacket{
		{MAC: "", Timestamp: time.Now().UTC().Unix(), RSSI: 123},
		{MAC: "01:02:03:04:05:06", RSSI: 123},
	}

	for _, notValidSnifferPacket := range notValidSnifferPackets {
		rec := sendTestRequestToHandler(defaultTestSnifferMAC, notValidSnifferPacket, s.packetAPI.CreatePacket)
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	}
}

func (s *PacketAPISuite) TestCreatePacketWithEmptyJSON() {
	responseStatusCode := sendTestRequestToHandlerWithEmptyJSON(s.packetAPI.CreatePacket)
	assert.Equal(s.T(), http.StatusBadRequest, responseStatusCode)
}

func (s *PacketAPISuite) TestCreatePacketWithCorruptedJSON() {
	responseStatusCode := sendTestRequestToHandlerWithCorruptedJSON(s.packetAPI.CreatePacket)
	assert.Equal(s.T(), http.StatusBadRequest, responseStatusCode)
}

func TestCreatePacketWithFailingDB(t *testing.T) {
	mockFailingPacketDB := createFailingMockPacketDB()
	packetAPI := PacketAPI{mockFailingPacketDB}

	snifferPacket := model.SnifferPacket{
		MAC: "22:44:66:88:AA:CC", Timestamp: time.Now().UTC().Unix(), RSSI: 123,
	}

	rec := sendTestRequestToHandler(defaultTestSnifferMAC, snifferPacket, packetAPI.CreatePacket)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func (s *PacketAPISuite) TestCreatePackets() {
	snifferPackets := []model.SnifferPacket{
		{MAC: "22:44:66:88:AA:CC", Timestamp: time.Now().UTC().Unix(), RSSI: 123},
		{MAC: "33:11:22:44:55:66", Timestamp: time.Now().UTC().Unix(), RSSI: 222},
	}

	rec := sendTestRequestToHandler(defaultTestSnifferMAC, snifferPackets, s.packetAPI.CreatePackets)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	var actualSnifferPackets []model.SnifferPacket
	json.NewDecoder(rec.Body).Decode(&actualSnifferPackets)
	assert.Equal(s.T(), snifferPackets, actualSnifferPackets)

	expectedPackets := []model.Packet{}
	for _, snifferPacket := range snifferPackets {
		expectedPackets = append(expectedPackets, *toPacket(&snifferPacket, defaultTestSnifferMAC))
	}
	assert.True(s.T(), assert.ObjectsAreEqual(expectedPackets, s.packetDB.Packets))
}

func (s *PacketAPISuite) TestCreatePacketsWithEmptyRequiredFields() {
	onlyNotValidPackets := []model.SnifferPacket{
		{MAC: "", Timestamp: time.Now().UTC().Unix(), RSSI: 123},
		{MAC: "01:02:03:04:05:06", RSSI: 123},
	}

	rec := sendTestRequestToHandler(defaultTestSnifferMAC, onlyNotValidPackets, s.packetAPI.CreatePackets)
	assert.Len(s.T(), s.packetDB.Packets, 0)
	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)

	validAndNotValidPackets := []model.SnifferPacket{
		{MAC: "", Timestamp: time.Now().UTC().Unix(), RSSI: 123},
		{MAC: "22:44:66:88:AA:CC", Timestamp: time.Now().UTC().Unix(), RSSI: 123},
		{MAC: "01:02:03:04:05:06", Timestamp: time.Now().UTC().Unix(), RSSI: 123},
		{MAC: "01:02:03:04:05:06", RSSI: 123},
		{MAC: "33:11:22:44:55:66", Timestamp: time.Now().UTC().Unix(), RSSI: 222},
	}

	rec = sendTestRequestToHandler(defaultTestSnifferMAC, validAndNotValidPackets, s.packetAPI.CreatePackets)
	assert.Len(s.T(), s.packetDB.Packets, 3)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)
}

func (s *PacketAPISuite) TestCreatePacketsWithEmptyJSON() {
	responseStatusCode := sendTestRequestToHandlerWithEmptyJSON(s.packetAPI.CreatePackets)
	assert.Equal(s.T(), http.StatusBadRequest, responseStatusCode)
}
func (s *PacketAPISuite) TestCreatePacketsWithCorruptedJSON() {
	responseStatusCode := sendTestRequestToHandlerWithCorruptedJSON(s.packetAPI.CreatePackets)
	assert.Equal(s.T(), http.StatusBadRequest, responseStatusCode)
}

func TestCreatePacketsWithFailingDB(t *testing.T) {
	mockPacketDB := createFailingMockPacketDB()
	packetAPI := PacketAPI{mockPacketDB}

	snifferPackets := []model.SnifferPacket{
		{MAC: "22:44:66:88:AA:CC", Timestamp: time.Now().UTC().Unix(), RSSI: 123},
		{MAC: "33:11:22:44:55:66", Timestamp: time.Now().UTC().Unix(), RSSI: 222},
	}

	rec := sendTestRequestToHandler(defaultTestSnifferMAC, snifferPackets, packetAPI.CreatePackets)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
