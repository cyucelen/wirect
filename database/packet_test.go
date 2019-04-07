package database

import (
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wirect/wirect-server/model"
)

func (s *DatabaseSuite) TestCreatePacket() {
	packets := []model.Packet{
		{
			MAC:        "00:11:22:33:44:55",
			Timestamp:  time.Now(),
			RSSI:       1234,
			SnifferMAC: "00:00:00:00:00:00",
		},
		{
			MAC:        "00:33:22:11:44:55",
			Timestamp:  time.Now().Add(10 * time.Second),
			RSSI:       333,
			SnifferMAC: "00:00:00:00:00:00",
		},
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
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  time.Now(),
			RSSI:       123,
			SnifferMAC: "01:02:03:04:05:06",
		},
		{
			MAC:        "00:11:CC:CC:44:55",
			Timestamp:  time.Now().Add(5 * time.Second),
			RSSI:       1234,
			SnifferMAC: "00:00:00:00:00:00",
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  time.Now().Add(10 * time.Second),
			RSSI:       333,
			SnifferMAC: "01:02:03:04:05:06",
		},
		{
			MAC:        "AA:BB:22:11:44:55",
			Timestamp:  time.Now().Add(20 * time.Second),
			RSSI:       333,
			SnifferMAC: "01:02:03:04:05:06",
		},
	}
	for _, packet := range packets {
		s.db.CreatePacket(&packet)
	}

	snifferPackets := s.db.GetPacketsBySniffer("01:02:03:04:05:06")

	assert.Len(s.T(), snifferPackets, 3)
	assert.Equal(s.T(), packets[0].Timestamp.Unix(), snifferPackets[0].Timestamp.Unix())
	assert.Equal(s.T(), packets[2].Timestamp.Unix(), snifferPackets[1].Timestamp.Unix())
	assert.Equal(s.T(), packets[3].Timestamp.Unix(), snifferPackets[2].Timestamp.Unix())
}

func (s *DatabaseSuite) TestGetPacketsBySnifferSince() {
	snifferOne := "01:02:03:04:05:06"
	snifferTwo := "00:00:00:00:00:00"
	createTwoSniffers(s, snifferOne, snifferTwo)

	since := time.Now().Add(10 * time.Minute)

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

func createTwoSniffers(s *DatabaseSuite, s1, s2 string) []model.Sniffer {
	s.db.DB.Delete(model.Sniffer{})

	sniffers := []model.Sniffer{
		{
			MAC:      s1,
			Name:     "library_sniffer",
			Location: "library",
		},
		{
			MAC:      s2,
			Name:     "copy_center_sniffer",
			Location: "copy_center",
		},
	}
	for _, sniffer := range sniffers {
		s.db.DB.Create(&sniffer)
	}
	return sniffers
}
