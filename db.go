package keygenerator

import (
	"github.com/gomodule/redigo/redis"
)

func NewRedisConn() (redis.Conn, error) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func CheckKeyExist(conn redis.Conn, key string) (bool, error) {
	result, err := conn.Do("EXISTS", "0"+key)
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

func MarkKeyAsUsed(conn redis.Conn, key string) error {
	_, err := conn.Do("SET", "0"+key, "0")
	if err != nil {
		return err
	}

	return nil
}
