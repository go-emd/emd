package leader

import (
	"time"
)

type WorkerCache struct {
	Timestamp time.Time // Leader controlled
	Metric interface{}
	Status string
	Health string
	State string // Leader controlled
}

type Cache struct {
	Workers map[string]WorkerCache
}
