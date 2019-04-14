package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/cyucelen/wirect/model"
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
	sniffer := new(model.Sniffer)
	if err := ctx.Bind(sniffer); err != nil || !isSnifferValid(sniffer) {
		ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}

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
