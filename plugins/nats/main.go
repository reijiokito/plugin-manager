// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"client3/config"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/exec"

	"client3/shared"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func main() {
	// Connect to NATS message broker
	broker, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer broker.Close()

	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	var server config.CPK
	server.LoadHandshakeConfig()
	server.LoadPluginMapConfig(map[string]plugin.Plugin{
		"client":     &shared.ClientPlugin{},
		"plugin":     &shared.GreeterPlugin{},
		"natsserver": &shared.NatsConnectPlugin{},
	})

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: server.GetHandshakeConfig(),
		Plugins:         server.GetPluginMapConfig(),
		Cmd:             exec.Command("../../plugin/plugin"),
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("client")
	if err != nil {
		log.Fatal(err)
	}

	greeter := raw.(shared.Client)
	fmt.Println(greeter.Init("CLIENT_C"))

	rawConnect, err := rpcClient.Dispense("natsserver")
	if err != nil {
		log.Fatal(err)
	}

	natsConnect := rawConnect.(shared.NatsConnect)

	fmt.Println(natsConnect.Subscript("CLIENT_C"))

	select {}
}
