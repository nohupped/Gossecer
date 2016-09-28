package main

import (
	"github.com/nohupped/GoLogger"
	"Ossecer/modules"
	"flag"
)

func main() {
	mylogger := GoLogger.New("/var/log/ossecer.log")
	defer mylogger.Close()
	configfileparam := flag.String("config", "/var/ossec/etc/ossec.conf", "Ossec main configuration file" +
		" where syslog_output is defined.")
	flag.Parse()
	mylogger.Info.Println("Parsing ", *configfileparam)
	host, ip :=modules.GetConfig(configfileparam)

	mylogger.Info.Println("Starting UDP server on ", host, ip)

	modules.StartUdpServer(host, ip, mylogger)
}