package connector

import (
	"testing"
	"time"
	"github.com/go-emd/emd/core"
	"github.com/go-emd/emd/log"
	"io/ioutil"
)

func TestLocal(t *testing.T) {
	log.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)

	mylocal := &Local{
		Base{
			core.Core{"Test"},
			make(chan interface{}, 0),
		},
	}

	// Sender
	go func() {
		mylocal.Open()
		mylocal.Channel() <- "HEY"
		mylocal.Close()
	}()

	// Receiver
	go func() {
		mylocal.Open()

		select {
		case <- time.After(time.Second * 2):
			t.Fail()
		case data := <- mylocal.Channel():
			if data != "HEY" {
				t.Fail()
			}
		}

		mylocal.Close()
	}()
}

func TestExternalUDP(t *testing.T) {
	log.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)

	myexternalIngress := &ExternalUDPIngress{
		External{
			Base{
				core.Core{"Test"},
				make(chan interface{}, 0),
			},
			"localhost",
			"60000",
		},
		Udp{nil},
		nil,
	}

	myexternalEgress := &ExternalUDPEgress{
		External{
			Base{
				core.Core{"Test"},
				make(chan interface{}, 0),
			},
			"localhost",
			"60000",
		},
		Udp{nil},
	}

	// Sender
	go func() {
		myexternalEgress.Open()
		myexternalEgress.Channel() <- "HEY"
		myexternalEgress.Close()
	}()

	// Receiver
	go func() {
		myexternalIngress.Open()

		select {
		case <- time.After(time.Second * 2):
			t.Fail()
		case data := <- myexternalIngress.Channel():
			if data != "HEY" {
				t.Fail()
			}
		}

		myexternalIngress.Close()
	}()
}
