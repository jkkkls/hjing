package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

var redisConn *redis.Client

// ConnectRedis 连接到redis服务器,在程序启动时调用
func ConnectRedis(connStr string) error {
	arr := strings.Split(connStr, ":")
	redisHost := arr[0]
	redisPort, _ := strconv.Atoi(arr[1])
	redisPwd := arr[2]
	redisDb, _ := strconv.Atoi(arr[3])

	redisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: redisPwd, // no password set
		DB:       redisDb,  // use default DB
	})
	if redisConn == nil {
		return fmt.Errorf("connect redis[%v] error", connStr)
	}
	return nil
}

func GetRedis() *redis.Client {
	return redisConn
}
