package ratelimiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLimiter(t *testing.T) {
	assert := assert.New(t)
	aa, has := GetLimiter("RATELIMITER")
	assert.True(has)
	assert.NotNil(aa)

	bb, has := GetLimiter("empty")
	assert.False(has)
	assert.Nil(bb)
}

func Test_GetDefaultLimiter(t *testing.T) {
	assert := assert.New(t)

	bb := GetDefaultLimiter()
	assert.NotNil(bb)
}

func Test_HasRegister(t *testing.T) {
	assert := assert.New(t)

	ok := HasRegister("RATELIMITER")
	assert.True(ok)

	ok = HasRegister("empty")
	assert.False(ok)
}
