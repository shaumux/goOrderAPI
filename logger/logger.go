package logger

import "github.com/sirupsen/logrus"

var Log = logrus.New()

func SetUpLogging() {
	Log.SetReportCaller(true)
	Log.SetLevel(logrus.ErrorLevel)
	Log.SetFormatter(&logrus.JSONFormatter{})
}
