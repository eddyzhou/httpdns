package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

var (
	initResolverOnce sync.Once
	errNoHost        = errors.New("resolver: not config this host")
	resolverLock     = new(sync.RWMutex)
	resolver         *Resolver
)

type Resolver struct {
	ipMap map[string][]string
	pv    uint64
}

func GetResolver() *Resolver {
	resolverLock.RLock()
	defer resolverLock.RUnlock()
	return resolver
}

func InitResolver(configFile string) {
	resolverLock.Lock()
	defer resolverLock.Unlock()
	resolver = newResolver(configFile)
}

func newResolver(configFile string) *Resolver {
	initResolverOnce.Do(func() { preInitResolver(configFile) })
	config := GetConfig()
	r := new(Resolver)
	r.ipMap = make(map[string][]string)
	for _, s := range config.Servers {
		host := s.Host
		var ips []string
		for _, v := range s.Ips {
			for i := 0; i < v.Weight; i++ {
				ips = append(ips, v.Ip)
			}
		}
		r.ipMap[host] = ips
	}
	return r
}

func (r *Resolver) GetIp(host string) (string, error) {
	if ips, ok := r.ipMap[host]; ok {
		pvFinal := atomic.AddUint64(&(r.pv), 1)
		return ips[pvFinal%uint64(len(ips))], nil
	}
	return "", errNoHost
}

func preInitResolver(configFile string) {
	loadConfig(configFile)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(configFile)
			resolverLock.Lock()
			resolver = newResolver(configFile)
			resolverLock.Unlock()
			log.Println("Reloaded")
		}
	}()
}
