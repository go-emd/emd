/*
	The worker package is the most important as its 
	what each component in the emd distribution will 
	need to implement.  Each work should do a specific 
	task and send its output to another worker and so 
	on and so forth.  This allows emd to take advantage 
	of multicore machines and multimachine distributions 
	giving embarrassingly distributed its name.
*/
package worker

import (
	"github.com/go-emd/emd/connector"
	"github.com/go-emd/emd/core"
)

// Perform's computation's and raw processing.
type Worker interface {
	Init()
	Run()
	Ports() map[string]connector.Connector
	Name() string
}

// The base worker.Work structure needs to be inherited by 
// each worker implementation.  It contains the workers core.Core 
// (name string) and all the ports pertaining to it including the 
// MDGMT_<worker name> port which allows the workers to speak to 
// its node leader.
type Work struct {
	core.Core
	Ports_ map[string]connector.Connector
}

// Returns the ports that the worker contains.
func (w Work) Ports() map[string]connector.Connector {
	return w.Ports_
}

// Returns the name that the worker is called (or its alias).
func (w Work) Name() string {
	return w.Name_
}
