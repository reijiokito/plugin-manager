// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter is the interface that we're exposing as a plugin.
type NatsConnect interface {
	Connect(clientName string, config plugin.Client) string
	SendData(clientName string, data string) string
	Subscript(subject string) string
	Close(clientName string) string
}

// Here is an implementation that talks over RPC
type NatsConnectRPC struct{ client *rpc.Client }

func (g *NatsConnectRPC) Connect(name string, config plugin.Client) string {
	var resp string

	data, errMarshal := json.Marshal(name)
	if errMarshal != nil {
		panic(errMarshal)
	}

	err := g.client.Call("Plugin.Connect", data, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	log.Println(fmt.Sprintf("KET QUA CONNECT TRA VE: %v", resp))
	return resp
}

func (g *NatsConnectRPC) SendData(name string, data string) string {
	var resp string
	err := g.client.Call("Plugin.SendData", []string{name, data}, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	log.Println(fmt.Sprintf("KET QUA SENDDATA TRA VE: %v", resp))
	return resp
}

func (g *NatsConnectRPC) Subscript(subject string) string {
	var resp string
	err := g.client.Call("Plugin.Subscript", subject, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	log.Println(fmt.Sprintf("KET QUA SENDDATA TRA VE: %v", resp))
	return resp
}

func (g *NatsConnectRPC) Close(name string) string {
	var resp string
	err := g.client.Call("Plugin.Close", name, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	log.Println(fmt.Sprintf("KET QUA SENDDATA TRA VE: %v", resp))
	return resp
}

// Here is the RPC client that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type NatsConnectRPCServer struct {
	// This is the real implementation
	Impl NatsConnect
}

func (s *NatsConnectRPCServer) Connect(args []byte, resp *string) error {
	log.Println(fmt.Sprintf("SERVER CONNECT: %v", args))
	var info string
	err := json.Unmarshal(args, &info)
	if err != nil {
		panic(err)
	}
	*resp = s.Impl.Connect(info, plugin.Client{})
	return nil
}

func (s *NatsConnectRPCServer) SendData(args []string, resp *string) error {
	log.Println(fmt.Sprintf("SERVER SEND DATA: %v", args))
	*resp = s.Impl.SendData(args[0], args[1])
	return nil
}

func (s *NatsConnectRPCServer) Subscript(arg string, resp *string) error {
	log.Println(fmt.Sprintf("SERVER SUBSCRIPT: %v", arg))
	*resp = s.Impl.Subscript(arg)
	return nil
}

func (s *NatsConnectRPCServer) Close(args string, resp *string) error {
	log.Println(fmt.Sprintf("SERVER CLOSED: %v", args))
	*resp = s.Impl.Close(args)
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
type NatsConnectPlugin struct {
	// Impl Injection
	Impl NatsConnect
}

func (p *NatsConnectPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &NatsConnectRPCServer{Impl: p.Impl}, nil
}

func (NatsConnectPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &NatsConnectRPC{client: c}, nil
}
