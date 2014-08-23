package connector

import (
	"github.com/go-emd/emd/log"
	"encoding/gob"
	"net"
)

// This must be inherited by all other 
// connector implementations that will 
// not force the workers to be on the 
// same node in the distribution.
//
// It will hold extra important information 
// that all external connector implementations 
// require such as the host and port that that 
// the other worker is on.
type External struct {
	Base
	Host string
	Port string
}

// Simply hold UDP specific information in order 
// to maintain a UDP connection.
type Udp struct {
	Conn *net.UDPConn
}

// Inherits the connector.External struct and 
// connector.External struct and adding the Buf 
// of type interface {}.  The Buf is very important 
// because the ExternalUDP connections use encoding/gob 
// to communicate reliably and with minimal bytes transferred.
// 
// The Buf interface is registered using the connector.ExternalUDPIngress.Register 
// function which will register the type that needs to be decoded with gob 
// allowing it to be deserialized and received over the wire.
type ExternalUDPIngress struct {
	External
	Udp
	Buf interface{}
}

// Registers the type to be decoded from the connector.ExternalUDPEgress data 
// that was serialized.
func (e ExternalUDPIngress) Register(t interface{}) {
	gob.Register(t)
	e.Buf = t
}

// Returns the underlying channel to read from.
func (e *ExternalUDPIngress) Channel() chan interface{} {
	return e.Channel_
}

// Opens the specified port to listen on for incoming gob 
// encoded data.
func (e *ExternalUDPIngress) Open() {
	addr, err := net.ResolveUDPAddr("udp", ":"+e.Port)
	if err != nil {
		log.ERROR.Println(err)
	}

	e.Conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.ERROR.Println(err)
	}

	go func(channel chan<- interface{}) {
		defer e.Conn.Close()
		decoder := gob.NewDecoder(e.Conn)

		for {
			err := decoder.Decode(&e.Buf)
			if err != nil {
				log.ERROR.Println(err)
			}

			channel <- e.Buf
		}
	}(e.Channel_)

	log.INFO.Println("ExternalUDPIngress: connector " + e.Name_ + " is opened.")
}

// Closes the UDP port being listened on.
func (e *ExternalUDPIngress) Close() {
	//e.Conn.Close() // Only in TCP
	//close(e.Channel) // Will be garbage collected
	log.INFO.Println("ExternalUDPIngress: connector " + e.Name_ + " is closed.")
}

// Client

// The base constructor of the ExternalUDPEgress connector implementation.  
// It's purpose is to send gob encoded data to the specified host:port.
type ExternalUDPEgress struct {
	External
	Udp
}

// Returns the base channel used under the hood.
func (e *ExternalUDPEgress) Channel() chan interface{} {
	return e.Channel_
}

// Sets up a connection to the specified host:port 
// and begins forwarding gob encoded data to it.
func (e *ExternalUDPEgress) Open() {
	addr, err := net.ResolveUDPAddr("udp", e.Host+":"+e.Port)
	if err != nil {
		log.ERROR.Println(err)
	}

	e.Conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.ERROR.Println(err)
	}

	go func(channel <-chan interface{}) {
		defer e.Conn.Close()
		encoder := gob.NewEncoder(e.Conn)

		for {
			select {
			case data := <-channel:
				err := encoder.Encode(&data)
				if err != nil {
					log.ERROR.Println(err)
				}
			}
		}
	}(e.Channel_)

	log.INFO.Println("ExternalUDPEgress: connector " + e.Name_ + " is opened.")
}

// Closes the host:port connection that was created.
func (e *ExternalUDPEgress) Close() {
	//e.Conn.Close() // Only in TCP
	//close(e.Channel) // Will be garbage collected
	log.INFO.Println("ExternalUDPEgress: connector " + e.Name_ + " is closed.")
}
