// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"client3/config"
	"fmt"
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Client is the interface that we're exposing as a plugin.
type Client interface {
	Init(name string) string
}

// Here is an implementation that talks over RPC
type ClientRPC struct{ client *rpc.Client }

func (g *ClientRPC) Init(name string) string {
	var resp string
	err := g.client.Call("Plugin.Init", name, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	log.Println(fmt.Sprintf("KET QUA TRA VE: %v", resp))
	return resp
}

// Here is the RPC client that ClientRPC talks to, conforming to
// the requirements of net/rpc
type ClientRPCServer struct {
	// This is the real implementation
	Impl         Client
	PluginServer *config.PluginServer
}

func (s *ClientRPCServer) Init(name string, resp *string) error {
	log.Println(fmt.Sprintf("SERVER INIT: %v", name))
	*resp = s.Impl.Init(name)

	m := s.PluginServer.GetPluginServer()
	v := make([]string, 0, len(m))

	for _, value := range m {
		v = append(v, value)
	}

	log.Println(fmt.Sprintf("ALL Client: %v", v))

	s.PluginServer.RegisterClient(name, name)
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC client for this plugin
// type. We construct a ClientRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC plugin. We return ClientRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type ClientPlugin struct {
	// Impl Injection
	Impl         Client
	PluginServer *config.PluginServer
}

func (p *ClientPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ClientRPCServer{Impl: p.Impl, PluginServer: p.PluginServer}, nil
}

func (ClientPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ClientRPC{client: c}, nil
}
