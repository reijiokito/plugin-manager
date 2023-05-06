package main

import (
	"flag"
	"fmt"
	go_pdk "github.com/reijiokito/go-pdk"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
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
	managerPort    = flag.String("manager_port", "8000", "Manager Port")
)

var configPlugins [2]map[string][]byte

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
	}
	return infos
}

func dumpBuiltInConfig(configDir string) {
	data, err := ioutil.ReadFile(configDir + "config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	type Plugin struct {
		BuiltIn []string `yaml:"builtIn"`
		Service []string `yaml:"service"`
	}
	type Config struct {
		Plugins []Plugin `yaml:"plugins"`
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	//Load built-in config
	for _, val := range cfg.Plugins[0].BuiltIn {
		config, err := ioutil.ReadFile(configDir + "/plugins/" + val + ".yaml")
		if err != nil {
			log.Fatal(err)
		}
		configPlugins[0][val] = config
	}

	//Load service config
	for _, val := range cfg.Plugins[1].Service {
		config, err := ioutil.ReadFile(configDir + "/plugins/" + val + ".yaml")
		if err != nil {
			log.Fatal(err)
		}
		configPlugins[1][val] = config
	}

}

func initBuildInPlugin(pdk *go_pdk.PDK) {
	for _, val := range go_pdk.Server.Plugins {
		if _, ok := configPlugins[0][val.Name]; ok {
			_, err := go_pdk.Server.StartInstance(go_pdk.PluginConfig{
				Name:   val.Name,
				Config: configPlugins[0][val.Name],
			})
			if err != nil {
				return
			}
		}
	}
	for _, val := range go_pdk.Server.Instances {
		if _, ok := configPlugins[0][val.Plugin.Name]; ok {
			exec(val.Handlers["access"], pdk)
		}
	}
}

func initServicePlugin(pdk *go_pdk.PDK) {
	for _, val := range go_pdk.Server.Plugins {
		if _, ok := configPlugins[1][val.Name]; ok {
			_, err := go_pdk.Server.StartInstance(go_pdk.PluginConfig{
				Name:   val.Name,
				Config: configPlugins[1][val.Name],
			})
			if err != nil {
				return
			}
		}
	}

	for _, val := range go_pdk.Server.Instances {
		if _, ok := configPlugins[1][val.Plugin.Name]; ok {
			exec(val.Handlers["access"], pdk)
		}
	}
}

func exec(f func(*go_pdk.PDK), pdk *go_pdk.PDK) {
	f(pdk)
}

func main() {
	/* Init module */
	pdk := go_pdk.Init(*pluginsDir)
	defer pdk.Release()

	//Dump all existed plugins info
	dumpAllInfo()

	//Read config from plugins
	for i := 0; i < 2; i++ {
		configPlugins[i] = make(map[string][]byte)
	}
	dumpBuiltInConfig(*configDir)

	//Initialize built-in plugin
	initBuildInPlugin(pdk)

	time.Sleep(time.Second)

	//Initialize service plugins
	initServicePlugin(pdk)

	//pdk.Start()

}
