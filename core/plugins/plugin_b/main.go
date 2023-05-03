package main

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/proto"
	proto2 "github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_b/service"
)

type Config struct {
	Name string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(pdk *go_pdk.PDK) {
	fmt.Print("Plugin: " + conf.Name)

	go_pdk.RegisterService(pdk.Nats.Connection, "/user/newB", service.CreateNew)

	//Send event from Nats
	pdk.PostEvent("manager.handshake", &proto.Handshake{Name: "HIHI"}, go_pdk.Scope{
		Local: true,
	})

	//Send from other plugin
	pdk.PostEvent("kkk", &proto2.HelloB{Name: fmt.Sprintf("kkk from plugin B ")}, go_pdk.Scope{
		Local: true,
	})

}
