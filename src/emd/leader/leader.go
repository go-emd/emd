package leader

import (
	"emd/connector"
	"emd/core"
	"emd/log"
	"emd/worker"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
	"time"
)

var (
	cache *Cache
)

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

type Lead struct {
	core.Core
	GUI_port string
	Workers  []worker.Worker
	Ports    map[string]connector.Connector
}

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

			Respond(rw, true, "Workers started :-)")
		}
	} else {
		log.INFO.Println("Leader: " + l.Name_ + " workers are already running.")
		Respond(rw, false, "Workers already started.")
	}

	return
}

func (l *Lead) Stop(rw http.ResponseWriter, r *http.Request) {
	if !allWorkersStopped() {
		log.INFO.Println("Leader: " + l.Name_ + " is stopping...")

		for k, v := range l.Ports {
			log.INFO.Println("Worker: " + k + "is stopping...")
			v.Channel() <- "STOP"

			tmp := cache.Workers[k]
			tmp.State = "Stopped"
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

func (l *Lead) Status(rw http.ResponseWriter, r *http.Request) {
	for k, v := range l.Ports {
		v.Channel() <- "STATUS"
		tmp := cache.Workers[k]

		select {
		case <-time.After(time.Second * 2):
			log.WARNING.Println("Unable to retrieve status of " + k)

			tmp.Health = "Unknown"
			tmp.Timestamp = time.Now()
			tmp.State = "Unknown"
			cache.Workers[k] = tmp

			Respond(rw, false, "Unknown")
			return
		case status := <-v.Channel():
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
	}

	Respond(rw, true, "Healthy")
	return
}

func (l *Lead) Metrics(rw http.ResponseWriter, r *http.Request) {
	metrics := make(map[string]interface{})

	for k, v := range l.Ports {
		v.Channel() <- "METRICS"

		select {
		case <-time.After(time.Second * 2):
			metrics[k] = "Unknown"
			log.WARNING.Println("Unable to retrieve metrics of " + k)
		case m := <-v.Channel():
			metrics[k] = m

			tmp := cache.Workers[k]
			tmp.Metric = metrics[k]
			tmp.Timestamp = time.Now()
			cache.Workers[k] = tmp

			log.INFO.Println("Received metrics from " + k)
		}
	}

	Respond(rw, true, metrics)
	return
}

func (l *Lead) Cache(rw http.ResponseWriter, r *http.Request) {
	Respond(rw, true, cache)
	return
}

func (l *Lead) Config(rw http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("config.json") // TODO make abs path part of lead.
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

func (l *Lead) Exit() {
	log.INFO.Println("Leader: " + l.Name_ + " is stopped.")
	os.Exit(0)
}

func allWorkersStopped() bool {
	for _, v := range cache.Workers {
		if v.State == "Stopped" {
			return true
		}
	}

	return false
}
