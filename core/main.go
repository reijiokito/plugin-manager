package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	go_pdk "github.com/reijiokito/go-pdk"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	dumpAllPlugins = flag.Bool("dump-all-plugins", true, "Dump info about all available plugins")
	pluginsDir     = flag.String("plugins-directory", "/usr/local/sigma/go-plugins", "Set directory `path` where to search plugins")
	configDir      = flag.String("config-plugin-directory", "/home/cong/Downloads/24_4/plugin-manager/core/config/", "Set config directory `path` where to load plugin configs")
	managerPort    = flag.String("manager_port", "localhost:8000", "Manager Port")
)

var configPlugins [2]map[string][]byte
var pluginInfos []go_pdk.PluginInfo

func dumpInfo(name string) {
	info, err := go_pdk.Server.GetPluginInfo(name)
	if err != nil {
		log.Printf("%s", err)
	}

	fmt.Println(fmt.Sprintf("Dump info plufin: %s", info.Name))
	for i, _ := range pluginInfos {
		if pluginInfos[i].Name == info.Name {
			pluginInfos[i] = *info
			return
		}
	}
	pluginInfos = append(pluginInfos, *info)
}

func dumpAllInfo() {
	fmt.Println("------Dump All plugin------")
	pluginPaths, err := filepath.Glob(path.Join(go_pdk.Server.PluginsDir, "/*.so"))
	if err != nil {
		log.Printf("can't get plugin names from %s: %s", go_pdk.Server.PluginsDir, err)
		return
	}

	for _, pluginPath := range pluginPaths {
		pluginName := strings.TrimSuffix(path.Base(pluginPath), ".so")

		x, err := go_pdk.Server.GetPluginInfo(pluginName)
		if err != nil {
			log.Printf("can't load Plugin %s: %s", pluginName, err)
			continue
		}
		pluginInfos = append(pluginInfos, *x)
	}
}

func dumpPluginConfig(name string) {
	config, err := os.ReadFile(*configDir + "/plugins/" + name + ".yaml")
	if err != nil {
		log.Fatal(err)
	}
	configPlugins[1][name] = config
}

func dumpAllPluginConfig() {
	data, err := os.ReadFile(*configDir + "config.yaml")
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
		config, err := os.ReadFile(*configDir + "/plugins/" + val + ".yaml")
		if err != nil {
			log.Fatal(err)
		}
		configPlugins[0][val] = config
	}

	//Load service config
	for _, val := range cfg.Plugins[1].Service {
		config, err := os.ReadFile(*configDir + "/plugins/" + val + ".yaml")
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

	//Read plugin configs
	for i := 0; i < 2; i++ {
		configPlugins[i] = make(map[string][]byte)
	}
	dumpAllPluginConfig()

	//Initialize built-in plugin
	initBuildInPlugin(pdk)

	time.Sleep(time.Second)

	//Initialize service plugins
	initServicePlugin(pdk)

	//pdk.Start()

	r := gin.Default()

	r.GET("/plugin/get-all-info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": pluginInfos,
		})
	})

	r.GET("/plugin/get-all-instance-info", func(c *gin.Context) {
		type InstanceInfo struct {
			Id                int
			Name              string
			Modtime           time.Time
			Loadtime          time.Time
			LastStartInstance time.Time
			LastCloseInstance time.Time
		}

		var list []InstanceInfo

		for key, value := range go_pdk.Server.Instances {
			list = append(list, InstanceInfo{
				Id:                key,
				Name:              value.Plugin.Name,
				Modtime:           value.Plugin.Modtime,
				Loadtime:          value.Plugin.Loadtime,
				LastStartInstance: value.Plugin.LastStartInstance,
				LastCloseInstance: value.Plugin.LastCloseInstance,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data": list,
		})
	})

	r.POST("/plugin/init", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}

		type Config struct {
			Name   string      `json:"name"`
			Config interface{} `json:"config"`
		}

		// Create a YAML file
		var data Config
		err = yaml.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal YAML"})
			return
		}
		yamlData, err := yaml.Marshal(data.Config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal YAML"})
			return
		}
		err = os.WriteFile(*configDir+"/plugins/"+data.Name+".yaml", yamlData, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create YAML file"})
			return
		}

		//Start instance
		dumpInfo(data.Name)
		dumpPluginConfig(data.Name)

		var instanceId int

		for _, val := range go_pdk.Server.Plugins {
			if val.Name == data.Name {
				if _, ok := configPlugins[1][data.Name]; ok {
					status, err := go_pdk.Server.StartInstance(go_pdk.PluginConfig{
						Name:   val.Name,
						Config: configPlugins[1][data.Name],
					})
					if err != nil {
						log.Println(fmt.Sprintf("Start instance err: %v", err))
						return
					}

					instanceId = status.Id
				}
			}

		}

		for _, val := range go_pdk.Server.Instances {
			if val.Plugin.Name == data.Name && val.Id == instanceId {
				if _, ok := configPlugins[1][data.Name]; ok {
					go exec(val.Handlers["access"], pdk)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": "Success",
		})
	})

	r.Run(*managerPort)

}
