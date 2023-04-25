package service

import (
	"fmt"
	"github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
)

func Hi(message *proto.HelloB) {
	fmt.Println(fmt.Sprintf("handle from plugin B - %v", message))
}
