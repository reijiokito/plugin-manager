package main

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	proto2 "github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_b/service"
	"google.golang.org/protobuf/proto"
)

func Access(pdk *go_pdk.PDK) {
	go_pdk.RegisterService(pdk.Connection, "/user/newB", service.CreateNew)

	go_pdk.RegisterChannelSubject("hi", service.Hi)

	for i := 0; i < 20; i++ {
		data, _ := proto.Marshal(&proto2.HelloB{Name: fmt.Sprintf("KKK from plugin B - %v", i)})
		go pdk.SendChannelData("kkk", data)
	}

	pdk.Start()
}
