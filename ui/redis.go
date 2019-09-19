package ui

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

// Redis functions
func redisConnection(target string) (redis.Conn, error) {

	c, err := redis.Dial("tcp", target)
	return c, err
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
