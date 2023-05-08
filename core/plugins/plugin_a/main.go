package main

import (
	"fmt"
	"github.com/reijiokito/go-pdk"
	proto2 "github.com/reijiokito/plugin-manager/core/plugins/plugin_a/proto"
	"google.golang.org/protobuf/proto"
)

type Config struct {
	Name string `yaml:"name"`
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
	fmt.Println("Plugin: ", conf.Name)

	go_pdk.Server.Plugins["nats"].Services["Subscribe"]("hello")

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

type SubjectHandler func(data proto.Message)

func Hello(message *proto2.Hello) {
	fmt.Println("Receive message: " + message.Name)
}
