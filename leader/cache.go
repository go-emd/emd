package leader

import (
	"time"
)

// Structure of each worker's cache.
type WorkerCache struct {
	Timestamp time.Time // Leader controlled
	Metric interface{}
	Status string
	Health string
	State string // Leader controlled
}

// Structure of the cache that the leader maintains 
// and sends back in response to a REST endpoint 
// request of cache.
type Cache struct {
	Workers map[string]WorkerCache
}
