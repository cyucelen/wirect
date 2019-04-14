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
	"github.com/cyucelen/wirect/api/mocks"
	"github.com/cyucelen/wirect/model"
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

func TestCreatePacket(t *testing.T) {
	now := time.Now().UTC()

	snifferPacket := model.SnifferPacket{
		MAC: "22:44:66:88:AA:CC", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00",
	}
	snifferPacketJSON, _ := json.Marshal(snifferPacket)

	req := httptest.NewRequest(http.MethodPost, "/packet", bytes.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := createMockPacketDB()
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePacket(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, string(snifferPacketJSON), strings.TrimRight(rec.Body.String(), "\n"))

	expectedPacket := *toPacket(&snifferPacket)

	assert.True(t, assert.ObjectsAreEqual(expectedPacket, mockPacketDB.CreatedPackets[0]))
}

func TestCreatePacketsWithEmptyRequiredFields(t *testing.T) {
	mockPacketDB := createMockPacketDB()

	now := time.Now().UTC()
	notValidSnifferPackets := []model.SnifferPacket{
		{MAC: "", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "01:02:03:04:05:06", Timestamp: now.Add(1 * time.Second), RSSI: 123, SnifferMAC: ""},
		{MAC: "01:02:03:04:05:06", RSSI: 123, SnifferMAC: "01:02:03:04:05:06"},
	}

	packetAPI := PacketAPI{mockPacketDB}

	for _, notValidSnifferPacket := range notValidSnifferPackets {
		notValidSnifferPacketJSON, _ := json.Marshal(notValidSnifferPacket)
		req := httptest.NewRequest(http.MethodPost, "/packet", bytes.NewReader(notValidSnifferPacketJSON))
		c, rec := createTestContext(req)

		packetAPI.CreatePacket(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreatePacketsWithEmptyJSON(t *testing.T) {
	snifferPacketJSON := `{}`

	req := httptest.NewRequest(http.MethodPost, "/packet", strings.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := &mocks.PacketDatabase{}
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePacket(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreatePacketsWithCorruptedJSON(t *testing.T) {
	snifferPacketJSON := `{"Tim}`

	req := httptest.NewRequest(http.MethodPut, "/packet", strings.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := &mocks.PacketDatabase{}
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePacket(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreatePacketsWithFailingDB(t *testing.T) {
	now := time.Now().UTC()

	snifferPacket := model.SnifferPacket{
		MAC: "22:44:66:88:AA:CC", Timestamp: now, RSSI: 123, SnifferMAC: "00:00:00:00:00:00",
	}
	snifferPacketJSON, _ := json.Marshal(snifferPacket)

	req := httptest.NewRequest(http.MethodPost, "/packet", bytes.NewReader(snifferPacketJSON))
	c, rec := createTestContext(req)

	mockPacketDB := createFailingMockPacketDB()
	packetAPI := PacketAPI{mockPacketDB}

	packetAPI.CreatePacket(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
