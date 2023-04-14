package config

import (
	"github.com/hashicorp/go-plugin"
)

type CPK struct {
	handshake plugin.HandshakeConfig
	pluginMap map[string]plugin.Plugin
}

func (c *CPK) LoadHandshakeConfig() {
	c.handshake = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "basic",
	}

}

func (c *CPK) LoadPluginMapConfig(mapConfig map[string]plugin.Plugin) {
	c.pluginMap = mapConfig
}

func (c *CPK) GetHandshakeConfig() plugin.HandshakeConfig {
	return c.handshake
}

func (c *CPK) GetPluginMapConfig() map[string]plugin.Plugin {
	return c.pluginMap
}

func (c *CPK) Init() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: c.GetHandshakeConfig(),
		Plugins:         c.GetPluginMapConfig(),
	})
}
