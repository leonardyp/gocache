package gocache

import (
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

type cache struct {
	sync.RWMutex
	gcInterval time.Duration
	Items      map[string]*Item
}

func NewCache() *cache {
	return &cache{
		Items: make(map[string]*Item),
	}
}
func NewCacheWithGcInterval(gcInterval time.Duration) *cache {
	return &cache{
		Items:      make(map[string]*Item),
		gcInterval: gcInterval,
	}
}
func (this *cache) Set(key string, val interface{}, ttl time.Duration) {
	this.Lock()
	t = time.Now().Add(ttl)
	this.Items[key] = &Item{
		Object: val,
		Expire: &t,
	}
	this.Unlock()
}
func (this *cache) Get(key string) interface{} {
	this.RLock()
	obj, _ := this.Items[key]
	this.RUnlock()
	if obj == nil {
		return nil
	}
	return obj.Object
}
func (this *cache) GetString(key string) string {
	this.RLock()
	obj, _ := this.Items[key]
	this.RUnlock()
	if obj == nil {
		return ""
	}
	return obj.Object.(string)
}
func (this *cache) GetInt(key string) int {
	this.RLock()
	obj, _ := this.Items[key]
	this.RUnlock()
	if obj == nil {
		return 0
	}
	return obj.Object.(int)
}
func (this *cache) GetInt64(key string) int64 {
	this.RLock()
	obj, _ := this.Items[key]
	this.RUnlock()
	if obj == nil {
		return 0
	}
	return obj.Object.(int64)
}
func (this *cache) Del(key string) {
	delete(this.Items, key)
}
func (this *cache) ItemsCount() int {
	this.RLock()
	l := len(this.Items)
	this.RUnlock()
	return l
}

var defaultGCInterval = 3 * time.Second

func (this *cache) SetGcInterval(inter time.Duration) {
	if inter > 0 {
		this.gcInterval = inter
		go this.startGc()
		return
	}
	if this.gcInterval > 0 {
		go this.startGc()
		return
	}

	go this.Clear()
	this.SetGcInterval(defaultGCInterval)

	return

}
func (this *cache) startGc() {
	for {
		select {
		case <-stopGC:
			return
		case <-time.Tick(this.gcInterval):
			this.Clear()
		}
	}
}
func (this *cache) StopGc() {
	close(stopGC)
}

func (this *cache) Clear() {
	for key, item := range this.Items {
		if (*item).Expired() {
			delete(this.Items, key)
		}
	}
}
