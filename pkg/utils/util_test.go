package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Min(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Min(1, 1), 1)
	assert.Equal(Min(1, 2), 1)
}

func Test_Max(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Max(1, 1), 1)
	assert.Equal(Max(1, 2), 2)
}

func Test_MimFloat64(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(MinFloat64(1.2, 5.1), 1.2)
	assert.Equal(MinFloat64(1.2, 1.20), 1.2)
}

func Test_MaxFloat64(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(MaxFloat64(1.2, 5.1), 5.1)
	assert.Equal(MaxFloat64(1.2, 1.20), 1.2)
}
