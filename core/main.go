package main

import (
	"flag"
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"github.com/reijiokito/plugin-manager/core/proxy"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	dump           = flag.String("dump-plugin-info", "", "Dump info about `plugin`")
	dumpAllPlugins = flag.Bool("dump-all-plugins", true, "Dump info about all available plugins")
	pluginsDir     = flag.String("plugins-directory", "/usr/local/sigma/go-plugins", "Set directory `path` where to search plugins")
	configDir      = flag.String("config-plugin-directory", "/home/cong/Downloads/24_4/plugin-manager/core/config/", "Set config directory `path` where to load plugin configs")
)

var configPlugins map[string][]byte

func dumpInfo() {
	info, err := go_pdk.Server.GetPluginInfo(*dump)
	if err != nil {
		log.Printf("%s", err)
	}

	fmt.Println("Dump info plufin: " + info.Name)

}

func dumpAllInfo() []go_pdk.PluginInfo {
	fmt.Println("------Dump All plugin------")
	pluginPaths, err := filepath.Glob(path.Join(go_pdk.Server.PluginsDir, "/*.so"))
	if err != nil {
		log.Printf("can't get plugin names from %s: %s", go_pdk.Server.PluginsDir, err)
		return nil
	}
	infos := make([]go_pdk.PluginInfo, len(pluginPaths))
	for i, pluginPath := range pluginPaths {
		pluginName := strings.TrimSuffix(path.Base(pluginPath), ".so")

		x, err := go_pdk.Server.GetPluginInfo(pluginName)
		if err != nil {
			log.Printf("can't load Plugin %s: %s", pluginName, err)
			continue
		}
		infos[i] = *x
		fmt.Println(fmt.Sprintf("Dump info plugin: %v", infos[i].Schema))
	}
	return infos
}

func dumpBuiltInConfig(configDir string) {
	data, err := ioutil.ReadFile(configDir + "config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	type Plugin struct {
		BuiltIn []string `yaml:"builtin"`
		Service []string `yaml:"service"`
	}
	type Config struct {
		Plugin []Plugin `yaml:"plugins"`
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

	//Load built-in plugin
	for _, val := range cfg.Plugin[0].BuiltIn {
		config, err := ioutil.ReadFile(configDir + "/build_in/" + val + ".yaml")
		if err != nil {
			log.Fatal(err)
		}

		configPlugins[val] = config
	}

	//Load service plugin
	for _, val := range cfg.Plugin[0].Service {
		config, err := ioutil.ReadFile(configDir + "/build_in/" + val + ".yaml")
		if err != nil {
			log.Fatal(err)
		}
		configPlugins[val] = config
	}

}

func main() {
	managerPort := flag.String("manager_port", "8000", "Manager Port")

	/* Init module */
	pdk := go_pdk.Init(*pluginsDir)
	defer pdk.Release()

	//Dump all existed plugins info
	dumpAllInfo()

	//Read config from plugins
	configPlugins = make(map[string][]byte)
	dumpBuiltInConfig(*configDir)

	//Initialize built-in plugin
	for _, val := range go_pdk.Server.Plugins {
		if val.Name == "nats" {
			type Config struct {
				NatsUrl      string
				NatsUsername string
				NatsPassword string
			}
			_, err := go_pdk.Server.StartInstance(go_pdk.PluginConfig{
				Name:   val.Name,
				Config: configPlugins[val.Name],
			})
			if err != nil {
				return
			}

		}
	}
	for _, val := range go_pdk.Server.Instances {
		if val.Plugin.Name == "nats" {
			exec(val.Handlers["access"], pdk)
		}
	}
	time.Sleep(time.Second)

	//Initialize plugins
	for _, val := range go_pdk.Server.Plugins {
		if val.Name != "nats" {
			_, err := go_pdk.Server.StartInstance(go_pdk.PluginConfig{
				Name:   val.Name,
				Config: configPlugins[val.Name],
			})
			if err != nil {
				return
			}
		}
	}

	for _, val := range go_pdk.Server.Instances {
		if val.Plugin.Name != "nats" {
			exec(val.Handlers["access"], pdk)
		}
	}

	//pdk.Start()

	proxy_ := proxy.NewProxy(pdk)
	log.Println("Sigma Plugin Manager Start with port " + *managerPort)
	if err := http.ListenAndServe(":"+*managerPort, proxy_); err != nil {
		log.Fatal(err)
	}
}

func exec(f func(*go_pdk.PDK), pdk *go_pdk.PDK) {
	f(pdk)
}
