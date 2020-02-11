package main

import (
	KGS "github.com/jscode017/go_key_generator_service"
	"log"
	"time"
)

func main() {
	go KGS.GenerateKeysToRedis()
	conn, err := KGS.NewRedisConn()
	if err != nil {
		log.Fatal(err)
	}
	getKeyNumTicker := time.NewTicker(time.Second * 1)
	getKeyTicker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-getKeyNumTicker.C:
			num, err := KGS.GetKeyNumsFromRedis(conn)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(num)

		case <-getKeyTicker.C:
			for i := 0; i < 5; i++ {
				key, err := KGS.GetKeyFromRedis(conn)
				if err != nil {
					log.Println(err)
					continue
				}

				log.Println(key)
			}
		}
	}
}
