package api

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/cyucelen/wirect/model"
	"github.com/labstack/echo"
)

type SnifferDatabase interface {
	CreateSniffer(sniffer *model.Sniffer) error
	GetSniffers() []model.Sniffer
	UpdateSniffer(sniffer *model.Sniffer) error
}

type SnifferAPI struct {
	DB SnifferDatabase
}

func (s *SnifferAPI) CreateSniffer(ctx echo.Context) error {
	sniffer := new(model.Sniffer)

	if err := ctx.Bind(sniffer); err != nil || !isSnifferValid(sniffer) {
		ctx.JSON(http.StatusBadRequest, nil)
		return err
	}

	if err := s.DB.CreateSniffer(sniffer); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return err
	}
	ctx.JSON(http.StatusCreated, sniffer)
	return nil
}

func (s *SnifferAPI) GetSniffers(ctx echo.Context) error {
	ctx.JSON(http.StatusOK, s.DB.GetSniffers())
	return nil
}

func (s *SnifferAPI) UpdateSniffer(ctx echo.Context) error {
	snifferMAC, _ := url.QueryUnescape(ctx.Param("snifferMAC"))
	sniffer := new(model.Sniffer)
	if err := ctx.Bind(sniffer); err != nil || snifferMAC == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}

	sniffer.MAC = snifferMAC

	if err := s.DB.UpdateSniffer(sniffer); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return errors.New("")
	}
	ctx.JSON(http.StatusOK, nil)
	return nil
}

func isSnifferValid(sniffer *model.Sniffer) bool {
	return sniffer.MAC != ""
}
