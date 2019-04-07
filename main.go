package main

import (
	"github.com/go-ffmt/ffmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gitlab.com/wirect/wirect-server/database"
	"gitlab.com/wirect/wirect-server/model"
)

// func getSingals(c echo.Context) error {
// 	var signals []model.Signal
// 	db.Find(&signals)

// 	return c.JSON(http.StatusOK, signals)
// }

// func saveSignals(c echo.Context) error {
// 	signals := new([]model.Signal)
// 	if err := c.Bind(signals); err != nil {
// 		return err
// 	}

// 	for _, signal := range *signals {
// 		db.Create(&signal)
// 	}

// 	return c.JSON(http.StatusOK, signals)

// }

func main() {
	db, err := database.New("sqlite3", "wirect.db", true)
	if err != nil {
		panic(err)
	}

	ffmt.Puts(db.DB.Find(new([]model.Sniffer)))

	// e := echo.New()

	// e.GET("/signal", getSingals)
	// e.POST("/signal", saveSignals)

	// e.Start(":3000")
}
