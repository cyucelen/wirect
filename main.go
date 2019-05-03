package main

import (
	"net/http"

	"github.com/cyucelen/wirect/database"
	"github.com/cyucelen/wirect/delivery/http"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/middleware"
)

func main() {
	db, err := database.New("sqlite3", "./wirect.db")

	if err != nil {
		panic(err)
	}

	e := server.Create(db)

	e.Use(middleware.Logger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	if err := e.Start(":1323"); err != nil {
		panic(err)
	}
}
