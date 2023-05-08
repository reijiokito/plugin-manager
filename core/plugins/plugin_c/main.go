package main

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
)

type Config struct {
	Name    string `yaml:"name"`
	Subject string `yaml:"subject"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(pdk *go_pdk.PDK) {
	fmt.Println("Plugin: ", conf.Name)
	for i := 0; i < 5; i++ {
		go_pdk.Server.Plugins["nats"].Services["Publish"](conf.Subject, []byte(fmt.Sprintf("Hello from plugin C %v", i)))
	}

}

func GetServices() map[string]func(...interface{}) {
	services := make(map[string]func(...interface{}))

	return services
}

func GetCallers() map[string]func(...interface{}) interface{} {
	callers := make(map[string]func(...interface{}) interface{})

	return callers
}
