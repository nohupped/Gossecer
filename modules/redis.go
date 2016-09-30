package modules

import (
	"gopkg.in/redis.v4"
	"strconv"
	"github.com/nohupped/GoLogger"
	"fmt"
	"crypto/md5"
	"time"
	"regexp"
)

var redisLogger *GoLogger.LogIt

// connectRedis will attempt a connection to the redis server and returns the client object
func connectRedis(redisServer string, redisPort string)  *redis.Client{
	client := redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		Addr: (redisServer + ":" + redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	CheckError(err)
	redisLogger.Info.Println(pong, "from", client.String())
	return client
}

var redisClient *redis.Client

// PutToRedis will call connectRedis to establish a connection, and use the
// returned client object to put the struct values into redis after applying filters and ttl.
func PutToRedis(redisServer string, redisPort string, filters []*regexp.Regexp, itemschan chan *Jsondata)  {
	redisLogger = GoLogger.New("/var/log/gossecer_redis.log")

	fmt.Println("filters are ", filters)
	if redisClient == nil {
		redisClient = connectRedis(redisServer, redisPort)
	}
	data := <- itemschan
	key :=  data.Component + " " + strconv.Itoa(data.Id)
	hexHashedKey := fmt.Sprintf("%x",md5.Sum([]byte(key)))
	redisLogger.Info.Println(key, " -> ", hexHashedKey)
	value := make(map[string]string)
	dummy := make(map[string]string)
	dummy["COUNTER"] = "1"
	value[data.Component] = data.Message
	status := redisClient.HMSet(hexHashedKey,value)
	Counter := redisClient.HIncrBy(hexHashedKey, "COUNTER", int64(1))
	SetTTL := redisClient.Expire(hexHashedKey, time.Second*60)
	redisLogger.Info.Println(status, Counter, SetTTL)


}

// To strip out the timestamp that comes like Sep 30 14:41:04
func ApplyFiltersOnMsg(msg string) string {

	return msg
}