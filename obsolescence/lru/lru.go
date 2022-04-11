package lru

import (
	"container/list"
)

//struct:字典+双向链表
type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element      //字典：key+链表指针
	OnEvicted func(key string, value Value) //回调函数

}

type entry struct {
	key   string
	value Value
}
type Value interface {
	Len() int
}

// 构造实例化cache函数
func New(maxBytes int64, OnExicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnExicted,
	}

}

//查找功能
func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)    //ele移动至队尾
		kv := ele.Value.(*entry) //断言
		return kv.value, true
	}
	return nil, false
}

//删除lru
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() //取首节点
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //长度更改
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) //被移除则回调
		}
	}
}

//添加（修改）
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value //修改
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes { //如果不够则lru
		c.RemoveOldest()
	}

}

func (c *Cache) Len() int {
	return c.ll.Len()
}
