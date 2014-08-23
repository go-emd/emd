/*
	A leader is a process that run's in the background 
	of a node in an emd distribution.  There can only be 
	one node leader per node per distribution.  The node 
	leader contains all the connections, workers, REST 
	endpoints, and monitors all of them.
 */
package leader

import (
	"github.com/go-emd/emd/connector"
	"github.com/go-emd/emd/core"
	"github.com/go-emd/emd/log"
	"github.com/go-emd/emd/worker"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
	"time"
)

// The cache.Cache variable keeps a constant rolling cache 
// of each job/connection's metrics, status, state, and when 
// the last time was it was updated.
var (
	cache *Cache
)

// Every leader must implement the leader.Leader interface 
// allowing it to initialize, run, exit and handle REST 
// requests.
type Leader interface {
	Init()
	Run()
	Exit()

	Start(http.ResponseWriter, *http.Request)
	Stop(http.ResponseWriter, *http.Request)
	Status(http.ResponseWriter, *http.Request)
	Metrics(http.ResponseWriter, *http.Request)
	Cache(http.ResponseWriter, *http.Request)
	Config(http.ResponseWriter, *http.Request)
}

// Each leader implementation must inherit the leader.Lead 
// struct to get information such as its name, the REST 
// ports to listen on, configuration file path, the workers 
// its supposed to monitor and maintain and the management 
// connections (or ports) the leader has with each worker.
type Lead struct {
	core.Core
	GUI_port string
	ConfigPath string
	Workers  []worker.Worker
	Ports    map[string]connector.Connector
}

// Initializes the leader and each of its workers, 
// creates a new cache and initializes each entry.
func (l *Lead) Init() {
	for _, w := range l.Workers {
		w.Init()
	}

	cache = new(Cache)
	cache.Workers = make(map[string]WorkerCache, len(l.Ports))
	for k, _ := range l.Ports {
		tmp := cache.Workers[k]
		tmp.Timestamp = time.Now()
		tmp.Health = "Unknown"
		tmp.State = "Initialized"
		tmp.Status = "Unknown"
		tmp.Metric = nil

		cache.Workers[k] = tmp
	}

	log.INFO.Println("Leader: " + l.Name_ + " is initialized.")
}

// Starts each worker in its own separate go routine and 
// spins up the REST server to handle monitoring and metrics 
// requests.
func (l *Lead) Run() {
	log.INFO.Println("Leader: " + l.Name_ + " is running...")

	// Start all the workers
	for _, w := range l.Workers {
		go w.Run()
		tmp := cache.Workers[w.Name()]
		tmp.State = "Running"
		tmp.Timestamp = time.Now()
		cache.Workers[w.Name()] = tmp
	}

	// Handle rest calls and continue managing nodes
	//   workers.
	http.HandleFunc("/start", l.Start)
	http.HandleFunc("/stop", l.Stop)
	http.HandleFunc("/status", l.Status)
	http.HandleFunc("/metrics", l.Metrics)
	http.HandleFunc("/cache", l.Cache)
	http.HandleFunc("/config", l.Config)

	http.ListenAndServe(":"+l.GUI_port, nil)
}

// A REST endpoint that will handle the status request and 
// respond with the health of each worker in the node.
func (l *Lead) Start(rw http.ResponseWriter, r *http.Request) {
	if allWorkersStopped() {
		log.INFO.Println("Leader: " + l.Name_ + " is starting it' workers...")

		for _, w := range l.Workers {
			w.Init()
			go w.Run()

			tmp := cache.Workers[w.Name()]
			tmp.State = "Running"
			tmp.Timestamp = time.Now()
			cache.Workers[w.Name()] = tmp
		}

		Respond(rw, true, "Workers started :-)")
	} else {
		log.INFO.Println("Leader: " + l.Name_ + " workers are already running.")
		Respond(rw, false, "Workers already started.")
	}

	return
}

// A REST endpoint that handles the stop request.  It will 
// stop all the workers in the node and when a second stop 
// request happens the node leader will exit if all the workers 
// are already stopped.
func (l *Lead) Stop(rw http.ResponseWriter, r *http.Request) {
	if !allWorkersStopped() {
		log.INFO.Println("Leader: " + l.Name_ + " is stopping...")

		for k, v := range l.Ports {
			log.INFO.Println("Worker: " + k + "is stopping...")
			
			tmp := cache.Workers[k]
			
			if writeChannel(v.Channel(), "STOP"){
				tmp.State = "Stopped"
			} else {
				log.WARNING.Println("Unable to stop worker " + k)
				tmp.State = "Unknown"
			}

			tmp.Timestamp = time.Now()
			cache.Workers[k] = tmp
		}
	} else {
		// Response won't return since the server is being shutdown.
		//Respond(rw, true, "Leader stopped :-(")
		l.Exit()
	}

	Respond(rw, true, "Workers stopped :-(")
	return
}

