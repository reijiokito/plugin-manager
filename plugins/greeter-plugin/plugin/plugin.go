package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/reijiokito/plugin-manager/plugins/greeter-plugin/shared"
	plugin "github.com/reijiokito/sigma-go-plugin"
	"log"
	"os"
	"time"
)

type GreeterHello struct {
	Logger hclog.Logger
}

func (g *GreeterHello) Calculate(a, b int32) int32 {
	g.Logger.Debug("message from GreeterHello.Calculate")
	count := 0
	for i := 0; i < 10; i++ {
		log.Println(fmt.Sprintf("Hello: %d", count))
		count++
		time.Sleep(time.Second * 1)
	}
	return a + b
}

func (g *GreeterHello) Greet(name string) string {
	g.Logger.Debug("message from GreeterHello.Greet")

	return fmt.Sprintf("Hello, %s!", name)
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	greeter := &GreeterHello{
		Logger: logger,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "basic",
		},
		Plugins: map[string]plugin.Plugin{
			"plugin": &shared.GreeterPlugin{Impl: greeter},
		},
		Logger: logger,
	})
}
