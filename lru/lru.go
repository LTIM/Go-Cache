package lru

import "container/list"

type LRUCache struct {
	MaxBytes  int64 //最大存储空间
	NBytes    int64 //已使用的容量
	ValueList *list.List
	Cache     map[string]*list.Element

	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数
}

type Entry struct {
	Key    string
	Evalue Value //任意实现len的结构体
}

type Value interface {
	Len() int //当前值所占用的内存大小
}

func NewLruCache(MaxBytes int64, onEvicted func(string, Value)) *LRUCache { //构建缓存 -》相当于构造函数
	return &LRUCache{
		MaxBytes:  MaxBytes,
		ValueList: list.New(),
		Cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (l *LRUCache) Get(key string) (value Value, ok bool) {

	if ele, ok := l.Cache[key]; ok {
		l.ValueList.MoveToFront(ele) //将该节点移动到队首
		kv := ele.Value.(*Entry)
		return kv.Evalue, true
	}

	return
}

func (l *LRUCache) RemoveOldest() { // 移除队尾的元素
	ele := l.ValueList.Back()

	if ele != nil {
		l.ValueList.Remove(ele)
		kv := ele.Value.(*Entry)
		delete(l.Cache, kv.Key)
		l.NBytes -= int64(len(kv.Key)) + int64(kv.Evalue.Len())

		if l.OnEvicted != nil {
			l.OnEvicted(kv.Key, kv.Evalue)
		}

	}

}

func (l *LRUCache) Add(key string, value Value) { //插入或者修改

	if ele, ok := l.Cache[key]; ok {
		l.ValueList.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		l.NBytes += int64(value.Len()) - int64(kv.Evalue.Len())
		kv.Evalue = value
	} else {
		ele := l.ValueList.PushFront(&Entry{key, value})
		l.Cache[key] = ele
		l.NBytes += int64(len(key)) + int64(value.Len())
	}

	for l.MaxBytes != 0 && l.MaxBytes < l.NBytes { //判断插入后容量是否饱和
		l.RemoveOldest()
	}

}

func (l *LRUCache) Len() int {
	return l.ValueList.Len()
}
