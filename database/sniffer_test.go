package database

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/wirect/wirect-server/model"
)

func (s *DatabaseSuite) TestCreateSniffer() {
	snifferMAC := "11:22:33:44:55:66"
	sniffer := model.Sniffer{MAC: snifferMAC, Name: "library_sniffer", Location: "library"}

	s.db.CreateSniffer(&sniffer)

	var sniffers []*model.Sniffer
	s.db.DB.Find(&sniffers)

	assert.Len(s.T(), sniffers, 2)
	assert.Equal(s.T(), sniffers[1].MAC, snifferMAC)
}

func (s *DatabaseSuite) TestGetSniffers() {
	sniffers := []model.Sniffer{
		{MAC: "11:22:33:44:55:66", Name: "library_sniffer", Location: "library"},
		{MAC: "33:44:55:88:99:33", Name: "copy_center_sniffer", Location: "copy_center"},
	}
	for _, sniffer := range sniffers {
		s.db.CreateSniffer(&sniffer)
	}

	sniffersInDB := s.db.GetSniffers() // plus default sniffer
	assert.Len(s.T(), sniffersInDB, 3)
	assert.Equal(s.T(), sniffers, sniffersInDB[1:])
}

func (s *DatabaseSuite) TestUpdateSniffer() {
	snifferToBeUpdatedMAC := "33:44:55:88:99:33"
	sniffers := []model.Sniffer{
		{MAC: "11:22:33:44:55:66", Name: "library_sniffer", Location: "library"},
		{MAC: snifferToBeUpdatedMAC, Name: "copy_center_sniffer", Location: "copy_center"},
	}
	for _, sniffer := range sniffers {
		s.db.CreateSniffer(&sniffer)
	}

	snifferUpdate := model.Sniffer{MAC: snifferToBeUpdatedMAC, Name: "room_sniffer", Location: "room"}

	err := s.db.UpdateSniffer(&snifferUpdate)
	assert.Nil(s.T(), err)

	sniffersInDB := s.db.GetSniffers()

	for _, sniffer := range sniffersInDB {
		if sniffer.MAC == snifferToBeUpdatedMAC {
			assert.Equal(s.T(), snifferUpdate, sniffer)
		}
	}

	assert.Len(s.T(), sniffersInDB, 3)
}
