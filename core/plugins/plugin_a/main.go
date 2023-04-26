package main

import (
	"fmt"
	"github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/event"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/service"
	proto2 "github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
)

func Access(pdk *go_pdk.PDK) {
	go_pdk.RegisterService(pdk.Nats.Connection, "/user/new", service.CreateNew)

	go_pdk.RegisterSubject("manager.handshake", event.HandShakeHandler)

	go_pdk.RegisterSubject("kkk", service.Hi)

	for i := 0; i < 500; i++ {
		go pdk.Chan.PostEvent("hkt", &proto2.HelloB{Name: fmt.Sprintf("hkt con chim non from plug in A- HKT - %v", i)})
	}

	pdk.Start()
}
