package api

import (
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/cyucelen/wirect/model"
	"github.com/labstack/echo"
)

type TimeAPI struct {
	Clock clock.Clock
}

func (t *TimeAPI) GetTime(ctx echo.Context) error {
	ctx.JSON(http.StatusOK, model.Time{Now: t.Clock.Now().Unix()})
	return nil
}
