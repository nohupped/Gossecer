package modules

import "fmt"

// SendAlert receives map[*Jsondata]*redis.StringStringMapCmd which contains the ossec data struct
// and the redis data that is stored, for furthur breakdown and write to a socket.
func SendALert(alertschan chan *AlertData)  {
	Alert := <- alertschan
	fmt.Println(Alert)
}
