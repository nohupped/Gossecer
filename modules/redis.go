package modules

import (
	"gopkg.in/redis.v4"
	"github.com/nohupped/GoLogger"
	"fmt"
	"crypto/md5"
	"time"
	"regexp"
	"strconv"
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
type Key map[int]int
type values map[string]string

// PutToRedis will call connectRedis to establish a connection, and use the
// returned client object to put the struct values into redis after applying filters and ttl.
// The key to redis would be an md5hash of the hostname + the normalized message, and the values
// are a COUNTER which is HIncrBy incremented, upon each occurance of the same key, and the value would be
// a key-value pair of Component and non-normalized Message.
func PutToRedis(redisServer string, redisPort string, filters []*regexp.Regexp, expire []Key, itemschan chan *Jsondata, alertschan chan *Jsondata)  {
	redisLogger = GoLogger.New("/var/log/gossecer_redis.log")
	if redisClient == nil {
		redisClient = connectRedis(redisServer, redisPort)
	}
	data := <- itemschan
	data.JsondataNormalize(filters)
	key := data.Component + " " + data.NormalizedMessage
	data.HashKey = fmt.Sprintf("%x",md5.Sum([]byte(key)))
	redisLogger.Info.Println(data.HashKey, " -> ", data)
	msg := values{}
//	COUNTER := values{}
	ruleset := values{}
//	COUNTER["COUNTER"] = "1"
	msg[data.Component] = data.Message
	ruleset["RULE"] = strconv.Itoa(data.Id)
	hashmsg := redisClient.HMSet(data.HashKey, msg)
	rule := redisClient.HMSet(data.HashKey, ruleset)
	Counter := redisClient.HIncrBy(data.HashKey, "COUNTER", int64(1))
	var SetTTL *redis.BoolCmd
	SetTTL = redisClient.Expire(data.HashKey, time.Second * 300) // default ttl
	// Checking length of [expire] section to zero
	if len(expire) != 0 {
		// Setting custom TTL based on rule ID.
		Outer:
		for _, i := range expire {
			for k, v := range i {
				if k == data.Id {
					SetTTL = redisClient.Expire(data.HashKey, time.Second * time.Duration(v))
					break Outer
				}
			}
		}
	}

	redisLogger.Info.Println("HASHMSG -> ", hashmsg,
		"\nCOUNTER ->", Counter,
		"\nSETTL ->", SetTTL,
		"\nRULE ->", rule, "\n\n")
	currentCount := redisClient.HMGet(data.HashKey, "COUNTER")
	countlist := currentCount.Val()
	var countint int
	for _, i := range countlist {
		countint, _ = strconv.Atoi(i.(string))
	}
	data.Counter = countint

	alertschan <- data
}