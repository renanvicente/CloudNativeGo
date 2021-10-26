package main

import (
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/renanvicente/grpc_sample/hashicorp-plugin/commons"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: run main/encode-json.go animal")
		os.Exit(1)
	}
	// Get the animal name, and build the path where we expect to
	// find the corresponding executable file.
	name := os.Args[1]
	module := fmt.Sprintf("./%s/%s", name, name)
	// Does the file exist?
	_, err := os.Stat(module)
	if os.IsNotExist(err) {
		log.Fatal("can't find an animal named ", name)
	}

	// pluginMap is the map of plug-ins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"sayer": &commons.SayerPlugin{},
	}
	// Launch the plugin process!
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: commons.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(module),
	})
	defer client.Kill()

	// Connect to the plugin via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}
	// Request the plug-in from the client
	raw, err := rpcClient.Dispense("sayer")
	if err != nil {
		log.Fatal(err)
	}
	// We should have a Sayer now! This feels like a normal interface
	// implementation, but is actually over an RPC connection.
	sayer := raw.(commons.Sayer)

	// Now we can use our loaded plug-in!
	fmt.Printf("A %s says: %q\n", name, sayer.Says())

}
