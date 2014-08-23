/*
	The config package manages the parsing of the 
	config.json configuration file required by each 
	emd distribution.  The file tells emd what the 
	nuts and bolts of the distribution are and how 
	they are connected.  It is affectively the 
	blueprint of the distribution.
*/
package config

import (
	"github.com/go-emd/emd/log"
	"encoding/json"
	"io/ioutil"
)

// Contains the variables related to a
// connector interface allowing the leader.template 
// file fill in these parameters.
type Connection struct {
	Type   string
	Worker string
	Alias  string
	Buffer string
}

// Basic configuration of a worker in 
// a distribution.  It contains only the name of 
// the worker and all of its connections.
type WorkConfig struct {
	Name        string
	Connections []Connection
}

// Contains all of the basic information 
// a node requires in emd.  It contains the hostname 
// this node leader will run and the workers that 
// run within it.
type NodeConfig struct {
	Hostname string
	Workers  []WorkConfig
}

// Contains misc things emd needs to know such 
// as if a NFS exists and what port to listen for REST 
// requests.  It contains all the nodes in the distribution.
type Config struct {
	Nfs      bool
	GUI_port string
	Nodes    []NodeConfig
}

// Processes the config.json and parses the file 
// into the Config structure using json.Unmarshal.
func Process(path string, config *Config) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.ERROR.Println("Unable to read config file.")
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.ERROR.Println("Unable to parse json config file.")
	}
}
