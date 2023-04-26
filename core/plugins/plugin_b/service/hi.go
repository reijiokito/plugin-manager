package service

import (
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
)

func Hi(ctx *go_pdk.Context, message *proto.HelloB) {
	fmt.Println(fmt.Sprintf("handle from plugin B - %v", message))
}
