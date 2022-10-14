package main

import (
	"fmt"
	lfuCache "geeCache/lfu"
	lruCache "geeCache/lru"
	"sync"
	"sync/atomic"
)

type String string //对应entry数组的value
var CALL_BACLK_VALUE []string
var INITIALIZED_CALL_BACK_VALUE uint32
var mu sync.Mutex

func (s String) Len() int {
	return len(s)
}

func TOnEvicted(key string, value lruCache.Value) {

	if atomic.LoadUint32((&INITIALIZED_CALL_BACK_VALUE)) != 1 { //原子化设置一个标志
		mu.Lock()
		defer mu.Unlock()
		if INITIALIZED_CALL_BACK_VALUE == 0 {
			CALL_BACLK_VALUE = make([]string, 0)
			atomic.StoreUint32(&INITIALIZED_CALL_BACK_VALUE, 1)
		}
	}

	CALL_BACLK_VALUE = append(CALL_BACLK_VALUE, key)

}

func main() {
	lru := lruCache.NewLruCache(int64(10), nil)
	lru.Add("key1", String("1234"))
	fmt.Println(lru.MaxBytes)
	fmt.Println(lru.NBytes)

	if v, ok := lru.Get("key1"); ok || string(v.(String)) != "1234" {
		fmt.Println("cache hit key1 = 1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		fmt.Println(" cache miss key2 failed")
	}

	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)

	lru1 := lruCache.NewLruCache(int64(cap), TOnEvicted)
	fmt.Println("容量：", lru1.MaxBytes)
	lru1.Add(k1, String(v1))
	lru1.Add(k2, String(v2))
	fmt.Println("current lru1 cap: ", lru1.NBytes)

	lru1.Add(k3, String(v3))

	if _, ok := lru1.Get("key1"); ok || lru1.Len() != 2 {
		fmt.Println("remove faile")
	} else {
		fmt.Println("remove key1")
	}

	vv, ok := lru1.Get("k3")

	if ok {
		fmt.Println(vv)
	} else {
		fmt.Println("not found")
	}

	for _, value := range CALL_BACLK_VALUE {
		fmt.Println("召回的值: ", value)
	}
	fmt.Println("----lfu----")
	// lfu := lfuCache.Constructor(2)
	// lfu.Put(1, 1)
	// fmt.Println("len :", lfu.ThisLen())
	// lfu.Put(2, 2)
	// fmt.Println("len :", lfu.ThisLen())
	// fmt.Println(lfu.Get(1))
	// lfu.Put(3, 3)
	// fmt.Println("len :", lfu.ThisLen())
	// fmt.Println(lfu.Get(2))
	// fmt.Println(lfu.Get(3))
	// lfu.Put(4, 4)
	// fmt.Println("len :", lfu.ThisLen())
	// fmt.Println(lfu.Get(1))
	// fmt.Println(lfu.Get(3))
	// fmt.Println(lfu.Get(4))

	lfu := lfuCache.Constructor(2)
	fmt.Println("get: ", lfu.Get(2))
	lfu.Put(2, 6)
	fmt.Println("len :", lfu.ThisLen())
	fmt.Println("get: ", lfu.Get(1))
	lfu.Put(1, 5)
	fmt.Println("len :", lfu.ThisLen())
	lfu.Put(1, 2)
	fmt.Println("len :", lfu.ThisLen())
	fmt.Println("get: ", lfu.Get(1))
	fmt.Println("get: ", lfu.Get(2))
}
