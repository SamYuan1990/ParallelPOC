package pipelinepoc

import log "github.com/sirupsen/logrus"

var instance *log.Logger

func GetLogger() *log.Logger {
	if instance == nil {
		instance = log.New()
		instance.SetLevel(log.TraceLevel)
	}
	return instance
}
