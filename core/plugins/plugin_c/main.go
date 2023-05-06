package main

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"time"
)

type Config struct {
	Add string `yaml:"add"`
	Url string `yaml:"url"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(pdk *go_pdk.PDK) {
	fmt.Println("Plugin: ", conf.Add)
	fmt.Println("Plugin: ", conf.Url)
	for i := 0; i < 5; i++ {
		go_pdk.Server.Plugins["nats"].Services["Publish"]("hello", []byte(fmt.Sprintf("Hello from plugin B %v", i)))
		time.Sleep(time.Second)
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
