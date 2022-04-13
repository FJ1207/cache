package sy

import (
	"fmt"
	"sync"
	"time"
)

var already = make(map[int]bool, 0) //slice 0为size
var m sync.Mutex                    //锁

func printonce(num int) {
	m.Lock() //上锁 如果一个goroutine进入则上锁，只能有一个进入，则避免重复打印
	if _, exist := already[num]; !exist {
		fmt.Println(num)
	}
	already[num] = true
	m.Unlock() //上锁
}

// func printOnce(num int) {
// 	m.Lock()
// 	defer m.Unlock()
// 	if _, exist := set[num]; !exist {
// 		fmt.Println(num)
// 	}
// 	set[num] = true
// }
func testprint() { //测试打印
	for i := 0; i < 10; i++ {
		go printonce(100) //并发
	}
	time.Sleep(time.Second)

}
