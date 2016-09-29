package modules

import (
	"github.com/nohupped/GoLogger"
	"net"
	"strconv"
	"sync"
	"encoding/json"
	"fmt"
	"bytes"
)

type Jsondata struct {
	Crit           int `json:"crit"`
	Id             int `json:"id"`
	Component      string `json:"component"`
	Classification string `json:"classification"`
	Description    string `json:"description"`
	Message        string `json:"message"`
}
var udplogger *GoLogger.LogIt
var syshostname string

func StartUdpServer(host string, port int, hostname string, itemschan chan *Jsondata, wg *sync.WaitGroup)  {
	syshostname = hostname
	udplogger = GoLogger.New("/var/log/gossecer_udp.log")
	defer udplogger.Close()
	serveraddr, err := net.ResolveUDPAddr("udp4", (host + ":" + strconv.Itoa(port)))
	CheckError(err)
	listner, err := net.ListenUDP("udp", serveraddr)
	CheckError(err)
	udplogger.Info.Println("Udp server listening on", serveraddr.String())
	defer listner.Close()
	for ; ;  {
		itemschan <- handleUDP(listner)

	}
	wg.Done()

}

func handleUDP(conn *net.UDPConn) *Jsondata{
	buffer := make([]byte, 65507) // TODO check for end of line and if not, append to the existing byte array
	jsonstring := new(Jsondata)
	//n, addr, err := conn.ReadFromUDP(buffer)
	n, addr, err := conn.ReadFromUDP(buffer)
	leanbuf := buffer[:n]

	// Doing the below shit because for the string "<132>Sep 29 10:04:10 myhostname ossec: {"crit":2,"..., it has to
	// be split with the pattern myhostname ossec: to get the proper json.
	splitbytes := bytes.Split(leanbuf, []byte((syshostname + " ossec: ")))
	udplogger.Info.Println("Client: ", addr)
	CheckError(err)
	fmt.Println(string(splitbytes[1]))
	udplogger.Info.Println(string(buffer[:n]))
	json.Unmarshal(splitbytes[1], &jsonstring)
	return jsonstring
}
