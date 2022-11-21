package temporary

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

/**
 * os.File 读写指针位置不区分
 * 实现内容不可修改, 读取位置自控制
 **/

var (
	ErrFileSeekInvalidWhence    error = fmt.Errorf("file seek invalid whence")
	ErrFileSeekNegativePosition error = fmt.Errorf("file seek negative position")
)

type temporaryFile struct {
	file         *os.File
	seekPosition int64
	size         int64
}

func newTemporaryFile(dir, pattern string) (tf *temporaryFile, err error) {
	var file *os.File
	if file, err = ioutil.TempFile(dir, pattern); err != nil {
		return
	}
	tf = &temporaryFile{
		file: file,
	}
	return
}

func (tf *temporaryFile) Size() int64 {
	return tf.size
}

func (tf *temporaryFile) Type() string {
	return TemporaryFile
}

func (tf *temporaryFile) Name() string {
	return tf.file.Name()
}

func (tf *temporaryFile) Bytes() []byte {
	return []byte{}
}

func (tf *temporaryFile) Sync() (err error) {
	return tf.file.Sync()
}

func (tf *temporaryFile) Write(p []byte) (n int, err error) {
	n, err = tf.file.Write(p)
	tf.size += int64(n)
	return
}

func (tf *temporaryFile) Read(p []byte) (n int, err error) {
	n, err = tf.file.ReadAt(p, tf.seekPosition)
	tf.seekPosition += int64(n)
	return
}

func (tf *temporaryFile) Seek(offset int64, whence int) (abs int64, err error) {
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = tf.seekPosition + offset
	case io.SeekEnd:
		abs = tf.size + offset
	default:
		err = ErrFileSeekInvalidWhence
		return
	}
	if abs < 0 {
		err = ErrFileSeekNegativePosition
		return
	}
	tf.seekPosition = abs
	return
}

func (tf *temporaryFile) Close() (err error) {
	os.Remove(tf.file.Name())
	return tf.file.Close()
}
