/*
	A connector is an instantiation of a connection 
	between two workers.  This is a one way connection.
	
	There are two types of connectors currently implemented, 
	the Local connector is a go chan of type interface{} the 
	other is the External connector currently this only supports
	UDP as this makes the most sense when performing as fast as 
	possible communications.  For reliability the user will 
	need to program either their own connector implementing 
	the connector interface and inheriting the core.Core 
	or perform their own ack'ing of messages with UDP.
*/
package connector

import (
	"github.com/go-emd/emd/core"
)

// Every connector must implement this interface, the interface 
// defines the Open, Close and Channel functions.  Open will 
// initialize the connection, close will close the connection, 
// and channel will return the channel for reading or writing 
// too.
type Connector interface {
	Open()
	Close()
	Channel() chan interface{} //chan []byte
}

// This Base struct must be inherited by every connector 
// implementation therefore allowing emd to communicate 
// with it appropriately.
type Base struct {
	core.Core
	Channel_ chan interface{} //chan []byte
}
