/*
	Load allows for the ability to forward data between 
	multiple input channels to multiple output channels.  
	There are a couple options on how to forward the input 
	data to the output channels.
	
	All the implementations of the handler func for forwarding 
	the tuples between the input and output channels must take 
	the parameters of the load.Kind type.
*/
package load

import (
	"github.com/go-emd/emd/log"
	"os"
)

// Function type that must be passed into the NtoN function call.  
// This can either be the RoundRobin function, Copy function or a 
// custom function in which the leader.template file will need 
// to be edited within the emd distribution.
type Kind func(outputs []chan interface{}, inputs []chan interface{})

// Allows N connections' channels to have ingress traffic and will 
// forward that traffic to the egress array of N connectors' channels.  
// There are multiple types of forwarding that can be used such as 
// round robin, and copy.
func NtoN(handler Kind, outputs []chan interface{}, inputs []chan interface{}) {
	if len(inputs) == 0 || len(outputs) == 0 {
		log.ERROR.Println("NtoN requires at least one input and one output")
		os.Exit(1)
	}

	go handler(outputs, inputs)
}
