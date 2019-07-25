package main

import (
	"github.com/sirupsen/logrus"
	"goOrderAPI/logger"
	"goOrderAPI/models"
	"sync"
)

var once sync.Once

func main() {
	once.Do(
		func() {
			logger.SetUpLogging()
		})
	defer models.GetDB().Close()
	err := NewHTTPServer()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Error Starting server")
	}
}
