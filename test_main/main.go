package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
		fmt.Println("main завершен")
		DB.Wait()
		fmt.Println("DB завершен")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	srv := startHttpServer()
	httpStopped := make(chan struct{})
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		DB.Close()
		httpStopped <- struct{}{}
	}()
	<-httpStopped
	fmt.Println("http stopped")
	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}

}

func startHttpServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		count, err := DB.Count()
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Println("Count: ", err)
			return
		}
		io.WriteString(w, strconv.Itoa(count))
	})
	http.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		log.Println("keys run")
		keys, err := DB.Keys()
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Println("Keys: ", err)
			return
		}
		log.Println("keys send")
		for _, val := range keys {
			io.WriteString(w, val+"\n")
		}
		log.Println("keys ok")
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	return srv
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
				fmt.Println("FPS: ", i, "Len: ", count)
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
