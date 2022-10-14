package lfu

import (
	"container/list"
)

type LFUCache struct {
	nodes    map[int]*list.Element
	lists    map[int]*list.List
	capacity int
	min      int
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		nodes:    make(map[int]*list.Element),
		lists:    make(map[int]*list.List),
		capacity: capacity,
		min:      0,
	}
}

func (this *LFUCache) Get(key int) int {
	currentEle, ok := this.nodes[key]
	if !ok {
		return -1
	}

	currentNode := currentEle.Value.(*node)

	this.lists[currentNode.frequency].Remove(currentEle)
	if this.lists[currentNode.frequency].Len() == 0 { //移除后清0 删除
		delete(this.lists, currentNode.frequency)
		if this.min == currentNode.frequency {
			this.min = currentNode.frequency + 1
		}
	}

	currentNode.frequency++

	if _, ok := this.lists[currentNode.frequency]; !ok {
		this.lists[currentNode.frequency] = list.New()
	}

	newList := this.lists[currentNode.frequency]
	newElement := newList.PushFront(currentNode)
	this.nodes[key] = newElement

	return currentNode.value
}

func (this *LFUCache) Put(key int, value int) {
	//检查是否已经存在
	if currentEle, ok := this.nodes[key]; ok {
		currentNode := currentEle.Value.(*node)
		currentNode.value = value
		this.Get(key)
		return
	}
	//如果不存在检查当前容量是否满了
	if this.capacity == len(this.nodes) {
		minList := this.lists[this.min]
		lastNode := minList.Back()
		delete(this.nodes, lastNode.Value.(*node).key)
		minList.Remove(lastNode)
		if minList.Len() == 0 {
			delete(this.lists, this.min)
			this.UpdateMin() //更新min
		}
	}

	//添加新节点

	addNode := &node{
		key:       key,
		value:     value,
		frequency: 1,
	}

	if _, ok := this.lists[1]; !ok { //存在为1的链表
		this.lists[1] = list.New()
		this.min = 1
	}

	newList := this.lists[1]
	newNode := newList.PushFront(addNode)
	this.nodes[key] = newNode
}

func (this *LFUCache) UpdateMin() {
	m := 0
	flg := true

	for key, _ := range this.lists {
		if flg {
			m = key
		}

		if m > key {
			m = key
		}
	}

	this.min = m
}

func (this *LFUCache) ThisLen() int {
	return len(this.nodes)
}
