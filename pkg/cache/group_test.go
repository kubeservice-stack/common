package cache

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoLFU(t *testing.T) {
	assert := assert.New(t)

	var g Group
	g.plugin = New(32).Setting()
	v, _, err := g.Do("key", func() (interface{}, error) {
		log.Println(g)
		return "bar", nil
	}, true)

	got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"
	log.Println(got, want)

	assert.Equal(want, got)
	assert.Nil(err)

}
func TestDoLRU(t *testing.T) {
	assert := assert.New(t)

	var g Group
	g.plugin = New(32).EvictType(LRU).Setting()
	v, _, err := g.Do("key", func() (interface{}, error) {
		log.Println(g, g.plugin)
		return "bar", nil
	}, true)

	got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"
	log.Println(got, want)

	assert.Equal(want, got)
	assert.Nil(err)

	v, _, err = g.Do("key", func() (interface{}, error) {
		log.Println(g, g.plugin)
		return g, nil
	}, true)

	assert.Equal(v, g)
	assert.Nil(err)

}

func TestDoErrLFU(t *testing.T) {
	assert := assert.New(t)

	var g Group
	g.plugin = New(32).Setting()
	someErr := errors.New("Some error")

	v, _, err := g.Do("key", func() (interface{}, error) {
		log.Println("dongjiang")
		return nil, someErr
	}, true)

	assert.Equal(err, someErr)
	assert.Nil(v)
}

func TestDoErrLRU(t *testing.T) {
	assert := assert.New(t)

	var g Group
	g.plugin = New(32).EvictType(LRU).Setting()
	someErr := errors.New("Some error")

	v, _, err := g.Do("key", func() (interface{}, error) {
		log.Println("dongjiang")
		return nil, someErr
	}, true)

	assert.Equal(err, someErr)
	assert.Nil(v)
}

func TestDoDupSuppressLFU(t *testing.T) {
	assert := assert.New(t)

	var g Group
	g.plugin = New(32).Setting()
	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var count = 0
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			log.Println("count:", count)
			count++
			v, _, err := g.Do("key", fn, true)
			assert.Nil(err)
			assert.Equal(v, "bar")
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	got := atomic.LoadInt32(&calls)

	assert.Equal(got, int32(1))
}

func TestDoDupSuppressLRU(t *testing.T) {
	assert := assert.New(t)

	var g Group
	g.plugin = New(32).EvictType(LRU).Setting()
	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var count = 0
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			log.Println("count:", count)
			count++
			v, _, err := g.Do("key", fn, true)
			assert.Nil(err)
			assert.Equal(v, "bar")
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	got := atomic.LoadInt32(&calls)

	assert.Equal(got, int32(1))
}
