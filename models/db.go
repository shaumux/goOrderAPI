package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"goOrderAPI/logger"
	"os"
)

var db *gorm.DB

func init() {

	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
	dbHost := os.Getenv("DBHOST")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	logger.Log.WithFields(logrus.Fields{"dbURI": dbUri}).Info("Connecting to database")

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Fatal(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Order{}) //Database migration
}

func GetDB() *gorm.DB {
	return db
}
