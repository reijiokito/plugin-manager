package main

import (
	"encoding/json"
	"flag"
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/proxy"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var version = "development"

const MODULE = "manager"

var (
	dump           = flag.String("dump-plugin-info", "", "Dump info about `plugin`")
	dumpAllPlugins = flag.Bool("dump-all-plugins", true, "Dump info about all available plugins")
	pluginsDir     = flag.String("plugins-directory", "/usr/local/sigma/go-plugins", "Set directory `path` where to search plugins")
)

func printVersion() {
	fmt.Printf("Version: %s\nRuntime Version: %s\n", version, runtime.Version())
}

func dumpInfo() {
	s := newServer()

	info, err := s.GetPluginInfo(*dump)
	if err != nil {
		log.Printf("%s", err)
	}

	fmt.Println("Dump info plufin: " + info.Name)

}

func isParentAlive() bool {
	return os.Getppid() != 1 // assume ppid 1 means process was adopted by init
}

func main() {
	managerPort := flag.String("manager_port", "8000", "Manager Port")

	natsUrl := flag.String("nats_url", "127.0.0.1", "Nats URL")
	natsUsername := flag.String("nats_username", "", "Nats Username")
	natsPassword := flag.String("nats_password", "", "Nats Password")

	//dir := flag.String("dir", "/usr/local/sigma/go-plugins", "Directory")

	flag.Parse()

	config := go_pdk.Configuration{
		NatsUrl:      *natsUrl,
		NatsUsername: *natsUsername,
		NatsPassword: *natsPassword,
	}

	/* Init module */
	pdk := go_pdk.Init(MODULE, &config)
	defer pdk.Release()

	s := newServer()
	fmt.Println("------Dump All plufin------")
	pluginPaths, err := filepath.Glob(path.Join(s.pluginsDir, "/*.so"))
	if err != nil {
		log.Printf("can't get plugin names from %s: %s", s.pluginsDir, err)
		return
	}
	infos := make([]PluginInfo, len(pluginPaths))
	for i, pluginPath := range pluginPaths {
		pluginName := strings.TrimSuffix(path.Base(pluginPath), ".so")

		x, err := s.GetPluginInfo(pluginName)
		if err != nil {
			log.Printf("can't load Plugin %s: %s", pluginName, err)
			continue
		}
		infos[i] = *x
		fmt.Println("Dump info plufin: " + infos[i].Name)

	}

	type Config struct {
		Address string
	}

	cf := Config{
		Address: "CONFIG",
	}

	c, err := json.Marshal(cf)
	if err != nil {
		return
	}

	for _, val := range s.plugins {
		instance, err := s.StartInstance(PluginConfig{
			Name:   val.name,
			Config: c,
		})
		if err != nil {
			return
		}
		defer s.CloseInstance(instance.Id)
	}

	for _, val := range s.instances {
		go exec(val.handlers["access"], pdk)
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
