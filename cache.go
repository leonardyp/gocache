package gocache

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

var (
	error_empty = errors.New("error empty")
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
	setted     bool
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
func (this *cache) GetString(key string) (string, error) {
	val := this.Get(key)
	if val == nil {
		return "", error_empty
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.String:
		return val.(string), nil
	}
	return "", errors.New("not string")
}
func (this *cache) GetInt(key string) (int, error) {
	val := this.Get(key)
	if val == nil {
		return 0, error_empty
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.Int:
		return val.(int), nil

	}
	return 0, errors.New("not int")
}
func (this *cache) GetInt64(key string) (int64, error) {
	val := this.Get(key)
	if val == nil {
		return 0, error_empty
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.Int64:
		return val.(int64), nil

	}
	return 0, errors.New("not int64")
}
func (this *cache) Del(key string) {
	this.Lock()
	delete(this.Items, key)
	this.Unlock()
}
func (this *cache) ItemsCount() int {
	this.RLock()
	l := len(this.Items)
	this.RUnlock()
	return l
}

var defaultGCInterval = 3 * time.Second

func (this *cache) SetGcInterval(inter time.Duration) {
	defer func() { this.setted = true }()
	if this.setted {
		return
	}
	if inter > 0 {
		this.gcInterval = inter
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
	fmt.Println("test webhook")
}

func (this *cache) Clear() {
	for key, item := range this.Items {
		if (*item).Expired() {
			delete(this.Items, key)
		}
	}
}
