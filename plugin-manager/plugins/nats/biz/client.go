package biz

import (
	"encoding/gob"
	"github.com/hashicorp/go-hclog"
	"os"
)

type Client struct {
	Logger hclog.Logger
}

func (c *Client) Init(name string) string {
	file, err := os.Create("client_" + name + ".gob")
	if err != nil {
		return "Connect error: Sent from plugin"
	}
	defer file.Close()

	// create a new encoder object
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(name)
	if err != nil {
		return "Setup config error: Sent from plugin"
	}

	return "Client Connected. Sent from Server"
}
