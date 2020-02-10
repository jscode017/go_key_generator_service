package keygenerator

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kjk/betterguid"
	"log"
	"sync"
	"time"
)

type KeyGeneratorService struct {
	Keys []string
	sync.Mutex
}

func NewKeyGenerator() *KeyGeneratorService {
	return &KeyGeneratorService{
		Keys: make([]string, 0),
	}
}
func GenerateKey(conn redis.Conn) (string, error) {
	key := betterguid.New()[5:15] //get a string by the length 10
	keyExist, err := CheckKeyExist(conn, key)
	if err != nil {
		return "", err
	}
	for keyExist {
		key = betterguid.New()[3:13] //get a string by the length 10
		keyExist, err = CheckKeyExist(conn, key)
		if err != nil {
			return "", err
		}
	}
	err = MarkKeyAsUsed(conn, key) // it is ok for key lost but not for using conflict keys, and it is more efficient not to do database operation when getting keys
	if err != nil {
		return "", err
	}

	return key, nil

}

func (kgs *KeyGeneratorService) Generate() {
	conn, err := NewRedisConn()
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(time.Millisecond * 10)
	for {
		select {
		case <-ticker.C:
			kgs.Lock()
			if len(kgs.Keys) >= 300 {
				kgs.Unlock()
				time.Sleep(100 * time.Millisecond)
				continue
			}
			key, err := GenerateKey(conn)
			if err != nil {
				log.Println(err)
				kgs.Unlock()
				continue
			}
			kgs.Keys = append(kgs.Keys, key)
			kgs.Unlock()

		}
	}
}

func (kgs *KeyGeneratorService) GetKey() string {
	kgs.Lock()
	defer kgs.Unlock()
	key := kgs.Keys[len(kgs.Keys)-1]
	kgs.Keys = kgs.Keys[:len(kgs.Keys)-1]
	return key
}
