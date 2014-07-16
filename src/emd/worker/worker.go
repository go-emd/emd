package worker

import (
	"emd/connector"
	"emd/core"
)

/*
 * Perform's computation's and raw
 * processing.
 */
type Worker interface {
	Init()
	Run()
	Ports() map[string]connector.Connector
	Name() string
}

/*
 *
 * Work - Worker implementations.
 *
 */
type Work struct {
	core.Core
	Ports_ map[string]connector.Connector
}

func (w Work) Ports() map[string]connector.Connector {
	return w.Ports_
}

func (w Work) Name() string {
	return w.Name_
}
