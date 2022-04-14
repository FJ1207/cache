package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FJ1207/cache/concurrent/mutexincludobsolesence"
	"github.com/FJ1207/cache/httpserver"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	mutexincludobsolesence.NewGroup("scores", 2<<10, mutexincludobsolesence.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:8888"
	peers := httpserver.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}