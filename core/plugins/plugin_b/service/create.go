package service

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
	"log"
)

func CreateNew(ctx *go_pdk.Service, message *proto.HandshakeB) {
	log.Println(fmt.Sprintf("Receive event Handshake: %v", message))

	//
	ctx.Done(&proto.HelloB{Name: "NgonB"})
}
