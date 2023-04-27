package main

import (
	"flag"
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"plugin"
)

const MODULE = "manager"

func main() {
	managerPort := flag.String("manager_port", "8000", "Manager Port")

	natsUrl := flag.String("nats_url", "127.0.0.1", "Nats URL")
	natsUsername := flag.String("nats_username", "", "Nats Username")
	natsPassword := flag.String("nats_password", "", "Nats Password")

	dir := flag.String("dir", "/usr/local/sigma/go-plugins", "Directory")

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
	pluginFiles, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}
	for _, f := range pluginFiles {
		pluginPath := filepath.Join(*dir, f.Name())
		p, err := plugin.Open(pluginPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to load plugin %s: %v\n", pluginPath, err))
		}

		initSymbol, err := p.Lookup("Access")
		if err != nil {
			log.Fatal(fmt.Errorf("failed to lookup MyFunction symbol: %v", err))
		}

		pluginAccess, ok := initSymbol.(func(*go_pdk.PDK))
		if !ok {
			log.Fatal(fmt.Errorf("failed to convert MyFunc" +
				"tion symbol to expected function signature"))
		}

		go exec(pluginAccess, pdk)
	}

	pdk.Start()

	proxy_ := proxy.NewProxy(pdk)
	log.Println("Sigma Plugin Manager Start with port " + *managerPort)
	if err := http.ListenAndServe(":"+*managerPort, proxy_); err != nil {
		log.Fatal(err)
	}

}

func exec(f func(*go_pdk.PDK), pdk *go_pdk.PDK) {
	f(pdk)
}
