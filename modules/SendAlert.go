package modules

import (
	"github.com/nohupped/GoLogger"
	"fmt"
	"net"
)

var alertLogger *GoLogger.LogIt


// CheckCounter receives *Jsondata which contains the ossec data struct
// and compares it against the redis hashed key's counter. Uses the same
// redisClient connection used in PutToRedis.
func CheckCounter(counterchan chan *Jsondata, threshold []Key, alertschan chan *Jsondata)  {
	alertLogger = GoLogger.New("/var/log/gossecer_alert.log")
	Alert := <-counterchan
	Alert.Threshold = 15 // default value

	Outer:
	for _, i := range threshold {
		for k, v := range i {
			if Alert.Id == k {
				Alert.Threshold = v
				break Outer
			}
		}
	}

	if Alert.Counter >= Alert.Threshold {
		Counter := redisClient.HIncrBy(Alert.HashKey, "Alerted", int64(1))
		Alert.Alerted, _ = Counter.Result()
		alertschan <- Alert
	}
	alertLogger.Info.Println("RuleID -->", Alert.Id)
	alertLogger.Info.Println("Message -->", Alert.Message)
	alertLogger.Info.Println("Counter -->", Alert.Counter)
	alertLogger.Info.Println("Threshold -->", Alert.Threshold)

}

func StartAlert(alertHost string, alertPort string) *net.UDPConn {
	host := alertHost + ":" + alertPort
	RemoteAddr, err := net.ResolveUDPAddr("udp", host)
	CheckError(err)
	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	CheckError(err)
	return conn

}

var conn *net.UDPConn

func SendAlert(alertschan chan *Jsondata,alertHost string, alertPort string) {
	Alert := <- alertschan

	if conn == nil {
		fmt.Println("Starting alert message")
		conn = StartAlert(alertHost, alertPort)
	}

	alertMessage := []byte(fmt.Sprintf("Hostname:%s SyslogCrit:%d Rule:%d Message:%s Times alerted:%d", Alert.Component, Alert.Crit, Alert.Id, Alert.Message, Alert.Alerted))
	length, err := conn.Write(alertMessage)
	if err != nil {
		alertLogger.Err.Println(err)
	}
	alertLogger.Info.Println(length, "bytes sent..")



}
