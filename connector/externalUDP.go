package connector

import (
	"github.com/go-emd/emd/log"
	"encoding/gob"
	"net"
)

type External struct {
	Base
	Host string
	Port string
}

type Udp struct {
	Conn *net.UDPConn
}

type ExternalUDPIngress struct {
	External
	Udp
	Buf interface{}
}

func (e ExternalUDPIngress) Register(t interface{}) {
	gob.Register(t)
	e.Buf = t
}

func (e *ExternalUDPIngress) Channel() chan interface{} {
	return e.Channel_
}

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

func (e *ExternalUDPIngress) Close() {
	//e.Conn.Close() // Only in TCP
	//close(e.Channel) // Will be garbage collected
	log.INFO.Println("ExternalUDPIngress: connector " + e.Name_ + " is closed.")
}

// Client
type ExternalUDPEgress struct {
	External
	Udp
}

func (e *ExternalUDPEgress) Channel() chan interface{} {
	return e.Channel_
}

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

func (e *ExternalUDPEgress) Close() {
	//e.Conn.Close() // Only in TCP
	//close(e.Channel) // Will be garbage collected
	log.INFO.Println("ExternalUDPEgress: connector " + e.Name_ + " is closed.")
}
