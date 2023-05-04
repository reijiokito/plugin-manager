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
	"time"
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
	info, err := go_pdk.Server.GetPluginInfo(*dump)
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

	/* Init module */
	pdk := go_pdk.Init(*pluginsDir)
	defer pdk.Release()

	pdk.Start()

	time.Sleep(time.Second)

	fmt.Println("------Dump All plugin------")
	pluginPaths, err := filepath.Glob(path.Join(go_pdk.Server.PluginsDir, "/*.so"))
	if err != nil {
		log.Printf("can't get plugin names from %s: %s", go_pdk.Server.PluginsDir, err)
		return
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
		fmt.Println("Dump info plugin: " + infos[i].Name)
		fmt.Println(fmt.Sprintf("Dump info plugin: %v", infos[i].Schema))

	}

	type Config struct {
		Name string
	}

	cf := Config{
		Name: "CONFIG",
	}

	c, err := json.Marshal(cf)
	if err != nil {
		return
	}

	for _, val := range go_pdk.Server.Plugins {
		instance, err := go_pdk.Server.StartInstance(go_pdk.PluginConfig{
			Name:   val.Name,
			Config: c,
		})
		if err != nil {
			return
		}
		defer go_pdk.Server.CloseInstance(instance.Id)
	}

	for _, val := range go_pdk.Server.Instances {
		go exec(val.Handlers["access"], pdk)
	}

	proxy_ := proxy.NewProxy(pdk)
	log.Println("Sigma Plugin Manager Start with port " + *managerPort)
	if err := http.ListenAndServe(":"+*managerPort, proxy_); err != nil {
		log.Fatal(err)
	}
}

func exec(f func(*go_pdk.PDK), pdk *go_pdk.PDK) {
	f(pdk)
}
