package load

import (
	"testing"
	"github.com/go-emd/emd/log"
	"io/ioutil"
)

func TestCopy(t *testing.T) {
	log.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)

	in := make([]chan interface{}, 1)
	in[0] = make(chan interface{}, 0)

	out := make([]chan interface{}, 2)
	out[0] = make(chan interface{}, 0)
	out[1] = make(chan interface{}, 0)

	NtoN(Copy, out, in)

	in[0] <- "TEST"
	out1 := <- out[0]
	out2 := <- out[1]

	if out1.(string) != "TEST" || out2.(string) != "TEST" {
		t.Fail()
	}

	in[0] <- nil
}

func TestRoundRobin(t *testing.T) {
	log.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)

	in := make([]chan interface{}, 1)
	in[0] = make(chan interface{}, 0)

	out := make([]chan interface{}, 2)
	out[0] = make(chan interface{}, 0)
	out[1] = make(chan interface{}, 0)

	NtoN(RoundRobin, out, in)

	in[0] <- "TEST1"
	out1 := <- out[0]
	
	in[0] <- "TEST2"
	out2 := <- out[1]

	if out1.(string) != "TEST1" || out2.(string) != "TEST2" {
		t.Fail()
	}

	in[0] <- nil
}
