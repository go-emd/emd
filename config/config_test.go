package config

import (
	"testing"
	"reflect"
)

func TestProcess(t *testing.T) {
	solution := Config{
		Nfs: false,
		GUI_port: "1234",
		Nodes: []NodeConfig{
			{
				Hostname: "example.com",
				Workers: []WorkConfig{
					{
						Name: "MyName",
						Connections: []Connection{
							{
								Type: "LocalEgress",
								Worker: "WorkerName",
								Alias: "WorkerAlias",
								Buffer: "0",
							},
						},
					},
				},
			},
		},
	}

	cfg := Config{}
	Process("config_test.json", &cfg)

	if !reflect.DeepEqual(cfg, solution) {
		t.Fail()
	}
}
