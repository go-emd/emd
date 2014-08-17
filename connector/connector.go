package connector

import (
	"github.com/go-emd/emd/core"
)

type Connector interface {
	Open()
	Close()
	Channel() chan interface{} //chan []byte
}

type Base struct {
	core.Core
	Channel_ chan interface{} //chan []byte
}
