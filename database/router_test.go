package database

import (
	"github.com/cyucelen/wirect/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestCreateRouter() {
	router := model.Router{MAC: "00:11:22:33:44:55", SSID: "2020", SnifferMAC: "11:11:22:22:33:33"}
	s.db.CreateRouter(&router)

	var actualRouter model.Router
	s.db.DB.Last(&actualRouter)

	assert.Equal(s.T(), router, actualRouter)
}

func (s *DatabaseSuite) TestGetRoutersBySniffer() {
	snifferMAC := "00:00:00:00:00:00"
	routers := []model.Router{
		{MAC: "11:22:33:44:55:66", SSID: "2020", SnifferMAC: snifferMAC},
		{MAC: "22:33:44:55:66:77", SSID: "1010", SnifferMAC: snifferMAC},
		{MAC: "AA:BB:CC:DD:EE:FF", SSID: "Arch", SnifferMAC: snifferMAC},
		{MAC: "FF:AA:BB:FF:CC:DD", SSID: "dont h@ck m3", SnifferMAC: snifferMAC},
	}

	for _, router := range routers {
		s.db.CreateRouter(&router)
	}

	actualRouters := s.db.GetRoutersBySniffer(snifferMAC)
	assert.Equal(s.T(), routers, actualRouters)
}
