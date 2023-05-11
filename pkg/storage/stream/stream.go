// Copyright (c) 2015,2016 Damian Gryski <damian@gryski.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package stream

import (
	"encoding/binary"
	"io"
)

// stream is a stream of bits.
type Stream struct {
	S     []byte // the data stream
	Count uint8  // how many bits are valid in current byte
}

func (b *Stream) Bytes() []byte {
	return b.S
}

// reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
func (b *Stream) Reset() {
	b.S = b.S[:0]
	b.Count = 0
}

type Bit bool

const (
	Zero Bit = false
	One  Bit = true
)

func (b *Stream) WriteBit(bit Bit) {
	if b.Count == 0 {
		b.S = append(b.S, 0)
		b.Count = 8
	}

	i := len(b.S) - 1

	if bit {
		b.S[i] |= 1 << (b.Count - 1)
	}

	b.Count--
}

func (b *Stream) WriteByte(byt byte) {
	if b.Count == 0 {
		b.S = append(b.S, 0)
		b.Count = 8
	}

	i := len(b.S) - 1

	// fill up b.b with b.Count bits from byt
	b.S[i] |= byt >> (8 - b.Count)

	b.S = append(b.S, 0)
	i++
	b.S[i] = byt << b.Count
}

func (b *Stream) WriteBits(u uint64, nbits int) {
	u <<= (64 - uint(nbits))
	for nbits >= 8 {
		byt := byte(u >> 56)
		b.WriteByte(byt)
		u <<= 8
		nbits -= 8
	}

	for nbits > 0 {
		b.WriteBit((u >> 63) == 1)
		u <<= 1
		nbits--
	}
}

type StreamReader struct {
	stream       []byte
	streamOffset int // The offset from which read the next byte from the stream.

	buffer uint64 // The current buffer, filled from the stream, containing up to 8 bytes from which read bits.
	valid  uint8  // The number of bits valid to read (from left) in the current buffer.
}

func NewReader(b []byte) StreamReader {
	return StreamReader{
		stream: b,
	}
}

func (b *StreamReader) ReadBit() (Bit, error) {
	if b.valid == 0 {
		if !b.LoadNextBuffer(1) {
			return false, io.EOF
		}
	}

	return b.ReadBitFast()
}

// readBitFast is like readBit but can return io.EOF if the internal buffer is empty.
// If it returns io.EOF, the caller should retry reading bits calling readBit().
// This function must be kept small and a leaf in order to help the compiler inlining it
// and further improve performances.
func (b *StreamReader) ReadBitFast() (Bit, error) {
	if b.valid == 0 {
		return false, io.EOF
	}

	b.valid--
	bitmask := uint64(1) << b.valid
	return (b.buffer & bitmask) != 0, nil
}

func (b *StreamReader) ReadBits(nbits uint8) (uint64, error) {
	if b.valid == 0 {
		if !b.LoadNextBuffer(nbits) {
			return 0, io.EOF
		}
	}

	if nbits <= b.valid {
		return b.ReadBitsFast(nbits)
	}

	// We have to read all remaining valid bits from the current buffer and a part from the next one.
	bitmask := (uint64(1) << b.valid) - 1
	nbits -= b.valid
	v := (b.buffer & bitmask) << nbits
	b.valid = 0

	if !b.LoadNextBuffer(nbits) {
		return 0, io.EOF
	}

	bitmask = (uint64(1) << nbits) - 1
	v = v | ((b.buffer >> (b.valid - nbits)) & bitmask)
	b.valid -= nbits

	return v, nil
}

// readBitsFast is like readBits but can return io.EOF if the internal buffer is empty.
// If it returns io.EOF, the caller should retry reading bits calling readBits().
// This function must be kept small and a leaf in order to help the compiler inlining it
// and further improve performances.
func (b *StreamReader) ReadBitsFast(nbits uint8) (uint64, error) {
	if nbits > b.valid {
		return 0, io.EOF
	}

	bitmask := (uint64(1) << nbits) - 1
	b.valid -= nbits

	return (b.buffer >> b.valid) & bitmask, nil
}

func (b *StreamReader) ReadByte() (byte, error) {
	v, err := b.ReadBits(8)
	if err != nil {
		return 0, err
	}
	return byte(v), nil
}

// loadNextBuffer loads the next bytes from the stream into the internal buffer.
// The input nbits is the minimum number of bits that must be read, but the implementation
// can read more (if possible) to improve performances.
func (b *StreamReader) LoadNextBuffer(nbits uint8) bool {
	if b.streamOffset >= len(b.stream) {
		return false
	}

	// Handle the case there are more then 8 bytes in the buffer (most common case)
	// in a optimized way. It's guaranteed that this branch will never read from the
	// very last byte of the stream (which suffers race conditions due to concurrent
	// writes).
	if b.streamOffset+8 < len(b.stream) {
		b.buffer = binary.BigEndian.Uint64(b.stream[b.streamOffset:])
		b.streamOffset += 8
		b.valid = 64
		return true
	}

	// We're here if the are 8 or less bytes left in the stream. Since this reader needs
	// to handle race conditions with concurrent writes happening on the very last byte
	// we make sure to never over more than the minimum requested bits (rounded up to
	// the next byte). The following code is slower but called less frequently.
	nbytes := int((nbits / 8) + 1)
	if b.streamOffset+nbytes > len(b.stream) {
		nbytes = len(b.stream) - b.streamOffset
	}

	buffer := uint64(0)
	for i := 0; i < nbytes; i++ {
		buffer = buffer | (uint64(b.stream[b.streamOffset+i]) << uint(8*(nbytes-i-1)))
	}

	b.buffer = buffer
	b.streamOffset += nbytes
	b.valid = uint8(nbytes * 8)

	return true
}
