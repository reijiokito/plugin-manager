package event

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/proto"

	"log"
)

func HandShakeHandler(ctx *go_pdk.Context, message *proto.Handshake) {
	log.Println(fmt.Sprintf("Receive event Handshake: %v", message))
}
