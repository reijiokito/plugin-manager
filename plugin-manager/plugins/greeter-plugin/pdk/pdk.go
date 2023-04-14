package pdk

import (
	"github.com/nats-io/nats.go"
	plugin "github.com/reijiokito/sigma-go-plugin"
	"greeter-plugin/config"
	"greeter-plugin/log"
	"greeter-plugin/shared"
	"os/exec"
)

type PDK struct {
	Client         *plugin.Client
	Name           string
	Logger         log.Log
	NATs           *nats.Conn
	ClientProtocol plugin.ClientProtocol
}

func Init(pluginPath string) *PDK {
	var log log.Log
	log.InitLog("basic")

	config := config.LoadConfig()

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  config.Module.HandShakeConfig.ProtocolVersion,
			MagicCookieKey:   config.Module.HandShakeConfig.MagicCookieKey,
			MagicCookieValue: config.Module.HandShakeConfig.MagicCookieValue,
		},
		Plugins: map[string]plugin.Plugin{
			"plugin": &shared.GreeterPlugin{},
		},
		Cmd:     exec.Command(pluginPath),
		Logger:  log.Logger,
		Managed: true,
	})
	defer client.Kill()

	return &PDK{
		Client: client,
		Logger: log,
	}
}
