package connpool

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_PopTimeout(t *testing.T) {
	assert := assert.New(t)

	maxActiveNum := 1
	pool := NewConnectionPool(
		maxActiveNum,
		0,
		time.Second*time.Duration(1),
		3,
		func() (interface{}, error) {
			log.Println("New handler")
			return true, nil
		},
		func(c interface{}) {
			log.Println("Destroy handler")
		},
		func(c interface{}) {
			log.Println("Clear handler")
		},
	)

	c, err := pool.Pop()
	assert.Nil(err)

	err = pool.Push(c)
	assert.Nil(err)

	c1, err := pool.Pop()
	assert.Nil(err)

	_, err = pool.Pop()
	assert.NotNil(err, err.Error())

	err = pool.Push(c1)
	assert.Nil(err)
}

func Test_Muti(t *testing.T) {
	assert := assert.New(t)

	maxActiveNum := 2
	pool := NewConnectionPool(
		maxActiveNum,
		0,
		time.Microsecond*time.Duration(100),
		3,
		func() (interface{}, error) {
			log.Println("New handler")
			return true, nil
		},
		func(c interface{}) {
			log.Println("Destroy handler")
		},
		func(c interface{}) {
			log.Println("Clear handler")
		},
	)

	c, err := pool.Pop()
	assert.Nil(err)

	err = pool.Push(c)
	assert.Nil(err)

	c1, err := pool.Pop()
	assert.Nil(err)

	time.Sleep(1 * time.Second)

	c2, err := pool.Pop()
	assert.Nil(err)

	c3, err := pool.Pop()
	assert.NotNil(err)

	err = pool.Push(c1)
	assert.Nil(err)

	err = pool.Push(c2)
	assert.Nil(err)

	err = pool.Push(c3)
	assert.NotNil(err)
}

func Test_AllConnect(t *testing.T) {
	assert := assert.New(t)

	maxActiveNum := 0 //不限制
	pool := NewConnectionPool(
		maxActiveNum,
		0,
		time.Microsecond*100,
		1,
		func() (interface{}, error) {
			log.Println("New handler")
			return true, nil
		},
		func(c interface{}) {
			log.Println("Destroy handler")
		},
		func(c interface{}) {
			log.Println("Clear handler")
		},
	)

	c, err := pool.Pop()
	assert.Nil(err)

	err = pool.Push(c)
	assert.Nil(err)

	c1, err := pool.Pop()
	assert.Nil(err)

	time.Sleep(1 * time.Second)

	c2, err := pool.Pop()
	assert.Nil(err)

	c3, err := pool.Pop()
	assert.Nil(err)

	err = pool.Push(c1)
	assert.Nil(err)

	err = pool.Push(c2)
	assert.Nil(err)

	err = pool.Push(c3)
	assert.Nil(err)

	time.Sleep(1 * time.Second)

	c, err = pool.Pop()
	assert.Nil(err)

	err = pool.Push(c)
	assert.Nil(err)
}

func Test_GetCounter(t *testing.T) {
	assert := assert.New(t)
	maxActiveNum := 0 //不限制
	pool := NewConnectionPool(
		maxActiveNum,
		0,
		time.Microsecond*100,
		1,
		func() (interface{}, error) {
			log.Println("New handler")
			return true, nil
		},
		func(c interface{}) {
			log.Println("Destroy handler")
		},
		func(c interface{}) {
			log.Println("Clear handler")
		},
	)

	assert.Equal(pool.GetActiveNum(), 0)
	assert.Equal(pool.GetIdleNum(), 0)
	assert.Equal(pool.GetWaitNum(), 0)

}

func Test_ClearPool(t *testing.T) {
	assert := assert.New(t)
	maxActiveNum := 2 //不限制
	pool := NewConnectionPool(
		maxActiveNum,
		0,
		time.Microsecond*100,
		1,
		func() (interface{}, error) {
			log.Println("New handler")
			return true, nil
		},
		func(c interface{}) {
			log.Println("Destroy handler")
		},
		func(c interface{}) {
			log.Println("Clear handler")
		},
	)

	c1, err := pool.Pop()
	assert.Nil(err)

	time.Sleep(1 * time.Second)

	c3, err := pool.Pop()
	assert.Nil(err)

	_, err = pool.Pop()
	assert.NotNil(err)

	err = pool.Push(c1)
	assert.Nil(err)

	assert.Equal(pool.GetActiveNum(), 1)
	assert.Equal(pool.GetIdleNum(), 1)
	assert.Equal(pool.GetWaitNum(), 0)

	err = pool.Push(c3)
	assert.Nil(err)

	assert.Equal(pool.GetActiveNum(), 0)
	assert.Equal(pool.GetIdleNum(), 2)
	assert.Equal(pool.GetWaitNum(), 0)

	pool.ClearPool()
}
