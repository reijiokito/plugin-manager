package main

import (
	"github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/event"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/service"
)

func Access(pdk *go_pdk.PDK) {
	go_pdk.RegisterService(pdk.Nats.Connection, "/user/new", service.CreateNew)

	go_pdk.RegisterSubject("kkk", service.Hi)

	go_pdk.RegisterSubject("manager.handshake", event.HandShakeHandler)
}
