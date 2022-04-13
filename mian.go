package main

//假设并发打印操作，则会出现重复打印
import (
	"fmt"
	"sync"
	"time"
)

var already = make(map[int]bool, 0) //slice 0为size
var m sync.Mutex

func printonce(num int) {
	m.Lock() //上锁
	if _, exist := already[num]; !exist {
		fmt.Println(num)
	}
	already[num] = true
	m.Unlock() //上锁
}

func main() {
	for i := 0; i < 10; i++ {
		go printonce(100) //并发
	}

	time.Sleep(time.Second)

}

//互斥锁Lock()和Unlock()
