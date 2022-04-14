package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//map 实现虚拟节点和现实节点之间的映射

type Hash func(data []byte) uint32 //计算哈希值函数

type Map struct {
	hash     Hash           //函数
	keys     []int          //哈希环
	replicas int            //虚拟节点倍数
	hashMap  map[int]string //虚拟节点与真实节点之间的映射，键是虚拟节点的哈希值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map { //自定义虚拟节点倍数和Hash函数
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//真实节点的add()方法
func (m *Map) Add(keys ...string) { //添加真实节点
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ { //对应每个真实节点key创建replicas个虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //strconv.Itoa(i) + key是指虚拟节点名称
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}

	}
	sort.Ints(m.keys) //对环上的哈希值排序
}

//选择节点的Get()方法
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 { //环为空
		return ""
	}

	hash := int(m.hash([]byte(key))) //

	idx := sort.Search(len(m.keys), func(i int) bool { //顺时针找到第一个虚拟节点下标
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]] //idx 确定最近的keys[i]哈希值，根据哈希值确定真实节点名称

}
