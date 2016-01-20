package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type IpConfig struct {
	Ip     string `json:"ip"`
	Weight int    `json:"weight"`
}

type ServerConfig struct {
	Host string     `json:"host"`
	Ips  []IpConfig `json:"ips"`
}

type Config struct {
	Servers   []ServerConfig `json:"servers"`
	DnsServer []string       `json:"dnsServers"`
}

var (
	config     *Config
	configLock = new(sync.RWMutex)
)

func loadConfig(configFile string) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("open config: ", err)
		return
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
		return
	}
	configLock.Lock()
	config = temp
	configLock.Unlock()
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}
