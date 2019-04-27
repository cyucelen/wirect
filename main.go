package main

import (
	"github.com/cyucelen/wirect/database"
	"github.com/cyucelen/wirect/router"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/middleware"
)

func main() {
	db, err := database.New("sqlite3", "./wirect.db", true)

	if err != nil {
		panic(err)
	}

	e := router.Create(db)

	e.Use(middleware.Logger())

	if err := e.Start(":1323"); err != nil {
		panic(err)
	}
}
