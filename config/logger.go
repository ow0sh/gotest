package config

import (
	"github.com/sirupsen/logrus"
)

type logger struct {
	log *logrus.Logger
}

func (log *logger) Log() *logrus.Logger {
	return logrus.New()
}
