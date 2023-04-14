// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter is the interface that we're exposing as a plugin.
type Greeter interface {
	Greet(name string) string
	Calculate(a, b int32) int32
}

// Here is an implementation that talks over RPC
type GreeterRPC struct{ client *rpc.Client }

func (g *GreeterRPC) Greet(name string) string {
	var resp string
	err := g.client.Call("Plugin.Greet", name, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	log.Println(fmt.Sprintf("KET QUA TRA VE: %v", resp))
	return resp
}

func (g *GreeterRPC) Calculate(a, b int32) int32 {
	var resp int32

	err := g.client.Call("Plugin.Calculate", [2]int32{a, b}, &resp)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("KET QUA TINH TOAN TRA VE: %v", resp))
	return -1
}

// Here is the RPC client that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type GreeterRPCServer struct {
	// This is the real implementation
	Impl Greeter
}

func (s *GreeterRPCServer) Greet(name string, resp *string) error {
	log.Println(fmt.Sprintf("SERVER XU LY: %v", name))
	*resp = s.Impl.Greet(name)
	return nil
}

func (s *GreeterRPCServer) Calculate(args [2]int32, resp *int32) error {
	log.Println(fmt.Sprintf("SERVER CALCULATE: %v", args))
	*resp = s.Impl.Calculate(args[0], args[1])
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC client for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC plugin. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type GreeterPlugin struct {
	// Impl Injection
	Impl Greeter
}

func (p *GreeterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl}, nil
}

func (GreeterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GreeterRPC{client: c}, nil
}
