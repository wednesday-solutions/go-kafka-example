package rediscache

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func redisDial() (redigo.Conn, error) {
	conn, err := redigo.Dial("tcp", os.Getenv("REDIS_ADDRESS"))
	return conn, err
}

var RedisPool = &redigo.Pool{
	MaxActive:   8640000,
	MaxIdle:     200,
	IdleTimeout: 240 * time.Second,
	Wait:        true,
	Dial:        redisDial,
}

// SetKeyValue ...
func SetKeyValue(key string, data interface{}) error {
	conn, err := redisDial()
	if err != nil {
		return fmt.Errorf("error in redis connection %s", err)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, string(b))
	return err
}

// GetKeyValue ...
func GetKeyValue(key string) (interface{}, error) {
	conn, err := redisDial()
	if err != nil {
		return nil, fmt.Errorf("error in redis connection %s", err)
	}

	reply, err := conn.Do("GET", key)
	return reply, err
}

func RemoveKeyValue(key string) (interface{}, error) {
	conn, err := redisDial()
	if err != nil {
		return nil, fmt.Errorf("error in redis connection %s", err)
	}

	reply, err := conn.Do("DEL", key)
	return reply, err
}
