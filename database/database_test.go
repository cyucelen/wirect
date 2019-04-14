package database

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/cyucelen/wirect/model"
)

type DatabaseSuite struct {
	suite.Suite
	db         *GormDatabase
	testDBPath string
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

func (s *DatabaseSuite) BeforeTest(suiteName, testName string) {
	s.testDBPath = "./testDBs/test_suite.db"
	db, err := New("sqlite3", s.testDBPath, true)
	os.Chmod(s.testDBPath, 777)
	assert.Nil(s.T(), err)
	s.db = db
}

func (s *DatabaseSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	os.RemoveAll(filepath.Dir(s.testDBPath))
}

func (s *DatabaseSuite) TestNewCreatesDefaultSniffer() {
	snifferCount := 0
	s.db.DB.Find(new(model.Sniffer)).Count(&snifferCount)
	s.NotZero(snifferCount)
}

func TestInvalidDialect(t *testing.T) {
	path := "./testDBs/test.db"

	_, err := New("what is this dialect", path, true)
	assert.Error(t, err)
}