// A REST endpoint that handles the status request, it will 
// return each workers status depending on if the timeout 
// of two seconds happens then it will send the Unknown 
// status.
func (l *Lead) Status(rw http.ResponseWriter, r *http.Request) {
	for k, v := range l.Ports {
		tmp := cache.Workers[k]
		var status interface{}

		if writeChannel(v.Channel(), "STATUS") {
			if status = readChannel(v.Channel()); status == nil {
				log.WARNING.Println("Unable to retrieve status of " + k)

				tmp.Health = "Unknown"
				tmp.Timestamp = time.Now()
				tmp.State = "Unknown"
				cache.Workers[k] = tmp

				Respond(rw, false, "Unknown")
				return
			} else {
				log.INFO.Println("Received status from " + k)

				if status == "Unhealthy" {
					tmp.Health = "Unhealthy"
					tmp.Timestamp = time.Now()
					cache.Workers[k] = tmp

					Respond(rw, true, "Unhealthy")
					return
				} else {
					tmp.Health = "Healthy"
					tmp.Timestamp = time.Now()
					cache.Workers[k] = tmp
				}
			}
		} else {
			log.WARNING.Println("Unable to retrieve status of " + k)

			tmp.Health = "Unknown"
			tmp.Timestamp = time.Now()
			tmp.State = "Unknown"
			cache.Workers[k] = tmp

			Respond(rw, false, "Unknown")
			return
		}
	}

	Respond(rw, true, "Healthy")
	return
}

// A REST endpoint that handles the metrics request.  It will 
// return a json serialized structure of a map containing an 
// interface.  All the metrics are specified by each worker 
// separately.
func (l *Lead) Metrics(rw http.ResponseWriter, r *http.Request) {
	metrics := make(map[string]interface{})

	for k, v := range l.Ports {
		if writeChannel(v.Channel(), "METRICS") {
			if metrics[k] = readChannel(v.Channel()); metrics[k] != nil {
				tmp := cache.Workers[k]
				tmp.Metric = metrics[k]
				tmp.Timestamp = time.Now()
				cache.Workers[k] = tmp

				log.INFO.Println("Received metrics from " + k)
			} else {
				metrics[k] = "Unknown"
				log.WARNING.Println("Unable to retrieve metrics of " + k)
			}
		} else {
			metrics[k] = "Unknown"
			log.WARNING.Println("Unable to retrieve metrics of " + k)
		}
	}

	Respond(rw, true, metrics)
	return
}

// A REST endpoint that will return the current cache that the 
// leader has.  This is useful to see if anything wrong is 
// happening in the distribution.
func (l *Lead) Cache(rw http.ResponseWriter, r *http.Request) {
	Respond(rw, true, cache)
	return
}

// A REST endpoint that will return the config file being 
// used by the node leader.  This is most useful for any 
// GUI's that want to gather information about the distribution 
// as a whole.
func (l *Lead) Config(rw http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile(l.ConfigPath)
	if err != nil {
		log.ERROR.Println("Unable to load config file.")
	}

	var tmp interface{}

	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.ERROR.Println("Unable to parse config file.")
	}

	Respond(rw, true, tmp)
	return
}

// The last function a leader will call.  Currently just 
// uses os.Exit to quit.
func (l *Lead) Exit() {
	log.INFO.Println("Leader: " + l.Name_ + " is stopped.")
	os.Exit(0)
}

// Private function used to detect is all the workers are 
// currently stopped or not.  This is used by the leader.Stop 
// function to tell if the workers need to be stopped or 
// the leader needs to exit.
func allWorkersStopped() bool {
	for _, v := range cache.Workers {
		if v.State == "Stopped" {
			return true
		}
	}

	return false
}

// Will write to channel but timeout if channel 
// is unavailable.
func writeChannel(ch chan<- interface{}, data interface{}) bool {
	select {
	case <- time.After(time.Second * 2):
		return false
	case ch <- data:
		return true
	}

	return false
}

// Will read from channel but will timeout if channel 
// is unavailable.
func readChannel(ch <-chan interface{}) interface{} {
	select {
	case <- time.After(time.Second * 2):
		return nil
	case data := <- ch:
		return data
	}
}
