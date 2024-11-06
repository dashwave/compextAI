package logger

import "github.com/sirupsen/logrus"

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	// logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
}

func GetLogger() *logrus.Logger {
	return logger
}
