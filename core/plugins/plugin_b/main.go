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

	for i := 0; i < 2000; i++ {
		go pdk.Chan.PostEvent("kkk", &proto2.HelloB{Name: fmt.Sprintf("kkk from plugin B - %v", i)})
	}

	pdk.Nats.PostEvent("manager.handshake", &proto.Handshake{Name: "HIHI"})

	pdk.Start()
}
