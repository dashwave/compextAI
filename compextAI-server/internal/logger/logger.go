package logger

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	// logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
}

func GetLogger() *logrus.Logger {
	return logger
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request: ", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
