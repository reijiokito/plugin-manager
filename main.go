package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Module Module
}

type Module struct {
	Name            string
	HandShakeConfig HandshakeConfig
}

type HandshakeConfig struct {
	ProtocolVersion  uint
	MagicCookieKey   string
	MagicCookieValue string
}

func (c *Config) LoadConfig(path string) error {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		log.Println("Cannot read config file: " + path)
		return errors.New("read config error")
	}

	err = json.Unmarshal(jsonFile, &c)
	if err != nil {
		log.Println("Cannot parse data config file: " + path)
		return errors.New("read config error")
	}
	return nil
}

func loadPlugin(path string) (Config, error) {
	//Load manifest file
	pluginConfigDir := path + "/" + "manifest.json"
	var c Config
	err := c.LoadConfig(pluginConfigDir)
	if err != nil {
		return c, err
	}
	return c, nil
}

func watchPlugins(pluginDir string, pluginChan chan string) {
	//Setup interval 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	fmt.Println(pluginDir)

	for range ticker.C {
		files, err := os.ReadDir(pluginDir)
		if err != nil {
			log.Println("error reading plugin directory:", err)
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				continue
			}

			pluginPath := pluginDir + "/" + file.Name()
			pluginChan <- pluginPath
		}
	}
}

func main() {
	pluginDir := "./plugins"
	pluginChan := make(chan string, 1)

	go watchPlugins(pluginDir, pluginChan)

	plugins := make(map[string]Config)
	var pluginPaths []string

	//Prevent terminate program
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case pluginPath := <-pluginChan:
			fmt.Println("Loading plugin:", pluginPath)
			pluginConfig, err := loadPlugin(pluginPath)

			if err != nil {
				log.Println("error loading plugin:", err)
				continue
			} else {
				plugins[pluginConfig.Module.Name] = pluginConfig

				buf, err := json.Marshal(&plugins)
				if err != nil {
					log.Println("error writing configs:", err)
				}

				os.WriteFile("manifest.json", buf, 0644)

				pluginPaths = append(pluginPaths, pluginPath)
			}

		//Received an OS signal, so kill the current process
		case <-signals:
			// Received an OS signal, so kill the current process
			fmt.Println("Received OS signal. Restarting...")

			// Wait for the current process to finish
			time.Sleep(2 * time.Second)

			// Start a new process to replace the current one
			cmd := exec.Command(os.Args[0], os.Args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			// Start the new process and exit the current one
			if err := cmd.Start(); err != nil {
				fmt.Printf("Failed to restart program: %s\n", err)
				return
			}

			os.Exit(0)
		}
	}
}
