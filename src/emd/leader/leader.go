package leader

import (
	"emd/connector"
	"emd/core"
	"emd/log"
	"emd/worker"
	"net/http"
	"os"
	"time"
)

type Leader interface {
	Init()
	Run()
	Exit()

	Start(http.ResponseWriter, *http.Request)
	Stop(http.ResponseWriter, *http.Request)
	Status(http.ResponseWriter, *http.Request)
	Metrics(http.ResponseWriter, *http.Request)
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

	log.INFO.Println("Leader: " + l.Name + " is initialized.")
}

func (l *Lead) Run() {
	log.INFO.Println("Leader: " + l.Name + " is running...")

	// Start all the workers
	for _, w := range l.Workers {
		go w.Run()
	}

	// Handle rest calls and continue managing nodes
	//   workers.
	http.HandleFunc("/start", l.Start)
	http.HandleFunc("/stop", l.Stop)
	http.HandleFunc("/status", l.Status)
	http.HandleFunc("/metrics", l.Metrics)

	http.ListenAndServe(":"+l.GUI_port, nil)
}

func (l *Lead) Start(rw http.ResponseWriter, r *http.Request) {
	log.INFO.Println("Leader: " + l.Name + " is starting it' workers...")

	for _, w := range l.Workers {
		w.Init()
		go w.Run()
	}

	Respond(rw, true, "Workers started :-)")
	return
}

func (l *Lead) Stop(rw http.ResponseWriter, r *http.Request) {
	log.INFO.Println("Leader: " + l.Name + " is stopping...")

	for k, v := range l.Ports {
		log.INFO.Println("Worker: " + k + "is stopping...")
		v.Channel() <- "STOP"
	}

	//l.Exit()

	Respond(rw, true, "Workers stopped :-(")
	return
}

func (l *Lead) Status(rw http.ResponseWriter, r *http.Request) {
	for k, v := range l.Ports {
		v.Channel() <- "STATUS"
		status := <-v.Channel()
		log.INFO.Println("Received status from " + k)

		if status == "Unhealthy" {
			Respond(rw, true, "Unhealthy")
			return
		}
	}

	Respond(rw, true, "Healthy")
	return
}

func (l *Lead) Metrics(rw http.ResponseWriter, r *http.Request) {
	var metrics interface{}

	for k, v := range l.Ports {
		v.Channel() <- "METRICS"
		metrics = <-v.Channel()

		log.INFO.Println("Received metrics from " + k)
	}

	Respond(rw, true, metrics)
	return
}

func (l *Lead) Exit() {
	// Give workers a chance to cleanup and exit properly.
	time.Sleep(time.Second * 2)
	log.INFO.Println("Leader: " + l.Name + " is stopped.")
	os.Exit(0)
}
