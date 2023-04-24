package main

import (
	"github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/event"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/proto"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/service"
)

const MODULE = "service_A"

func Access(pdk *go_pdk.PDK) {
	go_pdk.RegisterService(pdk.Connection, "/user/new", service.CreateNew)

	go_pdk.RegisterEvent("manager.handshake", event.HandShakeHandler)

	pdk.PostEvent("handshake", &proto.Handshake{
		Name: "XYA",
	})

	select {}
}
