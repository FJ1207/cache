package mutexincludobsolesence

import (
	"fmt"
	"log"
	"sync"

	"github.com/FJ1207/cache/obsolescence/lru"
	//	"github.com/FJ1207/cache/sy"
)

type ByteView struct { //slice[]大小
	b []byte
}

func (v ByteView) Len() int { //slice切片长度
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte { //返回拷贝切片
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v ByteView) String() string {
	return string(v.b)
}

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64 // 最大容量
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil) //延迟初始化，延迟至使用的时候
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

type Getter interface {
	Get(key string) ([]byte, error) //函数类型实现某一个接口，称之为接口型函数
}

type GetterFunc func(key string) ([]byte, error) //函数

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)//GetterFunc回调,接口型函数
}



//最核心的Group数据结构
type Group struct { //Group为缓存
	name string
	getter Getter //缓存未命中时获取源数据的回调
	mainCache cache //一开始实现的并发缓存
}	

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) //groups map string *Group

)

func NewGroup(name string,cacheBytes int64 ,getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes:cacheBytes },
	}
	groups[name]=g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	g := groups[name]
	mu.Unlock()
	return g
}

//核心方法
//Get 方法实现了上述所说的流程 ⑴ 和 ⑶。
//流程 ⑴ ：从 mainCache 中查找缓存，如果存在则返回缓存值。
//流程 ⑶ ：缓存不存在，则调用 load 方法，load 调用 getLocally，getLocally 调用用户回调函数 g.getter.Get() 
//获取源数据，并且将源数据添加到缓存 mainCache 中（通过 populateCache 方法）
func (g *Group) Get(key string)(ByteView,error) {
	if v, ok := g.mainCache.get(key); ok { //可以找得到
		log.Println("[GeeCache] hit")
		return v, nil
	}
	if key == "" { //key为空
		return ByteView{}, fmt.Errorf("key is required")
	} else {
		return g.load(key) //未缓存，则load回调获取源数据，存入
	}
}

func (g *Group) load(key string) (ByteView,error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView,error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{},err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string , value ByteView) {
	g.mainCache.add(key, value)
}

// func returnAge() (age int) {
// 	age = 30
// 	return
// }

// func returnAge2() int {
// 	age := 30
// 	return age
// }