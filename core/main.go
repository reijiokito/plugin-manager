package main

import (
	"flag"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/proxy"

	"log"
	"net/http"
)

const MODULE = "core"

func main() {
	managerPort := flag.String("manager_port", "8080", "Manager Port")

	natsUrl := flag.String("nats_url", "127.0.0.1", "Nats URL")
	natsUsername := flag.String("nats_username", "", "Nats Username")
	natsPassword := flag.String("nats_password", "", "Nats Password")

	flag.Parse()

	config := go_pdk.Configuration{
		NatsUrl:      *natsUrl,
		NatsUsername: *natsUsername,
		NatsPassword: *natsPassword,
	}

	/* Init module */
	go_pdk.Init(MODULE, &config)
	defer go_pdk.Release()

	//TODO: Implement APIs create/check/verify plugin

	proxy_ := proxy.NewProxy()
	log.Println("Sigma Plugin Manager Start with port " + *managerPort)
	if err := http.ListenAndServe(":"+*managerPort, proxy_); err != nil {
		log.Fatal(err)
	}

}
