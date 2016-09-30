package main

import (
	"github.com/nohupped/GoLogger"
	"Gossecer/modules"
	"flag"
	"sync"
	"os"
	"regexp"
)
func main() {
	// Flags and Variable declaration
	mylogger := GoLogger.New("/var/log/gossecer.log")
	defer mylogger.Close()
	hostname, err := os.Hostname()
	modules.CheckError(err)
	configfileparam := flag.String("config", "/etc/gossec.conf", "The program's main configuration file")
	flag.Parse()
	// Main config file
	configfile, err := modules.GetConfig(*configfileparam)
	modules.CheckError(err)
	// Get OSSec's configfile
	OSSecConfGlobal, err := configfile.GetSection("ossec")
	modules.CheckError(err)
	ossecConf, err := OSSecConfGlobal.GetKey("ConfFile")
	modules.CheckError(err)

	// Get Redis Config file
	RedisConfGlobal, err := configfile.GetSection("redis")
	redisServer, err := RedisConfGlobal.GetKey("Server")
	modules.CheckError(err)
	redisPort, err := RedisConfGlobal.GetKey("Port")
	modules.CheckError(err)

	// Filters
	FiltersGlobal, err := configfile.GetSection("filters")
	modules.CheckError(err)
	filters_keys := FiltersGlobal.Keys()
	var regexps []*regexp.Regexp

	for _, i := range filters_keys {
		regexps = append(regexps, regexp.MustCompile(i.MustString("")))
	}
	// End of variable declaration

	mylogger.Info.Println("Parsing ", ossecConf.String())
	host, ip :=modules.GetOSSecConfig(ossecConf.MustString("/var/ossec/etc/ossec.conf"))

	mylogger.Info.Println("Starting UDP server on ", host, ip, "for", hostname)
	itemschan := make(chan *modules.Jsondata, 1024)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go modules.StartUdpServer(host, ip, hostname, itemschan, wg)
	go func() {
		for ; ; {
			modules.PutToRedis(redisServer.MustString("localhost"), redisPort.MustString("6379"),  regexps, itemschan)
		}
	}()
	wg.Wait()
}