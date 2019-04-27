package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cyucelen/wirect/model"
	"github.com/stretchr/testify/assert"
)

type InMemoryDBSuite struct {
	suite.Suite
	db *InMemoryDB
}

func TestInMemoryDBSuite(t *testing.T) {
	suite.Run(t, new(InMemoryDBSuite))
}

func (s *InMemoryDBSuite) BeforeTest(string, string) {
	s.db = &InMemoryDB{}
}

func (s *InMemoryDBSuite) TestCreatePacket() {
	expectedPacket := model.Packet{
		MAC:        "00:00:11:11:22:22",
		RSSI:       1.23,
		SnifferMAC: "AA:AA:AA:BB:CC:DD",
		Timestamp:  time.Now().Unix(),
	}

	s.db.CreatePacket(&expectedPacket)

	actualPacket := s.db.Packets[0]

	assert.Equal(s.T(), expectedPacket, actualPacket)
}

func (s *InMemoryDBSuite) TestGetPacketsBySniffer() {
	snifferMAC := "AA:AA:AA:BB:CC:DD"

	packets := []model.Packet{
		{
			MAC:        "00:00:11:11:22:22",
			RSSI:       1.23,
			SnifferMAC: snifferMAC,
			Timestamp:  time.Now().Unix(),
		},
		{
			MAC:        "22:22:12:11:22:22",
			RSSI:       1.23,
			SnifferMAC: snifferMAC,
			Timestamp:  time.Now().Unix(),
		},
		{
			MAC:        "00:00:11:11:22:22",
			RSSI:       1.23,
			SnifferMAC: "11:11:11:11:11:11",
			Timestamp:  time.Now().Unix(),
		},
	}

	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	actualPackets := s.db.GetPacketsBySniffer(snifferMAC)
	expectedPackets := packets[0:2]

	assert.Equal(s.T(), expectedPackets, actualPackets)
}

func (s *InMemoryDBSuite) TestGetPacketsBySnifferSince() {
	snifferOne := "01:02:03:04:05:06"
	snifferTwo := "00:00:00:00:00:00"

	since := time.Now().UTC()

	packets := []model.Packet{
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  since.Add(-10 * time.Minute).Unix(),
			RSSI:       123,
			SnifferMAC: snifferOne,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  since.Add(1 * time.Second).Unix(),
			RSSI:       333,
			SnifferMAC: snifferOne,
		},
		{
			MAC:        "00:11:CC:CC:44:55",
			Timestamp:  since.Add(-5 * time.Minute).Unix(),
			RSSI:       1234,
			SnifferMAC: snifferTwo,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  since.Unix(),
			RSSI:       333,
			SnifferMAC: snifferOne,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  since.Add(5 * time.Minute).Unix(),
			RSSI:       333,
			SnifferMAC: snifferOne,
		},
	}
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	snifferPacketsSince := s.db.GetPacketsBySnifferSince(snifferOne, since.Unix())

	assert.Len(s.T(), snifferPacketsSince, 3)
	assert.Equal(s.T(), packets[3].Timestamp, snifferPacketsSince[0].Timestamp)
	assert.Equal(s.T(), packets[1].Timestamp, snifferPacketsSince[1].Timestamp)
	assert.Equal(s.T(), packets[4].Timestamp, snifferPacketsSince[2].Timestamp)
}

func (s *InMemoryDBSuite) TestCreateSniffer() {
	expectedSniffer := model.Sniffer{
		MAC:      "AA:AA:AA:BB:CC:DD",
		Name:     "lab_sniffer",
		Location: "lab",
	}

	s.db.CreateSniffer(&expectedSniffer)

	actualSniffer := s.db.Sniffers[0]

	assert.Equal(s.T(), expectedSniffer, actualSniffer)
}

func (s *InMemoryDBSuite) TestGetSniffers() {
	expectedSniffers := []model.Sniffer{
		{
			MAC:      "AA:AA:AA:BB:CC:DD",
			Name:     "lab_sniffer",
			Location: "lab",
		},
		{
			MAC:      "AA:AA:DD:AA:BB:CC",
			Name:     "room_sniffer",
			Location: "room",
		},
	}

	for _, sniffer := range expectedSniffers {
		s.db.CreateSniffer(&sniffer)
	}

	actualSniffers := s.db.GetSniffers()

	assert.Equal(s.T(), expectedSniffers, actualSniffers)
}

func (s *InMemoryDBSuite) TestUpdateSniffer() {
	sniffers := []model.Sniffer{
		{
			MAC:      "AA:AA:AA:BB:CC:DD",
			Name:     "lab_sniffer",
			Location: "lab",
		},
		{
			MAC:      "AA:AA:DD:AA:BB:CC",
			Name:     "room_sniffer",
			Location: "room",
		},
	}

	for _, sniffer := range sniffers {
		s.db.CreateSniffer(&sniffer)
	}

	snifferUpdate := model.Sniffer{
		MAC:      "AA:AA:AA:BB:CC:DD",
		Name:     "lib_sniffer",
		Location: "library",
	}

	s.db.UpdateSniffer(&snifferUpdate)

	actualSnifferAfterUpdate := s.db.Sniffers[0]

	assert.Equal(s.T(), snifferUpdate, actualSnifferAfterUpdate)
}
