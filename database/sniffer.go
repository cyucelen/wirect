package database

import "gitlab.com/wirect/wirect-server/model"

func (g *GormDatabase) CreateSniffer(sniffer *model.Sniffer) error {
	return g.DB.Create(sniffer).Error
}

func (g *GormDatabase) GetSniffers() []model.Sniffer {
	var sniffers []model.Sniffer
	g.DB.Find(&sniffers)

	return sniffers
}

func (g *GormDatabase) UpdateSniffer(sniffer *model.Sniffer) error {
	return g.DB.Save(sniffer).Error
}
