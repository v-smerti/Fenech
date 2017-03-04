package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"bitbucket.org/tdmv/fenech"
	uid "github.com/satori/go.uuid"
)

var DB *fenech.Fenech
var WG = new(sync.WaitGroup)

func init() {
	db, err := fenech.New("./dir")
	if err != nil {
		log.Fatal("Mew: ", err)
	}
	DB = db
	log.Println(`Run db`)
}
func main() {
	defer func() {
		WG.Wait()
		fmt.Println("main упал")
		//fmt.Println("ожидаю завершения всех операций DB")
		DB.Wait()
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		DB.Close()
	}()

	fps()
}

func fps() {
	var i int64 = 0
	r := new(sync.RWMutex)
	WG.Add(1)
	go func() {

		defer WG.Done()
		var localI int64 = 0
		counts, err := DB.Count()
		if err != nil {
			panic(err.Error())
		}
		count := int64(counts)
		for {
			time.Sleep(1 * time.Second)
			r.RLock()

			if localI == 0 {
				count = count + i
				fmt.Println("FPS 1: ", i, "Len: ", count)
			} else {
				fps := i - localI
				count = count + fps
				if fps == 0 {
					break
				}
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
		WG.Add(1)
		go func() {

			defer WG.Done()
			for {
				r.Lock()
				i = i + 1
				r.Unlock()
				if err := DB.Set(uid.NewV4().String(), SecureRandomBytes(10)); err != nil {
					log.Println(err)
					break
				}
			}
		}()
	}

	fmt.Println("main")

}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52 possibilities
	letterIdxBits = 6                                                      // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
)

func SecureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}
