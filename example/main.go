package main

import (
	KGS "github.com/jscode017/go_key_generator_service"
	"log"
	"time"
)

func main() {
	kgs := KGS.NewKeyGenerator()
	go kgs.Generate()
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			log.Println(kgs.Keys)
			log.Println(len(kgs.Keys))
			for i := 0; i < 5; i++ {
				key := kgs.GetKey()
				log.Printf(key)
			}
			log.Println("")
			log.Println(len(kgs.Keys))

		}
	}
}
