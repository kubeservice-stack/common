package temporary

import (
	"fmt"
	"io"
)

/**
 * 结合bytes.Buffer和bytes.Reader
 * bytes.Buffer数据只读一次性
 * bytes.Reader数据没有写入接口
 **/

const (
	smallBufferSize int = 64
	maxInt          int = int(^uint(0) >> 1)
)

var (
	ErrBufferTooLarge             error = fmt.Errorf("buffer too large")
	ErrBufferSeekInvalidWhence    error = fmt.Errorf("buffer seek invalid whence")
	ErrBufferSeekNegativePosition error = fmt.Errorf("buffer seek negative position")
)

type temporaryBuffer struct {
	buffer       []byte
	seekPosition int64
}

func newTemporaryBuffer() (tb *temporaryBuffer) {
	return &temporaryBuffer{}
}

func makeSlice(n int) (_ []byte, err error) {
	defer func() {
		if recover() != nil {
			err = ErrBufferTooLarge
			return
		}
	}()
	if n <= smallBufferSize {
		return make([]byte, n, smallBufferSize), nil
	}
	return make([]byte, n), nil
}

func (tb *temporaryBuffer) Size() int64 {
	return int64(len(tb.buffer))
}

func (tb *temporaryBuffer) Type() string {
	return TemporaryBuffer
}

func (tb *temporaryBuffer) Name() string {
	return ""
}

func (tb *temporaryBuffer) Bytes() []byte {
	return tb.buffer
}

func (tb *temporaryBuffer) Sync() (err error) {
	return
}

func (tb *temporaryBuffer) tryGrowByReslice(n int) (int, bool) {
	var l int = len(tb.buffer)
	if l+n <= cap(tb.buffer) {
		tb.buffer = tb.buffer[:l+n]
		return l, true
	}
	return 0, false
}

func (tb *temporaryBuffer) grow(n int) (_ int, err error) {
	var (
		tryGrow int
		tryOK   bool
	)
	if tryGrow, tryOK = tb.tryGrowByReslice(n); tryOK {
		return tryGrow, nil
	}
	if tb.buffer == nil && n <= smallBufferSize {
		if tb.buffer, err = makeSlice(n); err != nil {
			return
		}
		return 0, nil
	}
	var m, c int = len(tb.buffer), cap(tb.buffer)
	if n <= c/2-m { // 需要的长度比开辟的容量一半还少
	} else if c > maxInt-c-n {
		err = ErrBufferTooLarge
		return
	} else { // 新开容量
		var buffer []byte
		if buffer, err = makeSlice(2*c + n); err != nil {
			return
		}
		copy(buffer, tb.buffer)
		tb.buffer = buffer
	}
	tb.buffer = tb.buffer[:m+n]
	return m, nil
}

func (tb *temporaryBuffer) Write(p []byte) (n int, err error) {
	var (
		tryGrow int
		tryOK   bool
	)
	if tryGrow, tryOK = tb.tryGrowByReslice(len(p)); !tryOK {
		if tryGrow, err = tb.grow(len(p)); err != nil {
			return
		}
	}
	return copy(tb.buffer[tryGrow:], p), nil
}

func (tb *temporaryBuffer) Read(p []byte) (n int, err error) {
	if tb.seekPosition >= int64(len(tb.buffer)) {
		err = io.EOF
		return
	}
	n = copy(p, tb.buffer[tb.seekPosition:])
	tb.seekPosition += int64(n)
	return
}

func (tb *temporaryBuffer) Seek(offset int64, whence int) (abs int64, err error) {
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = tb.seekPosition + offset
	case io.SeekEnd:
		abs = int64(len(tb.buffer)) + offset
	default:
		err = ErrBufferSeekInvalidWhence
		return
	}
	if abs < 0 {
		err = ErrBufferSeekNegativePosition
		return
	}
	tb.seekPosition = abs
	return
}

func (tb *temporaryBuffer) Close() (err error) {
	return
}
