package main

import (
	"flag"
	"github.com/reijiokito/sigma-go-plugin/sigma"
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

	config := sigma.Configuration{
		NatsUrl:      *natsUrl,
		NatsUsername: *natsUsername,
		NatsPassword: *natsPassword,
	}

	/* Init module */
	sigma.Init(MODULE, &config)
	defer sigma.Release()

	//TODO: Implement APIs create/check/verify plugin

	log.Println("Sigma Plugin Manager Start with port " + *managerPort)
	if err := http.ListenAndServe(":"+*managerPort, nil); err != nil {
		log.Fatal(err)
	}

}
