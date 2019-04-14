package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
	"github.com/cyucelen/wirect/model"
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
		Timestamp:  time.Now(),
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
			Timestamp:  time.Now(),
		},
		{
			MAC:        "22:22:12:11:22:22",
			RSSI:       1.23,
			SnifferMAC: snifferMAC,
			Timestamp:  time.Now(),
		},
		{
			MAC:        "00:00:11:11:22:22",
			RSSI:       1.23,
			SnifferMAC: "11:11:11:11:11:11",
			Timestamp:  time.Now(),
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
			Timestamp:  since.Add(-10 * time.Minute),
			RSSI:       123,
			SnifferMAC: snifferOne,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  since.Add(1 * time.Second),
			RSSI:       333,
			SnifferMAC: snifferOne,
		},
		{
			MAC:        "00:11:CC:CC:44:55",
			Timestamp:  since.Add(-5 * time.Minute),
			RSSI:       1234,
			SnifferMAC: snifferTwo,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  since,
			RSSI:       333,
			SnifferMAC: snifferOne,
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  since.Add(5 * time.Minute),
			RSSI:       333,
			SnifferMAC: snifferOne,
		},
	}
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	snifferPacketsSince := s.db.GetPacketsBySnifferSince(snifferOne, since)

	assert.Len(s.T(), snifferPacketsSince, 3)
	assert.Equal(s.T(), packets[3].Timestamp.Unix(), snifferPacketsSince[0].Timestamp.Unix())
	assert.Equal(s.T(), packets[1].Timestamp.Unix(), snifferPacketsSince[1].Timestamp.Unix())
	assert.Equal(s.T(), packets[4].Timestamp.Unix(), snifferPacketsSince[2].Timestamp.Unix())
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
