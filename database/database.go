package database

import (
	"os"
	"path/filepath"

	"github.com/cyucelen/wirect/model"
	"github.com/jinzhu/gorm"
)

var mkdirAllFunc = os.MkdirAll

// GormDatabase is a wrapper for the gorm framework
type GormDatabase struct {
	DB *gorm.DB
}

// New creates a new wrapper for the gorm database framework
func New(dialect, connection string, createDefaultSnifferIfNotExist bool) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)
	db, err := gorm.Open(dialect, connection)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxOpenConns(1) // sqlite cannot handle concurrent writes
	db.AutoMigrate(&model.Packet{}, &model.Router{}, &model.Sniffer{})

	snifferCount := 0
	db.Find(new(model.Sniffer)).Count(&snifferCount)

	if createDefaultSnifferIfNotExist && snifferCount == 0 {
		db.Create(&model.Sniffer{MAC: "00:00:00:00:00:00", Name: "default", Location: "default"})
	}

	return &GormDatabase{DB: db}, nil
}

// Close closes the database connection
func (d *GormDatabase) Close() {
	d.DB.Close()
}

func createDirectoryIfSqlite(dialect string, connection string) {
	if dialect == "sqlite3" {
		if _, err := os.Stat(filepath.Dir(connection)); os.IsNotExist(err) {
			if err := mkdirAllFunc(filepath.Dir(connection), 0777); err != nil {
				panic(err)
			}
		}
	}
}
