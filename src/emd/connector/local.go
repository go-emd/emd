package connector

import (
	"emd/log"
)

type Local struct {
	Base
}

func (l *Local) Open() {
	log.INFO.Println("Local: " + l.Name + " is opened.")
}

func (l *Local) Close() {
	log.INFO.Println("Local: " + l.Name + " is closed.")
}

func (l *Local) Channel() chan interface{} {
	return l.Channel_
}
