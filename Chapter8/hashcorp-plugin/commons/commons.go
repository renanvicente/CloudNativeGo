package commons

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// Sayer says what an animal says.
type Sayer interface {
	Says() string
}

// Here is an implementation that talks over RPC
type SayerRPC struct {
	client *rpc.Client
}

func (g SayerRPC) Says() string {
	var resp string

	err := g.client.Call("Plugin.Says", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

// Here is the RPC server that SayerRPC talks to, conforming to
// the requirements of net/rpc
type SayerRPCServer struct {
	// This is the real implementation
	Impl Sayer
}

func (s *SayerRPCServer) Says(args interface{}, resp *string) error {
	*resp = s.Impl.Says()
	return nil
}

type SayerPlugin struct {
	Impl Sayer
}

func (SayerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &SayerRPC{client: c}, nil
}

func (p *SayerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &SayerRPCServer{Impl: p.Impl}, nil
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
