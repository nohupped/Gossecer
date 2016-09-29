package modules

import (
	"gopkg.in/redis.v4"
	"strconv"
	"github.com/nohupped/GoLogger"
	"fmt"
	"strings"
)

var redisLogger *GoLogger.LogIt

func connectRedis(redisServer *string, redisPort *int)  *redis.Client{
	client := redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		Addr: (*redisServer + ":" + strconv.Itoa(*redisPort)),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	CheckError(err)
	redisLogger.Info.Println(pong, "from", client.String())
	return client
}

var redisClient *redis.Client

func PutToRedis(redisServer *string, redisPort *int, itemschan chan *Jsondata)  {
	redisLogger = GoLogger.New("/var/log/gossecer_redis.log")
	if redisClient == nil {
		redisClient = connectRedis(redisServer, redisPort)
	}
	data := <- itemschan
	splitmsg := strings.Split(data.Message, " ")
	if len(splitmsg) > 10 {
		msg := data.Component + " " + (strings.Join(splitmsg[10:len(splitmsg)], " "))
		fmt.Println(msg)
	}






}