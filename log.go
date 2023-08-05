package slacker

import (
	"log"
	"os"
)

func (s *Slacker) CustomLogger(logger SlackLogger) {
	s.logger = logger
}

type SlackLogger interface {
	Printf(format string, v ...interface{})
	Output(calldepth int, s string) error
}

type Eventer interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
}

var defaultLogger = log.New(os.Stderr, "Slacker: ", log.LstdFlags)

func (s *Slacker) logf(format string, v ...interface{}) {
	if eventer, ok := s.logger.(Eventer); ok {
		eventer.Infof(format, v...)
		return
	}
	s.logger.Printf(format, v...)
}

func (s *Slacker) debugf(format string, v ...interface{}) {
	if !s.debug {
		return
	}
	if eventer, ok := s.logger.(Eventer); ok {
		eventer.Debugf(format, v...)
		return
	}
	format = "DEBUG: " + format
	s.logger.Printf(format, v...)
}
