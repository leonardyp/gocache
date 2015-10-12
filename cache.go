package gocache

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	Object interface{}
	Expire *time.Time
}

var (
	t      time.Time
	stopGC = make(chan bool, 1)
)

func (item *Item) Expired() bool {
	if item.Expire == nil {
		return false
	}
	return item.Expire.Before(time.Now())
}

type Cache struct {
	sync.RWMutex
	Items map[string]*Item
}

func NewCache() *Cache {
	c := &Cache{
		Items: make(map[string]*Item),
	}
	return c
}
func (this *Cache) Set(key string, val interface{}, ttl time.Duration) {
	this.Lock()
	t = time.Now().Add(ttl)
	this.Items[key] = &Item{
		Object: val,
		Expire: &t,
	}
	this.Unlock()
}
func (this *Cache) Get(key string) interface{} {
	this.RLock()
	obj, _ := this.Items[key]
	this.RUnlock()
	if obj == nil {
		return nil
	}
	return obj.Object
}
func (this *Cache) Del(key string) {
	delete(this.Items, key)
}
func (this *Cache) ItemsCount() int {
	this.RLock()
	l := len(this.Items)
	this.RUnlock()
	return l
}

var defaultGCInterval = 1 * time.Second

func SetGcInterval(inter time.Duration) {
	fmt.Println(inter)
	defaultGCInterval = inter
}
func (this *Cache) startGC(stop chan bool) {
	for {
		select {
		case <-stop:
			return
		case <-time.Tick(defaultGCInterval):
			clear(this)
		}
	}
}

func (this *Cache) StartGc() {
	go this.startGC(stopGC)
}
func (this *Cache) StopGc() {
	stopGC <- true
	this.startGC(stopGC)
}

func clear(c *Cache) {
	for key, item := range c.Items {
		if (*item).Expired() {
			delete(c.Items, key)
		}
	}
}
