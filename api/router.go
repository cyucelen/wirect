package api

import (
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

	for _, router := range routers {
		internalRouter := toInternalRouter(snifferMAC, &router)
		if err := r.DB.CreateRouter(internalRouter); err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return err
		}
	}

	ctx.JSON(http.StatusCreated, routers)
	return nil
}

func (r *RouterAPI) GetRouters(ctx echo.Context) error {
	snifferMAC, err := getSnifferMAC(ctx)
	if err != nil {
		return err
	}

	routers := r.DB.GetRoutersBySniffer(snifferMAC)
	externalRouters := []model.RouterExternal{}

	for _, router := range routers {
		externalRouters = append(externalRouters, *toExternal(&router))
	}

	ctx.JSON(http.StatusOK, externalRouters)
	return nil
}

func toInternalRouter(snifferMAC string, router *model.RouterExternal) *model.Router {
	return &model.Router{
		SSID:       router.SSID,
		SnifferMAC: snifferMAC,
		LastSeen:   router.LastSeen,
	}
}

func toExternal(router *model.Router) *model.RouterExternal {
	return &model.RouterExternal{
		SSID:     router.SSID,
		LastSeen: router.LastSeen,
	}
}
