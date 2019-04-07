package api

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/wirect/wirect-server/model"
)

type SnifferDatabase interface {
	CreateSniffer(sniffer *model.Sniffer) error
	GetSniffers() []model.Sniffer
	UpdateSniffer(sniffer *model.Sniffer) error
}

type SnifferAPI struct {
	DB SnifferDatabase
}

func (s *SnifferAPI) CreateSniffer(ctx echo.Context) {
	sniffer := new(model.Sniffer)
	if err := ctx.Bind(sniffer); err != nil || sniffer.MAC == "" {
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}
	s.DB.CreateSniffer(sniffer)
	ctx.JSON(http.StatusCreated, sniffer)
}

func (s *SnifferAPI) GetSniffers(ctx echo.Context) {
	ctx.JSON(http.StatusOK, s.DB.GetSniffers())
}

func (s *SnifferAPI) UpdateSniffer(ctx echo.Context) {
	sniffer := new(model.Sniffer)
	if err := ctx.Bind(sniffer); err != nil || sniffer.MAC == "" {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	if err := s.DB.UpdateSniffer(sniffer); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
