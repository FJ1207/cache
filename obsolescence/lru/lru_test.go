package lru

import (
	"fmt"
	"reflect"
	"testing"
)

// Str是新类型
type Str string

//
func (a Str) Len() int {
	return len(a)
}

//testGet
func TestGet(t *testing.T) {
	c := New(int64(0), nil)
	// fmt.Printf("error\n")
	c.Add("key01", Str("value01"))
	// TODO 了解类型断言
	if v, ok := c.Get("key01"); !ok || v.(Str) != "value01" {
		fmt.Printf("error")
	}
	if _, ok := c.Get("key02"); ok {
		fmt.Printf("error")
	}
}

//test Removeoldest
func TestRemoveOldest(t *testing.T) {
	//c := New(int64(cap),nil)
	k1, k2, k3 := "K1", "K2", "K3"
	v1, v2, v3 := "V1", "V2", "V3"
	d := len(k1 + k2 + v1 + v2)
	c := New(int64(d), nil)
	c.Add(k1, Str(v1))
	c.Add(k2, Str(v2))
	c.Add(k3, Str(v3))
	if _, ok := c.Get(k1); ok || c.Len() != 2 {
		fmt.Printf("error")
	}
}

// 测试回调
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", Str("123456"))
	lru.Add("k2", Str("k2"))
	lru.Add("k3", Str("k3"))
	lru.Add("k4", Str("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
