package modules

import (
	"github.com/nohupped/GoLogger"
	"net"
	"strconv"
	"sync"
	"encoding/json"
	"bytes"
	"regexp"
	"fmt"
)

// Jsondata struct is used to populate from the ossec forwarded udp datagrams
// into their respective fields. Message holds the actual string, whereas
// NormalizedMessage contains the message that is stripped out of timestamps, IPs, portnumbers, etc.
// based on the filters.
type Jsondata struct {
	Crit           int `json:"crit"`
	Id             int `json:"id"`
	Component      string `json:"component"`  //hostname
	Classification string `json:"classification"`
	Description    string `json:"description"`
	Message        string `json:"message"`
	NormalizedMessage string // normalised message
}

// JsondataNormalize will read the Message variable in the struct Jsondata, normalize it and
// puts the normalized string into NormalizedMessage field.
func (d *Jsondata) JsondataNormalize(regexps []*regexp.Regexp)  {
	for _, i := range regexps {
		// Todo
		fmt.Println(i)
	}
}

var udplogger *GoLogger.LogIt
var syshostname string

// StartUDPServer starts a UDP server on the port to which OSSEC is sending datagrams, reads it,
// unmarshals it into the Jsondata struct and writes it to the itemschan.
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
	buffer := make([]byte, 65507)
	jsonstring := new(Jsondata)
	n, addr, err := conn.ReadFromUDP(buffer)
	leanbuf := buffer[:n] //stripping off any unwanted bytes at the end

	// Doing the below shit because for the string "<132>Sep 29 10:04:10 myhostname ossec: {"crit":2,"..., it has to
	// be split with the pattern myhostname ossec: to get the proper json.
	splitbytes := bytes.Split(leanbuf, []byte((syshostname + " ossec: ")))
	udplogger.Info.Println("Client: ", addr)
	CheckError(err)
	udplogger.Info.Println(string(buffer[:n]))
	json.Unmarshal(splitbytes[1], &jsonstring)
	return jsonstring
}
