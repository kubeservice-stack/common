/*
Copyright 2022 The KubeService-Stack Authors.

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

package bufioutil

import (
	"errors"
)

var errOutOfRange = errors.New("index out of range")

type Buffer struct {
	buf    []byte // encode/decode byte stream
	index  int    // read point
	length int    // length of buf
}

// NewBuffer allocates a new Buffer and initializes its internal data to
// the contents of the argument slice.
func NewBuffer(e []byte) *Buffer {
	return &Buffer{buf: e, length: len(e)}
}

// SetBuf replaces the internal buffer with the slice,
// ready for unmarshalling the contents of the slice.
func (p *Buffer) SetBuf(s []byte) {
	p.buf = s
	p.index = 0
	p.length = len(s)
}

func (p *Buffer) SetIdx(idx int) {
	p.index = idx
}

func (p *Buffer) GetByte() (b byte, err error) {
	if p.index >= p.length {
		err = errOutOfRange
		return
	}
	b = p.buf[p.index]
	p.index++
	return
}
