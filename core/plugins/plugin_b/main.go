package main

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/proto"
	proto2 "github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_b/service"
)

func Access(pdk *go_pdk.PDK) {
	go_pdk.RegisterService(pdk.Nats.Connection, "/user/newB", service.CreateNew)

	//Send event from Nats
	pdk.Nats.PostEvent("manager.handshake", &proto.Handshake{Name: "HIHI"})

	//Send from other plugin
	for i := 0; i < 2; i++ {
		go pdk.Chan.PostEvent("kkk", &proto2.HelloB{Name: fmt.Sprintf("kkk from plugin B - %v", i)})
	}

	pdk.Start()
}
