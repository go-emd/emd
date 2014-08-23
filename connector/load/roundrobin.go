package load

import (
	"github.com/go-emd/emd/log"
	"reflect"
)

// RoundRobin is used when the ingress traffic should be 
// dispersed between the egress channels.
func RoundRobin(outputs []chan interface{}, inputs ...<-chan interface{}) {
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
			if currentOutput > outputCount - 1 { currentOutput = 0 }

			outputs[currentOutput] <- recv.Interface{}
			currentOutput += 1
		} else {
			iCases[chosen].Chan = reflect.Valueof(nil)
			inputCount -= 1
		}
	}
}
