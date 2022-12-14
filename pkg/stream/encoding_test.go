package stream

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UvariantLittleEndian(t *testing.T) {
	for _, i := range []int64{
		-1 << 63,
		-1<<63 + 1,
		-1,
		0,
		1,
		2,
		10,
		20,
		63,
		64,
		65,
		127,
		128,
		129,
		255,
		256,
		257,
		1<<63 - 1,
	} {
		assertValueEqual(t, uint64(i), UvariantSize(uint64(i)))
	}
}

func Test_UvarintLittleEndian_special_cases(t *testing.T) {
	// nil buffer
	value, size := UvarintLittleEndian(nil)
	assert.Zero(t, value)
	assert.Zero(t, size)

	// overflow
	var buf2 = []byte{
		1, 1,
		0x80, 0x80, 0x80, 0x80, 0x80,
		0x80, 0x80, 0x80, 0x80, 0x80,
	}
	value, size = UvarintLittleEndian(buf2)
	assert.Zero(t, value)
	assert.Equal(t, -11, size)
}

func assertValueEqual(t *testing.T, value uint64, size int) {
	var buf [binary.MaxVarintLen64]byte
	realSize := PutUvariantLittleEndian(buf[:], value)
	assert.Equal(t, size, realSize)

	// put into the tail
	decodedValue, decodedSize := UvarintLittleEndian(buf[:realSize])
	assert.Equal(t, decodedSize, realSize)
	assert.Equal(t, decodedValue, value)
}
