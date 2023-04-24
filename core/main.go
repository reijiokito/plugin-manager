package main

import (
	"flag"
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/proxy"
	"log"
	"net/http"
	"plugin"
)

const MODULE = "manager"

func main() {
	managerPort := flag.String("manager_port", "8000", "Manager Port")

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
	pdk := go_pdk.Init(MODULE, &config)
	defer pdk.Release()

	//READ
	p, err := plugin.Open("/home/cong/Downloads/24_4/plugin-manager/core/plugins/plugin_a/main")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to load plugin file: %v", err))
	}

	// Look up the `MyFunction` symbol in the plugin.
	initSymbol, err := p.Lookup("Access")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to lookup MyFunction symbol: %v", err))
	}

	pluginAccess, ok := initSymbol.(func(*go_pdk.PDK))
	if !ok {
		log.Fatal(fmt.Errorf("failed to convert MyFunction symbol to expected function signature"))
	}

	go exec(pluginAccess, pdk)

	proxy_ := proxy.NewProxy(pdk)
	log.Println("Sigma Plugin Manager Start with port " + *managerPort)
	if err := http.ListenAndServe(":"+*managerPort, proxy_); err != nil {
		log.Fatal(err)
	}
}

func exec(f func(*go_pdk.PDK), pdk *go_pdk.PDK) {
	f(pdk)
}
