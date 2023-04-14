package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Module Module
}

type Module struct {
	Name            string
	HandShakeConfig HandshakeConfig
}

type HandshakeConfig struct {
	ProtocolVersion  uint
	MagicCookieKey   string
	MagicCookieValue string
}

func LoadConfig() *Config {
	jsonFile, err := ioutil.ReadFile("../manifest.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
