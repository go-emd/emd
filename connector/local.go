package connector

import (
	"github.com/go-emd/emd/log"
)

// The most basic implementation of a connector, 
// this turns into just a go chan of type 
// interface{}.  It allows only one way communication 
// in order to keep all connectors in sync.
type Local struct {
	Base
}

// For a chan the Open method is useless since the 
// chan is already ready to go.  But this is nice 
// for logging the sequential life of the connector.
func (l *Local) Open() {
	log.INFO.Println("Local: " + l.Name_ + " is opened.")
}

// For a chan the close method is useful but problem is 
// which side of the communication should close the chan.  
// Therefore we rely on garbage collection to perform 
// these necessary actions.
func (l *Local) Close() {
	log.INFO.Println("Local: " + l.Name_ + " is closed.")
}

// Returns the chan interface{} that is in the underlying 
// inherited connector.Base class.
func (l *Local) Channel() chan interface{} {
	return l.Channel_
}
