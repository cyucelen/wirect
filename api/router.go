package api

import (
	"errors"
	"net/http"

	"github.com/cyucelen/wirect/model"
	"github.com/labstack/echo"
)

type RouterDatabase interface {
	CreateRouter(router *model.Router) error
	GetRoutersBySniffer(snifferMAC string) []model.Router
}

type RouterAPI struct {
	DB RouterDatabase
}

func (r *RouterAPI) CreateRouters(ctx echo.Context) error {
	var routers []model.RouterExternal
	if err := ctx.Bind(&routers); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return err
	}

	snifferMAC, err := getSnifferMAC(ctx)
	if err != nil {
		return err
	}

	validRouters := filterValidRouters(routers)

	if len(validRouters) == 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return errors.New("")
	}

	for _, router := range validRouters {
		internalRouter := toInternalRouter(snifferMAC, &router)
		if err := r.DB.CreateRouter(internalRouter); err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return err
		}
	}

	ctx.JSON(http.StatusCreated, validRouters)
	return nil
}

func (r *RouterAPI) GetRouters(ctx echo.Context) error {
	snifferMAC, err := getSnifferMAC(ctx)
	if err != nil {
		return err
	}
	routers := r.DB.GetRoutersBySniffer(snifferMAC)

	ctx.JSON(http.StatusOK, routers)
	return nil
}

func filterValidRouters(routers []model.RouterExternal) []model.RouterExternal {
	validRouters := []model.RouterExternal{}
	for _, router := range routers {
		if router.MAC != "" {
			validRouters = append(validRouters, router)
		}
	}

	return validRouters
}

func toInternalRouter(snifferMAC string, router *model.RouterExternal) *model.Router {
	return &model.Router{
		MAC:        router.MAC,
		SSID:       router.SSID,
		SnifferMAC: snifferMAC,
	}
}

func toExternal(router *model.Router) *model.RouterExternal {
	return &model.RouterExternal{
		MAC:  router.MAC,
		SSID: router.SSID,
	}
}
