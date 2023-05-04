package main

import (
	"fmt"
	"github.com/reijiokito/go-pdk"
	proto2 "github.com/reijiokito/plugin-manager/core/plugins/plugin_b/proto"
	"google.golang.org/protobuf/proto"
)

type Config struct {
	Name string
}

func New() interface{} {
	return &Config{}
}

func GetServices() map[string]func(...interface{}) {
	services := make(map[string]func(...interface{}))
	services["PostEvent"] = PostEvent

	return services
}

func GetCallers() map[string]func(...interface{}) interface{} {
	callers := make(map[string]func(...interface{}) interface{})
	callers["Request"] = Request

	return callers
}

func (conf Config) Access(pdk *go_pdk.PDK) {
	pdk.PostEvent("kkk", &proto2.HelloB{Name: fmt.Sprintf("kkk from plugin B ")}, go_pdk.Scope{
		Local: true,
	})

}

func PostEvent(args ...interface{}) { // account_created
	subject := args[0].(string)
	data := args[1].(proto.Message)

	fmt.Println(subject)
	fmt.Println(data)
}

func Request(args ...interface{}) interface{} {
	return "ok"
}
