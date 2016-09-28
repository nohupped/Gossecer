package modules

import (
	"github.com/nohupped/GoLogger"
	"net"
	"strconv"
	"sync"
	"encoding/json"
	"fmt"
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

func StartUdpServer(host string, port int, wg *sync.WaitGroup)  {
	udplogger = GoLogger.New("/var/log/gossecer_udp.log")
	defer udplogger.Close()
	serveraddr, err := net.ResolveUDPAddr("udp4", (host + ":" + strconv.Itoa(port)))
	CheckError(err)
	listner, err := net.ListenUDP("udp", serveraddr)
	CheckError(err)
	udplogger.Info.Println("Udp server listening on", serveraddr.String())
	defer listner.Close()
	for ; ;  {
		handleUDP(listner)
	}
	wg.Done()

}

func handleUDP(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	jsonstring := new(Jsondata)
	//jsonstring := make(map[string]Jsondata)
	n, addr, err := conn.ReadFromUDP(buffer)
	udplogger.Info.Println("Client: ", addr)
	CheckError(err)
	udplogger.Info.Println(string(buffer[:n]))
	fmt.Println(string(buffer[36:n]))
	json.Unmarshal(buffer[36:n], &jsonstring)
	fmt.Println(jsonstring)
}
