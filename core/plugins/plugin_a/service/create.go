package service

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/proto"
	"log"
)

func CreateNew(ctx *go_pdk.Service, message *proto.Handshake) {
	log.Println(fmt.Sprintf("Receive event Handshake: %v", message))
	//
	ctx.Done(&proto.Hello{Name: "Ngon"})
}
