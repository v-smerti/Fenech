package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"bitbucket.org/tdmv/fenech"
	uid "github.com/satori/go.uuid"
)

var DB *fenech.Fenech

func init() {
	db, err := fenech.New("./dir")
	if err != nil {
		log.Fatal("Mew: ", err)
	}
	DB = db
	log.Println(`Run db`)
}
func main() {
	fps()
}

func fps() {
	var i int64 = 0
	r := new(sync.RWMutex)
	go func() {
		var localI int64 = 0
		count := int64(DB.Count())
		for {
			time.Sleep(2 * time.Second)
			r.RLock()
			if localI == 0 {
				count = count + i
				fmt.Println("FPS 1: ", i, "Len: ", count)
			} else {
				fps := i - localI
				count = count + fps
				fmt.Println("FPS: ", fps, "Len: ", count)
			}
			localI = i
			r.RUnlock()
		}
	}()
	q := 0
	for {
		q++
		if q > 30 {
			break
		}
		log.Println(q)
		go func() {
			for {
				r.Lock()
				i = i + 1
				r.Unlock()
				if err := DB.Set(uid.NewV4().String(), uid.NewV4().Bytes()); err != nil {
					log.Fatal(err)
				}
			}
		}()
	}
	for {
		r.Lock()
		i = i + 1
		r.Unlock()
		if err := DB.Set(uid.NewV4().String(), uid.NewV4().Bytes()); err != nil {
			log.Fatal(err)
		}
	}
}
