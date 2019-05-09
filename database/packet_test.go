package database

import (
	"time"

	"github.com/benbjohnson/clock"

	"github.com/cyucelen/wirect/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestCreatePacket() {
	packets := []model.Packet{
		{MAC: "00:11:22:33:44:55", Timestamp: time.Now().UTC().Unix(), RSSI: 1234, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "00:33:22:11:44:55", Timestamp: time.Now().UTC().Add(10 * time.Second).Unix(), RSSI: 333, SnifferMAC: "00:00:00:00:00:00"},
	}

	var lastPacketInDB model.Packet
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
		s.db.DB.Last(&lastPacketInDB)
		assert.ObjectsAreEqualValues(packet, lastPacketInDB)
	}

	var packetsInDB []model.Packet
	s.db.DB.Find(&packetsInDB)
	assert.Len(s.T(), packetsInDB, 2)
}

func (s *DatabaseSuite) TestGetPacketsBySniffer() {
	createTwoSniffers(s, "01:02:03:04:05:06", "00:00:00:00:00:00")

	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: time.Now().UTC().Unix(), RSSI: 123, SnifferMAC: "01:02:03:04:05:06"},
		{MAC: "00:11:CC:CC:44:55", Timestamp: time.Now().UTC().Add(5 * time.Second).Unix(), RSSI: 1234, SnifferMAC: "00:00:00:00:00:00"},
		{MAC: "AA:BB:22:11:44:55", Timestamp: time.Now().UTC().Add(10 * time.Second).Unix(), RSSI: 333, SnifferMAC: "01:02:03:04:05:06"},
		{MAC: "AA:BB:22:11:44:55", Timestamp: time.Now().UTC().Add(20 * time.Second).Unix(), RSSI: 333, SnifferMAC: "01:02:03:04:05:06"},
	}
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	snifferPackets := s.db.GetPacketsBySniffer("01:02:03:04:05:06")

	assert.Len(s.T(), snifferPackets, 3)
	assert.Equal(s.T(), packets[0].Timestamp, snifferPackets[0].Timestamp)
	assert.Equal(s.T(), packets[2].Timestamp, snifferPackets[1].Timestamp)
	assert.Equal(s.T(), packets[3].Timestamp, snifferPackets[2].Timestamp)
}

func (s *DatabaseSuite) TestGetPacketsBySnifferSince() {
	snifferOne := "01:02:03:04:05:06"
	snifferTwo := "00:00:00:00:00:00"
	createTwoSniffers(s, snifferOne, snifferTwo)

	mockClock := clock.NewMock()

	now := mockClock.Now().UTC().Add(10 * time.Minute)
	since := now.Add(-20 * time.Minute)

	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-10 * time.Minute).Unix(), RSSI: 123, SnifferMAC: snifferOne},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(1 * time.Second).Unix(), RSSI: 333, SnifferMAC: snifferOne},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(-5 * time.Minute).Unix(), RSSI: 1234, SnifferMAC: snifferTwo},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Unix(), RSSI: 333, SnifferMAC: snifferOne},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(5 * time.Minute).Unix(), RSSI: 333, SnifferMAC: snifferOne},
	}
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	snifferPacketsSince := s.db.GetPacketsBySnifferSince(snifferOne, since.Unix())

	assert.Len(s.T(), snifferPacketsSince, 4)
	assert.Equal(s.T(), packets[0].Timestamp, snifferPacketsSince[0].Timestamp)
	assert.Equal(s.T(), packets[3].Timestamp, snifferPacketsSince[1].Timestamp)
	assert.Equal(s.T(), packets[1].Timestamp, snifferPacketsSince[2].Timestamp)
	assert.Equal(s.T(), packets[4].Timestamp, snifferPacketsSince[3].Timestamp)
}

func (s *DatabaseSuite) TestGetPacketsBySnifferBetween() {
	snifferOne := "01:02:03:04:05:06"
	snifferTwo := "00:00:00:00:00:00"
	createTwoSniffers(s, snifferOne, snifferTwo)

	mockClock := clock.NewMock()
	now := mockClock.Now().UTC().Add(20 * time.Minute)

	from := now.Add(-10 * time.Minute)
	until := now.Add(2 * time.Minute)

	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-10 * time.Minute).Unix(), RSSI: 123, SnifferMAC: snifferOne},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(-5 * time.Minute).Unix(), RSSI: 1234, SnifferMAC: snifferTwo},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(1 * time.Second).Unix(), RSSI: 333, SnifferMAC: snifferOne},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Unix(), RSSI: 333, SnifferMAC: snifferOne},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(5 * time.Minute).Unix(), RSSI: 333, SnifferMAC: snifferOne},
	}
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	snifferPacketsBetween := s.db.GetPacketsBySnifferBetweenDates(snifferOne, from.Unix(), until.Unix())

	assert.Len(s.T(), snifferPacketsBetween, 3)
	assert.Equal(s.T(), packets[0].Timestamp, snifferPacketsBetween[0].Timestamp)
	assert.Equal(s.T(), packets[3].Timestamp, snifferPacketsBetween[1].Timestamp)
	assert.Equal(s.T(), packets[2].Timestamp, snifferPacketsBetween[2].Timestamp)
}

func (s *DatabaseSuite) TestCountUniqueMACAddressesBySnifferBetweenDates() {
	snifferOne := "01:02:03:04:05:06"
	snifferTwo := "00:00:00:00:00:00"
	createTwoSniffers(s, snifferOne, snifferTwo)

	mockClock := clock.NewMock()
	now := mockClock.Now().UTC().Add(20 * time.Minute)

	from := now.Add(-10 * time.Minute)
	until := now.Add(2 * time.Minute)

	packets := []model.Packet{
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(-10 * time.Minute).Unix(), RSSI: 123, SnifferMAC: snifferOne},
		{MAC: "00:11:CC:CC:44:55", Timestamp: now.Add(-5 * time.Minute).Unix(), RSSI: 1234, SnifferMAC: snifferTwo},
		{MAC: "CC:BB:FA:AE:FC:6C", Timestamp: now.Unix(), RSSI: 333, SnifferMAC: snifferOne},
		{MAC: "FF:FB:44:21:64:25", Timestamp: now.Add(1 * time.Second).Unix(), RSSI: 333, SnifferMAC: snifferOne},
		{MAC: "AA:BB:22:11:44:55", Timestamp: now.Add(5 * time.Minute).Unix(), RSSI: 333, SnifferMAC: snifferOne},
	}
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	uniqueMACCount := s.db.GetUniqueMACCountBySnifferBetweenDates(snifferOne, from.Unix(), until.Unix())

	assert.Equal(s.T(), 3, uniqueMACCount)
}

func createTwoSniffers(s *DatabaseSuite, s1, s2 string) []model.Sniffer {
	s.db.DB.Delete(model.Sniffer{})

	sniffers := []model.Sniffer{
		{MAC: s1, Name: "library_sniffer", Location: "library"},
		{MAC: s2, Name: "copy_center_sniffer", Location: "copy_center"},
	}
	for _, sniffer := range sniffers {
		s.db.DB.Create(&sniffer)
	}
	return sniffers
}
