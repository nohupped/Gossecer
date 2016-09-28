package modules

import (
	"github.com/nohupped/GoLogger"
	"net"
	"strconv"
	"fmt"
)

func StartUdpServer(host string, port int, mylogger *GoLogger.LogIt)  {
	serveraddr, err := net.ResolveUDPAddr(host, strconv.Itoa(port))
	CheckError(err)
	fmt.Println(serveraddr)

}
