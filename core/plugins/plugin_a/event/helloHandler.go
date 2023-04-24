package event

import (
	"fmt"
	"github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_a/proto"

	"log"
)

func HelloHandler(ctx *go_pdk.Context, message *proto.Hello) {
	log.Println(fmt.Sprintf("Plugin A: Receive event Hello: %v", message))
}
