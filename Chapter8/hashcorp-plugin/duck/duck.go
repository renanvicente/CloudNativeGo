package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/renanvicente/grpc_sample/hashicorp-plugin/commons"
)

// Here is a real implementation of Sayer
type Duck struct {
}

func (g *Duck) Says() string {
	return "Quack!"
}

func main() {
	sayer := &Duck{}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"sayer": &commons.SayerPlugin{Impl: sayer},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: commons.HandshakeConfig,
		Plugins:         pluginMap,
	})
}
