package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wirect/wirect-server/api/mocks"
	"gitlab.com/wirect/wirect-server/model"
)

func createMockPacketDB() *mocks.PacketDatabase {
	mockPacketDB := &mocks.PacketDatabase{}
	mockPacketDB.On("CreatePacket", mock.AnythingOfType("*model.Packet")).Return(nil)

	return mockPacketDB
}

func createFailingMockPacketDB() *mocks.PacketDatabase {
	mockPacketDB := &mocks.PacketDatabase{}
	mockPacketDB.On("CreatePacket", mock.AnythingOfType("*model.Packet")).Return(errors.New(""))

	return mockPacketDB
}

func TestCreatePackets(t *testing.T) {
	now := time.Now().UTC()

	snifferPackets := []model.SnifferPacket{
		{MAC: "22:44:66:88:AA:CC", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "33:11:22:44:55:66", Timestamp: now.Add(5 * time.Second), RSSI: 222, SnifferMAC: "00:00:00:00:00:00"},
	}
	snifferPacketJSON, _ := json.Marshal(snifferPackets)

	req := httptest.NewRequest(http.MethodPost, "/packet", bytes.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := createMockPacketDB()
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePackets(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, string(snifferPacketJSON), strings.TrimRight(rec.Body.String(), "\n"))

	expectedPackets := []model.Packet{}
	for _, snifferPacket := range snifferPackets {
		expectedPackets = append(expectedPackets, *toPacket(&snifferPacket))
	}
	assert.True(t, assert.ObjectsAreEqual(expectedPackets, mockPacketDB.CreatedPackets))
}

func TestCreatePacketsWithEmptyRequiredFields(t *testing.T) {
	now := time.Now().UTC()
	onlyNotValidPackets := []model.SnifferPacket{
		{MAC: "", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "01:02:03:04:05:06", Timestamp: now.Add(1 * time.Second), RSSI: 123, SnifferMAC: ""},
		{MAC: "01:02:03:04:05:06", RSSI: 123, SnifferMAC: "01:02:03:04:05:06"},
	}
	onlyNotValidPacketsJSON, _ := json.Marshal(onlyNotValidPackets)

	req := httptest.NewRequest(http.MethodPost, "/packet", bytes.NewReader(onlyNotValidPacketsJSON))
	c, rec := createTestContext(req)

	mockPacketDB := createMockPacketDB()
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePackets(c)

	assert.Len(t, mockPacketDB.CreatedPackets, 0)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	validAndNotValidPackets := []model.SnifferPacket{
		{MAC: "", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "22:44:66:88:AA:CC", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "01:02:03:04:05:06", Timestamp: now.Add(1 * time.Second), RSSI: 123, SnifferMAC: ""},
		{MAC: "01:02:03:04:05:06", RSSI: 123, SnifferMAC: "01:02:03:04:05:06"},
		{MAC: "33:11:22:44:55:66", Timestamp: now.Add(5 * time.Second), RSSI: 222, SnifferMAC: "00:00:00:00:00:00"},
	}
	validAndNotValidPacketsJSON, _ := json.Marshal(validAndNotValidPackets)
	req = httptest.NewRequest(http.MethodPost, "/packet", bytes.NewReader(validAndNotValidPacketsJSON))
	c, rec = createTestContext(req)

	packetAPI.CreatePackets(c)

	assert.Len(t, mockPacketDB.CreatedPackets, 2)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreatePacketsWithEmptyJSON(t *testing.T) {
	snifferPacketJSON := `{}`

	req := httptest.NewRequest(http.MethodPost, "/packet", strings.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := &mocks.PacketDatabase{}
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePackets(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreatePacketsWithCorruptedJSON(t *testing.T) {
	snifferPacketJSON := `{"Tim}`

	req := httptest.NewRequest(http.MethodPut, "/packet", strings.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := &mocks.PacketDatabase{}
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePackets(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreatePacketsWithFailingDB(t *testing.T) {
	now := time.Now().UTC()

	snifferPackets := []model.SnifferPacket{
		{MAC: "22:44:66:88:AA:CC", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "33:11:22:44:55:66", Timestamp: now.Add(5 * time.Second), RSSI: 222, SnifferMAC: "00:00:00:00:00:00"},
	}
	snifferPacketJSON, _ := json.Marshal(snifferPackets)

	req := httptest.NewRequest(http.MethodPost, "/packet", bytes.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := createFailingMockPacketDB()
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePackets(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
