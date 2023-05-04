package main

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/service"
)

type Config struct {
	Name string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(pdk *go_pdk.PDK) {
	go_pdk.RegisterSubject("kkk", service.Hi)

	subject := "Here"

	s := go_pdk.Server.Plugins["plugin_a"].Callers["Request"](subject, "")
	fmt.Println(s)

}

func GetServices() map[string]func(...interface{}) {
	services := make(map[string]func(...interface{}))

	return services
}

func GetCallers() map[string]func(...interface{}) interface{} {
	callers := make(map[string]func(...interface{}) interface{})

	return callers
}
