package main

import (
	"github.com/nohupped/GoLogger"
	"Gossecer/modules"
	"flag"
	"sync"
	"os"
//	"fmt"
)
func main() {
	mylogger := GoLogger.New("/var/log/gossecer.log")
	defer mylogger.Close()
	hostname, err := os.Hostname()
	modules.CheckError(err)

	configfileparam := flag.String("config", "/var/ossec/etc/ossec.conf", "Ossec main configuration file" +
		" where syslog_output is defined.")
	redisServer := flag.String("redisServer", "localhost", "redis server hostname or IP")
	redisPort := flag.Int("redisPort", 6379, "Port at which redis server is listening")
	flag.Parse()
	mylogger.Info.Println("Parsing ", *configfileparam)
	host, ip :=modules.GetConfig(configfileparam)

	mylogger.Info.Println("Starting UDP server on ", host, ip, "for", hostname)
	itemschan := make(chan *modules.Jsondata, 1024)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go modules.StartUdpServer(host, ip, hostname, itemschan, wg)
	go func() {
		for ; ; {
			modules.PutToRedis(redisServer, redisPort, itemschan)
		}
	}()
	wg.Wait()
}