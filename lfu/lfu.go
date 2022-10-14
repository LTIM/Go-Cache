package lfu

import (
	"container/list"
	"unsafe"
)

//  Least Frequently Used 即最不经常最少使用

type LfuCache struct {
	nodes map[int]*list.Element
	lists map[int]*list.List

	MaxBytes     int
	currentBytes int
	Min          int
}

type node struct {
	key       int
	value     int
	frequency int
}

func NewLfuCache(cap int) *LfuCache {
	return &LfuCache{
		MaxBytes: cap,
		lists:    make(map[int]*list.List),
		nodes:    make(map[int]*list.Element),
		Min:      0,
	}
}

func (this *LfuCache) Get(key int) int {
	value, ok := this.nodes[key]

	if !ok {
		return -1
	}

	currentNode := value.Value.(*node)
	this.lists[currentNode.frequency].Remove(value) // 将当前node从原频次链表中移除
	currentNode.frequency++

	if _, ok := this.lists[currentNode.frequency]; !ok { //判断当前节点的频次是否在map中存在
		this.lists[currentNode.frequency] = list.New() //如果没有则生成
	}

	newList := this.lists[currentNode.frequency] //获取当前频次结果的的链表
	newNode := newList.PushFront(currentNode)    //把当前节点移动到链表的前面
	this.nodes[key] = newNode                    //更新该链表在nodes hashmap内element的地址

	if currentNode.frequency-1 == this.Min && this.lists[currentNode.frequency-1].Len() == 0 { //判断原frequency是否为最小值
		this.Min++
	}

	return currentNode.value

}

func (this *LfuCache) Put(key int, value int) error {
	var err error

	if currentValue, ok := this.nodes[key]; ok { //判断当前节点是否存在
		currentNode := currentValue.Value.(*node)
		currentNode.value = value
		this.Get(key)
		return err
	}

	for this.currentBytes == this.MaxBytes || (this.currentBytes+int(unsafe.Sizeof(key))) > this.MaxBytes { // 判断当前节点是否已满
		//如果满了需要删除
		this.RemoveOldest()
	}

	return err
}

func (this *LfuCache) RemoveOldest() {
	currentList := this.lists[this.Min]
	lastNode := currentList.Back()
	delete(this.nodes, lastNode.Value.(*node).key)
	currentList.Remove(lastNode)

	if currentList == nil { //删除后是否更新min值 删完了链表为空怎么处理
		delete(this.lists, this.Min) //移除

		if len(this.lists) == 0 { //如果为空将min置为原始值0
			this.Min = 0
			return
		}

		thisMin := 0
		flg := true

		for key, _ := range this.lists {
			if flg { //考虑竞态
				thisMin = key
			} else {
				if thisMin > key {
					thisMin = key
				}
			}
		}
		this.Min = thisMin

	}

}
