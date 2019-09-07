package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func main() {
	con := redisConnection()
	redisSet("hoge", "10", con)
	pos, _ := strconv.Atoi(redisGet("hoge", con))
	fmt.Println(pos)
}

func redisConnection() redis.Conn {
	const IPPORT = "127.0.0.1:6379"

	c, err := redis.Dial("tcp", IPPORT)
	if err != nil {
		panic(err)
	}
	return c
}

func redisSet(key string, value string, c redis.Conn) {
	c.Do("SET", key, value)
}

func redisGet(key string, c redis.Conn) string {
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return s
}
