package keygenerator

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

func GenerateKeysToRedis() {
	conn, err := NewRedisConn()
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	generateKeyTicker := time.NewTicker(time.Millisecond * 10)
	getKeyNumsTicker := time.NewTicker(time.Second * 2)
	var keyNums int64
	keyNums = 0
	for {
		select {
		case <-generateKeyTicker.C:
			if keyNums >= 1000 { // I know there may be some concurrency problem, but the number of keys do not need to be so strict
				time.Sleep(100 * time.Millisecond)
				continue
			}
			err = GenerateKeyToRedis(conn)
			if err != nil {
				log.Println(err)
				continue
			}
		case <-getKeyNumsTicker.C:
			preLen := keyNums
			keyNums, err = GetKeyNumsFromRedis(conn)
			if err != nil {
				log.Println(err)
				keyNums = preLen
				continue
			}

		}
	}
}

func GenerateKeyToRedis(conn redis.Conn) error {
	key, err := GenerateKey(conn)
	if err != nil {
		return err
	}

	_, err = conn.Do("RPUSH", "keys", key)
	if err != nil {
		return err
	}

	return nil
}
func GetKeyFromRedis(conn redis.Conn) (string, error) {
	key, err := conn.Do("RPOP", "keys")
	if err != nil {
		return "", err
	}

	return string(key.([]uint8)), nil
}

func GetKeyNumsFromRedis(conn redis.Conn) (int64, error) {
	result, err := conn.Do("LLEN", "keys")
	if err != nil {
		return -1, err
	}

	return result.(int64), err
}
