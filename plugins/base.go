package main

import "github.com/sirupsen/logrus"

type BasePlugin struct {
}

func (b *BasePlugin) Logger() *logrus.Logger {
	return logrus.New()
}
