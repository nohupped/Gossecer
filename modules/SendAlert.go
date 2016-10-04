package modules

import (
	"fmt"
)


// CheckCounter receives *Jsondata which contains the ossec data struct
// and compares it against the redis hashed key's counter. Uses the same
// redisClient connection used in PutToRedis.
func CheckCounter(alertschan chan *Jsondata, threshold []Key)  {
	Alert := <- alertschan
	fmt.Println(Alert.Counter)

}
