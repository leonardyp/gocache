package gocache

import (
	"fmt"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	c := NewCache()
	c.Set("a", struct{ Name string }{"leo"}, 2*time.Second)
	fmt.Printf("%#v\n", c.Get("a"))
}
func TestGet(t *testing.T) {
	c := NewCache()
	fmt.Printf("%#v\n", c.Get("a"))
}

func BenchmarkCacheGet(b *testing.B) {
	b.StopTimer()
	c := NewCache()
	c.Set("foo", "bar", time.Second*0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Get("foo")
	}
}
func BenchmarkCacheSet(b *testing.B) {
	b.StopTimer()
	c := NewCache()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Set("foo", "bar", time.Second*0)
	}
}

func TestGc(t *testing.T) {
	c := NewCache()
	c.Set("a", 1332424, time.Second*1)
	c.StartGc()
	fmt.Printf("%v\n", c.Get("a"))
	time.Sleep(time.Second * 2)
	fmt.Printf("%v\n", c.Get("a"))
}

func TestGcInterval(t *testing.T) {
	c := NewCache()
	c.Set("a", 1332424, time.Second*3)
	SetGcInterval(time.Second * 2)
	c.StartGc()
	fmt.Printf("%v\n", c.Get("a"))
	time.Sleep(time.Second * 2)
	fmt.Printf("%v\n", c.Get("a"))
}
func TestGcStop(t *testing.T) {
	c := NewCache()
	c.Set("a", 1332424, time.Second*1)
	c.StopGc()
	time.Sleep(2 * time.Second)
	fmt.Printf("%v\n", c.Get("a"))
}
