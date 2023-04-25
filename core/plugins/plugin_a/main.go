package main

import (
	"fmt"
	"github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/event"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/service"
	proto2 "github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
	"google.golang.org/protobuf/proto"
)

func Access(pdk *go_pdk.PDK) {
	go_pdk.RegisterService(pdk.Connection, "/user/new", service.CreateNew)

	go_pdk.RegisterEvent("manager.handshake", event.HandShakeHandler)

	go_pdk.RegisterChannelSubject("kkk", service.Hi)

	go_pdk.RegisterChannelSubject("hkt", service.Hi)

	for i := 0; i < 20; i++ {
		data, _ := proto.Marshal(&proto2.HelloB{Name: fmt.Sprintf("con chim non from plug in A - %v", i)})
		go pdk.SendChannelData("hi", data)
	}

	for i := 0; i < 5; i++ {
		data, _ := proto.Marshal(&proto2.HelloB{Name: fmt.Sprintf("con chim non from plug in A- HKT - %v", i)})
		go pdk.SendChannelData("hkt", data)
	}

	pdk.Start()
}
