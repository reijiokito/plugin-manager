package main

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
)

type Config struct {
	Name string `yaml:"name"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(pdk *go_pdk.PDK) {
	fmt.Println("Plugin: ", conf.Name)

}

func GetServices() map[string]func(...interface{}) {
	services := make(map[string]func(...interface{}))

	return services
}

func GetCallers() map[string]func(...interface{}) interface{} {
	callers := make(map[string]func(...interface{}) interface{})

	return callers
}
