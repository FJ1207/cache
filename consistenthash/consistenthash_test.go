package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3,func(key []byte) uint32 { //自定义一个哈希算法
		i,_ := strconv.Atoi(string(key))
		return uint32(i)  //哈希值i
	})

	hash.Add("6","4","2") //添加真实节点

	testCase := map[string]string { //虚拟节点：真实节点
		"2": "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k,v := range testCase {
		if hash.Get(k) != v {
			t.Errorf("请求%s ，真实%s", k ,v )
		}
	}

}

