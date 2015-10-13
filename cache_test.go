package gocache

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	c := NewCache()
	c.Set("a", struct{ Name string }{"leo"}, 2*time.Second)
	c.SetGcInterval(1 * time.Second)

	Convey("TestSet", t, func() {
		So(fmt.Sprint(c.Get("a")), ShouldEqual, `{leo}`)
		time.Sleep(3 * time.Second)
		So(c.Get("a"), ShouldBeNil)
	})
}
func TestGet(t *testing.T) {
	c := NewCache()

	Convey("TestGet", t, func() {
		So(c.Get("a"), ShouldBeNil)
	})
}

func BenchmarkCacheGet(b *testing.B) {
	b.StopTimer()
	c := NewCache()
	c.Set("foo", "bar", 0*time.Second)
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
		c.Set("foo", "bar", 0*time.Second)
	}
}

func TestGc(t *testing.T) {
	c := NewCache()
	c.Set("a", "gc demo", 1*time.Second)
	c.SetGcInterval(1 * time.Second)

	Convey("TestGc", t, func() {
		So(c.Get("a"), ShouldEqual, "gc demo")
		time.Sleep(2 * time.Second)
		So(c.Get("a"), ShouldBeNil)
	})
}

func TestGcInterval(t *testing.T) {
	c := NewCache()
	c.Set("a", "gcInterval demo", 3*time.Second)
	c.SetGcInterval(2 * time.Second)

	Convey("TestGcInterval", t, func() {
		So(c.Get("a"), ShouldEqual, "gcInterval demo")
		time.Sleep(2 * time.Second)
		So(c.Get("a"), ShouldEqual, "gcInterval demo")
	})
}
func TestGcStop(t *testing.T) {
	c := NewCache()
	c.Set("a", "gc stop demo", 1*time.Second)
	c.SetGcInterval(2 * time.Second)
	c.StopGc()
	time.Sleep(2 * time.Second)

	Convey("TestGcStop", t, func() {
		So(c.Get("a"), ShouldEqual, "gc stop demo")
	})
}

func TestNewCacheWithGcInterval(t *testing.T) {
	c := NewCacheWithGcInterval(2 * time.Second)
	c.SetGcInterval(0 * time.Second)
	c.Set("a", "NewCacheWithGcInterval demo", 1*time.Second)
	time.Sleep(3 * time.Second)

	Convey("TestNewCacheWithGcInterval", t, func() {
		So(c.Get("a"), ShouldBeNil)
		c.Set("a", "NewCacheWithGcInterval demo", 1*time.Second)
		c.StopGc()
		time.Sleep(2 * time.Second)
		So(c.Get("a"), ShouldEqual, "NewCacheWithGcInterval demo")
	})
}
func TestDefaultGcInterval(t *testing.T) {
	c := NewCache()
	c.SetGcInterval(0)
	c.Set("a", "DefaultGcInterval demo", 2*time.Second)
	Convey("TestDefaultGcInterval", t, func() {
		So(c.Get("a"), ShouldEqual, "DefaultGcInterval demo")
		time.Sleep(4 * time.Second)
		So(c.Get("a"), ShouldBeNil)
	})
}
