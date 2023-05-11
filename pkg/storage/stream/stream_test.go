/*
Copyright 2023 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stream

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBstreamReader(t *testing.T) {
	// Write to the bit stream.
	w := Stream{}
	for _, bit := range []Bit{true, false} {
		w.WriteBit(bit)
	}
	for nbits := 1; nbits <= 64; nbits++ {
		w.WriteBits(uint64(nbits), nbits)
	}
	for v := 1; v < 10000; v += 123 {
		w.WriteBits(uint64(v), 29)
	}

	// Read back.
	r := NewReader(w.Bytes())
	for _, bit := range []Bit{true, false} {
		v, err := r.ReadBitFast()
		if err != nil {
			v, err = r.ReadBit()
		}
		require.NoError(t, err)
		require.Equal(t, bit, v)
	}
	for nbits := uint8(1); nbits <= 64; nbits++ {
		v, err := r.ReadBitsFast(nbits)
		if err != nil {
			v, err = r.ReadBits(nbits)
		}
		require.NoError(t, err)
		require.Equal(t, uint64(nbits), v, "nbits=%d", nbits)
	}
	for v := 1; v < 10000; v += 123 {
		actual, err := r.ReadBitsFast(29)
		if err != nil {
			actual, err = r.ReadBits(29)
		}
		require.NoError(t, err)
		require.Equal(t, uint64(v), actual, "v=%d", v)
	}
}
