package database

import (
	"github.com/cyucelen/wirect/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestCreateRouter() {
	router := model.Router{SSID: "2020", SnifferMAC: "11:11:22:22:33:33", LastSeen: 1000}
	s.db.CreateRouter(&router)

	var actualRouter model.Router
	s.db.DB.Last(&actualRouter)

	assert.Equal(s.T(), router, actualRouter)
}

func (s *DatabaseSuite) TestGetRoutersBySniffer() {
	snifferMAC := "00:00:00:00:00:00"
	routers := []model.Router{
		{SSID: "2020", SnifferMAC: snifferMAC},
		{SSID: "1010", SnifferMAC: snifferMAC},
		{SSID: "Arch", SnifferMAC: snifferMAC},
		{SSID: "dont h@ck m3", SnifferMAC: snifferMAC},
	}

	for _, router := range routers {
		s.db.CreateRouter(&router)
	}

	actualRouters := s.db.GetRoutersBySniffer(snifferMAC)
	assert.Equal(s.T(), routers, actualRouters)
}
