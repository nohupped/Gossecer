package main

import (
	"github.com/nohupped/GoLogger"
	"Gossecer/modules"
	"flag"
	"sync"
)

func main() {
	mylogger := GoLogger.New("/var/log/gossecer.log")
	defer mylogger.Close()
	configfileparam := flag.String("config", "/var/ossec/etc/ossec.conf", "Ossec main configuration file" +
		" where syslog_output is defined.")
	flag.Parse()
	mylogger.Info.Println("Parsing ", *configfileparam)
	host, ip :=modules.GetConfig(configfileparam)

	mylogger.Info.Println("Starting UDP server on ", host, ip)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go modules.StartUdpServer(host, ip, wg)
	wg.Wait()
}