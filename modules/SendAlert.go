package modules

import (
	"github.com/nohupped/GoLogger"
	"net"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var alertLogger *GoLogger.LogIt


// CheckCounter receives *Jsondata which contains the ossec data struct
// and compares it against the redis hashed key's counter. Uses the same
// redisClient connection used in PutToRedis.
func CheckCounter(counterchan chan *Jsondata, threshold []Key, alertschan chan *Jsondata)  {
	alertLogger = GoLogger.New("/var/log/gossecer_alert.log")
	Alert := <-counterchan

	lrange := redisClient.LRange(Alert.RPush, 0, -1)
	ExitTimeCheck:
	for _, i := range lrange.Val() {
		histime, _ := strconv.ParseInt(i, 10, 64)
		if (Alert.CurrentEventOccurrenceTime - histime) >= Alert.TTL.Nanoseconds() {
			fmt.Println(time.Unix(0, histime), "exceeded TTL time of", Alert.TTL.Nanoseconds(), "for", Alert.HashKey, Alert.RPush)
			fmt.Println(redisClient.LPop(Alert.RPush)) // Poping the oldest timestamp because it expired, has to decrement COUNTER
			fmt.Println(redisClient.HIncrBy(Alert.HashKey, "COUNTER", int64(-1))) //Decrementing the COUNTER by 1 each time.

		} else {
			fmt.Println("Exiting to timecheck")
			break ExitTimeCheck
		}
	}

	// Update the struct's Counter with the current counter value
	currentCount := redisClient.HMGet(Alert.HashKey, "COUNTER")
	countlist := currentCount.Val()
	var countint int
	for _, i := range countlist {
		countint, _ = strconv.Atoi(i.(string))
	}
	Alert.Counter = countint

	// Setting default threshold
	Alert.Threshold = 15

	Outer:
	for _, i := range threshold {
		for k, v := range i {
			if Alert.Id == k {
				Alert.Threshold = v
				break Outer
			}
		}
	}

	// Check of Counter has breached the configured Threshold
	if Alert.Counter >= Alert.Threshold {
		Counter := redisClient.HIncrBy(Alert.HashKey, "Alerted", int64(1))
		Alert.Alerted, _ = Counter.Result()
		alertschan <- Alert
	}
/*	alertLogger.Info.Println("RuleID -->", Alert.Id)
	alertLogger.Info.Println("Message -->", Alert.Message)
	alertLogger.Info.Println("Counter -->", Alert.Counter)
	alertLogger.Info.Println("Threshold -->", Alert.Threshold)
	*/

}

// StartAlert will create a udp connection object and returns it.
func StartAlert(alertHost string, alertPort string) *net.UDPConn {
	host := alertHost + ":" + alertPort
	RemoteAddr, err := net.ResolveUDPAddr("udp", host)
	CheckError(err)
	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	CheckError(err)
	return conn

}

// conn is defined globally because if it is defined inside the function, it will open a new socket for every alert.
var conn *net.UDPConn

// SendAlert will use StartAlert() to send alerts that met the configured threshold to the configured UDP port.
// A different program can be used to read these datagrams from the port to send to any alerting mechanism.
func SendAlert(alertschan chan *Jsondata,alertHost string, alertPort string) {
	Alert := <- alertschan

	if conn == nil {
		conn = StartAlert(alertHost, alertPort)
	}
	alert := map[string]interface{}{"Hostname" : Alert.Component, "Syslogcrit": Alert.Crit, "TotalEventOccurance": Alert.TotalCount,
		"EventOccurance": Alert.Counter, "EventThreshold": Alert.Threshold, "TimesAlerted": Alert.Alerted,
		"RuleID": Alert.Id, "Message": Alert.Message}
	alertjson, _ := json.Marshal(&alert)

	//length, err := conn.Write(alertMessage)
	_, err := conn.Write(alertjson)
	if err != nil {
		alertLogger.Err.Println(err)
	}
//	reset := redisClient.HSet(Alert.HashKey, "COUNTER", "0")
	redisClient.HSet(Alert.HashKey, "COUNTER", "0")
//	alertLogger.Info.Println("RESET:", reset)
//	alertLogger.Info.Println(length, "bytes sent..")



}
