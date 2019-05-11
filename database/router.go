package database

import "github.com/cyucelen/wirect/model"

func (g *GormDatabase) CreateRouter(router *model.Router) error {
	return g.DB.Create(router).Error
}

func (g *GormDatabase) GetRoutersBySniffer(snifferMAC string) []model.Router {
	var routers []model.Router
	g.DB.Where("sniffer_mac = ?", snifferMAC).Find(&routers)
	return routers
}
