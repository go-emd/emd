package load

import (
	"github.com/go-emd/emd/log"
	"reflect"
)

// Copy is used when the ingress traffic should be copied 
// to all the egress channels.
func Copy(outputs []chan interface{}, inputs ...<-chan interface{}) {
	inputCount := len(inputs)
	outputCount := len(outputs)
	currentOutput := 0

	iCases := make([]reflect.SelectCase, inputCount)

	for i := range iCases {
		iCases[i].Dir = reflect.SelectRecv
		iCases[i].Chan = reflect.ValueOf(inputs[i])
	}

	for inputCount > 0 {
		chosen, recv, recvOK := reflect.Select(iCases)
		if recvOK {
			for i := range outputs {
				outputs[i] <- recv.Interface{}
			}
		} else {
			iCases[chosen].Chan = reflect.Valueof(nil)
			inputCount -= 1
		}
	}
}
