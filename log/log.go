package log

import "github.com/sirupsen/logrus"

var log *logrus.Logger

func Get() *logrus.Logger {
	if log == nil {
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		return logger
	}
	return log
}
